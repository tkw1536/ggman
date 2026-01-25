package cmd

//spellchecker:words github cobra ggman internal
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
)

//spellchecker:words worktree wrapcheck

func NewHereCommand() *cobra.Command {
	impl := new(here)

	cmd := &cobra.Command{
		Use:   "here",
		Short: "Print the root path to the current ggman-controlled repository",
		Long: `Here prints the current ggman-controlled repository to standard output.
This is the repository in the working directory, or a parent thereof.

The '--tree' flag additionally prints the path relative to the git worktree root.`,
		Args: cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Tree, "tree", "t", false, "also print the current HEAD reference and relative path to the root of the git worktree")

	return cmd
}

type here struct {
	Tree bool
}

func (h *here) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	root, worktree, err := environment.At(cmd.Context(), ".")
	if err != nil {
		return fmt.Errorf("%w: %w", env.ErrUnableLocalPath, err)
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), root); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	if h.Tree {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), worktree); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}

	return nil
}
