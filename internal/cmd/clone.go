package cmd

//spellchecker:words errors essio shellescape github cobra ggman internal pkglib exit
import (
	"errors"
	"fmt"
	"os"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/ggman/internal/git"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
)

//spellchecker:words canonicalize GGROOT

func NewCloneCommand() *cobra.Command {
	impl := new(clone)

	cmd := &cobra.Command{
		Use:   "clone URL [ARGS...]",
		Short: "Clone a repository into the local directory structure",
		Long: `Clone clones a repository into its location within '$GGROOT'.

For example

    ggman clone git@github.com:hello/world.git

clones into '$GGROOT/github.com/hello/world'.
Any URL format works; the canonical URL is used for cloning.

For example

    ggman clone https://github.com/hello/world.git

produces the same result.

The '--exact-url' flag uses the provided URL without canonicalization:

    ggman clone --exact-url https://github.com/hello/world.git

Additional arguments can be passed to git after '--':

    ggman clone --exact-url https://github.com/hello/world.git -- --branch dev --depth 2

This executes 'git clone git@github.com:hello/world.git --branch dev --depth 2'.
The '--' separator distinguishes ggman flags from git flags.`,
		Args: cobra.ArbitraryArgs,

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Force, "force", "f", false, "do not complain when a repository already exists in the target directory. Incompatible with '--overwrite'")
	flags.BoolVarP(&impl.Overwrite, "overwrite", "o", false, "if the local directory already exists delete it before attempting to clone again. Incompatible with '--force'")
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
	Force     bool
	Overwrite bool
	Local     bool
	Exact     bool
	Plain     bool
	To        string
}

func (c *clone) ParseArgs(cmd *cobra.Command, args []string) error {
	if c.Local {
		c.Plain = true
	}
	if (c.Plain) && c.To != "" {
		return errCloneInvalidDestFlags
	}

	if c.Overwrite && c.Force {
		return errCloneInvalidForceFlags
	}

	c.Positional.URL = args[0]
	c.Positional.Args = args[1:]

	return nil
}

var (
	errCloneInvalidDestFlags  = exit.NewErrorWithCode(`invalid destination: "--to" and "--plain" may not be used together`, env.ExitCommandArguments)
	errCloneInvalidForceFlags = exit.NewErrorWithCode(`"--overwrite" and "--force" are incompatible`, env.ExitCommandArguments)
	errCloneInvalidDest       = exit.NewErrorWithCode("failed to determine local destination", env.ExitGeneralArguments)
	errCloneCheckDest         = exit.NewErrorWithCode("failed to check if destination is a directory", env.ExitGeneric)
	errCloneDeleteDest        = exit.NewErrorWithCode("failed to delete existing directory", env.ExitGeneric)
	errCloneLocalURI          = exit.NewErrorWithCode("invalid remote URI: invalid scheme, not a remote path", env.ExitCommandArguments)
	errCloneAlreadyExists     = exit.NewErrorWithCode("failed to clone repository: another git repository already exists in target location", env.ExitGeneric)
	errCloneNoArguments       = exit.NewErrorWithCode(`failed to pass arguments: external "git" not found`, env.ExitGeneric)
	errCloneOther             = exit.NewErrorWithCode("", env.ExitGeneric)

	errCloneNoComps = errors.New("failed to find components of URI")
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

	if c.Overwrite {
		isDir, err := fsx.IsDirectory(local, false)
		if err != nil {
			return fmt.Errorf("%w: %w", errCloneCheckDest, err)
		}
		if isDir {
			_, err := fmt.Fprintf(cmd.OutOrStdout(), "Deleting existing directory %q\n", local)
			if err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
			if err := os.RemoveAll(local); err != nil {
				return fmt.Errorf("%w: %w", errCloneDeleteDest, err)
			}
		}
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
