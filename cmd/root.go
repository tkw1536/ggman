package cmd

import (
	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Root is the 'ggman root' command.
//
// The 'ggman root' command prints the ggman root directory to standard output.
// It does not require the root directory to exist.
var Root program.Command = root{}

type root struct{}

func (root) Name() string {
	return "root"
}

func (root) Options(flagset *pflag.FlagSet) program.Options {
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
