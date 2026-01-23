package cmd

//spellchecker:words errors essio shellescape github cobra ggman internal pkglib exit
import (
	"errors"
	"fmt"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/ggman/internal/git"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words canonicalize GGROOT

func NewCloneCommand() *cobra.Command {
	impl := new(clone)

	cmd := &cobra.Command{
		Use:   "clone URL [ARGS...]",
		Short: "Clone a repository into the local directory structure",
		Long: `Clone a new repository into the respective location using 'ggman clone' with the repository URL as the argument, for example:

  ggman clone git@github.com:hello/world.git

which will clone the the hello world repository into '$GGROOT/github.com/hello/world'.
This cloning not only works for the canonical repository url, but for any other url as well.
For example:

  ggman clone https://github.com/hello/world.git

will do the same as the above command.

When it is not desired that the canonical URL should be used, pass the '--exact-url' flag:

  ggman clone --exact-url https://github.com/hello/world.git

This will clone using the exact url into the same folder as above.

If ggman has access to a real 'git' executable, it is also possible to pass additional arguments to it.
For example:

  ggman clone --exact-url https://github.com/hello/world.git -- --branch dev --depth 2

will execute the command ` + "`" + `git clone git@github.com:hello/world.git --branch dev --depth 2` + "`" + ` under the hood.
The extra '--' is needed to allow ggman to separate the internal flags from the external flags.`,
		Args: cobra.ArbitraryArgs,

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Force, "force", "f", false, "do not complain when a repository already exists in the target directory")
	flags.BoolVarP(&impl.Local, "local", "l", false, "alias of \"--plain\"")
	flags.BoolVarP(&impl.Exact, "exact-url", "e", false, "don't canonicalize URL before cloning and use exactly the passed URL")
	flags.BoolVar(&impl.Plain, "plain", false, "clone like a standard git would: into an appropriately named subdirectory of the current directory")
	flags.StringVarP(&impl.To, "to", "t", "", "clone repository into specified directory")

	return cmd
}

type clone struct {
	Positional struct {
		URL  string
		Args []string
	}
	Force bool
	Local bool
	Exact bool
	Plain bool
	To    string
}

func (c *clone) ParseArgs(cmd *cobra.Command, args []string) error {
	if c.Local {
		c.Plain = true
	}
	if (c.Plain) && c.To != "" {
		return errCloneInvalidDestFlags
	}

	c.Positional.URL = args[0]
	c.Positional.Args = args[1:]

	return nil
}

var (
	errCloneInvalidDestFlags = exit.NewErrorWithCode(`invalid destination: "--to" and "--plain" may not be used together`, env.ExitCommandArguments)
	errCloneInvalidDest      = exit.NewErrorWithCode("unable to determine local destination", env.ExitGeneralArguments)
	errCloneLocalURI         = exit.NewErrorWithCode("invalid remote URI: invalid scheme, not a remote path", env.ExitCommandArguments)
	errCloneAlreadyExists    = exit.NewErrorWithCode("unable to clone repository: another git repository already exists in target location", env.ExitGeneric)
	errCloneNoArguments      = exit.NewErrorWithCode("external `git` not found, can not pass any additional arguments to `git clone`", env.ExitGeneric)
	errCloneOther            = exit.NewErrorWithCode("", env.ExitGeneric)

	errCloneNoComps = errors.New("unable to find components of URI")
)

func (c *clone) Exec(cmd *cobra.Command, args []string) error {
	// get the environment
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot:    true,
		NeedsCanFile: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	// grab the url to clone and make sure it is not local
	url := env.ParseURL(c.Positional.URL)
	if url.IsLocal() {
		return fmt.Errorf("%q: %w", c.Positional.URL, errCloneLocalURI)
	}

	// find the remote and local paths to clone to / from
	remote := c.Positional.URL
	if !c.Exact {
		remote = environment.Canonical(url)
	}
	local, err := c.dest(environment, url)
	if err != nil {
		return fmt.Errorf("%q: %w: %w", c.Positional.URL, errCloneInvalidDest, err)
	}

	// do the actual cloning!
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Cloning %q into %q ...\n", remote, local); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	switch err := environment.Git.Clone(cmd.Context(), streamFromCommand(cmd), remote, local, c.Positional.Args...); {
	case err == nil:
		return nil
	case errors.Is(err, git.ErrCloneAlreadyExists):
		if c.Force {
			_, err := fmt.Fprintln(cmd.OutOrStdout(), "Clone already exists in target location, done.")
			if err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
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
func (c *clone) dest(environment *env.Env, url env.URL) (path string, err error) {
	switch {
	case c.Plain: // clone into directory named automatically
		comps := url.Components()
		if len(comps) == 0 {
			return "", errCloneNoComps
		}
		path, err = environment.Abs(comps[len(comps)-1])
	case c.To != "": // clone directory into a directory
		path, err = environment.Abs(c.To)
	default: // normal clone!
		path, err = environment.Local(url)
	}

	if err != nil {
		return "", fmt.Errorf("%q: failed to get destination: %w", url.String(), err)
	}
	return path, nil
}
