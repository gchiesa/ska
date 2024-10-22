package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/gchiesa/ska/pkg/skaffolder"
	"github.com/gchiesa/ska/pkg/util"
	"github.com/manifoldco/promptui"
)

type ConfigDeleteCmd struct {
	Name        string `arg:"-n,--name,required" help:"The name of the named configuration to delete"`
	AutoApprove bool   `arg:"-y,--auto-approve" help:"Skip the confirmation prompt"`
}

func (c *ConfigDeleteCmd) Execute(ctx context.Context) error {
	ska := skaffolder.NewSkaConfigTask(ctx.Value(configFolderPath("path")).(string))

	lastUpdate, err := ska.QueryNamedConfigJSON(c.Name, "{.State.LastUpdate}") // nolint:govet
	if err != nil {
		return err
	}
	var result = make([]ConfigListResultItem, 0)
	result = append(result, ConfigListResultItem{
		NamedConfig: c.Name,
		LastUpdate:  lastUpdate,
	})
	if !c.AutoApprove {
		output, err := util.RenderWithOutputFormat(result, "table")
		if err != nil {
			return err
		}
		fmt.Println(string(output))
		p := promptui.Prompt{
			Label:     fmt.Sprintf("Do you really want to delete this configuration named %s", c.Name),
			IsConfirm: true,
		}
		response, err := p.Run()
		if err != nil {
			if errors.Is(err, promptui.ErrAbort) {
				log.Infof("You responded: %s, so not proceeding further.", response)
				return nil
			}
			return err
		}
	}
	if err := ska.DeleteConfig(c.Name); err != nil {
		return err
	}
	log.Infof("Deleted configuration: %s", c.Name)
	return nil
}
