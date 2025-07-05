package cmd

//spellchecker:words github ggman goprogram exit
import (
	"fmt"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
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

var errPullCustom = exit.NewErrorWithCode("", exit.ExitGeneric)

func (pull) Run(context ggman.Context) error {
	hasError := false
	for _, repo := range context.Environment.Repos(true) {
		if _, err := context.Printf("Pulling %q\n", repo); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if e := context.Environment.Git.Pull(context.IOStream, repo); e != nil {
			if _, err := context.EPrintf("%s\n", e); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
			hasError = true
		}
	}

	if hasError {
		return errPullCustom
	}

	return nil
}
