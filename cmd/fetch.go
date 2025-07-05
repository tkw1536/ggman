package cmd

//spellchecker:words github ggman goprogram exit
import (
	"fmt"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

//spellchecker:words nolint wrapcheck

// Fetch is the 'ggman fetch' command.
//
// 'ggman fetch' is the equivalent of running 'git fetch --all' on all locally cloned repositories.
var Fetch ggman.Command = fetch{}

type fetch struct{}

func (fetch) Description() ggman.Description {
	return ggman.Description{
		Command:     "fetch",
		Description: "run \"git fetch --all\" on locally cloned repositories",

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

var errFetchCustom = exit.NewErrorWithCode("", exit.ExitGeneric)

func (fetch) Run(context ggman.Context) error {
	hasError := false

	// iterate over all the repositories, and run git fetch
	for _, repo := range context.Environment.Repos(true) {
		if _, err := context.Printf("Fetching %q\n", repo); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if e := context.Environment.Git.Fetch(context.IOStream, repo); e != nil {
			if _, err := context.EPrintln(e.Error()); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
			hasError = true
		}
	}

	if hasError {
		return errFetchCustom
	}
	return nil
}
