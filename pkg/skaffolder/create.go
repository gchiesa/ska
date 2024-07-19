package skaffolder

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/configuration"
	"github.com/gchiesa/ska/internal/contentprovider"
	"github.com/gchiesa/ska/internal/processor"
	"github.com/gchiesa/ska/internal/tui"
)

type SkaCreate struct {
	TemplateURI     string
	DestinationPath string
	Variables       map[string]string
	Log             *log.Entry
}

func NewSkaCreate(templateURI string, destinationPath string, variables map[string]string) *SkaCreate {
	logCtx := log.WithFields(log.Fields{
		"pkg": "skaffolder",
	})
	return &SkaCreate{
		TemplateURI:     templateURI,
		DestinationPath: destinationPath,
		Variables:       variables,
		Log:             logCtx,
	}
}

func (s *SkaCreate) Create() error {
	templateProvider, err := contentprovider.ByURI(s.TemplateURI)
	defer func(templateProvider contentprovider.RemoteContentProvider) {
		_ = templateProvider.Cleanup()
	}(templateProvider)

	if err != nil {
		return err
	}

	configService := configuration.NewConfigService()

	if err = templateProvider.DownloadContent(); err != nil { //nolint:govet //not a bit deal
		return err
	}

	fileTreeProcessor := processor.NewFileTreeProcessor(templateProvider.WorkingDir(), s.DestinationPath)
	defer func(fileTreeProcessor *processor.FileTreeProcessor) {
		_ = fileTreeProcessor.Cleanup()
	}(fileTreeProcessor)

	// check if interactive mode is required
	var interactiveServiceVariables map[string]string

	interactiveService := tui.NewSkaInteractiveService(
		templateProvider.WorkingDir(),
		fmt.Sprintf("Variables for blueprint: %s", s.TemplateURI))
	if interactiveService.ShouldRun() {
		if err = interactiveService.Run(); err != nil {
			return err
		}
		interactiveServiceVariables = interactiveService.Variables()
	}

	vars := mapStringToMapInterface(s.Variables)

	// merge the known variables with overrides
	for k, v := range mapStringToMapInterface(interactiveServiceVariables) {
		vars[k] = v
	}

	if err = fileTreeProcessor.Render(vars); err != nil { //nolint:govet //not a bit deal
		return err
	}

	// save the config
	err = configService.
		WithVariables(vars).
		WithBlueprintUpstream(templateProvider.RemoteURI()).
		WriteConfig(s.DestinationPath)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"method": "Create", "path": s.DestinationPath}).Infof("template created under destination path: %s", s.DestinationPath)
	return nil
}
