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
		Long: `While creating a folder structure for cloning new repositories, ggman can run operations on any other folder structure contained within the GGROOT directory.
For this purpose the 'ggman ls' command lists all repositories that have been found in this structure to standard output.

For easier integration into scripts, 'ggman ls' supports an '--exit-code' argument.
If this is given, the command will return exit code 0 iff at least one repository is found, and exit code 1 otherwise.

Furthermore, the flag '--one' or the flags '--count' / '-n' can be given to limit the number of results.
This is useful in specific scripting circumstances.

#### Filtering repositories

When running multi-repository operations, it is possible to limit the operations to a specific subset of repositories.
This is achieved by using the 'for' keyword along with a pattern.
For example, 'ggman --for "github.com/*/example" ls' will list all repositories from 'github.com' that are named 'example'.

Examples for simple supported patterns:

- "world" matches: git@github.com:hello/world.git, https://github.com/hello/world
- "hello/*" matches: git@github.com:hello/earth.git, git@github.com:hello/mars.git
- "hello/m*" matches: git@github.com:hello/mars.git, git@github.com:hello/mercury.git
- "github.com/*/*" matches: git@github.com:hello/world.git, git@github.com:bye/world.git
- "github.com/hello" matches: git@github.com:hello/world.git, git@github.com:hello/mars.git

Patterns are generally applied against URL components (see below for details on how the splitting works).
For example, to match the pattern 'hello/*', it is first split into the patterns 'hello' and '*'.
These are then matched individually against the components of the URL.
The matched components have to be sequential, but don't have to be at either end of the URL.

Each component pattern can be one of the following:

- A case-insensitive fnmatch.3 pattern with '*', '?' and '[]' holding their usual meanings;
- A fuzzy string match, meaning the characters in the pattern have to occur in order of the characters in the string.

When no special fnmatch characters are found, the implementation assumes a fuzzy match.
Fuzzy matching can also be explicitly disabled by passing the global '--no-fuzzy-filter' argument.

A special case is when a pattern begins with '^' or ends with '$' (or both).
Then any fuzzy matching is disabled, and any matches must start at the beginning (in the case '^') or end at the end (in the case '$') of the URL (or both).
For example 'hello/world' matches both 'git@github.com:hello/world.git' and 'hello.com/world/example.git', but 'hello/world$' only matches the former.

Note that the '--for' argument also works for exact repository urls, e.g. 'ggman --for "https://github.com/tkw1536/ggman" ls'.
'--for' also works with absolute or relative filepaths to locally installed repositories.

In addition, the '--for' argument by default uses a fuzzy matching algorithm.
For example, the pattern 'wrld' will also match a repository called 'world'.
Fuzzy matching only works for patterns that do not contain a special glob characters ('*' and friends).
It is also possible to turn off fuzzy matching entirely by passing the '--no-fuzzy-filter' / '-n' argument.

In addition to the '--for' argument, you can also use the '--path' argument.
Instead of taking a pattern, it takes a (relative or absolute) filesystem path and matches all repositories under it.
This also works when not inside 'GGROOT'.
The '--path' argument can be provided multiple times.

The '--here' argument is an alias for '--path .', meaning it matches only the repository located in the current working directory, or repositories under it.

These flags can also be used against other multi-repository operations, such as 'ggman pull' or 'ggman fix'.`,
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
	errLsOnlyOneOfOneAndLimit = exit.NewErrorWithCode("only one of `--one` and `--count` may be provided", env.ExitCommandArguments)
	errLsLimitNegative        = exit.NewErrorWithCode("`--count` may not be negative", env.ExitCommandArguments)
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
