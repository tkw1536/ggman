package cmd

//spellchecker:words runtime github cobra ggman
import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
)

//spellchecker:words nolint wrapcheck

func NewVersionCommand() *cobra.Command {
	impl := new(version)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "print version information and exit",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type version struct{}

func (version) Exec(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s, built %s, using %s\n", ggman.BuildVersion, ggman.BuildTime, runtime.Version())
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}
