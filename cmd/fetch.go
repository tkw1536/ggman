package cmd

//spellchecker:words github cobra ggman pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/pkglib/exit"
)

func NewFetchCommand() *cobra.Command {
	impl := new(fetch)

	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "run \"git fetch --all\" on locally cloned repositories",
		Long:  `'ggman fetch' is the equivalent of running 'git fetch --all' on all locally cloned repositories.`,
		Args:  cobra.NoArgs,

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	return cmd
}

type fetch struct{}

func (fetch) AfterParse(cmd *cobra.Command, args []string) error {
	return nil
}

var errFetchCustom = exit.NewErrorWithCode("", exit.ExitGeneric)

func (fetch) Exec(cmd *cobra.Command, args []string) error {
	environment, err := ggman.GetEnv(cmd, env.Requirement{
		AllowsFilter: true,
		NeedsRoot:    true,
	})
	if err != nil {
		return err
	}

	hasError := false

	// iterate over all the repositories, and run git fetch
	for _, repo := range environment.Repos(true) {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Fetching %q\n", repo); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if e := environment.Git.Fetch(streamFromCommand(cmd), repo); e != nil {
			if _, err := fmt.Fprintln(cmd.ErrOrStderr(), e.Error()); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
			hasError = true
		}
	}

	if hasError {
		return errFetchCustom
	}
	return nil
}
