package skaffolder //nolint:typecheck

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

type SkaUpdate struct {
	BaseURI     string
	NamedConfig string
	Variables   map[string]string
	Options     *SkaOptions
	Log         *log.Entry
}

func NewSkaUpdate(baseURI, namedConfig string, variables map[string]string, options SkaOptions) *SkaUpdate {
	logCtx := log.WithFields(log.Fields{
		"pkg": "skaffolder",
	})
	return &SkaUpdate{
		BaseURI:     baseURI,
		NamedConfig: namedConfig,
		Variables:   variables,
		Options:     &options,
		Log:         logCtx,
	}
}

func (s *SkaUpdate) Update() error {
	localConfig := configuration.NewLocalConfigService(s.NamedConfig)

	// read the config from the folder
	if err := localConfig.ReadValidConfig(s.BaseURI); err != nil {
		return err
	}

	// allocate the template based on the configured upstream
	blueprintProvider, err := contentprovider.ByURI(localConfig.BlueprintUpstream())
	if err != nil {
		return err
	}

	defer func(templateProvider contentprovider.RemoteContentProvider) {
		_ = templateProvider.Cleanup()
	}(blueprintProvider)

	if err := blueprintProvider.DownloadContent(); err != nil { //nolint:govet //not a bit deal
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
		templateService = templateprovider.NewSprigTemplate(s.BaseURI)
	case templateprovider.JinjaTemplateType:
		templateService = templateprovider.NewJinjaTemplate(s.BaseURI)
	default:
		return fmt.Errorf("unknown template engine")
	}

	fileTreeProcessor := filetreeprocessor.NewFileTreeProcessor(blueprintProvider.WorkingDir(), s.BaseURI,
		filetreeprocessor.WithTemplateService(templateService),
		filetreeprocessor.WithSourceIgnorePaths(upstreamConfig.UpstreamIgnorePaths()),
		filetreeprocessor.WithDestinationIgnorePaths(localConfig.IgnorePaths()))

	defer func(fileTreeProcessor *filetreeprocessor.FileTreeProcessor) {
		_ = fileTreeProcessor.Cleanup()
	}(fileTreeProcessor)

	// merge the known variables from the yaml with overrides from command line
	vars := localConfig.Variables()
	for k, v := range mapStringToMapInterface(s.Variables) {
		vars[k] = v
	}

	var interactiveServiceVariables map[string]string

	interactiveService := tui.NewSkaInteractiveService(
		fmt.Sprintf("Variables for blueprint: %s", localConfig.BlueprintUpstream()),
		upstreamConfig.GetInputs())

	// check if interactive mode is required
	if !s.Options.NonInteractive && interactiveService.ShouldRun() {
		// overrides the variables from remote service with already saved variables
		interactiveService.SetDefaults(mapInterfaceToString(vars))

		if err = interactiveService.Run(); err != nil {
			return err
		}
		// retrieve the collected variables
		interactiveServiceVariables = interactiveService.Variables()
	}

	// update the variables with the interactive variables
	for k, v := range mapStringToMapInterface(interactiveServiceVariables) {
		vars[k] = v
	}

	// render
	if err := fileTreeProcessor.Render(vars); err != nil { //nolint:govet //not a bit deal
		return err
	}

	// render the ignore entries in the upstream configuration
	var skaConfigIgnorePaths []string
	sp := stringprocessor.NewStringProcessor(stringprocessor.WithTemplateService(templateService))
	skaConfigIgnorePaths, err = sp.RenderSliceOfStrings(upstreamConfig.SkaConfigIgnorePaths(), vars)
	if err != nil {
		return err
	}

	// save the config
	err = localConfig.
		WithVariables(vars).
		WithBlueprintUpstream(blueprintProvider.RemoteURI()).
		WithExtendIgnorePaths(localConfig.IgnorePaths()).
		WithExtendIgnorePaths(skaConfigIgnorePaths).
		WriteConfig(s.BaseURI)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"method": "Update", "path": s.BaseURI, "blueprintURI": localConfig.BlueprintUpstream()}).Info("local path updated with blueprint.")
	return nil
}
