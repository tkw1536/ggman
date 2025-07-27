package cmd

//spellchecker:words errors path filepath essio shellescape github cobra ggman internal dirs goprogram exit pkglib
import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/ggman/internal/dirs"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/fsx"
)

//spellchecker:words nolint wrapcheck

func NewRelocateCommand() *cobra.Command {
	impl := new(relocate)

	cmd := &cobra.Command{
		Use:   "relocate",
		Short: "move locally cloned repositories into locations as per \"ggman where\"",
		Long:  "Relocate moves all repositories to the location where they should be moved to if they had been cloned with 'ggman clone'.",
		Args:  cobra.NoArgs,

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.Simulate, "simulate", "s", false, "only print unix-like commands to move repositories around")

	return cmd
}

type relocate struct {
	Simulate bool
}

func (relocate) Description() ggman.Description {
	return ggman.Description{
		Command:     "relocate",
		Description: "move locally cloned repositories into locations as per \"ggman where\"",

		Requirements: env.Requirement{
			NeedsRoot:    true,
			NeedsCanFile: true,
			AllowsFilter: true,
		},
	}
}

var (
	errRelocateCreateParent = exit.NewErrorWithCode("unable to create parent directory for destination", exit.ExitGeneric)
	errRelocateMove         = exit.NewErrorWithCode("unable to move repository", exit.ExitGeneric)

	errRelocateRepoExists = exit.NewErrorWithCode("repository already exists", exit.ExitGeneric)
	errRelocatePathExists = exit.NewErrorWithCode("path already exists", exit.ExitGeneric)
)

func (r *relocate) AfterParse(cmd *cobra.Command, args []string) error {
	return nil
}

func (r *relocate) Exec(cmd *cobra.Command, args []string) error {
	environment, err := ggman.GetEnv(cmd)
	if err != nil {
		return err
	}

	for _, gotPath := range environment.Repos(false) {
		// determine the remote path and where it should go
		remote, err := environment.Git.GetRemote(gotPath, "")
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
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "mv %s %s\n", shellescape.Quote(gotPath), shellescape.Quote(shouldPath)); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
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
			got, err := environment.AtRoot(shouldPath)
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
