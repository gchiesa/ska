package cmd

import (
	"context"
	"fmt"
	"github.com/apex/log"
)

type ConfigCmd struct {
	*ConfigListCmd `arg:"subcommand:list"`
	FolderPath     string `arg:"-p,--path,required" help:"Local path where the .ska-config folder is located"`
}

type configFolderPath string

func (c *ConfigCmd) Execute(ctx context.Context) error {
	configCtx := context.WithValue(ctx, configFolderPath("path"), c.FolderPath)
	switch {
	case c.ConfigListCmd != nil:
		if err := args.ConfigCmd.ConfigListCmd.Execute(configCtx); err != nil {
			log.Fatalf("error executing config list command: %v", err)
		}
	default:
		fmt.Println("no subcommand specified, please use the --help flag to check available commands")
	}
	return nil
}
