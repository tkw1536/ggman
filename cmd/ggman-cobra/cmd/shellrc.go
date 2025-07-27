package cmd

//spellchecker:words embed github cobra ggman
import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
)

//spellchecker:words shellrc

func NewShellrcCommand() *cobra.Command {
	impl := new(shellrc)

	cmd := &cobra.Command{
		Use:   "shellrc",
		Short: "print additional aliases to be used in shell profiles in conjunction with ggman",
		Long:  "The 'ggman shellrc' command prints aliases to be used for shell profiles in conjunction with ggman.",
		Args:  cobra.NoArgs,

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	return cmd
}

type shellrc struct{}

func (shellrc) Description() ggman.Description {
	return ggman.Description{
		Command:     "shellrc",
		Description: "print additional aliases to be used in shell profiles in conjunction with ggman",
	}
}

// TODO: Make this private again
//
//go:embed shellrc.sh
var ShellrcSh string

func (*shellrc) AfterParse(cmd *cobra.Command, args []string) error {
	return nil
}

func (shellrc) Exec(cmd *cobra.Command, args []string) error {
	_, err := fmt.Fprint(cmd.OutOrStdout(), ShellrcSh)
	if err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	return nil
}
