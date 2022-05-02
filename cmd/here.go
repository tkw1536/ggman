package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

// Here is the 'ggman here' command.
//
// 'ggman here' prints the path to the root of the repository in the current working directory to standard output.
//   --tree
// When provided, also print the relative path from the root of the repository to the current path.
var Here ggman.Command = here{}

type here struct {
	Tree bool `short:"t" long:"tree" description:"also print the current HEAD reference and relative path to the root of the git worktree"`
}

func (here) Description() ggman.Description {
	return ggman.Description{
		Command:     "here",
		Description: "print the root path to the repository in the current repository",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (h here) Run(context ggman.Context) error {
	root, worktree, err := context.Environment.At(".")
	if err != nil {
		return err
	}

	context.Println(root)
	if h.Tree {
		context.Println(worktree)
	}

	return nil
}
