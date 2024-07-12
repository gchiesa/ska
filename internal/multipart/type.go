package multipart

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/part"
	"os"
)

type Multipart struct {
	id              string
	refFileUri      string
	contentOriginal []byte
	contentCompiled []byte
	parts           []part.Part
	workindDir      string
	log             *log.Entry
}

func NewMultipartFromFile(filePath string, id string) (*Multipart, error) {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, part.MultipartError)
	}
	if ok := isValidContent(fileData); !ok {
		return nil, fmt.Errorf("invalid contentOriginal for partial container: %w", part.InvalidContent)
	}

	logCtx := log.WithFields(log.Fields{
		"pkg": "multipart",
	})

	return &Multipart{
		id:              id,
		refFileUri:      filePath,
		contentOriginal: fileData,
		log:             logCtx,
	}, nil
}

func (mp *Multipart) Id() string {
	return mp.id
}

func (mp *Multipart) Parts() []part.Part {
	return mp.parts
}
