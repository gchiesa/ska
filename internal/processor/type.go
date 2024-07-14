package processor

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/multipart"
	"os"
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

func NewFileTreeProcessor(sourcePath, destinationPathRoot string) *FileTreeProcessor {
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

func (tp *FileTreeProcessor) Cleanup() error {
	tp.log.WithFields(log.Fields{"workingDir": tp.workingDir}).Debug("removing working dir.")
	return os.RemoveAll(tp.workingDir)
}
