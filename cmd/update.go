package cmd

import (
	"context"
	"github.com/gchiesa/ska/pkg/skaffolder"
)

type UpdateCmd struct {
	FolderPath     string            `arg:"-p,--path,required" help:"Local path where the .ska-config.yml file is located"`
	Variables      map[string]string `arg:"-v,separate" help:"Variables to use in the template. Can be specified multiple times"`
	NonInteractive bool              `arg:"-n,--non-interactive" help:"Run in non-interactive mode"`
}

func (c *UpdateCmd) Execute(ctx context.Context) error {
	options := &skaffolder.SkaOptions{
		NonInteractive: c.NonInteractive,
		Engine:         ctx.Value("engine").(string),
	}
	ska := skaffolder.NewSkaUpdate(
		c.FolderPath,
		c.Variables,
		*options,
	)
	return ska.Update()
}
