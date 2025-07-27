package cmd

import (
	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/goprogram/exit"
)

var errUnableToFindAliasedCommand = exit.NewErrorWithCode("unable to find aliased command", exit.ExitCommandArguments)

// NewAlias configures a command to act as an alias for another command.
// root is the root command the alias is to be added to.
// expansion is the expansion, including flags, it will act as an alias for.
func NewAlias(root *cobra.Command, alias *cobra.Command, expansion ...string) *cobra.Command {
	alias.DisableFlagParsing = true
	if alias.Long == "" {
		alias.Long = "alias for '" + shellescape.QuoteCommand(append([]string{root.Name()}, expansion...)) + "'"
	}
	alias.RunE = func(cmd *cobra.Command, args []string) error {
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
	}
	return alias
}
