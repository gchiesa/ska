package cmd

import (
	"context"
	"github.com/gchiesa/ska/internal/templateprovider"
	"github.com/gchiesa/ska/pkg/skaffolder"
)

type CreateCmd struct {
	TemplateURI     string            `arg:"-b,--blueprint,required" help:"URI of the template blueprint to use"`
	DestinationPath string            `arg:"-o,--output,required" help:"Destination path where to expand the blueprint"`
	Variables       map[string]string `arg:"-v,separate" help:"Variables to use in the template. Can be specified multiple times"`
	NonInteractive  bool              `arg:"-n,--non-interactive" help:"Run in non-interactive mode"`
}

func (c *CreateCmd) Execute(ctx context.Context) error {
	options := &skaffolder.SkaOptions{
		NonInteractive: c.NonInteractive,
		Engine:         ctx.Value(contextEngineKey("engine")).(templateprovider.TemplateType),
	}
	ska := skaffolder.NewSkaCreate(
		c.TemplateURI,
		c.DestinationPath,
		c.Variables,
		*options,
	)
	return ska.Create()
}
