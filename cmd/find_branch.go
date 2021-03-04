package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// FindBranch is the 'ggman find-branch' command.
//
// The 'find-branch' command lists all repositories that contain a branch with the provided name.
// The remotes will be listed in dictionary order of their local installation paths.
//   --exit-code
// When provided, exit with code 1 if no repositories are found.
var FindBranch program.Command = &findBranch{}

type findBranch struct {
	ExitCode bool `short:"e" long:"exit-code" description:"Exit with Status Code 1 when no repositories with provided branch exist"`
}

func (findBranch) BeforeRegister() {}

func (f *findBranch) Description() program.Description {
	return program.Description{
		Name:        "find-branch",
		Description: "List repositories containing a specific branch",

		PosArgsMin: 1,
		PosArgsMax: 1,

		PosArgName: "BRANCH",

		PosArgDescription: "Name of branch to find",

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
		hasBranch, err := context.Git.ContainsBranch(repo, context.Args[0])
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
