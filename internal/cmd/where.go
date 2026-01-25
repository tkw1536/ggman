package cmd

//spellchecker:words github cobra ggman internal
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
)

//spellchecker:words positionals wrapcheck GGROOT GGNORM

func NewWhereCommand() *cobra.Command {
	impl := new(where)

	cmd := &cobra.Command{
		Use:   "where URL",
		Short: "Print the location where a repository would be cloned to",
		Long: `Where prints the path where a repository from URL would be cloned to.

Each path segment corresponds to a URL component.
The root defaults to '~/Projects' and can be customized via '$GGROOT'.
The 'ggman root' alias prints the root directory.

For example, 'https://github.com/hello/world.git' clones to '$GGROOT/github.com/hello/world'.
This works for any URL, not just 'github.com'.

Since ggman 1.12, path resolution considers existing directories.
Existing sub-paths differing only by case are reused.

For example, if '$GGROOT/github.com/hello' exists and 'https://github.com/HELLO/world.git' is cloned:

- Before 1.12: cloned to '$GGROOT/github.com/HELLO/world', creating duplicate directories
- After 1.12: cloned to '$GGROOT/github.com/hello/world', reusing the existing directory

The first matching directory (alphanumerically) is used.
Exact name matches are preferred over case-insensitive matches.

The '$GGNORM' environment variable controls normalization:

- 'smart' => use first matching path, prefer exact matches (default)
- 'fold' => fold paths, do not prefer exact matches
- 'none' => always use exact paths (legacy behavior)`,
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
