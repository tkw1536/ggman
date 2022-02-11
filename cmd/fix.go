package cmd

import (
	"sync"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program/exit"
)

// Fix is the 'ggman fix' command.
//
// The 'ggman fix' command canonicalizes the urls of all remotes of a repository.
//   --simulate
// Instead of writing out the changes to disk, only print what would be done.
var Fix ggman.Command = &fix{}

type fix struct {
	Simulate bool `short:"s" long:"simulate" description:"Do not perform any canonicalization. Only print what would be done"`
}

func (fix) BeforeRegister(program *ggman.Program) {}

func (f *fix) Description() ggman.Description {
	return ggman.Description{
		Name:        "fix",
		Description: "Canonicalizes remote URLs for cloned repositories. ",
		Requirements: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
			AllowsFilter: true,
		},
	}
}

func (fix) AfterParse() error {
	return nil
}

var errFixCustom = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (f fix) Run(context ggman.Context) error {
	simulate := f.Simulate

	hasError := false
	for _, repo := range context.Environment.Repos() {
		var initialMessage sync.Once // send an initial log message to the user, once

		if e := context.Environment.Git.UpdateRemotes(repo, func(url, remoteName string) (string, error) {
			canon := context.Environment.Canonical(env.ParseURL(url))

			if url == canon {
				return url, nil
			}

			initialMessage.Do(func() {
				if !simulate {
					context.Printf("Fixing remote of %q", repo)
				} else {
					context.Printf("Simulate fixing remote of %q", repo)
				}
			})

			context.Printf("Updating %s: %s -> %s\n", remoteName, url, canon)

			// either return the canonical url, or (if we're simulating) the old url
			if simulate {
				return url, nil
			}

			return canon, nil
		}); e != nil {
			context.EPrintln(e.Error())
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
