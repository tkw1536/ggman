package cmd

import (
	"os"
	"path/filepath"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// Link is the 'ggman link' command.
//
// The 'ggman link' symlinks the repository in the path passed as the first argument where it would have been cloned to inside 'ggman root'.
var Link program.Command = link{}

type link struct{}

func (link) Name() string {
	return "link"
}

func (link) Options() program.Options {
	return program.Options{
		MinArgs: 1,
		MaxArgs: 1,

		Metavar: "PATH",

		UsageDescription: "Path of repository to symlink. ",

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (link) AfterParse() error {
	return nil
}

var errLinkDoesNotExist = ggman.Error{
	ExitCode: ggman.ExitGeneric,
	Message:  "Unable to link repository: Can not open source repository. ",
}

var errLinkSamePath = ggman.Error{
	ExitCode: ggman.ExitGeneric,
	Message:  "Unable to link repository: Link source and target are identical. ",
}

var errLinkAlreadyExists = ggman.Error{
	ExitCode: ggman.ExitGeneric,
	Message:  "Unable to link repository: Another directory already exists in target location. ",
}

var errLinkUnknown = ggman.Error{
	ExitCode: ggman.ExitGeneric,
	Message:  "Unknown linking error: %s",
}

func (link) Run(context program.Context) error {
	// make sure that the path is absolute
	// to avoid relative symlinks
	from, e := context.Abs(context.Args[0])
	if e != nil {
		return errLinkDoesNotExist
	}

	// open the source repository and get the remotre
	r, e := context.Git.GetRemote(from)
	if e != nil {
		return errLinkDoesNotExist
	}

	// find the target path
	to := context.Local(env.ParseURL(r))
	parentTo := filepath.Dir(to)

	// if it's the same path, we throw an error
	if from == to {
		return errLinkSamePath
	}

	// make sure it doesn't exist
	if _, e := os.Stat(to); !os.IsNotExist(e) {
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
