package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/stream"
)

var errNoArgumentsProvided = exit.NewErrorWithCode("need at least one argument. use `ggman license` to view licensing information", exit.ExitGeneralArguments)

// Command returns the main ggman command
//
//nolint:contextcheck // don't need to pass down the context
func NewCommand(ctx context.Context, parameters env.Parameters, stream stream.IOStream) *cobra.Command {
	root := &cobra.Command{
		Use:     "ggman",
		Aliases: []string{os.Args[0]},
		Short:   "A golang tool that can manage all your git repositories. ",

		RunE: func(cmd *cobra.Command, args []string) error {
			return errNoArgumentsProvided
		},
	}
	root.SetContext(ctx)

	var flags env.Flags
	ggman.SetFlags(root, &flags)
	ggman.SetParameters(root, &parameters)

	// setup flags
	pflags := root.PersistentFlags()

	pflags.StringArrayVarP(&flags.For, "for", "f", flags.For, "filter list of repositories. Argument can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url")
	pflags.StringArrayVarP(&flags.FromFile, "from-file", "i", flags.FromFile, "filter list of repositories to only those matching filters from the given file. File should contain one filter per line, with common comment chars being ignored")
	pflags.BoolVarP(&flags.NoFuzzyFilter, "no-fuzzy-filter", "n", flags.NoFuzzyFilter, "disable fuzzy matching for filters")

	pflags.BoolVarP(&flags.Here, "here", "H", flags.Here, "filter list of repositories to only contain those that are in the current directory or subtree. alias for \"-p .\"")
	pflags.StringArrayVarP(&flags.Path, "path", "P", flags.Path, "filter list of repositories to only contain those that are in or under the specified path. may be used multiple times")

	pflags.BoolVarP(&flags.Dirty, "dirty", "d", flags.Dirty, "filter list of repositories to only contain repositories with uncommitted changes")
	pflags.BoolVarP(&flags.Clean, "clean", "c", flags.Clean, "filter list of repositories to only contain repositories without uncommitted changes")

	pflags.BoolVarP(&flags.Synced, "synced", "s", flags.Synced, "filter list of repositories to only contain those which are up-to-date with remote")
	pflags.BoolVarP(&flags.UnSynced, "unsynced", "u", flags.UnSynced, "filter list of repositories to only contain those not up-to-date with remote")

	pflags.BoolVarP(&flags.Tarnished, "tarnished", "t", flags.Tarnished, "filter list of repositories to only contain those that are dirty or unsynced")
	pflags.BoolVarP(&flags.Pristine, "pristine", "p", flags.Pristine, "filter list of repositories to only contain those that are clean and synced")

	// add all the commands
	root.AddCommand(
		NewLSCommand(),
		NewCanonCommand(),
	)

	// setup the command

	root.SetIn(stream.Stdin)
	root.SetOut(stream.Stdout)
	root.SetErr(stream.Stderr)

	root.SilenceErrors = true
	root.SilenceUsage = true

	return root
}
