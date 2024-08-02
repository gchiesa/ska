package cmd

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"github.com/apex/log/handlers/json"
	"os"
)

const programName = "ska"
const githubRepo = "https://gchiesa/ska"

var commandVersion = "development"

type arguments struct {
	CreateCmd  *CreateCmd `arg:"subcommand:create"`
	UpdateCmd  *UpdateCmd `arg:"subcommand:update"`
	Debug      bool       `arg:"-d"`
	JSONOutput bool       `arg:"-j,--json" help:"Enable JSON output for logging"`
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

func Execute(version string) error {
	commandVersion = version
	log.SetHandler(cli.New(os.Stderr))
	log.SetLevel(log.InfoLevel)

	var args arguments
	arg.MustParse(&args)

	if args.Debug {
		log.SetLevel(log.DebugLevel)
	}
	if args.JSONOutput {
		log.SetHandler(json.New(os.Stderr))
	}

	switch {
	case args.CreateCmd != nil:
		if err := args.CreateCmd.Execute(); err != nil {
			log.Fatalf("error executing create command: %v", err)
		}
	case args.UpdateCmd != nil:
		if err := args.UpdateCmd.Execute(); err != nil {
			log.Fatalf("error executing update command: %v", err)
		}
	default:
		fmt.Println("no subcommand specified, please use the --help flag to check available commands")
	}

	return nil
}
