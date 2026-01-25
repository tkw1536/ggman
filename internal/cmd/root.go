package cmd

//spellchecker:words context github cobra ggman internal pkglib exit stream
import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/stream"
)

//spellchecker:words contextcheck unsynced pflags shellrc GGROOT
var (
	errInvalidFlags       = exit.NewErrorWithCode("unknown flags passed", env.ExitGeneralArguments)
	errGenericOutput      = exit.NewErrorWithCode("unknown output error", env.ExitGeneric)
	errGenericEnvironment = exit.NewErrorWithCode("failed to initialize environment", env.ExitInvalidEnvironment)
)

// Command returns the main ggman command
//
//nolint:contextcheck // don't need to pass down the context
func NewCommand(ctx context.Context, parameters env.Parameters) *cobra.Command {
	var flags env.Flags

	root := &cobra.Command{
		Use:     "ggman",
		Version: ggman.BuildVersion,
		Aliases: []string{os.Args[0]},
		Short:   "A golang tool that can manage all your git repositories. ",
		Long: `ggman is a golang tool that can manage all your git repositories.

#### What is ggman?

When you only have a couple of git repositories that you work on it is perfectly feasible to manage them by using git clone, git pull and friends.
However once the number of repositories grows beyond a small number this can become tedious:

- It is hard to find which folder a repository has been cloned to
- Getting an overview of what is cloned and what is not is hard
- It's not easily possible to perform actions on more than one repo at once, e.g. git pull

This is the problem ggman is designed to solve.
It allows one to:

- Maintain and expand a local directory structure of multiple repositories
- Run actions (such as git clone, git pull) on groups of repositories

#### Why ggman?

While similar tools exist these commonly have a lot of downsides:

- they enforce a flat directory structure;
- they are limited to one repository provider (such as GitHub or GitLab); or
- they are only available from within an IDE or GUI.

ggman considers these as major downsides.
The goals and principles of ggman are:

- to be command-line first;
- to be simple to install, configure and use;
- to encourage an obvious hierarchical directory structure, but remain fully functional with any directory structure;
- to remain free of forge- or provider-specific code; and
- to not store any repository-specific data outside of the repositories themselves (enabling the user to switch back to only git at any point).

#### How does ggman work?

ggman is split into several sub-commands, which are described in the relevant documentation pages.
ggman has the following general exit behavior:

- Exit Code 0: Everything went ok
- Exit Code 1: Command Parsing went ok, but a subcommand-specific error occurred
- Exit Code 2: The user asked for an unknown subcommand
- Exit Code 3: Command-independent argument parsing failed, e.g. an invalid '--for'
- Exit Code 4: Command-dependent argument parsing failed
- Exit Code 5: Invalid configuration
- Exit Code 6: Unable to parse a repository name`,
	}

	// setup flags
	{
		pflags := root.PersistentFlags()

		pflags.StringArrayVarP(&flags.For, "for", "F", flags.For, "filter list of repositories. Argument can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url")
		pflags.StringArrayVarP(&flags.FromFile, "from-file", "I", flags.FromFile, "filter list of repositories to only those matching filters from the given file. File should contain one filter per line, with common comment chars being ignored")
		pflags.BoolVarP(&flags.NoFuzzyFilter, "no-fuzzy-filter", "N", flags.NoFuzzyFilter, "disable fuzzy matching for filters")

		pflags.BoolVarP(&flags.Here, "here", "H", flags.Here, "filter list of repositories to only contain those that are in the current directory or subtree. alias for \"-p .\"")
		pflags.StringArrayVarP(&flags.Path, "path", "P", flags.Path, "filter list of repositories to only contain those that are in or under the specified path. may be used multiple times")

		pflags.BoolVarP(&flags.Dirty, "dirty", "D", flags.Dirty, "filter list of repositories to only contain repositories with uncommitted changes")
		pflags.BoolVarP(&flags.Clean, "clean", "C", flags.Clean, "filter list of repositories to only contain repositories without uncommitted changes")

		pflags.BoolVarP(&flags.Synced, "synced", "S", flags.Synced, "filter list of repositories to only contain those which are up-to-date with remote")
		pflags.BoolVarP(&flags.UnSynced, "unsynced", "U", flags.UnSynced, "filter list of repositories to only contain those not up-to-date with remote")

		pflags.BoolVarP(&flags.Tarnished, "tarnished", "T", flags.Tarnished, "filter list of repositories to only contain those that are dirty or unsynced")
		pflags.BoolVarP(&flags.Pristine, "pristine", "R", flags.Pristine, "filter list of repositories to only contain those that are clean and synced")
	}

	root.SetContext(ctx)
	root.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return fmt.Errorf("%w: %w", errInvalidFlags, err)
	})

	root.SilenceErrors = true
	root.SilenceUsage = true

	env.SetFlags(root, &flags)
	env.SetParameters(root, &parameters)

	// add all the commands
	root.AddCommand(
		NewCanonCommand(),
		NewCloneCommand(),
		NewCompsCommand(),
		NewEnvCommand(),
		NewExecCommand(),
		NewFetchCommand(),
		NewFindBranchCommand(),
		NewFindFileCommand(),
		NewFixCommand(),
		NewHereCommand(),
		NewLicenseCommand(),
		NewLinkCommand(),
		NewLsCommand(),
		NewLsrCommand(),
		NewPullCommand(),
		NewRelocateCommand(),
		NewShellrcCommand(),
		NewSweepCommand(),
		NewWhereCommand(),
		NewWebCommand(),
		NewURLCommand(),
		NewVersionCommand(),
		NewDocCommand(),
		NewCompletionCmd(),
	)

	for _, alias := range []struct {
		Command   *cobra.Command
		Expansion []string
	}{
		{
			Command: &cobra.Command{
				Use:   "git",
				Short: "Execute a git command using a native 'git' executable. ",
			},
			Expansion: []string{"exec", "--", "git"},
		},
		{
			Command: &cobra.Command{
				Use:   "root",
				Short: "Print the ggman root folder. ",
			},
			Expansion: []string{"env", "--raw", "GGROOT"},
		},
		{
			Command: &cobra.Command{
				Use:   "require",
				Short: "Require a remote git repository to be installed locally. ",
			},
			Expansion: []string{"clone", "--", "--force"},
		},
		{
			Command: &cobra.Command{
				Use:   "for",
				Short: "Filter repositories by a given filter. ",
			},
			Expansion: []string{"--for"},
		},
	} {
		addAlias(root, alias.Command, alias.Expansion...)
	}

	// wrap all the argument errors
	var wrapAllArgs func(cmd *cobra.Command)
	wrapAllArgs = func(cmd *cobra.Command) {
		cmd.Args = wrapArgs(cmd.Args)
		for _, child := range cmd.Commands() {
			wrapAllArgs(child)
		}
	}
	wrapAllArgs(root)

	return root
}

var errInvalidArguments = exit.NewErrorWithCode("invalid arguments passed", env.ExitCommandArguments)

// wrapArgs wraps a [cobra.PositionalArgs] error with an invalid arguments error.
// The wrapping occurs by calling [fmt.Errorf] with a string of "%w: %w" and [errInvalidArguments].
// If pos is nil, it is passed through as-is.
func wrapArgs(pos cobra.PositionalArgs) cobra.PositionalArgs {
	if pos == nil {
		return pos
	}

	return func(cmd *cobra.Command, args []string) error {
		err := pos(cmd, args)
		if err == nil {
			return nil
		}
		return fmt.Errorf("%w: %w", errInvalidArguments, err)
	}
}

// streamFromCommand returns a stream.IOStream from the given command.
func streamFromCommand(cmd *cobra.Command) stream.IOStream {
	return stream.IOStream{
		Stdout: cmd.OutOrStdout(),
		Stderr: cmd.ErrOrStderr(),
		Stdin:  cmd.InOrStdin(),
	}
}
