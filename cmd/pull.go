package cmd

//spellchecker:words github ggman goprogram exit
import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

//spellchecker:words nolint wrapcheck

// Pull is the 'ggman pull' command.
//
// 'ggman pull' is the equivalent of running 'git pull' on all locally installed repositories.
var Pull ggman.Command = pull{}

type pull struct{}

func (pull) Description() ggman.Description {
	return ggman.Description{
		Command:     "pull",
		Description: "run \"git pull\" on locally cloned repositories",

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

var errPullCustom = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (pull) Run(context ggman.Context) error {
	hasError := false
	for _, repo := range context.Environment.Repos(true) {
		if _, err := context.Printf("Pulling %q\n", repo); err != nil {
			return ggman.ErrGenericOutput.WrapError(err) //nolint:wrapcheck
		}
		if e := context.Environment.Git.Pull(context.IOStream, repo); e != nil {
			if _, err := context.EPrintf("%s\n", e); err != nil {
				return ggman.ErrGenericOutput.WrapError(err) //nolint:wrapcheck
			}
			hasError = true
		}
	}

	if hasError {
		return errPullCustom
	}

	return nil
}
