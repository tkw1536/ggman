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
		Long: `Relocate moves all repositories to the location where they should be moved to if they had been cloned with 'ggman clone'.

This can be useful to move one or multiple repositories into ggman's default structure.
It can also be useful when a repository changes its' remote url and should be moved to a new location.`,
		Args: cobra.NoArgs,

		RunE: impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.OnlyCurrentRemote, "only-current-remote", "o", false, "consider only the current remote (as opposed to all remotes) when checking if a repository is in the correct location")
	flags.BoolVarP(&impl.Simulate, "simulate", "s", false, "only print unix-like commands to move repositories around")

	return cmd
}

type relocate struct {
	Simulate          bool
	OnlyCurrentRemote bool
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
		// check if we are in a valid location
		valid, err := isValidLocation(gotPath, r.OnlyCurrentRemote, cmd, environment)
		if err != nil || valid {
			continue
		}

		// find the path it should go to!
		shouldPath, err := getCanonicalLocation(gotPath, cmd, environment)
		if err != nil {
			return fmt.Errorf("failed to determine canonical location: %w", err)
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

// isValidLocation checks if the repository at path is in the correct location.
func isValidLocation(path string, onlyCurrentRemote bool, cmd *cobra.Command, environment *env.Env) (bool, error) {
	if onlyCurrentRemote {
		remote, err := getCanonicalLocation(path, cmd, environment)
		if err != nil {
			return false, fmt.Errorf("failed to get canonical location: %w", err)
		}
		return fsx.Same(path, remote), nil
	}
	remotes, err := environment.Git.GetAllRemotes(cmd.Context(), path)
	if err != nil {
		return false, fmt.Errorf("failed to get remotes: %w", err)
	}

	// if we don't have any remotes any path is fine.
	if len(remotes) == 0 {
		return true, nil
	}

	// check for each remote if it is fine.
	for _, remote := range remotes {
		shouldPath, err := environment.Local(env.ParseURL(remote))
		if err != nil {
			return false, fmt.Errorf("%w: %w", env.ErrUnableLocalPath, err)
		}

		// if it is the same, don't move it
		if fsx.Same(path, shouldPath) {
			return true, nil
		}
	}

	// none of the remotes were fine!
	return false, nil
}

// getCanonicalLocation gets the canonical location for a given repository.
func getCanonicalLocation(path string, cmd *cobra.Command, environment *env.Env) (string, error) {
	remote, err := environment.Git.GetRemote(cmd.Context(), path, "")
	if err != nil || remote == "" { // ignore remotes that don't exist
		return "", fmt.Errorf("failed to get remotes: %w", err)
	}
	shouldPath, err := environment.Local(env.ParseURL(remote))
	if err != nil {
		return "", fmt.Errorf("%w: %w", env.ErrUnableLocalPath, err)
	}
	return shouldPath, nil
}
