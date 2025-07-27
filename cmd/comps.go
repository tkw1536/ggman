package cmd

//spellchecker:words github cobra ggman
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
)

func NewCompsCommand() *cobra.Command {
	impl := new(comps)

	cmd := &cobra.Command{
		Use:   "comps URL",
		Short: "print the components of a URL",
		Long: `When invoked, it prints the components of the first argument passed to it.
Each component is printed on a separate line of standard output.`,
		Args: cobra.ExactArgs(1),

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	return cmd
}

//spellchecker:words nolint wrapcheck

type comps struct {
	Positional struct {
		URL env.URL
	}
}

func (c *comps) AfterParse(cmd *cobra.Command, args []string) error {
	c.Positional.URL = env.ParseURL(args[0])
	return nil
}

func (c *comps) Exec(cmd *cobra.Command, args []string) error {
	for _, comp := range c.Positional.URL.Components() {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), comp); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}

	return nil
}
