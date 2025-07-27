package cmd

//spellchecker:words runtime github cobra ggman constants
import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/constants"
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
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s, built %s, using %s\n", constants.BuildVersion, constants.BuildTime, runtime.Version())
	if err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	return nil
}
