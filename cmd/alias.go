package cmd

import (
	"github.com/spf13/cobra"
	"go.tkw01536.de/goprogram/exit"
)

var errUnableToFindAliasedCommand = exit.NewErrorWithCode("unable to find aliased command", exit.ExitCommandArguments)

// NewAlias adds a new command that acts like an alias.
func NewAlias(root *cobra.Command, alias string, description string, expansion ...string) *cobra.Command {
	return &cobra.Command{
		Use:                alias,
		Short:              description,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// build the complete expanded command line
			full := make([]string, 0, len(expansion)+len(args))
			full = append(full, expansion...)
			full = append(full, args...)

			// setup state
			root.SetContext(cmd.Context())
			root.SetIn(cmd.InOrStdin())
			root.SetOut(cmd.OutOrStdout())
			root.SetErr(cmd.ErrOrStderr())
			root.SetArgs(full)

			// and execute it!
			return root.Execute()
		},
	}
}
