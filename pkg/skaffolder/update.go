package skaffolder //nolint:typecheck

import (
	"fmt"

	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/contentprovider"
	"github.com/gchiesa/ska/internal/filetreeprocessor"
	"github.com/gchiesa/ska/internal/localconfigservice"
	"github.com/gchiesa/ska/internal/stringprocessor"
	"github.com/gchiesa/ska/internal/tui"
	"github.com/gchiesa/ska/internal/upstreamconfigservice"
	"github.com/gchiesa/ska/pkg/templateprovider"
)

// SkaUpdateTask updates an existing destination directory from its recorded
// blueprint upstream. Use NewSkaUpdateTask to construct it and call Update
// to apply changes.
type SkaUpdateTask struct {
	// BaseURI is the path to the destination directory containing .ska-config.
	BaseURI string
	// NamedConfig selects which configuration (if multiple) to use in .ska-config.
	NamedConfig string
	// Variables contains key/value overrides merged with the stored variables.
	Variables map[string]string
	// Options controls interactive prompts and template engine selection.
	Options *SkaTaskOptions
	// Log is a contextual logger used by the task.
	Log *log.Entry
}

// NewSkaUpdateTask constructs a SkaUpdateTask with the provided parameters.
// The returned task can be executed via (*SkaUpdateTask).Update.
func NewSkaUpdateTask(baseURI, namedConfig string, variables map[string]string, options SkaTaskOptions) *SkaUpdateTask {
	logCtx := log.WithFields(log.Fields{
		"pkg": "skaffolder",
	})
	return &SkaUpdateTask{
		BaseURI:     baseURI,
		NamedConfig: namedConfig,
		Variables:   variables,
		Options:     &options,
		Log:         logCtx,
	}
}

// Update downloads the upstream blueprint recorded in the destination's
// .ska-config and applies changes to the destination directory. It optionally
// prompts for variables unless NonInteractive is set.
func (s *SkaUpdateTask) Update() error {
	localConfig := localconfigservice.NewLocalConfigService(s.NamedConfig)

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
	upstreamConfig, err := upstreamconfigservice.NewUpstreamConfigService().LoadFromPath(blueprintProvider.WorkingDir())
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
		fmt.Sprintf("Blueprint: %s", localConfig.BlueprintUpstream()),
		upstreamConfig.GetInputs()).SetShowBanner(s.Options.ShowBanner).
		SetWriteOnce(true)

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
	sp := stringprocessor.NewStringProcessor(stringprocessor.WithTemplateService(templateService))
	skaConfigIgnorePaths, err := sp.RenderSliceOfStrings(upstreamConfig.SkaConfigIgnorePaths(), vars)
	if err != nil {
		return err
	}

	// save the config
	err = localConfig.WithVariables(vars).
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
