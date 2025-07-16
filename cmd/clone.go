package cmd

//spellchecker:words errors github ggman goprogram exit parser
import (
	"errors"
	"fmt"

	"al.essio.dev/pkg/shellescape"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/ggman/git"
	"go.tkw01536.de/goprogram/exit"
)

//spellchecker:words nolint wrapcheck canonicalize canonicalized

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
		Args []string `description:"additional arguments to pass to \"git clone\". "             positional-arg-name:"ARG"`
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

		Requirements: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
		},
	}
}

func (c *clone) AfterParse() error {
	if (c.Here || c.Local) && c.To != "" {
		return errCloneInvalidDestFlags
	}
	return nil
}

var (
	errCloneInvalidDestFlags = exit.NewErrorWithCode(`invalid destination: "--to" and "--here" may not be used together`, exit.ExitCommandArguments)
	errCloneInvalidDest      = exit.NewErrorWithCode("unable to determine local destination", exit.ExitGeneralArguments)
	errCloneLocalURI         = exit.NewErrorWithCode("invalid remote URI: invalid scheme, not a remote path", exit.ExitCommandArguments)
	errCloneAlreadyExists    = exit.NewErrorWithCode("unable to clone repository: another git repository already exists in target location", exit.ExitGeneric)
	errCloneNoArguments      = exit.NewErrorWithCode("external `git` not found, can not pass any additional arguments to `git clone`", exit.ExitGeneric)
	errCloneOther            = exit.NewErrorWithCode("", exit.ExitGeneric)

	errCloneNoComps = errors.New("unable to find components of URI")
)

func (c clone) Run(context ggman.Context) error {
	// grab the url to clone and make sure it is not local
	url := env.ParseURL(c.Positional.URL)
	if url.IsLocal() {
		return fmt.Errorf("%q: %w", c.Positional.URL, errCloneLocalURI)
	}

	// find the remote and local paths to clone to / from
	remote := c.Positional.URL
	if !c.Exact {
		remote = context.Environment.Canonical(url)
	}
	local, err := c.dest(context, url)
	if err != nil {
		return fmt.Errorf("%q: %w: %w", c.Positional.URL, errCloneInvalidDest, err)
	}

	// do the actual cloning!
	if _, err := context.Printf("Cloning %q into %q ...\n", remote, local); err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	switch err := context.Environment.Git.Clone(context.IOStream, remote, local, c.Positional.Args...); {
	case err == nil:
		return nil
	case errors.Is(err, git.ErrCloneAlreadyExists):
		if c.Force {
			_, err := context.Println("Clone already exists in target location, done.")
			if err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
			return nil
		}
		return errCloneAlreadyExists
	case errors.Is(err, git.ErrArgumentsUnsupported):
		return fmt.Errorf("%w: %v", errCloneNoArguments, shellescape.QuoteCommand(c.Positional.Args))
	default:
		return fmt.Errorf("%w%w", errCloneOther, err)
	}
}

// dest returns the destination path to clone the repository into.
func (c clone) dest(context ggman.Context, url env.URL) (path string, err error) {
	switch {
	case c.Here || c.Local: // clone into directory named automatically
		comps := url.Components()
		if len(comps) == 0 {
			return "", errCloneNoComps
		}
		path, err = context.Environment.Abs(comps[len(comps)-1])
	case c.To != "": // clone directory into a directory
		path, err = context.Environment.Abs(c.To)
	default: // normal clone!
		path, err = context.Environment.Local(url)
	}

	if err != nil {
		return "", fmt.Errorf("%q: failed to get destination: %w", url.String(), err)
	}
	return path, nil
}
