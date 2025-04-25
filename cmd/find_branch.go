package cmd

//spellchecker:words github ggman goprogram exit
import (
	"fmt"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

//spellchecker:words positionals nolint wrapcheck

// FindBranch is the 'ggman find-branch' command.
//
// The 'find-branch' command lists all repositories that contain a branch with the provided name.
// The remotes will be listed in dictionary order of their local installation paths.
//
//	--exit-code
//
// When provided, exit with code 1 if no repositories are found.
var FindBranch ggman.Command = findBranch{}

type findBranch struct {
	Positionals struct {
		Branch string `description:"name of branch to find" positional-arg-name:"BRANCH" required:"1-1"`
	} `positional-args:"true"`
	ExitCode bool `description:"exit with status code 1 when no repositories with provided branch exist" long:"exit-code" short:"e"`
}

func (findBranch) Description() ggman.Description {
	return ggman.Description{
		Command:     "find-branch",
		Description: "list repositories containing a specific branch",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

var errFindBranchCustom = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (f findBranch) Run(context ggman.Context) error {
	foundRepo := false
	for _, repo := range context.Environment.Repos(true) {
		// check if the repository has the branch!
		hasBranch, err := context.Environment.Git.ContainsBranch(repo, f.Positionals.Branch)
		if err != nil {
			panic(err)
		}
		if !hasBranch {
			continue
		}

		foundRepo = true
		if _, err := context.Println(repo); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if f.ExitCode && !foundRepo {
		return errFindBranchCustom
	}

	return nil
}
