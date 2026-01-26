package cmd

//spellchecker:words github cobra ggman internal pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words canonicalized wrapcheck

func NewLsrCommand() *cobra.Command {
	impl := new(lsr)

	cmd := &cobra.Command{
		Use:   "lsr",
		Short: "List remote URLs to all repositories",
		Long: `Lsr prints the remote URLs of all repositories.

The '--canonical' flag prints canonical URLs instead of the original URLs.`,
		Args: cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Canonical, "canonical", "c", false, "print canonicalized URLs")

	return cmd
}

type lsr struct {
	Canonical bool
}

var errLSRInvalidCanfile = exit.NewErrorWithCode("failed to parse CANFILE", env.ExitInvalidEnvironment)

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
	for _, repo := range environment.Repos(cmd.Context(), true) {
		remote, err := environment.Git.GetRemote(cmd.Context(), repo, "")
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
