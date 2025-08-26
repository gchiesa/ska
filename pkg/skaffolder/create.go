// Package skaffolder exposes the high-level API to create and update a destination
// directory from a remote or local blueprint using different template engines.
package skaffolder

import (
	"fmt"

	"github.com/gchiesa/ska/pkg/templateprovider"

	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/contentprovider"
	"github.com/gchiesa/ska/internal/filetreeprocessor"
	"github.com/gchiesa/ska/internal/localconfigservice"
	"github.com/gchiesa/ska/internal/stringprocessor"
	"github.com/gchiesa/ska/internal/tui"
	"github.com/gchiesa/ska/internal/upstreamconfigservice"
)

// SkaCreateTask defines the parameters and options for expanding a blueprint
// into a destination path. Use NewSkaCreateTask to construct it and call Create
// to render files under DestinationPath.
type SkaCreateTask struct {
	// TemplateURI is the URI to the blueprint (e.g. local path or remote github/gitlab URL).
	TemplateURI string
	// DestinationPath is the folder where the blueprint will be rendered.
	DestinationPath string
	// NamedConfig is the optional configuration name stored under .ska-config.
	NamedConfig string
	// Variables contains key/value pairs used by the template engine.
	Variables map[string]string
	// Options controls interactive mode and template engine.
	Options *SkaTaskOptions
	// Log is a contextual logger used by the task.
	Log *log.Entry
}

// SkaTaskOptions controls the behavior of create/update tasks.
//   - NonInteractive: when true, skip TUI prompts and use provided variables only
//   - ShowBanner: when true, display the interactive banner when using TUI
//   - Engine: select the template engine (sprig or jinja)
type SkaTaskOptions struct {
	NonInteractive bool
	ShowBanner     bool
	Engine         templateprovider.TemplateType // jinja or sprig
}

// NewSkaCreateTask constructs a SkaCreateTask with the provided parameters.
// The returned task can be executed via (*SkaCreateTask).Create.
func NewSkaCreateTask(templateURI, destinationPath, namedConfig string, variables map[string]string, options SkaTaskOptions) *SkaCreateTask {
	logCtx := log.WithFields(log.Fields{
		"pkg": "skaffolder",
	})
	return &SkaCreateTask{
		TemplateURI:     templateURI,
		DestinationPath: destinationPath,
		NamedConfig:     namedConfig,
		Variables:       variables,
		Options:         &options,
		Log:             logCtx,
	}
}

// Create expands the blueprint referred by TemplateURI into DestinationPath.
// It optionally prompts for variables unless NonInteractive is set.
func (s *SkaCreateTask) Create() error {
	// blueprint provider
	blueprintProvider, err := contentprovider.ByURI(s.TemplateURI)
	if err != nil {
		return err
	}

	defer func(templateProvider contentprovider.RemoteContentProvider) {
		_ = templateProvider.Cleanup()
	}(blueprintProvider)

	// configservice
	localConfig := localconfigservice.NewLocalConfigService(s.NamedConfig)

	// check if localconfig already exist, if yes we fail
	if localConfig.ConfigExists(s.DestinationPath) {
		log.Infof("Default config or specified named config already exists, please use a different name for the configuration")
		return fmt.Errorf("configuration already exists")
	}

	if err = blueprintProvider.DownloadContent(); err != nil { //nolint:govet //not a bit deal
		return err
	}

	// load the config for upstream blueprint
	upstreamConfig, err := upstreamconfigservice.NewUpstreamConfigService().LoadFromPath(blueprintProvider.WorkingDir())
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
		fmt.Sprintf("Blueprint: %s", s.TemplateURI),
		upstreamConfig.GetInputs()).SetShowBanner(s.Options.ShowBanner)

	// check if interactive mode is required
	if !s.Options.NonInteractive && interactiveService.ShouldRun() {
		// set any potential variables set by the user
		interactiveService.SetDefaults(s.Variables)
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
