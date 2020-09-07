package cmd

import (
	"sync"

	flag "github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Fix is the 'ggman fix' command
var Fix program.Command = fix{}

type fix struct{}

func (fix) Name() string {
	return "fix"
}

func (fix) Options(flagset *flag.FlagSet) program.Options {
	return program.Options{
		FlagValue:       "--simulate",
		FlagDescription: "If set, only print what would be done. ",

		Environment: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
			AllowsFilter: true,
		},
	}
}

func (fix) AfterParse() error {
	return nil
}

var errFixCustom = ggman.Error{
	ExitCode: ggman.ExitGeneric,
}

func (fix) Run(context program.Context) error {
	simulate := context.Flag

	hasError := false
	for _, repo := range context.Repos() {
		var initialMessage sync.Once // send an initial log message to the user, once

		if e := context.Git.UpdateRemotes(repo, func(url, remoteName string) (string, error) {
			canon := context.Canonical(context.ParseURL(url))

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
