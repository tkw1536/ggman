package cmd

//spellchecker:words github cobra ggman internal pkglib exit
import (
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words wrapcheck wrld fnmatch GGROOT

func NewLsCommand() *cobra.Command {
	impl := new(ls)

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List local paths of cloned repositories",
		Long: `Ls lists all repositories found within the '$GGROOT' directory to standard output.

The '--exit-code' flag causes exit code 0 when at least one repository is found, and exit code 1 otherwise.
The '--count' flag limits output to at most the specified number of repositories.
The '--one' flag is equivalent to '--count 1' and limits output to at most one repository.


#### Filtering repositories

The '--for' flag limits operations to repositories matching a pattern.
For example, 'ggman --for "github.com/*/example" ls' lists all repositories from 'github.com' named 'example'.

Pattern examples:

- "world" => git@github.com:hello/world.git, https://github.com/hello/world
- "hello/*" => git@github.com:hello/earth.git, git@github.com:hello/mars.git
- "hello/m*" => git@github.com:hello/mars.git, git@github.com:hello/mercury.git
- "github.com/*/*" => git@github.com:hello/world.git, git@github.com:bye/world.git
- "github.com/hello" => git@github.com:hello/world.git, git@github.com:hello/mars.git

Patterns are applied against URL components.
The pattern 'hello/*' is split into 'hello' and '*', then matched sequentially against URL components.
Matched components must be sequential but need not be at the URL boundaries.

Each component pattern is one of:

- A case-insensitive fnmatch.3 pattern with '*', '?' and '[]' holding their usual meanings
- A fuzzy string match where pattern characters must occur in order within the string

Without special fnmatch characters, fuzzy matching is assumed.
The '--no-fuzzy-filter' flag disables fuzzy matching.

Patterns beginning with '^' or ending with '$' disable fuzzy matching and require matches at URL boundaries.
For example, 'hello/world' matches both 'git@github.com:hello/world.git' and 'hello.com/world/example.git', but 'hello/world$' matches only the former.

The '--for' flag also accepts exact repository URLs or filesystem paths.
Fuzzy matching applies by default: 'wrld' matches 'world'.
Fuzzy matching only applies to patterns without glob characters.

The '--path' flag matches all repositories under a filesystem path.
It works outside '$GGROOT' and can be specified multiple times.
The '--here' flag is equivalent to '--path .'.

These flags work with other multi-repository commands such as 'ggman pull' or 'ggman fix'.`,
		Args: cobra.NoArgs,

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
	errLSExitFlag             = exit.NewErrorWithCode("", env.ExitGeneric)
	errLsOnlyOneOfOneAndLimit = exit.NewErrorWithCode(`only one of "--one" and "--count" may be provided`, env.ExitCommandArguments)
	errLsLimitNegative        = exit.NewErrorWithCode(`"--count" may not be negative`, env.ExitCommandArguments)
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
	repos, scores := environment.RepoScores(cmd.Context(), true)
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
