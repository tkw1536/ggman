package cmd

//spellchecker:words github cobra ggman goprogram exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

//spellchecker:words nolint wrapcheck

func NewPullCommand() *cobra.Command {
	impl := new(pull)

	cmd := &cobra.Command{
		Use:   "pull",
		Short: "run \"git pull\" on locally cloned repositories",
		Long:  "'ggman pull' is the equivalent of running 'git pull' on all locally installed repositories.",
		Args:  cobra.NoArgs,

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	return cmd
}

type pull struct{}

func (pull) Description() ggman.Description {
	return ggman.Description{
		Command:     "pull",
		Description: "run \"git pull\" on locally cloned repositories",

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

var errPullCustom = exit.NewErrorWithCode("", exit.ExitGeneric)

func (*pull) AfterParse(cmd *cobra.Command, args []string) error {
	return nil
}

func (pull) Exec(cmd *cobra.Command, args []string) error {
	environment, err := ggman.GetEnv(cmd, env.Requirement{
		AllowsFilter: true,
		NeedsRoot:    true,
	})
	if err != nil {
		return err
	}

	hasError := false
	for _, repo := range environment.Repos(true) {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Pulling %q\n", repo); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if e := environment.Git.Pull(streamFromCommand(cmd), repo); e != nil {
			if _, err := fmt.Fprintln(cmd.ErrOrStderr(), e); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
			hasError = true
		}
	}

	if hasError {
		return errPullCustom
	}

	return nil
}
