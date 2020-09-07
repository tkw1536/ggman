package cmd

import (
	"flag"

	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Here is the 'ggman here' command
var Here program.Command = here{}

type here struct{}

func (here) Name() string {
	return "here"
}

func (here) Options(flagset *flag.FlagSet) program.Options {
	return program.Options{
		FlagValue:       "--tree",
		FlagDescription: "If provided, also print the current HEAD reference and relative path to the root of the git worktree. ",

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (here) AfterParse() error {
	return nil
}

func (here) Run(context program.Context) error {
	root, worktree, err := context.At(".")
	if err != nil {
		return err
	}

	context.Println(root)
	if context.Flag {
		context.Println(worktree)
	}

	return nil
}
