package cmd

//spellchecker:words github ggman goprogram exit
import (
	"fmt"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

//spellchecker:words canonicalized CANFILE nolint wrapcheck

// When provided, instead of printing the urls directly, prints the canonical remotes of all repositories.
var Lsr ggman.Command = lsr{}

type lsr struct {
	Canonical bool `description:"print canonicalized URLs" long:"canonical" short:"c"`
}

func (lsr) Description() ggman.Description {
	return ggman.Description{
		Command:     "lsr",
		Description: "list remote URLs to all locally cloned repositories",

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

var errLSRInvalidCanfile = exit.NewErrorWithCode("invalid CANFILE found", env.ExitInvalidEnvironment)

func (l lsr) Run(context ggman.Context) error {
	var lines env.CanFile
	if l.Canonical {
		var err error
		if lines, err = context.Environment.LoadDefaultCANFILE(); err != nil {
			return fmt.Errorf("%w: %w", errLSRInvalidCanfile, err)
		}
	}

	// and print them
	for _, repo := range context.Environment.Repos(true) {
		remote, err := context.Environment.Git.GetRemote(repo, "")
		if err != nil {
			continue
		}
		if l.Canonical {
			remote = env.ParseURL(remote).CanonicalWith(lines)
		}
		if _, err := context.Println(remote); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}

	return nil
}
