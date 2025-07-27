package cmd

//spellchecker:words ggman goprogram exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

func NewFindBranchCommand() *cobra.Command {
	impl := new(findBranch)

	cmd := &cobra.Command{
		Use:   "find-branch BRANCH",
		Short: "list repositories containing a specific branch",
		Long: `The 'find-branch' command lists all repositories that contain a branch with the provided name.
The remotes will be listed in dictionary order of their local installation paths.`,
		Args: cobra.ExactArgs(1),

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.ExitCode, "exit-code", "e", false, "exit with status code 1 when no repositories with provided branch exist")

	return cmd
}

//spellchecker:words positionals nolint wrapcheck

type findBranch struct {
	Positionals struct {
		Branch string
	}
	ExitCode bool
}

func (*findBranch) Description() ggman.Description {
	return ggman.Description{
		Command:     "find-branch",
		Description: "list repositories containing a specific branch",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (f *findBranch) AfterParse(cmd *cobra.Command, args []string) error {
	f.Positionals.Branch = args[0]
	return nil
}

var errFindBranchCustom = exit.NewErrorWithCode("", exit.ExitGeneric)

func (f *findBranch) Exec(cmd *cobra.Command, args []string) error {
	environment, err := ggman.GetEnv(cmd)
	if err != nil {
		return nil
	}

	foundRepo := false
	for _, repo := range environment.Repos(true) {
		// check if the repository has the branch!
		hasBranch, err := environment.Git.ContainsBranch(repo, f.Positionals.Branch)
		if err != nil {
			panic(err)
		}
		if !hasBranch {
			continue
		}

		foundRepo = true
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), repo); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if f.ExitCode && !foundRepo {
		return errFindBranchCustom
	}

	return nil
}
