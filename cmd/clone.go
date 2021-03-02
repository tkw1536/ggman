package cmd

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/program"
)

// Clone is the 'ggman clone' command.
//
// Clone clones the repository in the firsty argument into the path described to by 'ggman where'.
// It canonizes the url before cloning it.
// It optionally takes any argument that would be passed to the normal invocation of a git command.
//
// When 'git' is not available on the system ggman is running on, additional arguments may not be supported.
var Clone program.Command = clone{}

type clone struct{}

func (clone) Name() string {
	return "clone"
}

func (clone) Options() program.Options {
	return program.Options{
		SkipUnknownFlags: true,

		MinArgs: 1,
		MaxArgs: -1,
		Metavar: "ARG",

		UsageDescription: "URL of repository and arguments to pass to 'git clone'",

		Environment: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
		},
	}
}

func (clone) AfterParse() error {
	return nil
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

func (clone) Run(context program.Context) error {
	// find the remote
	remote := context.URLV(0)
	remoteURI := context.Canonical(remote)
	clonePath := context.Local(remote)

	// do the clone command!
	context.Printf("Cloning %q into %q ...\n", remoteURI, clonePath)
	switch err := context.Git.Clone(context.IOStream, remoteURI, clonePath, context.Args[1:]...); err {
	case nil:
		return nil
	case git.ErrCloneAlreadyExists:
		return errCloneAlreadyExists
	case git.ErrArgumentsUnsupported:
		return errCloneNoArguments
	default:
		return errCloneOther.WithMessage(err.Error())
	}
}
