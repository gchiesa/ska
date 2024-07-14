package cmd

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/configuration"
	"github.com/gchiesa/ska/internal/contentprovider"
	"github.com/gchiesa/ska/internal/processor"
)

type CreateCmd struct {
	TemplateURI     string            `arg:"-t,--template,required" help:"URI of the template"`
	DestinationPath string            `arg:"-d,--destination,required" help:"Destination path"`
	Variables       map[string]string `arg:"-v,separate" help:"Variables to use in the template. Can be specified multiple times"`
}

func (c *CreateCmd) Execute() error {
	templateProvider, err := contentprovider.ByURI(c.TemplateURI)
	defer func(templateProvider contentprovider.RemoteContentProvider) {
		_ = templateProvider.Cleanup()
	}(templateProvider)

	if err != nil {
		return err
	}

	configService := configuration.NewConfigService()

	if err := templateProvider.DownloadContent(); err != nil { //nolint:govet //not a bit deal
		return err
	}

	fileTreeProcessor := processor.NewFileTreeProcessor(templateProvider.WorkingDir(), c.DestinationPath)
	defer func(fileTreeProcessor *processor.FileTreeProcessor) {
		_ = fileTreeProcessor.Cleanup()
	}(fileTreeProcessor)

	vars := mapStringToMapInterface(c.Variables)
	if err := fileTreeProcessor.Render(vars); err != nil { //nolint:govet //not a bit deal
		return err
	}

	// save the config
	err = configService.
		WithVariables(vars).
		WithBlueprintUpstream(templateProvider.RemoteURI()).
		WriteConfig(c.DestinationPath)
	if err != nil {
		return err
	}

	log.Infof("template created under file path: %s", c.DestinationPath)
	return nil
}
