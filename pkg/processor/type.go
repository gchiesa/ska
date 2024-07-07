package processor

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/pkg/multipart"
)

type FileTreeProcessor struct {
	sourcePath          string
	destinationPathRoot string
	workingDir          string
	multiparts          []*multipart.Multipart
	log                 *log.Entry
}

type TreeRendererOptions struct {
}

func NewFileTreeProcessor(sourcePath, destinationPathRoot string, options TreeRendererOptions) *FileTreeProcessor {
	logCtx := log.WithFields(log.Fields{
		"pkg": "processor",
	})
	return &FileTreeProcessor{
		sourcePath:          sourcePath,
		destinationPathRoot: destinationPathRoot,
		workingDir:          "",
		log:                 logCtx,
	}
}

func (tp *FileTreeProcessor) WorkingDir() string {
	return tp.workingDir
}
