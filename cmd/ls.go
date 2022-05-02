package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

// Ls is the 'ggman ls' command.
//
// When called, the ggman ls command prints a list of paths to all locally cloned repositories to standard output.
//   --exit-code
// When provided, exit with code 1 if no repositories are found.
//   --one
// List at most one repository
var Ls ggman.Command = ls{}

type ls struct {
	ExitCode bool `short:"e" long:"exit-code" description:"return exit code 1 if no repositories are found"`
	Scores   bool `short:"s" long:"scores" description:"show scores returned from filter along with repositories"`
	One      bool `short:"o" long:"one" description:"list at most one repository, for use in shell scripts"`
}

func (ls) Description() ggman.Description {
	return ggman.Description{
		Command:     "ls",
		Description: "list local paths to all locally cloned repositories",

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

var errLSExitFlag = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (l ls) Run(context ggman.Context) error {
	repos, scores := context.Environment.RepoScores()
	if l.One && len(repos) > 0 {
		repos = repos[:1]
	}
	for i, repo := range repos {
		if l.Scores {
			context.Printf("%f %s\n", scores[i], repo)
			continue
		}
		context.Println(repo)
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if l.ExitCode && len(repos) == 0 {
		return errLSExitFlag
	}

	return nil
}
