package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Pull is the 'ggman pull' command
var Pull program.Command = pull{}

type pull struct{}

func (pull) Name() string {
	return "pull"
}

func (pull) Options(flagset *pflag.FlagSet) program.Options {
	return program.Options{
		Environment: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

func (pull) AfterParse() error {
	return nil
}

var errPullCustom = ggman.Error{
	ExitCode: ggman.ExitGeneric,
}

func (pull) Run(context program.Context) error {
	hasError := false
	for _, repo := range context.Env.Repos() {
		context.Printf("Pulling %q\n", repo)
		if e := context.Git.Pull(ggman.NewEnvIOStream(), repo); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			hasError = true
		}
	}

	if hasError {
		return errPullCustom
	}

	return nil
}
