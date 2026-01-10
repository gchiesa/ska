package cmd

import (
	"context"

	"github.com/gchiesa/ska/pkg/skaffolder"
)

type ConfigRenameCmd struct {
	Name    string `arg:"-o,--name,required" help:"The name of the named configuration to rename"`
	NewName string `arg:"-n,--new-name,required" help:"The new name to give to the named configuration"`
}

func (c *ConfigRenameCmd) Execute(ctx context.Context) error {
	ska := skaffolder.NewSkaConfigTask(ctx.Value(configFolderPath("path")).(string))

	if err := ska.RenameNamedConfig(c.Name, c.NewName); err != nil {
		return err
	}
	return nil
}
