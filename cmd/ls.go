package cmd

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/gg"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/program/exit"
)

// Ls is the 'ggman ls' command.
//
// When called, the ggman ls command prints a list of paths to all locally cloned repositories to standard output.
//   --exit-code
// When provided, exit with code 1 if no repositories are found.
//   --one
// List at most one repository
var Ls program.Command = &ls{}

type ls struct {
	ExitCode bool `short:"e" long:"exit-code" description:"Return exit code 1 if no repositories are found"`
	One      bool `short:"o" long:"one" description:"List at most one repository, for use in shell scripts"`
}

func (*ls) BeforeRegister(program *program.Program) {}

func (l *ls) Description() program.Description {
	return program.Description{
		Name:        "ls",
		Description: "List local paths to all locally cloned repositories. ",

		Environment: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

func (*ls) AfterParse() error {
	return nil
}

var errLSExitFlag = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (l *ls) Run(context program.Context) error {
	repos := gg.C2E(context).Repos()
	if l.One && len(repos) > 0 {
		repos = repos[:1]
	}
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
