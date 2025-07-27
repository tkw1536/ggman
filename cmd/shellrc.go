package cmd

//spellchecker:words embed github cobra ggman
import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//spellchecker:words shellrc

func NewShellrcCommand() *cobra.Command {
	impl := new(shellrc)

	cmd := &cobra.Command{
		Use:   "shellrc",
		Short: "print additional aliases to be used in shell profiles in conjunction with ggman",
		Long:  "The 'ggman shellrc' command prints aliases to be used for shell profiles in conjunction with ggman.",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type shellrc struct{}

//go:embed shellrc.sh
var shellrcSh string

func (shellrc) Exec(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprint(cmd.OutOrStdout(), shellrcSh)
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}
