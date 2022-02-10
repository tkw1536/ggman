package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/program/exit"
)

// Pull is the 'ggman pull' command.
//
// 'ggman pull' is the equivalent of running 'git pull' on all locally installed repositories.
var Pull ggman.Command = pull{}

type pull struct{}

func (pull) BeforeRegister(program *ggman.Program) {}

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

var errPullCustom = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (pull) Run(context ggman.Context) error {
	hasError := false
	for _, repo := range context.Runtime().Repos() {
		context.Printf("Pulling %q\n", repo)
		if e := context.Runtime().Git.Pull(context.IOStream, repo); e != nil {
			context.EPrintf("%s\n", e)
			hasError = true
		}
	}

	if hasError {
		return errPullCustom
	}

	return nil
}
