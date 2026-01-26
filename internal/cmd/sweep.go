package cmd

//spellchecker:words github cobra ggman internal walker pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/ggman/internal/walker"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words GGROOT wrapcheck

func NewSweepCommand() *cobra.Command {
	impl := new(sweep)

	cmd := &cobra.Command{
		Use:   "sweep",
		Short: "find empty folders in the root folder",
		Long: `Sweep identifies empty non-git directories within '$GGROOT'.
A directory is empty if it contains only recursively empty subdirectories.
These directories remain after 'ggman relocate' or manual repository deletion.

Output is ordered such that

    ggman sweep | xargs rmdir

can remove all directories.
`,
		Args: cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type sweep struct{}

var errSweepScan = exit.NewErrorWithCode("failed to scan for empty directories", env.ExitGeneric)

func (sweep) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	results, err := walker.Sweep(func(path string, root walker.FS, depth int) (stop bool) {
		return environment.Git.IsRepository(cmd.Context(), path)
	}, walker.Params{
		Root: walker.NewRealFS(environment.Root, false),
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errSweepScan, err)
	}

	for _, r := range results {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), r); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}
	return nil
}
