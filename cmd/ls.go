package cmd

//spellchecker:words ggman goprogram exit
import (
	"fmt"

	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

//spellchecker:words nolint wrapcheck

// Ls is the 'ggman ls' command.
//
// When called, the ggman ls command prints a list of paths to all locally cloned repositories to standard output.
//
//	--exit-code
//
// When provided, exit with code 1 if no repositories are found.
//
//	--one
//
// # List at most one repository
//
// --count NUMBER
//
// List as most COUNT repositories. May not be combined with the one flag.
var Ls ggman.Command = ls{}

type ls struct {
	ExitCode bool `description:"return exit code 1 if no repositories are found"                   long:"exit-code" short:"e"`
	Scores   bool `description:"show scores returned from filter along with repositories"          long:"scores"    short:"s"`
	One      bool `description:"list at most one repository, for use in shell scripts"             long:"one"       short:"o"`
	Limit    int  `description:"list at most this many repositories. May not be combined with one" long:"count"     short:"n"`
}

func (ls) Description() ggman.Description {
	return ggman.Description{
		Command:     "ls",
		Description: "list local paths to all locally cloned repositories",

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

var (
	errLSExitFlag             = exit.NewErrorWithCode("", exit.ExitGeneric)
	errLsOnlyOneOfOneAndLimit = exit.NewErrorWithCode("only one of `--one` and `--count` may be provided", exit.ExitCommandArguments)
	errLsLimitNegative        = exit.NewErrorWithCode("`--count` may not be negative", exit.ExitCommandArguments)
)

func (ls ls) AfterParse() error {
	if ls.Limit < 0 {
		return errLsLimitNegative
	}
	if ls.Limit != 0 && ls.One {
		return errLsOnlyOneOfOneAndLimit
	}
	return nil
}

func (l ls) Run(context ggman.Context) error {
	if l.One {
		l.Limit = 1
	}
	repos, scores := context.Environment.RepoScores(true)
	if l.Limit > 0 && len(repos) > l.Limit {
		repos = repos[:l.Limit]
	}
	for i, repo := range repos {
		if l.Scores {
			if _, err := context.Printf("%f %s\n", scores[i], repo); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
			continue
		}
		if _, err := context.Println(repo); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if l.ExitCode && len(repos) == 0 {
		return errLSExitFlag
	}

	return nil
}
