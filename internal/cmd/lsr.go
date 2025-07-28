package cmd

//spellchecker:words github cobra ggman internal pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words canonicalized CANFILE nolint wrapcheck

func NewLsrCommand() *cobra.Command {
	impl := new(lsr)

	cmd := &cobra.Command{
		Use:   "lsr",
		Short: "list remote URLs to all locally cloned repositories",
		Long:  "When provided, instead of printing the urls directly, prints the canonical remotes of all repositories.",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Canonical, "canonical", "c", false, "print canonicalized URLs")

	return cmd
}

type lsr struct {
	Canonical bool
}

var errLSRInvalidCanfile = exit.NewErrorWithCode("invalid CANFILE found", env.ExitInvalidEnvironment)

func (l *lsr) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		AllowsFilter: true,
		NeedsRoot:    true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	var lines env.CanFile
	if l.Canonical {
		var err error
		if lines, err = environment.LoadDefaultCANFILE(); err != nil {
			return fmt.Errorf("%w: %w", errLSRInvalidCanfile, err)
		}
	}

	// and print them
	for _, repo := range environment.Repos(true) {
		remote, err := environment.Git.GetRemote(repo, "")
		if err != nil {
			continue
		}
		if l.Canonical {
			remote = env.ParseURL(remote).CanonicalWith(lines)
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), remote); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}

	return nil
}
