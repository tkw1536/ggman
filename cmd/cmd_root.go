package cmd

import (
	"flag"

	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Root is the 'ggman root' command
var Root program.Command = root{}

type root struct{}

func (root) Name() string {
	return "root"
}

func (root) Options(flagset *flag.FlagSet) program.Options {
	return program.Options{
		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (root) AfterParse() error {
	return nil
}

func (root) Run(context program.Context) error {
	context.Println(context.Env.Root)
	return nil
}
