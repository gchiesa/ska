package skaffolder

import (
	"fmt"
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/configuration"
	"github.com/gchiesa/ska/internal/contentprovider"
	"github.com/gchiesa/ska/internal/filetreeprocessor"
	"github.com/gchiesa/ska/internal/stringprocessor"
	"github.com/gchiesa/ska/internal/templateprovider"
	"github.com/gchiesa/ska/internal/tui"
)

type SkaCreate struct {
	TemplateURI     string
	DestinationPath string
	NamedConfig     string
	Variables       map[string]string
	Options         *SkaOptions
	Log             *log.Entry
}

type SkaOptions struct {
	NonInteractive bool
	Engine         templateprovider.TemplateType // jinja or sprig
}

func NewSkaCreate(templateURI, destinationPath, namedConfig string, variables map[string]string, options SkaOptions) *SkaCreate {
	logCtx := log.WithFields(log.Fields{
		"pkg": "skaffolder",
	})
	return &SkaCreate{
		TemplateURI:     templateURI,
		DestinationPath: destinationPath,
		NamedConfig:     namedConfig,
		Variables:       variables,
		Options:         &options,
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
	localConfig := configuration.NewLocalConfigService(s.NamedConfig)

	// check if localconfig already exist, if yes we fail
	if localConfig.ConfigExists(s.DestinationPath) {
		log.Infof("Default config or specified named config already exists, please use a different name for the configuration")
		return fmt.Errorf("configuration already exists")
	}

	if err = blueprintProvider.DownloadContent(); err != nil { //nolint:govet //not a bit deal
		return err
	}

	// load the config for upstream blueprint
	upstreamConfig, err := configuration.NewUpstreamConfigService().LoadFromPath(blueprintProvider.WorkingDir())
	if err != nil {
		return err
	}

	// template engine
	var templateService templateprovider.TemplateService
	switch s.Options.Engine {
	case templateprovider.SprigTemplateType:
		templateService = templateprovider.NewSprigTemplate(s.TemplateURI)
	case templateprovider.JinjaTemplateType:
		templateService = templateprovider.NewJinjaTemplate(s.TemplateURI)
	default:
		return fmt.Errorf("unknown template engine")
	}

	fileTreeProcessor := filetreeprocessor.NewFileTreeProcessor(blueprintProvider.WorkingDir(), s.DestinationPath,
		filetreeprocessor.WithTemplateService(templateService),
		filetreeprocessor.WithSourceIgnorePaths(upstreamConfig.UpstreamIgnorePaths()),
		filetreeprocessor.WithDestinationIgnorePaths(localConfig.IgnorePaths()))

	defer func(fileTreeProcessor *filetreeprocessor.FileTreeProcessor) {
		_ = fileTreeProcessor.Cleanup()
	}(fileTreeProcessor)

	var interactiveServiceVariables map[string]string

	interactiveService := tui.NewSkaInteractiveService(
		fmt.Sprintf("Variables for blueprint: %s", s.TemplateURI),
		upstreamConfig.GetInputs())

	// check if interactive mode is required
	if !s.Options.NonInteractive && interactiveService.ShouldRun() {
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

	// render the ignore entries in the upstream configuration
	sp := stringprocessor.NewStringProcessor(stringprocessor.WithTemplateService(templateService))
	skaConfigIgnorePaths, err := sp.RenderSliceOfStrings(upstreamConfig.SkaConfigIgnorePaths(), vars)
	if err != nil {
		return err
	}

	// save the config
	err = localConfig.
		WithVariables(vars).
		WithBlueprintUpstream(blueprintProvider.RemoteURI()).
		WithIgnorePaths(skaConfigIgnorePaths).
		WriteConfig(s.DestinationPath)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"method": "Create", "path": s.DestinationPath, "blueprintUri": blueprintProvider.RemoteURI()}).Info("blueprint expanded under destination path.")
	return nil
}
