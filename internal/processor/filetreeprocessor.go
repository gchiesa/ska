package processor

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/templateservice"
	"os"
)

func NewFileTreeProcessor(sourcePath, destinationPathRoot string, options ...func(*FileTreeProcessor)) *FileTreeProcessor {
	logCtx := log.WithFields(log.Fields{
		"pkg": "processor",
	})

	tp := &FileTreeProcessor{
		sourcePath:          sourcePath,
		destinationPathRoot: destinationPathRoot,
		workingDir:          "",
		templateService:     templateservice.NewSprigTemplate("default"),
		log:                 logCtx,
	}
	// configure options
	for _, opt := range options {
		opt(tp)
	}
	return tp
}

func (tp *FileTreeProcessor) WorkingDir() string {
	return tp.workingDir
}

func (tp *FileTreeProcessor) Cleanup() error {
	tp.log.WithFields(log.Fields{"workingDir": tp.workingDir}).Debug("removing working dir.")
	return os.RemoveAll(tp.workingDir)
}

func (tp *FileTreeProcessor) Render(withVariables map[string]interface{}) error {
	if tp.workingDir == "" {
		var err error
		tp.workingDir, err = os.MkdirTemp("", "swansonRenderer")
		if err != nil {
			return err
		}
	}

	if err := tp.buildStagingFileTree(withVariables); err != nil {
		return err
	}

	// WAVE 2 - decompose the swanson managed partials
	// create a set of partials that are related to the files in the staging directory
	if err := tp.loadMultiparts(); err != nil {
		return err
	}

	// WAVE 3 - expand template
	// render all the templates, but if a partial exists for a file then expands only the partials
	if err := tp.renderStagingFileTree(withVariables); err != nil {
		return err
	}

	// WAVE 4 - copy to destination the staging directory
	// copy the staging directory to the destination with the process
	// for each file (non-swanson) copy the file first
	// then replace the partials with the expanded content
	// **IF the file mustBeSkipped then skip, otherwise copy
	// **IF the file already exists in the destination then
	// only replace the partials with the expanded content
	if err := tp.updateDestinationFileTree(); err != nil {
		return err
	}
	return nil
}

func WithTemplateService(templateService *templateservice.SprigTemplate) func(tp *FileTreeProcessor) {
	return func(tp *FileTreeProcessor) {
		tp.templateService = templateService
	}
}

func WithErrorOnMissingKey(errorOnMissingKey bool) func(tp *FileTreeProcessor) {
	return func(tp *FileTreeProcessor) {
		tp.templateService.WithErrorOnMissingKey(errorOnMissingKey)
	}
}