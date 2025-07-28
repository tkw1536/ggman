package cmd

//spellchecker:words github cobra ggman internal
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
)

//spellchecker:words positionals nolint wrapcheck

func NewWhereCommand() *cobra.Command {
	impl := new(where)

	cmd := &cobra.Command{
		Use:   "where URL",
		Short: "print the location where a repository would be cloned to",
		Long: `When invoked, the ggman where command prints to standard output the location where the remote repository described by the first argument would be cloned to. 
This location is a subfolder of the directory outputted by 'ggman root'. 
Each segment of the path corresponding to a component of the repository url.

This command does not perform any interactions with the remote repository or the local disk, in particular it does not require access to the remote repository or require it to be installed.`,
		Args: cobra.ExactArgs(1),

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type where struct {
	Positionals struct {
		URL string
	}
}

func (w *where) ParseArgs(cmd *cobra.Command, args []string) error {
	w.Positionals.URL = args[0]
	return nil
}

func (w *where) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	localPath, err := environment.Local(env.ParseURL(w.Positionals.URL))
	if err != nil {
		return fmt.Errorf("%w: %w", env.ErrUnableLocalPath, err)
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), localPath)
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}
