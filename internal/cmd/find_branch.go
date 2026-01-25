package cmd

//spellchecker:words github cobra ggman internal pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

func NewFindBranchCommand() *cobra.Command {
	impl := new(findBranch)

	cmd := &cobra.Command{
		Use:   "find-branch BRANCH",
		Short: "List repositories containing a specific branch",
		Long: `Find-branch lists repositories containing a branch with the specified name.

git 2.28 introduced 'init.defaultBranch' to set the default branch for new repositories.
Existing repositories are not affected.
Find-branch can locate repositories still using an old branch name.`,
		Args: cobra.ExactArgs(1),

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.ExitCode, "exit-code", "e", false, "exit with status code 1 when no repositories with provided branch exist")

	return cmd
}

//spellchecker:words positionals wrapcheck

type findBranch struct {
	Positionals struct {
		Branch string
	}
	ExitCode bool
}

func (f *findBranch) ParseArgs(cmd *cobra.Command, args []string) error {
	f.Positionals.Branch = args[0]
	return nil
}

var errFindBranchCustom = exit.NewErrorWithCode("", env.ExitGeneric)

func (f *findBranch) Exec(cmd *cobra.Command, args []string) error {
	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot: true,
	})
	if err != nil {
		return nil
	}

	foundRepo := false
	for _, repo := range environment.Repos(cmd.Context(), true) {
		// check if the repository has the branch!
		hasBranch, err := environment.Git.ContainsBranch(cmd.Context(), repo, f.Positionals.Branch)
		if err != nil {
			panic(err)
		}
		if !hasBranch {
			continue
		}

		foundRepo = true
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), repo); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if f.ExitCode && !foundRepo {
		return errFindBranchCustom
	}

	return nil
}
