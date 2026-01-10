package cmd

//spellchecker:words essio shellescape github cobra
import (
	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
)

// addAlias adds the alias command to act as an alias for another command of root.
// root is the root command the alias is to be added to.
// expansion is the expansion, including flags, it will act as an alias for.
func addAlias(root *cobra.Command, alias *cobra.Command, expansion ...string) {
	alias.DisableFlagParsing = true
	if alias.Long == "" {
		alias.Long = "alias for '" + shellescape.QuoteCommand(append([]string{root.Name()}, expansion...)) + "'"
	}
	addModifier(root, alias, func(args []string) ([]string, error) {
		argv := make([]string, 0, len(expansion)+len(args))
		argv = append(argv, expansion...)
		argv = append(argv, args...)
		return argv, nil
	})
}

// addModifier adds cmd as a modifier command using the transform function.
// Instead of being invoked normally, a transform command does not parse arguments and instead transforms them using the given function.
// It then invokes the root command with these arguments.
func addModifier(root *cobra.Command, cmd *cobra.Command, transform func(args []string) ([]string, error)) {
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
	root.AddCommand(cmd)
}
