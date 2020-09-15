package cmd

import (
	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// FindBranch is the 'ggman find-branch' command
var FindBranch program.Command = &findBranch{}

type findBranch struct {
	ExitCode bool
}

func (findBranch) Name() string {
	return "find-branch"
}

func (f *findBranch) Options(flagset *pflag.FlagSet) program.Options {
	flagset.BoolVarP(&f.ExitCode, "exit-code", "e", f.ExitCode, "Exit with Code 1 when no repositories with provided branch exist. ")

	return program.Options{
		MinArgs: 1,
		MaxArgs: 1,

		Metavar: "BRANCH",

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

func (f findBranch) Run(context program.Context) error {
	count := 0
	for _, repo := range context.Repos() {
		hasBranch, err := context.Git.ContainsBranch(repo, context.Argv[0])
		if err != nil {
			panic(err)
		}
		if !hasBranch {
			continue
		}

		count++
		context.Println(repo)
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if f.ExitCode && count == 0 {
		return errFindBranchCustom
	}

	return nil
}
