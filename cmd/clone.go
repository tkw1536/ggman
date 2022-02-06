package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/program"
)

// Clone is the 'ggman clone' command.
//
// Clone clones the remote repository in the first argument into the path described to by 'ggman where'.
// It canonizes the url before cloning it.
// It optionally takes any argument that would be passed to the normal invocation of a git command.
//
// When 'git' is not available on the system ggman is running on, additional arguments may not be supported.
var Clone program.Command = &clone{}

type clone struct {
	Force bool `short:"f" long:"force" description:"Don't complain when a repository already exists in the target directory"`
}

func (*clone) BeforeRegister(program *program.Program) {}

func (*clone) Description() program.Description {
	return program.Description{
		Name:        "clone",
		Description: "Clone a repository into a path described by 'ggman where'",

		SkipUnknownOptions: true,

		PosArgsMin: 1,
		PosArgsMax: -1,
		PosArgName: "ARG",

		PosArgDescription: "URL of repository and arguments to pass to 'git clone'",

		Environment: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
		},
	}
}

func (*clone) AfterParse() error {
	return nil
}

var errCloneInvalidURI = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Invalid remote URI %q: Invalid scheme, not a remote path. ",
}

var errCloneAlreadyExists = ggman.Error{
	ExitCode: ggman.ExitGeneric,
	Message:  "Unable to clone repository: Another git repository already exists in target location. ",
}

var errCloneNoArguments = ggman.Error{
	ExitCode: ggman.ExitGeneric,
	Message:  "External 'git' not found, can not pass any additional arguments to 'git clone'. ",
}

var errCloneOther = ggman.Error{
	ExitCode: ggman.ExitGeneric,
}

func (c *clone) Run(context program.Context) error {
	// grab the url to clone and make sure it is not local
	url := context.URLV(0)
	if url.IsLocal() {
		return errCloneInvalidURI.WithMessageF(context.Args[0])
	}

	// find the remote and local paths to clone to / from
	remote := context.Canonical(url)
	local, err := context.Local(url)
	if err != nil {
		return err
	}

	// do the actual cloning!
	context.Printf("Cloning %q into %q ...\n", remote, local)
	switch err := context.Git.Clone(context.IOStream, remote, local, context.Args[1:]...); err {
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
