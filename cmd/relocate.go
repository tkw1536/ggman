package cmd

//spellchecker:words errors path filepath essio shellescape github ggman internal dirs goprogram exit pkglib
import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"al.essio.dev/pkg/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/dirs"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/fsx"
)

// Relocate is the 'ggman relocate' command.
//
// Relocate moves all repositories to the location where they should be moved to if they had been cloned with 'ggman clone'.
var Relocate ggman.Command = relocate{}

type relocate struct {
	Simulate bool `short:"s" long:"simulate" description:"only print unix-like commands to move repositories around"`
}

func (relocate) Description() ggman.Description {
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
	Message:  "unable to create parent directory for destination",
	ExitCode: exit.ExitGeneric,
}

var errUnableToMoveRepo = exit.Error{
	Message:  "unable to move repository",
	ExitCode: exit.ExitGeneric,
}

var errRepositoryAlreadyExists = exit.Error{
	Message:  "a repository already exists at %q",
	ExitCode: exit.ExitGeneric,
}

var errPathAlreadyExists = exit.Error{
	Message:  "path %q already exists",
	ExitCode: exit.ExitGeneric,
}

func (r relocate) Run(context ggman.Context) error {
	for _, gotPath := range context.Environment.Repos(false) {
		// determine the remote path and where it should go
		remote, err := context.Environment.Git.GetRemote(gotPath, "")
		if err != nil || remote == "" { // ignore remotes that don't exist
			continue
		}
		shouldPath, err := context.Environment.Local(env.ParseURL(remote))
		if err != nil {
			return err
		}

		// if it is the same, don't move it
		if fsx.Same(gotPath, shouldPath) {
			continue
		}

		parentPath := filepath.Dir(shouldPath)

		// print what is being done
		if _, err := context.Printf("mkdir -p %s\n", shellescape.Quote(parentPath)); err != nil {
			return ggman.ErrGenericOutput.WrapError(err)
		}
		if _, err := context.Printf("mv %s %s\n", shellescape.Quote(gotPath), shellescape.Quote(shouldPath)); err != nil {
			return ggman.ErrGenericOutput.WrapError(err)
		}
		if r.Simulate {
			continue
		}

		// do it!
		if err := os.MkdirAll(parentPath, dirs.NewModBits); err != nil {
			return errUnableMoveCreateParent.WrapError(err)
		}

		// if there already is a target repository at the path
		{
			got, err := context.Environment.AtRoot(shouldPath)
			if err != nil {
				return errUnableToMoveRepo.WrapError(err)
			}
			if got != "" {
				return errRepositoryAlreadyExists.WithMessageF(got)
			}
		}

		// do the rename
		{
			err := os.Rename(gotPath, shouldPath)

			// check if an error was returned because the path already existed
			// (fs.ErrPermission is returned by Windows)
			if errors.Is(err, fs.ErrExist) || errors.Is(err, fs.ErrPermission) {
				if exists, _ := fsx.Exists(shouldPath); exists {
					return errPathAlreadyExists.WithMessageF(shouldPath)
				}
			}

			if err != nil {
				return errUnableToMoveRepo.WrapError(err)
			}
		}
	}

	return nil
}
