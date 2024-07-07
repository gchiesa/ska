package main

import (
	"github.com/apex/log"
	"github.com/apex/log/handlers/logfmt"
	"github.com/gchiesa/ska/pkg/processor"
	"github.com/gchiesa/ska/pkg/provider"
	"os"
)

var version = ""

const (
	templatePath    = "/Users/gchiesa/git/swanson/swanson/tools/root-project"
	destinationPath = "/tmp/test"
)

func main() {
	log.SetHandler(logfmt.New(os.Stderr))
	log.SetLevel(log.DebugLevel)
	templateProvider := provider.NewLocalPath(templatePath, false)
	if err := templateProvider.DownloadContent(); err != nil {
		log.Fatalf("error downloading template: %v", err)
	}
	log.Infof("working dir: %s", templateProvider.WorkingDir())

	processor := processor.NewFileTreeProcessor(templateProvider.WorkingDir(), destinationPath, processor.TreeRendererOptions{})

	variables := map[string]string{
		"appName":      "ThisApp",
		"testFileName": "test-file",
	}
	if err := processor.Render(variables); err != nil {
		log.Fatalf("error rendering template: %v", err)
	}
	log.Infof("destination dir: %s", processor.WorkingDir())
	println("done")
}
