package cmd

//spellchecker:words github cobra ggman internal foreach pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/ggman/internal/foreach"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words wrapcheck

func NewPullCommand() *cobra.Command {
	impl := new(pull)

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Run \"git pull\" on locally cloned repositories",
		Long: `Pull runs 'git pull' on all repositories.

The '--parallel' flag sets how many pulls to run concurrently.
The default is 1 (sequential). Use 0 for no limit.

After all repositories have been processed, Pull prints a summary line
with the number of successful and failed pulls.`,
		Args: cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.IntVarP(&impl.Parallel, "parallel", "p", 1, "number of pulls to run in parallel, 0 for no limit")

	return cmd
}

type pull struct {
	Parallel int
}

var errPullCustom = exit.NewErrorWithCode("", env.ExitGeneric)

func (p *pull) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		AllowsFilter: true,
		NeedsRoot:    true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	ok, failed, err := foreach.Repos(
		cmd.Context(),
		streamFromCommand(cmd),
		environment.Repos(cmd.Context(), true),
		p.Parallel,
		"Pulling",
		environment.Git.Pull,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Pulled %d repositories, %d failed\n", ok, failed); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}

	if failed > 0 {
		return errPullCustom
	}

	return nil
}
