package cmd

//spellchecker:words github cobra ggman
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
)

//spellchecker:words worktree nolint wrapcheck

func NewHereCommand() *cobra.Command {
	impl := new(here)

	cmd := &cobra.Command{
		Use:   "here",
		Short: "print the root path to the repository in the current repository",
		Long:  `'ggman here' prints the path to the root of the repository in the current working directory to standard output.`,
		Args:  cobra.NoArgs,

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
	environment, err := ggman.GetEnv(cmd, env.Requirement{
		NeedsRoot: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericEnvironment, err)
	}

	root, worktree, err := environment.At(".")
	if err != nil {
		return fmt.Errorf("%w: %w", env.ErrUnableLocalPath, err)
	}

	if _, err := fmt.Fprintln(cmd.OutOrStdout(), root); err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	if h.Tree {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), worktree); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}

	return nil
}
