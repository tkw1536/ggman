package cmd

import (
	"flag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Ls is the 'ggman ls' command
var Ls program.Command = ls{}

type ls struct{}

func (ls) Name() string {
	return "ls"
}

func (ls) Options(flagset *flag.FlagSet) program.Options {
	return program.Options{
		FlagValue:       "--exit-code",
		FlagDescription: "If provided, return exit code 1 if no repositories are found. ",

		Environment: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

func (ls) AfterParse() error {
	return nil
}

var errLSExitFlag = ggman.Error{
	ExitCode: ggman.ExitGeneric,
}

func (ls) Run(context program.Context) error {
	repos := context.Repos()
	for _, repo := range repos {
		context.Println(repo)
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if context.Flag && len(repos) == 0 {
		return errLSExitFlag
	}

	return nil
}
