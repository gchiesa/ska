package cmd

import (
	"context"
	"fmt"
	"github.com/gchiesa/ska/pkg/skaffolder"
	"github.com/gchiesa/ska/pkg/util"
)

type ConfigListCmd struct {
}

type ConfigListResultItem struct {
	NamedConfig string `json:"NamedConfig" csv:"NamedConfig"`
	LastUpdate  string `json:"LastUpdate" csv:"LastUpdate"`
}

func (c *ConfigListCmd) Execute(ctx context.Context) error {
	ska := skaffolder.NewSkaConfigTask(ctx.Value(configFolderPath("path")).(string))

	namedConfigs, err := ska.ListNamedConfigs()
	if err != nil {
		return err
	}
	result := make([]ConfigListResultItem, 0, len(namedConfigs))
	for _, namedConfig := range namedConfigs {
		lastUpdate, err := ska.QueryNamedConfigJSON(namedConfig, "{.State.LastUpdate}") // nolint:govet
		if err != nil {
			return err
		}
		result = append(result, ConfigListResultItem{
			NamedConfig: namedConfig,
			LastUpdate:  lastUpdate,
		})
	}

	outputFormat := ctx.Value(consoleOutputFormat("output-format")).(string)
	output, err := util.RenderWithOutputFormat(result, outputFormat)
	if err != nil {
		return err
	}
	fmt.Println(string(output))
	return nil
}
