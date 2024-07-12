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
		return fmt.Errorf("invalid contentOriginal for partial container (file: %s): %w", mp.refFileUri, part.InvalidContent)
	}
	reSection := getPartialsRegexp()

	matches := reSection.FindAllSubmatch(mp.contentOriginal, -1)

	for _, match := range matches {
		cleanId := strings.TrimSpace(string(match[1]))
		cleanMatch := string(match[2])
		partial := part.NewPart(mp.refFileUri, cleanId).WithContent([]byte(cleanMatch))
		mp.parts = append(mp.parts, *partial)
	}
	return nil
}

func (mp *Multipart) PartsToFiles() ([]string, error) {
	workDir := filepath.Dir(mp.refFileUri)
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
		compiledPartial, err := os.ReadFile(filepath.Join(filepath.Dir(mp.refFileUri), part.RefFileBasename()))
		if err != nil {
			mp.log.Errorf("error reading file. %s: %v", part.RefFileBasename(), err)
		}
		re := buildReplaceRegexp(part.Id())
		dataContent = re.ReplaceAll(dataContent, []byte(`${1}:${2}`+string("\n"+string(compiledPartial))+`${4}`))
	}
	mp.contentCompiled = dataContent
	return dataContent
}

func (mp *Multipart) CompileToFile(filePath string, forceRecompile bool) error {
	dataContent := mp.Compile(forceRecompile)
	return os.WriteFile(filePath, dataContent, 0644)
}
