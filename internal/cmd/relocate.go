package cmd

//spellchecker:words errors path filepath essio shellescape github cobra ggman internal dirs pkglib exit
import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/dirs"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
)

//spellchecker:words wrapcheck

func NewRelocateCommand() *cobra.Command {
	impl := new(relocate)

	cmd := &cobra.Command{
		Use:   "relocate",
		Short: "Move locally cloned repositories into locations as per \"ggman where\"",
		Long:  "Relocate moves all repositories to the location where they should be moved to if they had been cloned with 'ggman clone'.",
		Args:  cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Simulate, "simulate", "s", false, "only print unix-like commands to move repositories around")

	return cmd
}

type relocate struct {
	Simulate bool
}

var (
	errRelocateCreateParent = exit.NewErrorWithCode("unable to create parent directory for destination", env.ExitGeneric)
	errRelocateMove         = exit.NewErrorWithCode("unable to move repository", env.ExitGeneric)

	errRelocateRepoExists = exit.NewErrorWithCode("repository already exists", env.ExitGeneric)
	errRelocatePathExists = exit.NewErrorWithCode("path already exists", env.ExitGeneric)
)

func (r *relocate) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot:    true,
		NeedsCanFile: true,
		AllowsFilter: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	for _, gotPath := range environment.Repos(cmd.Context(), false) {
		// determine the remote path and where it should go
		remote, err := environment.Git.GetRemote(cmd.Context(), gotPath, "")
		if err != nil || remote == "" { // ignore remotes that don't exist
			continue
		}
		shouldPath, err := environment.Local(env.ParseURL(remote))
		if err != nil {
			return fmt.Errorf("%w: %w", env.ErrUnableLocalPath, err)
		}

		// if it is the same, don't move it
		if fsx.Same(gotPath, shouldPath) {
			continue
		}

		parentPath := filepath.Dir(shouldPath)

		// print what is being done
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "mkdir -p %s\n", shellescape.Quote(parentPath)); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "mv %s %s\n", shellescape.Quote(gotPath), shellescape.Quote(shouldPath)); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
		if r.Simulate {
			continue
		}

		// do it!
		if err := os.MkdirAll(parentPath, dirs.NewModBits); err != nil {
			return fmt.Errorf("%q: %w: %w", parentPath, errRelocateCreateParent, err)
		}

		// if there already is a target repository at the path
		{
			got, err := environment.AtRoot(cmd.Context(), shouldPath)
			if err != nil {
				return fmt.Errorf("%w: %w", errRelocateMove, err)
			}
			if got != "" {
				return fmt.Errorf("%w at %q", errRelocateRepoExists, got)
			}
		}

		// do the rename
		{
			err := os.Rename(gotPath, shouldPath)

			// check if an error was returned because the path already existed
			// (fs.ErrPermission is returned by Windows)
			if errors.Is(err, fs.ErrExist) || errors.Is(err, fs.ErrPermission) {
				if exists, _ := fsx.Exists(shouldPath); exists {
					return fmt.Errorf("%q: %w", shouldPath, errRelocatePathExists)
				}
			}

			if err != nil {
				return fmt.Errorf("%w: %w", errRelocateMove, err)
			}
		}
	}

	return nil
}
