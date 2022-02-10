package cmd

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program/exit"
)

// Link is the 'ggman link' command.
//
// The 'ggman link' symlinks the repository in the path passed as the first argument where it would have been cloned to inside 'ggman root'.
var Link ggman.Command = link{}

type link struct{}

func (link) BeforeRegister(program *ggman.Program) {}

func (link) Description() ggman.Description {
	return ggman.Description{
		Name:        "link",
		Description: "Symlink a repository into the local repository structure. ",

		PosArgsMin: 1,
		PosArgsMax: 1,

		PosArgName: "PATH",

		PosArgDescription: "Path of repository to symlink",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (link) AfterParse() error {
	return nil
}

var errLinkDoesNotExist = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unable to link repository: Can not open source repository. ",
}

var errLinkSamePath = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unable to link repository: Link source and target are identical. ",
}

var errLinkAlreadyExists = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unable to link repository: Another directory already exists in target location. ",
}

var errLinkUnknown = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unknown linking error: %s",
}

func (link) Run(context ggman.Context) error {
	// make sure that the path is absolute
	// to avoid relative symlinks
	from, e := context.Runtime().Abs(context.Args[0])
	if e != nil {
		return errLinkDoesNotExist
	}

	// open the source repository and get the remotre
	r, e := context.Runtime().Git.GetRemote(from)
	if e != nil {
		return errLinkDoesNotExist
	}

	// find the target path
	to, err := context.Runtime().Local(env.ParseURL(r))
	if err != nil {
		return err
	}
	parentTo := filepath.Dir(to)

	// if it's the same path, we throw an error
	if from == to {
		return errLinkSamePath
	}

	// make sure it doesn't exist
	if _, e := os.Stat(to); !errors.Is(e, fs.ErrNotExist) {
		return errLinkAlreadyExists
	}

	context.Printf("Linking %q -> %q\n", to, from)

	// make the parent folder
	if e := os.MkdirAll(parentTo, os.ModePerm); e != nil {
		return errLinkUnknown.WithMessageF(e)
	}

	// and make the symlink
	if e := os.Symlink(from, to); e != nil {
		return errLinkUnknown.WithMessageF(e)
	}

	return nil
}
