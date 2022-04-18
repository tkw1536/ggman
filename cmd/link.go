package cmd

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

// Link is the 'ggman link' command.
//
// The 'ggman link' symlinks the repository in the path passed as the first argument where it would have been cloned to inside 'ggman root'.
var Link ggman.Command = &link{}

type link struct {
	Positionals struct {
		Path string `positional-arg-name:"PATH" required:"1-1" description:"path of repository to symlink"`
	} `positional-args:"true"`
}

func (link) BeforeRegister(program *ggman.Program) {}

func (link) Description() ggman.Description {
	return ggman.Description{
		Command:     "link",
		Description: "symlink a repository into the local repository structure",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

var errLinkDoesNotExist = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to link repository: can not open source repository",
}

var errLinkSamePath = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to link repository: link source and target are identical",
}

var errLinkAlreadyExists = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to link repository: another directory already exists in target location",
}

var errLinkUnknown = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unknown linking error: %s",
}

func (l *link) Run(context ggman.Context) error {
	// make sure that the path is absolute
	// to avoid relative symlinks
	from, e := context.Environment.Abs(l.Positionals.Path)
	if e != nil {
		return errLinkDoesNotExist
	}

	// open the source repository and get the remotre
	r, e := context.Environment.Git.GetRemote(from)
	if e != nil {
		return errLinkDoesNotExist
	}

	// find the target path
	to, err := context.Environment.Local(env.ParseURL(r))
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
