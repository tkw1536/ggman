package cmd

//spellchecker:words runtime github cobra ggman
import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
)

//spellchecker:words wrapcheck

func NewVersionCommand() *cobra.Command {
	impl := new(version)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information and exit",
		Long:  "Version prints version information.",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type version struct{}

func (version) Exec(cmd *cobra.Command, args []string) error {
	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "ggman %s\n", ggman.BuildVersion); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "GOOS=%s GOARCH=%s %s\n", runtime.GOOS, runtime.GOARCH, runtime.Version()); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}

	if hash, err := ggman.BuildHash(); err == nil {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "sha256 %s\n", hash); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}

	return nil
}
