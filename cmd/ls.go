package cmd

//spellchecker:words github cobra ggman pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words nolint wrapcheck

func NewLsCommand() *cobra.Command {
	impl := new(ls)

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "list local paths to all locally cloned repositories",
		Long:  `When called, the ggman ls command prints a list of paths to all locally cloned repositories to standard output.`,
		Args:  cobra.NoArgs,

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.ExitCode, "exit-code", "e", false, "return exit code 1 if no repositories are found")
	flags.BoolVarP(&impl.Scores, "scores", "s", false, "show scores returned from filter along with repositories")
	flags.BoolVarP(&impl.One, "one", "o", false, "list at most one repository, for use in shell scripts")
	flags.IntVarP(&impl.Limit, "count", "n", 0, "list at most this many repositories. May not be combined with one")

	return cmd
}

type ls struct {
	ExitCode bool
	Scores   bool
	One      bool
	Limit    int
}

var (
	errLSExitFlag             = exit.NewErrorWithCode("", exit.ExitGeneric)
	errLsOnlyOneOfOneAndLimit = exit.NewErrorWithCode("only one of `--one` and `--count` may be provided", exit.ExitCommandArguments)
	errLsLimitNegative        = exit.NewErrorWithCode("`--count` may not be negative", exit.ExitCommandArguments)
)

func (l *ls) ParseArgs(cmd *cobra.Command, args []string) error {
	if l.Limit < 0 {
		return errLsLimitNegative
	}
	if l.Limit != 0 && l.One {
		return errLsOnlyOneOfOneAndLimit
	}
	return nil
}

func (l *ls) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		AllowsFilter: true,
		NeedsRoot:    true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	if l.One {
		l.Limit = 1
	}
	repos, scores := environment.RepoScores(true)
	if l.Limit > 0 && len(repos) > l.Limit {
		repos = repos[:l.Limit]
	}
	for i, repo := range repos {
		if l.Scores {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%f %s\n", scores[i], repo); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
			continue
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), repo); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if l.ExitCode && len(repos) == 0 {
		return errLSExitFlag
	}

	return nil
}
