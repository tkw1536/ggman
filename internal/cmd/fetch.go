package cmd

//spellchecker:words github cobra ggman internal pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

func NewFetchCommand() *cobra.Command {
	impl := new(fetch)

	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "Run \"git fetch --all\" on all repositories",
		Long:  `Fetch is the equivalent of running 'git fetch --all' on all repositories.`,
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type fetch struct{}

var errFetchCustom = exit.NewErrorWithCode("", env.ExitGeneric)

func (fetch) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		AllowsFilter: true,
		NeedsRoot:    true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	hasError := false

	// iterate over all the repositories, and run git fetch
	for _, repo := range environment.Repos(cmd.Context(), true) {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Fetching %q\n", repo); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
		if e := environment.Git.Fetch(cmd.Context(), streamFromCommand(cmd), repo); e != nil {
			if _, err := fmt.Fprintln(cmd.ErrOrStderr(), e.Error()); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
			hasError = true
		}
	}

	if hasError {
		return errFetchCustom
	}
	return nil
}
