package cmd

import (
	flag "github.com/spf13/pflag"

	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Where is the 'ggman where' command
var Where program.Command = where{}

type where struct{}

func (where) Name() string {
	return "where"
}

func (where) Options(flagset *flag.FlagSet) program.Options {
	return program.Options{
		MinArgs: 1,
		MaxArgs: 1,

		Metavar: "REPO",

		UsageDescription: "Repository URI to find location of. ",

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (where) AfterParse() error {
	return nil
}

func (where) Run(context program.Context) error {
	localPath := context.Local(context.URLV(0))
	context.Println(localPath)
	return nil
}
