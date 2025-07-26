package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
)

func NewCompsCommand() *cobra.Command {
	comps := &cobra.Command{
		Use:   "comps",
		Short: "print the components of a URL",
		Args:  cobra.ExactArgs(1),

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ggman.SetRequirements(cmd, &env.Requirement{})
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			url := args[0]
			for _, comp := range env.ComponentsOf(url) {
				if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s\n", comp); err != nil {
					return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
				}
			}
			return nil
		},
	}

	return comps
}
