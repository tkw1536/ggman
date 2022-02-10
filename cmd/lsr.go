package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program/exit"
)

// Lsr is the 'ggman lsr' command.
//
// When called, the ggman ls command prints a list of remotes of all locally cloned repositories to standard output.
// The remotes will be listed in dictionary order of their local installation paths.
//  --canonical
// When provided, instead of printing the urls directly, prints the canonical remotes of all repositories
var Lsr ggman.Command = &lsr{}

type lsr struct {
	Canonical bool `short:"c" long:"canonical" description:"Print canonicalized URLs"`
}

func (lsr) BeforeRegister(program *ggman.Program) {}

func (l *lsr) Description() ggman.Description {
	return ggman.Description{
		Name:        "lsr",
		Description: "List remote URLs to all locally cloned repositories. ",

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

func (lsr) AfterParse() error {
	return nil
}

var errInvalidCanfile = exit.Error{
	Message:  "Invalid CANFILE found. ",
	ExitCode: exit.ExitInvalidEnvironment,
}

func (l lsr) Run(context ggman.Context) error {
	var lines env.CanFile
	if l.Canonical {
		var err error
		if lines, err = context.Runtime().LoadDefaultCANFILE(); err != nil {
			return errInvalidCanfile
		}
	}

	// and print them
	for _, repo := range context.Runtime().Repos() {
		remote, err := context.Runtime().Git.GetRemote(repo)
		if err != nil {
			continue
		}
		if l.Canonical {
			remote = env.ParseURL(remote).CanonicalWith(lines)
		}
		context.Println(remote)
	}

	return nil
}
