package cmd

//spellchecker:words github cobra ggman internal pkglib exit
import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words compinit

func NewCompletionCmd() *cobra.Command {
	impl := new(completion)

	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate a ggman tab completion script for a shell",
		Long: fmt.Sprintf(`Completion generates a tab completion script.
The shell name must be passed as the first argument.

Supported shells:

- bash
- zsh
- fish
- powershell

Installation is shell-dependent.

Bash:

    $ source <(%[1]s completion bash)

    # To load completions for each session, execute once:
    # Linux:
    $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
    # macOS:
    $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

Zsh:

    # Enable completion if not already enabled:
    $ echo "autoload -U compinit; compinit" >> ~/.zshrc

    # Install:
    $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

    # Restart shell for effect.

fish:

    $ %[1]s completion fish | source

    # Persistent:
    $ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

PowerShell:

    PS> %[1]s completion powershell | Out-String | Invoke-Expression

    # Persistent:
    PS> %[1]s completion powershell > %[1]s.ps1
    # Source this file from the PowerShell profile.
`, "ggman"),
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),

		PreRunE: func(cmd *cobra.Command, args []string) error {
			return impl.ParseArgs(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return impl.Exec(cmd, args)
		},
	}
}

type completion struct {
	Positional struct {
		Shell string
	}
}

var errNoShell = exit.NewErrorWithCode("no shell name provided", env.ExitCommandArguments)

func (c *completion) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errNoShell
	}

	c.Positional.Shell = args[0]
	return nil
}

func (c *completion) Exec(cmd *cobra.Command, args []string) error {
	var err error
	switch args[0] {
	case "bash":
		err = cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		err = cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		err = cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		panic("never reached")
	}
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}
