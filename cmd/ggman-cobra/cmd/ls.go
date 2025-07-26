package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

var (
	errLSExitFlag             = exit.NewErrorWithCode("", exit.ExitGeneric)
	errLsOnlyOneOfOneAndLimit = exit.NewErrorWithCode("only one of `--one` and `--count` may be provided", exit.ExitCommandArguments)
	errLsLimitNegative        = exit.NewErrorWithCode("`--count` may not be negative", exit.ExitCommandArguments)
)

func NewLSCommand() *cobra.Command {
	var (
		ExitCode bool
		Scores   bool
		One      bool
		Limit    int
	)

	ls := &cobra.Command{
		Use:   "ls",
		Short: "list local paths to all locally cloned repositories",
		Args:  cobra.NoArgs,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			ggman.SetRequirements(cmd, &env.Requirement{
				AllowsFilter: true,
				NeedsRoot:    true,
			})

			if Limit < 0 {
				return errLsLimitNegative
			}
			if Limit != 0 && One {
				return errLsOnlyOneOfOneAndLimit
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			env, err := ggman.GetEnv(cmd)
			if err != nil {
				return fmt.Errorf("failed to get environment: %w", err)
			}

			if One {
				Limit = 1
			}
			repos, scores := env.RepoScores(true)
			if Limit > 0 && len(repos) > Limit {
				repos = repos[:Limit]
			}

			stdout := cmd.OutOrStdout()
			for i, repo := range repos {
				if Scores {
					if _, err := fmt.Fprintf(stdout, "%f %s\n", scores[i], repo); err != nil {
						return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
					}
					continue
				}
				if _, err := fmt.Fprintln(stdout, repo); err != nil {
					return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
				}
			}

			// if we have --exit-code set and no results
			// we need to exit with an error code
			if ExitCode && len(repos) == 0 {
				return errLSExitFlag
			}

			return nil
		},
	}

	flags := ls.Flags()
	flags.BoolVarP(&ExitCode, "exit-code", "e", ExitCode, "return exit code 1 if no repositories are found")
	flags.BoolVarP(&Scores, "scores", "s", Scores, "show scores returned from filter along with repositories")
	flags.BoolVarP(&One, "one", "o", One, "list at most one repository, for use in shell scripts")
	flags.IntVarP(&Limit, "count", "n", Limit, "list at most this many repositories. May not be combined with one")

	return ls
}
