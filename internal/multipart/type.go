package multipart

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/part"
	"os"
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
