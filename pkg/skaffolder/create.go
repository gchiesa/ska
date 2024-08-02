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

func NewSkaCreate(templateURI, destinationPath string, variables map[string]string) *SkaCreate {
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
	// blueprint provider
	blueprintProvider, err := contentprovider.ByURI(s.TemplateURI)
	if err != nil {
		return err
	}

	defer func(templateProvider contentprovider.RemoteContentProvider) {
		_ = templateProvider.Cleanup()
	}(blueprintProvider)

	// configservice
	localConfig := configuration.NewLocalConfigService()

	if err = blueprintProvider.DownloadContent(); err != nil { //nolint:govet //not a bit deal
		return err
	}

	// load the config for upstream blueprint
	upstreamConfig, err := configuration.NewUpstreamConfigService().LoadFromPath(blueprintProvider.WorkingDir())
	if err != nil {
		return err
	}

	fileTreeProcessor := processor.NewFileTreeProcessor(blueprintProvider.WorkingDir(), s.DestinationPath,
		processor.WithErrorOnMissingKey(true),
		processor.WithSourceIgnorePaths(upstreamConfig.GetIgnorePaths()),
		processor.WithDestinationIgnorePaths(localConfig.GetIgnorePaths()))

	defer func(fileTreeProcessor *processor.FileTreeProcessor) {
		_ = fileTreeProcessor.Cleanup()
	}(fileTreeProcessor)

	var interactiveServiceVariables map[string]string

	interactiveService := tui.NewSkaInteractiveService(
		fmt.Sprintf("Variables for blueprint: %s", s.TemplateURI),
		upstreamConfig.GetInputs())

	// check if interactive mode is required
	if interactiveService.ShouldRun() {
		if err = interactiveService.Run(); err != nil {
			return err
		}
		// retrieve the collected variables
		interactiveServiceVariables = interactiveService.Variables()
	}

	// variables for templating
	vars := mapStringToMapInterface(s.Variables)

	// merge the known variables with overrides
	for k, v := range mapStringToMapInterface(interactiveServiceVariables) {
		vars[k] = v
	}

	if err = fileTreeProcessor.Render(vars); err != nil { //nolint:govet //not a bit deal
		return err
	}

	// save the config
	err = localConfig.
		WithVariables(vars).
		WithBlueprintUpstream(blueprintProvider.RemoteURI()).
		WriteConfig(s.DestinationPath)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"method": "Create", "path": s.DestinationPath}).Infof("template created under destination path: %s", s.DestinationPath)
	return nil
}
