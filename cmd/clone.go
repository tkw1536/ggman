package cmd

import (
	"errors"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/program/exit"
)

// Clone is the 'ggman clone' command.
//
// Clone clones the remote repository in the first argument into the path described to by 'ggman where'.
// It canonizes the url before cloning it.
// It optionally takes any argument that would be passed to the normal invocation of a git command.
//
// When 'git' is not available on the system ggman is running on, additional arguments may not be supported.
var Clone ggman.Command = &clone{}

type clone struct {
	Force bool   `short:"f" long:"force" description:"Don't complain when a repository already exists in the target directory"`
	Local bool   `short:"l" long:"local" description:"Clone into an appropriately named subdirectory of the current directory"`
	To    string `short:"t" long:"to" description:"Clone repository into specified directory"`
}

func (*clone) BeforeRegister(program *ggman.Program) {}

func (*clone) Description() ggman.Description {
	return ggman.Description{
		Name:        "clone",
		Description: "Clone a repository into a path described by 'ggman where'",

		SkipUnknownOptions: true,

		PosArgsMin: 1,
		PosArgsMax: -1,
		PosArgName: "ARG",

		PosArgDescription: "URL of repository and arguments to pass to 'git clone'",

		Requirements: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
		},
	}
}

var errInvalidDest = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Invalid destination: '--to' and '--local' may not be used together. ",
}

func (c *clone) AfterParse() error {
	if c.Local && c.To != "" {
		return errInvalidDest
	}
	return nil
}

var errCloneInvalidDest = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to determine local destination for %q: %s",
}

var errCloneLocalURI = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Invalid remote URI %q: Invalid scheme, not a remote path. ",
}

var errCloneAlreadyExists = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "Unable to clone repository: Another git repository already exists in target location. ",
}

var errCloneNoArguments = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "External 'git' not found, can not pass any additional arguments to 'git clone'. ",
}

var errCloneOther = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (c *clone) Run(context ggman.Context) error {
	// grab the url to clone and make sure it is not local
	url := ggman.URLV(context, 0)
	if url.IsLocal() {
		return errCloneLocalURI.WithMessageF(context.Args.Pos[0])
	}

	// find the remote and local paths to clone to / from
	remote := context.Environment.Canonical(url)
	local, err := c.dest(context, url)
	if err != nil {
		return errCloneInvalidDest.WithMessageF(context.Args.Pos[0], err)
	}

	// do the actual cloning!
	context.Printf("Cloning %q into %q ...\n", remote, local)
	switch err := context.Environment.Git.Clone(context.IOStream, remote, local, context.Args.Pos[1:]...); err {
	case nil:
		return nil
	case git.ErrCloneAlreadyExists:
		if c.Force {
			context.Println("Clone already exists in target location, done.")
			return nil
		}
		return errCloneAlreadyExists
	case git.ErrArgumentsUnsupported:
		return errCloneNoArguments
	default:
		return errCloneOther.WithMessage(err.Error())
	}
}

var errCloneNoComps = errors.New("unable to find components of URI")

// dest returns the destination path to clone the repository into
func (c clone) dest(context ggman.Context, url env.URL) (string, error) {
	if c.Local { // clone into directory named automatically
		comps := url.Components()
		if len(comps) == 0 {
			return "", errCloneNoComps
		}
		return context.Environment.Abs(comps[len(comps)-1])
	}

	if c.To != "" { // clone directory into a directory
		return context.Environment.Abs(c.To)
	}

	// normal clone!
	return context.Environment.Local(url)
}
