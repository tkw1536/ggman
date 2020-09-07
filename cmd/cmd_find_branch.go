package cmd

import (
	"flag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// FindBranch is the 'ggman find-branch' command
var FindBranch program.Command = findBranch{}

type findBranch struct{}

func (findBranch) Name() string {
	return "find-branch"
}

func (findBranch) Options(flagset *flag.FlagSet) program.Options {
	return program.Options{
		MinArgs: 1,
		MaxArgs: 1,

		Metavar: "BRANCH",

		FlagValue: "--exit-code",

		UsageDescription: "Name of branch to find in repositories. ",

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (findBranch) AfterParse() error {
	return nil
}

var errFindBranchCustom = ggman.Error{
	ExitCode: ggman.ExitGeneric,
}

func (findBranch) Run(context program.Context) error {
	repos := context.Repos()
	for _, repo := range repos {
		hasBranch, err := context.Git.ContainsBranch(repo, context.Argv[0])
		if err != nil {
			panic(err)
		}
		if !hasBranch {
			continue
		}

		context.Println(repo)
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if context.Flag && len(repos) == 0 {
		return errFindBranchCustom
	}

	return nil
}
