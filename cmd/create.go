package cmd

import (
	"context"

	"github.com/gchiesa/ska/pkg/skaffolder"
	"github.com/gchiesa/ska/pkg/templateprovider"
)

type CreateCmd struct {
	TemplateURI     string            `arg:"-b,--blueprint,required" help:"URI of the template blueprint to use"`
	DestinationPath string            `arg:"-o,--output" default:"." help:"Destination path where to expand the blueprint"`
	NamedConfig     string            `arg:"-n,--name" help:"The ska configuration name in case there are multiple templates configurations in the same root"`
	Variables       map[string]string `arg:"-v,separate" help:"Variables to use in the template. Can be specified multiple times"`
	NonInteractive  bool              `arg:"-n,--non-interactive" help:"Run in non-interactive mode"`
}

func (c *CreateCmd) Execute(ctx context.Context) error {
	options := &skaffolder.SkaTaskOptions{
		NonInteractive: c.NonInteractive,
		ShowBanner:     true,
		Engine:         ctx.Value(contextEngineKey("engine")).(templateprovider.TemplateType),
	}
	ska := skaffolder.NewSkaCreateTask(
		c.TemplateURI,
		c.DestinationPath,
		c.NamedConfig,
		c.Variables,
		*options,
	)
	return ska.Create()
}
