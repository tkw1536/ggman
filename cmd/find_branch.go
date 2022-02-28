package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/goprogram/exit"
	"github.com/tkw1536/ggman/goprogram/meta"
)

// FindBranch is the 'ggman find-branch' command.
//
// The 'find-branch' command lists all repositories that contain a branch with the provided name.
// The remotes will be listed in dictionary order of their local installation paths.
//   --exit-code
// When provided, exit with code 1 if no repositories are found.
var FindBranch ggman.Command = &findBranch{}

type findBranch struct {
	ExitCode bool `short:"e" long:"exit-code" description:"Exit with Status Code 1 when no repositories with provided branch exist"`
}

func (findBranch) BeforeRegister(program *ggman.Program) {}

func (f *findBranch) Description() ggman.Description {
	return ggman.Description{
		Command:     "find-branch",
		Description: "List repositories containing a specific branch",

		Positional: meta.Positional{
			Value:       "BRANCH",
			Description: "Name of branch to find",
			Min:         1,
			Max:         1,
		},

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (findBranch) AfterParse() error {
	return nil
}

var errFindBranchCustom = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (f findBranch) Run(context ggman.Context) error {
	foundRepo := false
	for _, repo := range context.Environment.Repos() {
		// check if the repository has the branch!
		hasBranch, err := context.Environment.Git.ContainsBranch(repo, context.Args.Pos[0])
		if err != nil {
			panic(err)
		}
		if !hasBranch {
			continue
		}

		foundRepo = true
		context.Println(repo)
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if f.ExitCode && !foundRepo {
		return errFindBranchCustom
	}

	return nil
}
