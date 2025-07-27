package cmd

//spellchecker:words github cobra ggman internal walker pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/ggman/internal/walker"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words GGROOT nolint wrapcheck

func NewSweepCommand() *cobra.Command {
	impl := new(sweep)

	cmd := &cobra.Command{
		Use:   "sweep",
		Short: "find empty folders in the root folder",
		Long: `The sweep command can be used to identify non-git directories within the GGROOT directory which are empty, or contain only subdirectories which are empty recursively. 
Such directories are left behind after running the 'ggman relocate' command, or after manually deleting repositories. 
The command takes no arguments, and produces them in an order such that they can be passed to 'rmdir' and be deleted.`,
		Args: cobra.NoArgs,

		RunE: impl.Exec,
	}

	return cmd
}

type sweep struct{}

var errSweepScan = exit.NewErrorWithCode("error scanning for empty directories", exit.ExitGeneric)

func (sweep) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	results, err := walker.Sweep(func(path string, root walker.FS, depth int) (stop bool) {
		return environment.Git.IsRepository(path)
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
