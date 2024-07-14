package cmd

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/apex/log"
	"github.com/apex/log/handlers/cli"
	"os"
)

const programName = "ska"
const version = "development"

type arguments struct {
	CreateCmd *CreateCmd `arg:"subcommand:create"`
	Debug     bool       `arg:"-d"`
}

func (arguments) Version() string { return fmt.Sprintf("%s version %s", programName, version) }

func Execute() error {
	log.SetHandler(cli.New(os.Stderr))
	log.SetLevel(log.InfoLevel)

	var args arguments
	arg.MustParse(&args)

	if args.Debug {
		log.SetLevel(log.DebugLevel)
	}

	switch {
	case args.CreateCmd != nil:
		if err := args.CreateCmd.Execute(); err != nil {
			log.Fatalf("error executing create command: %v", err)
		}
	default:
		fmt.Println("no subcommand specified")
	}

	return nil
}
