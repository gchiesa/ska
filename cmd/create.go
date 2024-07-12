package cmd

import (
	"github.com/apex/log"
	"github.com/gchiesa/ska/internal/configuration"
	"github.com/gchiesa/ska/internal/content_provider"
	"github.com/gchiesa/ska/internal/processor"
)

type CreateCmd struct {
	TemplateURI     string            `arg:"-t,--template,required" help:"URI of the template"`
	DestinationPath string            `arg:"-d,--destination,required" help:"Destination path"`
	Variables       map[string]string `arg:"-v,separate" help:"Variables to use in the template. Can be specified multiple times"`
}

func (c *CreateCmd) Execute() error {
	templateProvider, err := content_provider.ContentProviderByURI(c.TemplateURI)
	defer templateProvider.RemoveWorkingDir()

	if err != nil {
		log.Fatalf("error creating template provider: %v", err)
	}

	configService := configuration.NewConfigService()

	if err := templateProvider.DownloadContent(); err != nil {
		log.Fatalf("error downloading template: %v", err)
	}

	fileTreeProcessor := processor.NewFileTreeProcessor(templateProvider.WorkingDir(), c.DestinationPath, processor.TreeRendererOptions{})
	defer fileTreeProcessor.RemoveWorkingDir()

	vars := mapStringToMapInterface(c.Variables)
	if err := fileTreeProcessor.Render(vars); err != nil {
		log.Fatalf("error rendering template: %v", err)
	}

	// save the config
	err = configService.
		WithVariables(vars).
		WithBlueprintUpstream(templateProvider.RemoteURI()).
		WriteConfig(c.DestinationPath)
	if err != nil {
		log.Fatalf("error writing config: %v", err)
	}

	log.Infof("template created under file path: %s", c.DestinationPath)
	return nil
}
