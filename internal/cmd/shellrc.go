package cmd

//spellchecker:words embed github cobra
import (
	_ "embed"
	"fmt"

	"github.com/spf13/cobra"
)

//spellchecker:words shellrc ggcd ggdo ggclone ggshow ggcd ggcode ggcursor

func NewShellrcCommand() *cobra.Command {
	impl := new(shellrc)

	cmd := &cobra.Command{
		Use:   "shellrc",
		Short: "Print additional aliases to be used in shell profiles in conjunction with ggman",
		Long: `Shellrc prints shell aliases for use with ggman.

To install, add the following to '.zshrc' or '.bashrc':

    eval "$(ggman shellrc)"

Provided aliases:

- 'ggcd PATTERN' => cd into the first repository matching PATTERN
- 'gg PATTERN CMD' / 'ggdo PATTERN CMD' => run CMD in the first repository matching PATTERN
- 'ggclone URL' => clone repository and cd into it; skips clone if already exists
- 'ggshow PATTERN' => like ggcd, but also prints the HEAD commit
- 'ggcode PATTERN' => open repository in vscode
- 'ggcursor PATTERN' => open repository in cursor`,
		Args: cobra.NoArgs,

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
