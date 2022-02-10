package cmd

import (
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
	var hasError bool
	for _, repo := range context.Env.Repos() {
		context.Printf("Pulling %q\n", repo)
		if e := context.Env.Git.Pull(context.IOStream, repo); e != nil {
			context.EPrintf("%s\n", e)
			hasError = true
		}
	}

	if hasError {
		return errPullCustom
	}

	return nil
}
