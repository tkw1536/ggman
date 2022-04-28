package cmd

import (
	"os"
	"path/filepath"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

// Relocate is the 'ggman relocate' command.
//
// Relocate moves all repositories to the location where they should be moved to if they had been cloned with 'ggman clone'.
var Relocate ggman.Command = &relocate{}

type relocate struct {
	Simulate bool `short:"s" long:"simulate" description:"only print unix-like commands to move repositories around"`
}

func (relocate) BeforeRegister(program *ggman.Program) {}

func (r *relocate) Description() ggman.Description {
	return ggman.Description{
		Command:     "relocate",
		Description: "move locally cloned repositories into locations as per \"ggman where\"",

		Requirements: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
			AllowsFilter: true,
		},
	}
}

var errUnableMoveCreateParent = exit.Error{
	Message:  "unable to create parent directory for destination: %s",
	ExitCode: exit.ExitGeneric,
}

var errUnableToMoveRepo = exit.Error{
	Message:  "unable to move repository: %s",
	ExitCode: exit.ExitGeneric,
}

func (r relocate) Run(context ggman.Context) error {
	for _, gotPath := range context.Environment.Repos() {
		// determine the remote path and where it should go
		remote, err := context.Environment.Git.GetRemote(gotPath)
		if err != nil || remote == "" { // ignore remotes that don't exist
			continue
		}
		shouldPath, err := context.Environment.Local(env.ParseURL(remote))
		if err != nil {
			return err
		}

		// if it is the same, don't move it
		if gotPath == shouldPath {
			continue
		}

		parentPath := filepath.Join(shouldPath, "..")

		// print what is being done
		context.Printf("mkdir -p %s\n", shellescape.Quote(parentPath))
		context.Printf("mv %s %s\n", shellescape.Quote(gotPath), shellescape.Quote(shouldPath))
		if r.Simulate {
			continue
		}

		// do it!
		if err := os.MkdirAll(parentPath, os.ModePerm); err != nil {
			return errUnableMoveCreateParent.WithMessageF(err)
		}

		if err := os.Rename(gotPath, shouldPath); err != nil {
			return errUnableToMoveRepo.WithMessageF(err)
		}
	}

	return nil
}
