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
		Long:  "Version prints version information about this program.",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type version struct{}

func (version) Exec(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s, built using %s\n", ggman.BuildVersion, runtime.Version())
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}
