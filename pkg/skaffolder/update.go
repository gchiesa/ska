package skaffolder //nolint:typecheck

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/configuration"
	"github.com/gchiesa/ska/internal/contentprovider"
	"github.com/gchiesa/ska/internal/processor"
)

type SkaUpdate struct {
	BaseURI   string
	Variables map[string]string
	Log       *log.Entry
}

func NewSkaUpdate(baseURI string, variables map[string]string) *SkaUpdate {
	logCtx := log.WithFields(log.Fields{
		"pkg": "skaffolder",
	})
	return &SkaUpdate{
		BaseURI:   baseURI,
		Variables: variables,
		Log:       logCtx,
	}
}

func (s *SkaUpdate) Update() error {
	configService := configuration.NewConfigService()

	// read the config from the folder
	if err := configService.ReadConfig(s.BaseURI); err != nil {
		return err
	}

	// allocate the template based on the configured upstream
	templateProvider, err := contentprovider.ByURI(configService.BlueprintUpstream())
	defer func(templateProvider contentprovider.RemoteContentProvider) {
		_ = templateProvider.Cleanup()
	}(templateProvider)

	if err != nil {
		return err
	}

	if err := templateProvider.DownloadContent(); err != nil { //nolint:govet //not a bit deal
		return err
	}

	fileTreeProcessor := processor.NewFileTreeProcessor(templateProvider.WorkingDir(), s.BaseURI)
	defer func(fileTreeProcessor *processor.FileTreeProcessor) {
		_ = fileTreeProcessor.Cleanup()
	}(fileTreeProcessor)

	// merge the known variables with overrides
	vars := configService.Variables()
	for k, v := range mapStringToMapInterface(s.Variables) {
		vars[k] = v
	}

	if err := fileTreeProcessor.Render(vars); err != nil { //nolint:govet //not a bit deal
		return err
	}

	// save the config
	err = configService.
		WithVariables(vars).
		WithBlueprintUpstream(templateProvider.RemoteURI()).
		WriteConfig(s.BaseURI)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{"method": "Update", "path": s.BaseURI}).Infof("template updated under destination path: %s", s.BaseURI)
	return nil
}
