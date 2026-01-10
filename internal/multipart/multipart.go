package multipart

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/part"
)

type Multipart struct {
	id              string
	refFileURI      string
	contentOriginal []byte
	contentCompiled []byte
	parts           []part.Part
	log             *log.Entry
}

func NewMultipartFromFile(filePath, id string) (*Multipart, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, part.ErrMultipartError)
	}
	if ok := isValidContent(fileData); !ok {
		return nil, fmt.Errorf("invalid contentOriginal for partial container: %w", part.ErrInvalidContent)
	}

	logCtx := log.WithFields(log.Fields{
		"pkg": "multipart",
	})

	return &Multipart{
		id:              id,
		refFileURI:      filePath,
		contentOriginal: fileData,
		log:             logCtx,
	}, nil
}

func (mp *Multipart) ID() string {
	return mp.id
}

func (mp *Multipart) Parts() []part.Part {
	return mp.parts
}

func (mp *Multipart) ParseParts() error {
	if !isValidContent(mp.contentOriginal) {
		return fmt.Errorf("invalid contentOriginal for partial container (file: %s): %w", mp.refFileURI, part.ErrInvalidContent)
	}
	parts, err := part.ParseParts(mp.contentOriginal, mp.refFileURI)
	if err != nil {
		return err
	}
	mp.parts = parts
	return nil
}

// PartsToFiles generates individual files for each part in the current multipart and returns a list of the generated file paths.
func (mp *Multipart) PartsToFiles() ([]string, error) {
	workDir := filepath.Dir(mp.refFileURI)
	fileList := make([]string, 0)
	for _, part := range mp.parts {
		partialPath := filepath.Join(workDir, part.RefFileBasename())
		err := part.SetFilePath(partialPath).CreateFile()
		if err != nil {
			return nil, fmt.Errorf("error writing file: %w", err)
		}
		fileList = append(fileList, part.RefFilePath())
	}
	return fileList, nil
}

func (mp *Multipart) HasParts() bool {
	return len(mp.parts) > 0
}

// Compile processes the original content by injecting or replacing parts, optionally forcing re-compilation.
func (mp *Multipart) Compile(forceRecompile bool) []byte {
	if mp.contentCompiled != nil && !forceRecompile {
		return mp.contentCompiled
	}

	dataContent := mp.contentOriginal
	for _, part := range mp.parts {
		// find and replace the part in the content
		compiledPartial, err := os.ReadFile(filepath.Join(filepath.Dir(mp.refFileURI), part.RefFileBasename()))
		if err != nil {
			mp.log.Errorf("error reading file. %s: %v", part.RefFileBasename(), err)
		}
		re := buildReplaceRegexp(part.ID())
		dataContent = re.ReplaceAll(dataContent, []byte(`${1}:${2}`+string("\n"+string(compiledPartial))+`${4}`)) //nolint:unconvert //keeping for better readability
	}
	mp.contentCompiled = dataContent
	return dataContent
}

func (mp *Multipart) CompileToFile(filePath string, forceRecompile bool) error { // noling:gosec
	// if local file exists, we need to use that as content
	if _, err := os.Stat(filePath); err == nil {
		dataContent, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		mp.contentOriginal = dataContent
	}
	// adopt-aware compile: start from current contentOriginal and process parts
	dataContent := mp.contentOriginal
	for _, p := range mp.parts {
		// read compiled partial content
		compiledPartial, err := os.ReadFile(filepath.Join(filepath.Dir(mp.refFileURI), p.RefFileBasename()))
		if err != nil {
			mp.log.Errorf("error reading file. %s: %v", p.RefFileBasename(), err)
			continue
		}
		// First check if a managed block for this section already exists (for idempotency)
		re := buildReplaceRegexp(p.ID())
		if re.Match(dataContent) {
			// Replace existing managed block with updated content
			dataContent = re.ReplaceAll(dataContent, []byte(`${1}:${2}`+"\n"+string(compiledPartial)+`${4}`))
			continue
		}
		// No existing block - if a special adopt type is requested, process it
		if p.AdoptType() != "" {
			switch p.AdoptType() {
			case "ska-inject-after":
				dataContent = injectRelative(dataContent, p.AdoptArg(), string(buildManagedBlock(p.ID(), compiledPartial)), true)
			case "ska-inject-before":
				dataContent = injectRelative(dataContent, p.AdoptArg(), string(buildManagedBlock(p.ID(), compiledPartial)), false)
			case "ska-replace-match":
				dataContent = replaceMatch(dataContent, p.AdoptArg(), string(buildManagedBlock(p.ID(), compiledPartial)))
			default:
				// no directive, nothing to inject
			}
			continue
		}
	}
	mp.contentCompiled = dataContent
	return os.WriteFile(filePath, dataContent, 0o644) // noling:gosec
}

