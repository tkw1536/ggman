package cmd

import (
	"os"
	"path/filepath"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/program/exit"
)

// Relocate is the 'ggman relocate' command.
//
// Relocate moves all repositories to the location where they should be moved to if they had been cloned with 'ggman clone'.
var Relocate ggman.Command = &relocate{}

type relocate struct {
	Simulate bool `short:"s" long:"simulate" description:"Only print unix-like commands to move repositories around"`
}

func (relocate) BeforeRegister(program *ggman.Program) {}

func (r *relocate) Description() program.Description {
	return program.Description{
		Name:        "relocate",
		Description: "Move locally cloned repositories into locations as per 'ggman where'. ",

		Environment: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
			AllowsFilter: true,
		},
	}
}

func (relocate) AfterParse() error {
	return nil
}

var errUnableMoveCreateParent = exit.Error{
	Message:  "Unable to create parent directory for destination: %s",
	ExitCode: exit.ExitGeneric,
}

var errUnableToMoveRepo = exit.Error{
	Message:  "Unable to move repository: %s",
	ExitCode: exit.ExitGeneric,
}

func (r relocate) Run(context ggman.Context) error {
	for _, gotPath := range context.Runtime().Repos() {
		// determine the remote path and where it should go
		remote, err := context.Runtime().Git.GetRemote(gotPath)
		if err != nil || remote == "" { // ignore remotes that don't exist
			continue
		}
		shouldPath, err := context.Runtime().Local(env.ParseURL(remote))
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
