package cmd

//spellchecker:words github cobra ggman internal pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words wrapcheck

func NewPullCommand() *cobra.Command {
	impl := new(pull)

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "Run \"git pull\" on locally cloned repositories",
		Long:  "Pull is the equivalent of running 'git pull' on all locally installed repositories.",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type pull struct{}

var errPullCustom = exit.NewErrorWithCode("", exit.ExitGeneric)

func (pull) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		AllowsFilter: true,
		NeedsRoot:    true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	hasError := false
	for _, repo := range environment.Repos(true) {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Pulling %q\n", repo); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
		if e := environment.Git.Pull(streamFromCommand(cmd), repo); e != nil {
			if _, err := fmt.Fprintln(cmd.ErrOrStderr(), e); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
			hasError = true
		}
	}

	if hasError {
		return errPullCustom
	}

	return nil
}
