package cmd

import (
	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Fetch is the 'ggman fetch' command.
//
// 'ggman fetch' is the equivalent of running 'git fetch --all' on all locally installed repositories.
var Fetch program.Command = fetch{}

type fetch struct{}

func (fetch) Name() string {
	return "fetch"
}

func (fetch) Options(flagset *pflag.FlagSet) program.Options {
	return program.Options{
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
