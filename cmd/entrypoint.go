package cmd

import (
	"context"
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"github.com/gchiesa/ska/pkg/templateprovider"
	"os"
)

const programName = "ska"
const githubRepo = "https://github.com/gchiesa/ska"

var commandVersion = "development"

type arguments struct {
	CreateCmd    *CreateCmd `arg:"subcommand:create"`
	UpdateCmd    *UpdateCmd `arg:"subcommand:update"`
	ConfigCmd    *ConfigCmd `arg:"subcommand:config"`
	OutputFormat string     `arg:"--format" default:"table" help:"The format for output. It can be: csv,json,markdown,table"`
	Debug        bool       `arg:"-d"`
	JSONOutput   bool       `arg:"-j,--json" help:"Enable JSON output for logging"`
	Engine       string     `arg:"--engine" default:"sprig" help:"Template engine to use (sprig or jinja)"`
}

type contextEngineKey string
type consoleOutputFormat string

var args arguments

func Execute(version string) error {
	commandVersion = version
	log.SetHandler(cli.New(os.Stderr))
	log.SetLevel(log.InfoLevel)

	arg.MustParse(&args)

	if args.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if args.JSONOutput {
		log.SetHandler(json.New(os.Stderr))
	}

	if args.Engine != "sprig" && args.Engine != "jinja" {
		log.Fatalf("invalid template engine: %s", args.Engine)
	}

	ctx := context.TODO()
	ctx = context.WithValue(ctx, contextEngineKey("engine"), templateprovider.GetTypeFromString(args.Engine))
	ctx = context.WithValue(ctx, consoleOutputFormat("output-format"), args.OutputFormat)
	switch {
	case args.CreateCmd != nil:
		if err := args.CreateCmd.Execute(ctx); err != nil {
			log.Fatalf("error executing create command: %v", err)
		}
	case args.UpdateCmd != nil:
		if err := args.UpdateCmd.Execute(ctx); err != nil {
			log.Fatalf("error executing update command: %v", err)
		}
	case args.ConfigCmd != nil:
		if err := args.ConfigCmd.Execute(ctx); err != nil {
			log.Fatalf("error executing config command: %v", err)
		}
	default:
		fmt.Println("no subcommand specified, please use the --help flag to check available commands")
	}

	return nil
}

func (arguments) Version() string {
	return fmt.Sprintf("version: %s\n", commandVersion)
}

func (arguments) Description() string {
	return fmt.Sprintf(`
%s is a tool for scaffolding your directories based on blueprint templates available locally or removely on GitHub and GitLab.
`, programName)
}

func (arguments) Epilogue() string {
	return fmt.Sprintf(`
For more information check the repository on %s.

Made with love by https://github.com/gchiesa.
`, githubRepo)
}
