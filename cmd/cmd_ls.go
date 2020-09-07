package cmd

import (
	flag "github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Ls is the 'ggman ls' command
var Ls program.Command = &ls{}

type ls struct {
	ExitCode bool
}

func (ls) Name() string {
	return "ls"
}

func (l *ls) Options(flagset *flag.FlagSet) program.Options {
	flagset.BoolVarP(&l.ExitCode, "exit-code", "e", l.ExitCode, "If provided, return exit code 1 if no repositories are found. ")
	return program.Options{
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

func (l ls) Run(context program.Context) error {
	repos := context.Repos()
	for _, repo := range repos {
		context.Println(repo)
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if l.ExitCode && len(repos) == 0 {
		return errLSExitFlag
	}

	return nil
}
