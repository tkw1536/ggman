package cmd

//spellchecker:words essio shellescape github cobra goprogram exit
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
	return newModifier(root, alias, func(args []string) ([]string, error) {
		argv := make([]string, 0, len(expansion)+len(args))
		argv = append(argv, expansion...)
		argv = append(argv, args...)
		return argv, nil
	})
}

// newModifier configures cmd to be a modifier command using the transform function.
// Instead of being invoked normally, a transform command does not parse arguments and instead transforms them using the given function.
// It then invokes the root command with these arguments.
func newModifier(root *cobra.Command, cmd *cobra.Command, transform func(args []string) ([]string, error)) *cobra.Command {
	cmd.DisableFlagParsing = true
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		argv, err := transform(args)
		if err != nil {
			return err
		}

		// setup state
		root.SetContext(cmd.Context())
		root.SetIn(cmd.InOrStdin())
		root.SetOut(cmd.OutOrStdout())
		root.SetErr(cmd.ErrOrStderr())
		root.SetArgs(argv)

		// and execute it!
		return root.Execute()
	}
	return cmd
}
