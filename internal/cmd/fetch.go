package cmd

//spellchecker:words github cobra ggman internal foreach pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/ggman/internal/foreach"
	"go.tkw01536.de/pkglib/exit"
)

func NewFetchCommand() *cobra.Command {
	impl := new(fetch)

	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Run \"git fetch --all\" on all repositories",
		Long: `Fetch runs 'git fetch --all' on all repositories.

The '--parallel' flag sets how many fetches to run concurrently.
The default is 1 (sequential). Use 0 for no limit.

After all repositories have been processed, Fetch prints a summary line
with the number of successful and failed fetches.`,
		Args: cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.IntVarP(&impl.Parallel, "parallel", "p", 1, "number of fetches to run in parallel, 0 for no limit")

	return cmd
}

type fetch struct {
	Parallel int
}

var errFetchCustom = exit.NewErrorWithCode("", env.ExitGeneric)

func (f *fetch) Exec(cmd *cobra.Command, args []string) error {
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
		f.Parallel,
		"Fetching",
		environment.Git.Fetch,
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Fetched %d repositories, %d failed\n", ok, failed); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}

	if failed > 0 {
		return errFetchCustom
	}
	return nil
}
