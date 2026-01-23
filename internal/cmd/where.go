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
		Long: `
ggman manages all git repositories inside a given root directory, and automatically sets up new repositories relative to the URLs they are cloned from. 
Where prints to standard output the location where the remote repository described by the first argument would be cloned to. 

Each segment of the path corresponds to a component of the repository url.
The root folder defaults to '~/Projects' but can be customized using the 'GGROOT' environment variable. 
The root directory can be echoed using the command alias 'ggman root'. 

For example, when ggman clones a repository 'https://github.com/hello/world.git', this would automatically end up in '$GGROOT/github.com/hello/world'. 
This works not only for 'github.com' urls, but for any kind of url. 
To see where a repository would be cloned to (but not actually cloning it), use 'ggman where <REPO>'. 

As of ggman 1.12, this translation of URLs into paths takes existing paths into account.
In particular, it re-uses existing sub-paths if they differ from the requested path only by casing.

For example, say the directory '$GGROOT/github.com/hello' exists and the user requests to clone 'https://github.com/HELLO/world.git'.
Before 1.12, this clone would end up in '$GGROOT/github.com/HELLO/world', resulting in two directories '$GGROOT/github.com/HELLO' and '$GGROOT/github.com/hello'. 
After 1.12, this clone will end up in '$GGROOT/github.com/hello/world'.
While this means placing of repositories needs to touch the disk (and check for existing directories), it results in less directory clutter.

By default, the first matching directory (in alphanumerical order) is used as opposed to creating a new one.
If a directory with the exact name exists, this is preferred over a case-insensitive match.

This normalization behavior can be controlled using the 'GGNORM' environment variable.
It has three values:
  - 'smart' (use first matching path, prefer exact matches, default behavior);
  - 'fold' (fold paths, but do not prefer exact matches); and
  - 'none' (always use exact paths, legacy behavior)
`,
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
