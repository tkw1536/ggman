package cmd

//spellchecker:words sync ggman goprogram exit
import (
	"fmt"
	"sync"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

//spellchecker:words canonicalizes canonicalization nolint wrapcheck

// Fix is the 'ggman fix' command.
//
// The 'ggman fix' command canonicalizes the urls of all remotes of a repository.
//
//	--simulate
//
// Instead of writing out the changes to disk, only print what would be done.
var Fix ggman.Command = fix{}

type fix struct {
	Simulate bool `description:"do not perform any canonicalization, instead only print what would be done" long:"simulate" short:"s"`
}

func (fix) Description() ggman.Description {
	return ggman.Description{
		Command:     "fix",
		Description: "canonicalizes remote URLs for cloned repositories",
		Requirements: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
			AllowsFilter: true,
		},
	}
}

var errFixCustom = exit.NewErrorWithCode("", exit.ExitGeneric)

func (f fix) Run(context ggman.Context) error {
	simulate := f.Simulate

	hasError := false

	var innerError error
	for _, repo := range context.Environment.Repos(true) {
		var initialMessage sync.Once // send an initial log message to the user, once

		if e := context.Environment.Git.UpdateRemotes(repo, func(url, remoteName string) (string, error) {
			canon := context.Environment.Canonical(env.ParseURL(url))

			if url == canon {
				return url, nil
			}

			initialMessage.Do(func() {
				if !simulate {
					if _, err := context.Printf("Fixing remote of %q", repo); err != nil {
						innerError = fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
					}
				} else {
					if _, err := context.Printf("Simulate fixing remote of %q", repo); err != nil {
						innerError = fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
					}
				}
			})

			if innerError != nil {
				return "", innerError
			}

			if _, err := context.Printf("Updating %s: %s -> %s\n", remoteName, url, canon); err != nil {
				return "", fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}

			// either return the canonical url, or (if we're simulating) the old url
			if simulate {
				return url, nil
			}

			return canon, nil
		}); e != nil {
			_, _ = context.EPrintln(e.Error()) // no way to report error
			hasError = true
		}
	}

	// if we had an error, indicate that to the user
	if hasError {
		return errFixCustom
	}

	// and finish
	return nil
}
