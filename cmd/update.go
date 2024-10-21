package cmd

import (
	"context"
	"github.com/gchiesa/ska/internal/templateprovider"
	"github.com/gchiesa/ska/pkg/skaffolder"
)

type UpdateCmd struct {
	FolderPath     string            `arg:"-p,--path,required" help:"Local path where the .ska-config.yml file is located"`
	NamedConfig    string            `arg:"-n,--name" help:"The ska configuration name in case there are multiple templates configurations in the same root"`
	Variables      map[string]string `arg:"-v,separate" help:"Variables to use in the template. Can be specified multiple times"`
	NonInteractive bool              `arg:"-n,--non-interactive" help:"Run in non-interactive mode"`
}

func (c *UpdateCmd) Execute(ctx context.Context) error {
	options := &skaffolder.SkaTaskOptions{
		NonInteractive: c.NonInteractive,
		Engine:         ctx.Value(contextEngineKey("engine")).(templateprovider.TemplateType),
	}
	ska := skaffolder.NewSkaUpdateTask(
		c.FolderPath,
		c.NamedConfig,
		c.Variables,
		*options,
	)
	return ska.Update()
}
