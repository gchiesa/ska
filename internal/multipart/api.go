package multipart

import (
	"fmt"
	"github.com/gchiesa/ska/internal/part"
	"os"
	"path/filepath"
	"strings"
)

func (mp *Multipart) ParseParts() error {
	if !isValidContent(mp.contentOriginal) {
		return fmt.Errorf("invalid contentOriginal for partial container (file: %s): %w", mp.refFileURI, part.ErrInvalidContent)
	}
	reSection := getPartialsRegexp()

	matches := reSection.FindAllSubmatch(mp.contentOriginal, -1)

	for _, match := range matches {
		cleanID := strings.TrimSpace(string(match[1]))
		cleanMatch := string(match[2])
		partial := part.NewPart(mp.refFileURI, cleanID).WithContent([]byte(cleanMatch))
		mp.parts = append(mp.parts, *partial)
	}
	return nil
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
	// otherwise we use the multipart content
	dataContent := mp.Compile(forceRecompile)
	return os.WriteFile(filePath, dataContent, 0o644) // noling:gosec
}
