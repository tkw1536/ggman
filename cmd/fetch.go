package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Fetch is the 'ggman fetch' command.
//
// 'ggman fetch' is the equivalent of running 'git fetch --all' on all locally cloned repositories.
var Fetch program.Command = fetch{}

type fetch struct{}

func (fetch) BeforeRegister() {}

func (fetch) Description() program.Description {
	return program.Description{
		Name:        "fetch",
		Description: "Run 'git fetch --all' on locally cloned repositories",

		Environment: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

func (fetch) AfterParse() error {
	return nil
}

var errFetchCustom = ggman.Error{
	ExitCode: ggman.ExitGeneric,
}

func (fetch) Run(context program.Context) error {
	// find all the repos
	hasError := false

	// and fetch them
	for _, repo := range context.Repos() {
		context.Printf("Fetching %q\n", repo)
		if e := context.Git.Fetch(context.IOStream, repo); e != nil {
			context.EPrintln(e.Error())
			hasError = true
		}
	}

	if hasError {
		return errFetchCustom
	}
	return nil
}
