package cmd

import (
	"github.com/gchiesa/ska/pkg/skaffolder"
)

type UpdateCmd struct {
	FolderPath string            `arg:"required" help:"Folder path where the .ska-config.yml file is located"`
	Variables  map[string]string `arg:"-v,separate" help:"Variables to use in the template. Can be specified multiple times"`
}

func (c *UpdateCmd) Execute() error {
	ska := skaffolder.NewSkaUpdate(
		c.FolderPath,
		c.Variables,
	)
	return ska.Update()
}
