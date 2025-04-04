package cmd

//spellchecker:words errors github ggman goprogram exit parser
import (
	"errors"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
)

// Clone is the 'ggman clone' command.
//
// Clone clones the remote repository in the first argument into the path described to by 'ggman where'.
// It canonizes the url before cloning it.
// It optionally takes any argument that would be passed to the normal invocation of a git command.
//
// When 'git' is not available on the system ggman is running on, additional arguments may not be supported.
var Clone ggman.Command = clone{}

type clone struct {
	Positional struct {
		URL  string   `description:"URL of repository clone. Will be canonicalized by default. " positional-arg-name:"URL" required:"1-1"`
		Args []string `description:"additional arguments to pass to \"git clone\"."              positional-arg-name:"ARG"`
	} `positional-args:"true"`
	Force bool   `description:"do not complain when a repository already exists in the target directory" long:"force"     short:"f"`
	Local bool   `description:"alias of \"--here\""                                                      long:"local"     short:"l"`
	Exact bool   `description:"don't canonicalize URL before cloning and use exactly the passed URL"     long:"exact-url" short:"e"`
	Here  bool   `description:"clone into an appropriately named subdirectory of the current directory"  long:"here"`
	To    string `description:"clone repository into specified directory"                                long:"to"        short:"t"`
}

func (clone) Description() ggman.Description {
	return ggman.Description{
		Command:     "clone",
		Description: "clone a repository into a path described by \"ggman where\"",

		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},

		Requirements: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
		},
	}
}

var errInvalidDest = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  `invalid destination: "--to" and "--here" may not be used together`,
}

func (c *clone) AfterParse() error {
	if (c.Here || c.Local) && c.To != "" {
		return errInvalidDest
	}
	return nil
}

var errCloneInvalidDest = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "unable to determine local destination for %q",
}

var errCloneLocalURI = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "invalid remote URI %q: invalid scheme, not a remote path",
}

var errCloneAlreadyExists = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "unable to clone repository: another git repository already exists in target location",
}

var errCloneNoArguments = exit.Error{
	ExitCode: exit.ExitGeneric,
	Message:  "external `git` not found, can not pass any additional arguments to `git clone`",
}

var errCloneOther = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (c clone) Run(context ggman.Context) error {
	// grab the url to clone and make sure it is not local
	url := env.ParseURL(c.Positional.URL)
	if url.IsLocal() {
		return errCloneLocalURI.WithMessageF(c.Positional.URL)
	}

	// find the remote and local paths to clone to / from
	remote := c.Positional.URL
	if !c.Exact {
		remote = context.Environment.Canonical(url)
	}
	local, err := c.dest(context, url)
	if err != nil {
		return errCloneInvalidDest.WithMessageF(c.Positional.URL).WrapError(err)
	}

	// do the actual cloning!
	if _, err := context.Printf("Cloning %q into %q ...\n", remote, local); err != nil {
		return ggman.ErrGenericOutput.WrapError(err)
	}
	switch err := context.Environment.Git.Clone(context.IOStream, remote, local, c.Positional.Args...); {
	case err == nil:
		return nil
	case errors.Is(err, git.ErrCloneAlreadyExists):
		if c.Force {
			_, err := context.Println("Clone already exists in target location, done.")
			return ggman.ErrGenericOutput.WrapError(err)
		}
		return errCloneAlreadyExists
	case errors.Is(err, git.ErrArgumentsUnsupported):
		return errCloneNoArguments
	default:
		return errCloneOther.WithMessage(err.Error())
	}
}

var errCloneNoComps = errors.New("unable to find components of URI")

// dest returns the destination path to clone the repository into.
func (c clone) dest(context ggman.Context, url env.URL) (string, error) {
	if c.Here || c.Local { // clone into directory named automatically
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
