package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program/exit"
)

// Fetch is the 'ggman fetch' command.
//
// 'ggman fetch' is the equivalent of running 'git fetch --all' on all locally cloned repositories.
var Fetch ggman.Command = fetch{}

type fetch struct{}

func (fetch) BeforeRegister(program *ggman.Program) {}

func (fetch) Description() ggman.Description {
	return ggman.Description{
		Command:     "fetch",
		Description: "Run 'git fetch --all' on locally cloned repositories",

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

func (fetch) AfterParse() error {
	return nil
}

var errFetchCustom = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (fetch) Run(context ggman.Context) error {
	hasError := false

	// iterate over all the repositories, and run git fetch
	for _, repo := range context.Environment.Repos() {
		context.Printf("Fetching %q\n", repo)
		if e := context.Environment.Git.Fetch(context.IOStream, repo); e != nil {
			context.EPrintln(e.Error())
			hasError = true
		}
	}

	if hasError {
		return errFetchCustom
	}
	return nil
}
