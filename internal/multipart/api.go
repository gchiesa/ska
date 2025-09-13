package multipart

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gchiesa/ska/internal/part"
)

func (mp *Multipart) ParseParts() error {
	if !isValidContent(mp.contentOriginal) {
		return fmt.Errorf("invalid contentOriginal for partial container (file: %s): %w", mp.refFileURI, part.ErrInvalidContent)
	}
	reSection := getPartialsRegexp()

	matches := reSection.FindAllSubmatch(mp.contentOriginal, -1)

	for _, match := range matches {
		fullHeader := strings.TrimSpace(string(match[1]))
		cleanMatch := string(match[2])
		sectionName, adoptType, adoptArg := parseHeader(fullHeader)
		partial := part.NewPart(mp.refFileURI, sectionName).WithContent([]byte(cleanMatch))
		if adoptType != "" {
			partial = partial.SetAdopt(adoptType, adoptArg)
		}
		mp.parts = append(mp.parts, *partial)
	}
	return nil
}

// parseHeader splits the header "<section> [+ directive:arg]" into components
func parseHeader(header string) (sectionName, adoptType, adoptArg string) {
	// default: no directive, whole header is sectionName
	sectionName = strings.TrimSpace(header)
	adoptType, adoptArg = "", ""
	if strings.Contains(header, "+") {
		parts := strings.SplitN(header, "+", 2)
		sectionName = strings.TrimSpace(parts[0])
		directive := strings.TrimSpace(parts[1])
		// normalize
		directive = strings.TrimSpace(directive)
		// supported directives
		supported := []string{"ska-inject-after:", "ska-inject-before:", "ska-replace-match:"}
		for _, pref := range supported {
			if strings.HasPrefix(directive, pref) {
				adoptType = strings.TrimSuffix(pref, ":")
				adoptArg = strings.TrimSpace(strings.TrimPrefix(directive, pref))
				// strip optional angle brackets
				if strings.HasPrefix(adoptArg, "<") && strings.HasSuffix(adoptArg, ">") {
					adoptArg = strings.TrimSuffix(strings.TrimPrefix(adoptArg, "<"), ">")
				}
				break
			}
		}
	}
	return
}

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
