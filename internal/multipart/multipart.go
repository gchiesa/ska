package multipart

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

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
		re := buildReplaceRegexp(p.ID())
		if re.Match(dataContent) {
			// standard managed replace
			dataContent = re.ReplaceAll(dataContent, []byte(`${1}:${2}`+"\n"+string(compiledPartial)+`${4}`))
			continue
		}
		// if no match and adopt directive is present, perform injection/replacement
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
		if after {
			return append([]byte(payload), base...)
		}
		return append([]byte(payload), base...)
	}
	if anchor == "@end" {
		return append(base, []byte(payload)...)
	}
	// treat as regex
	re := regexp.MustCompile(anchor)
	loc := re.FindIndex(base)
	if loc == nil {
		// not found, append at end
		return append(base, []byte(payload)...)
	}
	if after {
		return append(append(base[:loc[1]], []byte(payload)...), base[loc[1]:]...)
	}
	// before
	return append(append(base[:loc[0]], []byte(payload)...), base[loc[0]:]...)
}

// replaceMatch replaces the first regex match (or its single capture group if present) with the payload
func replaceMatch(base []byte, regex string, payload string) []byte {
	re := regexp.MustCompile(regex)
	idxs := re.FindSubmatchIndex(base)
	if idxs == nil {
		return base
	}
	// if there's exactly one capture group, replace the group range; else replace whole match
	if len(idxs) == 4 {
		start, end := idxs[2], idxs[3]
		return append(append(base[:start], []byte(payload)...), base[end:]...)
	}
	start, end := idxs[0], idxs[1]
	return append(append(base[:start], []byte(payload)...), base[end:]...)
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
