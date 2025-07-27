package cmd

//spellchecker:words github cobra ggman constants legal
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

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	return cmd
}

type version struct{}

func (version) Description() ggman.Description {
	return ggman.Description{
		Command:     "version",
		Description: "print a version message and exit",
	}
}

func (*version) AfterParse(cmd *cobra.Command, args []string) error {
	return nil
}

func (version) Exec(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprintf(cmd.OutOrStdout(), "%s, built %s, using %s\n", constants.BuildVersion, constants.BuildTime, runtime.Version())
	if err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	return nil
}
