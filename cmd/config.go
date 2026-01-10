package cmd

import (
	"context"
	"fmt"

	"github.com/apex/log"
)

type ConfigCmd struct {
	*ConfigListCmd   `arg:"subcommand:list"`
	*ConfigRenameCmd `arg:"subcommand:rename"`
	*ConfigDeleteCmd `arg:"subcommand:delete"`
	FolderPath       string `arg:"-p,--path" default:"." help:"Local path where the .ska-config folder is located"`
}

type configFolderPath string

func (c *ConfigCmd) Execute(ctx context.Context) error {
	configCtx := context.WithValue(ctx, configFolderPath("path"), c.FolderPath)
	switch {
	case c.ConfigListCmd != nil:
		if err := args.ConfigCmd.ConfigListCmd.Execute(configCtx); err != nil {
			log.Fatalf("error executing config list command: %v", err)
		}
	case c.ConfigRenameCmd != nil:
		if err := args.ConfigCmd.ConfigRenameCmd.Execute(configCtx); err != nil {
			log.Fatalf("error executing config rename command: %v", err)
		}
	case c.ConfigDeleteCmd != nil:
		if err := args.ConfigCmd.ConfigDeleteCmd.Execute(configCtx); err != nil {
			log.Fatalf("error executing config delete command: %v", err)
		}
	default:
		fmt.Println("no subcommand specified, please use the --help flag to check available commands")
	}
	return nil
}