// buildManagedBlock creates a canonical managed block without directive in the header for idempotency
func buildManagedBlock(sectionName string, content []byte) []byte {
	// use '#' marker by default
	return []byte("# " + part.DelimiterStart + ":" + sectionName + "\n" + string(content) + "# " + part.DelimiterEnd + "\n")
}

// injectRelative injects payload before/after an anchor (@start, @end, or regex). If anchor not found (regex), appends at end.
func injectRelative(base []byte, anchor string, payload string, after bool) []byte {
	if anchor == "@start" {
		result := make([]byte, 0, len(payload)+len(base))
		result = append(result, payload...)
		result = append(result, base...)
		return result
	}
	if anchor == "@end" {
		result := make([]byte, 0, len(base)+len(payload))
		result = append(result, base...)
		result = append(result, payload...)
		return result
	}
	// treat as regex
	re := regexp.MustCompile(anchor)
	loc := re.FindIndex(base)
	if loc == nil {
		// not found, append at end
		result := make([]byte, 0, len(base)+len(payload))
		result = append(result, base...)
		result = append(result, payload...)
		return result
	}
	// Build result in a new slice to avoid overwriting base's underlying array
	var pos int
	if after {
		pos = loc[1]
	} else {
		pos = loc[0]
	}
	result := make([]byte, 0, len(base)+len(payload))
	result = append(result, base[:pos]...)
	result = append(result, payload...)
	result = append(result, base[pos:]...)
	return result
}

// replaceMatch replaces the first regex match (or its single capture group if present) with the payload.
// If the regex contains ^ or $ anchors but no (?m) flag, multiline mode is automatically enabled
// so that ^ and $ match line boundaries instead of just the start/end of the entire string.
func replaceMatch(base []byte, regex string, payload string) []byte {
	// Enable multiline mode if regex uses ^ or $ anchors but doesn't already have (?m)
	if (strings.Contains(regex, "^") || strings.Contains(regex, "$")) && !strings.Contains(regex, "(?m)") {
		regex = "(?m)" + regex
	}
	re := regexp.MustCompile(regex)
	idxs := re.FindSubmatchIndex(base)
	if idxs == nil {
		return base
	}
	// if there's exactly one capture group, replace the group range; else replace whole match
	var start, end int
	if len(idxs) == 4 {
		start, end = idxs[2], idxs[3]
	} else {
		start, end = idxs[0], idxs[1]
	}
	// Build result in a new slice to avoid overwriting base's underlying array
	result := make([]byte, 0, start+len(payload)+len(base)-end)
	result = append(result, base[:start]...)
	result = append(result, payload...)
	result = append(result, base[end:]...)
	return result
}

// validatePartialsInContent return error if the contentOriginal contains invalid partials
func isValidContent(content []byte) bool {
	reStart := regexp.MustCompile(part.DelimiterStart)
	reEnd := regexp.MustCompile(part.DelimiterEnd)
	reSection := regexp.MustCompile("(?s)" + part.DelimiterStart + "(.*?)" + part.DelimiterEnd) //nolint:goconst //keeping for better readability

	// is if the number of section matches with the number of starts and ends valid
	startEntries := reStart.FindAll(content, -1)
	endEntries := reEnd.FindAll(content, -1)
	sections := reSection.FindAll(content, -1)
	return len(startEntries) == len(sections) && len(endEntries) == len(sections)
}
