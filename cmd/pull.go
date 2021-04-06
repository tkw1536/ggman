package cmd

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Pull is the 'ggman pull' command.
//
// 'ggman pull' is the equivalent of running 'git pull' on all locally installed repositories.
var Pull program.Command = pull{}

type pull struct{}

func (pull) BeforeRegister(program *program.Program) {}

func (pull) Description() program.Description {
	return program.Description{
		Name:        "pull",
		Description: "Run 'git pull' on locally cloned repositories",

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
		if e := context.Git.Pull(context.IOStream, repo); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			hasError = true
		}
	}

	if hasError {
		return errPullCustom
	}

	return nil
}
