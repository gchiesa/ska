package cmd

import (
	"github.com/gchiesa/ska/pkg/skaffolder"
)

type CreateCmd struct {
	TemplateURI     string            `arg:"-t,--template,required" help:"URI of the template"`
	DestinationPath string            `arg:"-d,--destination,required" help:"Destination path"`
	Variables       map[string]string `arg:"-v,separate" help:"Variables to use in the template. Can be specified multiple times"`
}

func (c *CreateCmd) Execute() error {
	ska := skaffolder.NewSkaCreate(
		c.TemplateURI,
		c.DestinationPath,
		c.Variables,
	)
	return ska.Create()
}
