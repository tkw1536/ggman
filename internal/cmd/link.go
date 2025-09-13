package cmd

//spellchecker:words path filepath github cobra ggman internal dirs pkglib exit
import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/dirs"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
)

//spellchecker:words positionals wrapcheck

func NewLinkCommand() *cobra.Command {
	impl := new(link)

	cmd := &cobra.Command{
		Use:   "link PATH",
		Short: "Symlink a repository into the local repository structure",
		Long:  "Link symlinks the repository in the path passed as the first argument where it would have been cloned to inside 'ggman root'.",
		Args:  cobra.ExactArgs(1),

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	return cmd
}

type link struct {
	Positionals struct {
		Path string
	}
}

var (
	errLinkDoesNotExist  = exit.NewErrorWithCode("can not open source repository", env.ExitGeneric)
	errLinkSamePath      = exit.NewErrorWithCode("link source and target are identical", env.ExitGeneric)
	errLinkAlreadyExists = exit.NewErrorWithCode("another directory already exists in target location", env.ExitGeneric)
	errLinkCheck         = exit.NewErrorWithCode("unable to check directory", env.ExitGeneric)
	errLinkUnknown       = exit.NewErrorWithCode("unknown linking error", env.ExitGeneric)
)

func (l *link) ParseArgs(cmd *cobra.Command, args []string) error {
	l.Positionals.Path = args[0]
	return nil
}

func (l *link) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	// make sure that the path is absolute
	// to avoid relative symlinks
	from, e := environment.Abs(l.Positionals.Path)
	if e != nil {
		return errLinkDoesNotExist
	}

	// open the source repository and get the remote
	r, e := environment.Git.GetRemote(cmd.Context(), from, "")
	if e != nil {
		return errLinkDoesNotExist
	}

	// find the target path
	to, err := environment.Local(env.ParseURL(r))
	if err != nil {
		return fmt.Errorf("%w: %w", env.ErrUnableLocalPath, err)
	}
	parentTo := filepath.Dir(to)

	// if it's the same path, we throw an error
	if from == to {
		return errLinkSamePath
	}

	// make sure it doesn't exist
	{
		exists, err := fsx.Exists(to)
		if err != nil {
			return fmt.Errorf("%w: %q: %w", errLinkCheck, to, err)
		}
		if exists {
			return errLinkAlreadyExists
		}
	}

	if _, err := fmt.Fprintf(cmd.OutOrStdout(), "Linking %q -> %q\n", to, from); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}

	// make the parent folder
	if e := os.MkdirAll(parentTo, dirs.NewModBits); e != nil {
		return fmt.Errorf("%w: %w", errLinkUnknown, e)
	}

	// and make the symlink
	if e := os.Symlink(from, to); e != nil {
		return fmt.Errorf("%w: %w", errLinkUnknown, e)
	}

	return nil
}
