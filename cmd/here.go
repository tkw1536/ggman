package cmd

import (
	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Here is the 'ggman here' command.
//
// 'ggman here' prints the path to the root of the repository in the current working directory to standard output.
//   --tree
// When provided, also print the relative path from the root of the repository to the current path.
var Here program.Command = &here{}

type here struct {
	Tree bool
}

func (here) Name() string {
	return "here"
}

func (h *here) Options(flagset *pflag.FlagSet) program.Options {
	flagset.BoolVarP(&h.Tree, "tree", "t", h.Tree, "If provided, also print the current HEAD reference and relative path to the root of the git worktree. ")
	return program.Options{
		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (here) AfterParse() error {
	return nil
}

func (h here) Run(context program.Context) error {
	root, worktree, err := context.At(".")
	if err != nil {
		return err
	}

	context.Println(root)
	if h.Tree {
		context.Println(worktree)
	}

	return nil
}
