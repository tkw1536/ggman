package cmd

//spellchecker:words encoding json path filepath sync essio shellescape github cobra ggman internal pkglib collection exit
import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/collection"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words wrapcheck wrld fnmatch GGROOT canonicalized

func NewLsCommand() *cobra.Command {
	impl := new(ls)

	cmd := &cobra.Command{
		Use:   "ls",
		Short: "List local paths of cloned repositories",
		Long: `Ls lists all repositories found within the '$GGROOT' directory to standard output.

The '--exit-code' flag causes exit code 0 when at least one repository is found, and exit code 1 otherwise.
The '--count' flag limits output to at most the specified number of repositories.
The '--one' flag is equivalent to '--count 1' and limits output to at most one repository.

The '--relative' flag prints paths relative to '$GGROOT' instead of absolute paths.
The '--remote' flag prints remote URLs instead of local paths.
The '--canonical' flag prints canonicalized remote URLs instead of the original ones.

The '--scores' flag shows filtering scores in addition to any paths in the output.

By default, output consists of one repository (and possibly score) per line.
The '--json' flag outputs JSON instead of plain text.
The '--export' flag generates a bash script to re-clone all repositories.

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

The filter flags work with other multi-repository commands such as 'ggman pull' or 'ggman fix'.`,
		Args: cobra.NoArgs,

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&impl.ExitCode, "exit-code", "e", false, "return exit code 1 if no repositories are found")
	flags.BoolVarP(&impl.Scores, "scores", "s", false, "show scores returned from filter along with repositories")
	flags.BoolVarP(&impl.One, "one", "o", false, "list at most one repository, for use in shell scripts")
	flags.IntVarP(&impl.Limit, "count", "n", 0, "list at most this many repositories. May not be combined with one")
	flags.BoolVarP(&impl.Relative, "relative", "l", false, "compute relative paths instead of absolute ones")
	flags.BoolVarP(&impl.Remote, "remote", "r", false, "gather remote URLs instead of local ones")
	flags.BoolVarP(&impl.Canonical, "canonical", "c", false, "gather canonicalized remote URLs")
	flags.BoolVarP(&impl.JSON, "json", "j", false, "output JSON")
	flags.BoolVarP(&impl.Export, "export", "x", false, `generate a bash script to re-clone repositories. Implies "--remote" and "--relative"`)

	return cmd
}

type ls struct {
	ExitCode bool

	Scores bool
	One    bool
	Limit  int

	Relative bool

	Remote    bool
	Canonical bool

	JSON   bool
	Export bool
}

var (
	errLSExitFlag                = exit.NewErrorWithCode("", env.ExitGeneric)
	errLsInvalidCanfile          = exit.NewErrorWithCode("failed to parse CANFILE", env.ExitInvalidEnvironment)
	errLsOnlyOneOfOneAndLimit    = exit.NewErrorWithCode(`only one of "--one" and "--count" may be provided`, env.ExitCommandArguments)
	errLsLimitNegative           = exit.NewErrorWithCode(`"--count" may not be negative`, env.ExitCommandArguments)
	errLsCanonicalOnlyWithRemote = exit.NewErrorWithCode(`"--canonical" may only be used with "--remote"`, env.ExitCommandArguments)
	errLsRelativeNoRemote        = exit.NewErrorWithCode(`"--relative" may not be used with "--remote" unless "--json" or "--export" are set`, env.ExitCommandArguments)
	errLsJSONAndExport           = exit.NewErrorWithCode(`"--json" and "--export" cannot be used together`, env.ExitCommandArguments)
)

func (l *ls) ParseArgs(cmd *cobra.Command, args []string) error {
	if l.Limit < 0 {
		return errLsLimitNegative
	}
	if l.Limit != 0 && l.One {
		return errLsOnlyOneOfOneAndLimit
	}

	if l.One {
		l.Limit = 1
	}

	if l.Export {
		l.Remote = true
		l.Relative = true
	}

	if l.JSON && l.Export {
		return errLsJSONAndExport
	}

	if l.Canonical && !l.Remote {
		return errLsCanonicalOnlyWithRemote
	}
	if (!l.JSON && !l.Export) && l.Relative && l.Remote {
		return errLsRelativeNoRemote
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

	repos, err := l.getRepositories(cmd, environment)
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	switch {
	case l.JSON:
		if err := l.outputJSON(cmd, repos); err != nil {
			return err
		}
	case l.Export:
		if err := l.outputExport(cmd, repos); err != nil {
			return err
		}
	default:
		if err := l.outputPlain(cmd, repos); err != nil {
			return err
		}
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if l.ExitCode && len(repos) == 0 {
		return errLSExitFlag
	}

	return nil
}

func (l *ls) outputExport(cmd *cobra.Command, repos []Repo) error {
	w := cmd.OutOrStdout()
	if _, err := fmt.Fprintln(w, "#!/bin/bash"); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	if _, err := fmt.Fprintln(w, "set -e"); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	if _, err := fmt.Fprintln(w, ""); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	if _, err := fmt.Fprintln(w, "# Generated by ggman export"); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}

	for _, repo := range repos {
		if _, err := fmt.Fprintf(w, "mkdir -p %s\n", shellescape.Quote(repo.Relative)); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
		if _, err := fmt.Fprintf(w, "git clone %s %s\n", shellescape.Quote(repo.Remote), shellescape.Quote(repo.Relative)); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}

	return nil
}

func (l *ls) outputJSON(cmd *cobra.Command, repos []Repo) error {
	encoder := json.NewEncoder(cmd.OutOrStdout())
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(repos); err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}

func (l *ls) outputPlain(cmd *cobra.Command, repos []Repo) error {
	for _, repo := range repos {
		value := repo.Path
		switch {
		case l.Canonical:
			value = repo.Canonical
		case l.Remote:
			value = repo.Remote
		case l.Relative:
			value = repo.Relative
		default:
		}

		if l.Scores {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%f %s\n", repo.Score, value); err != nil {
				return fmt.Errorf("%w: %w", errGenericOutput, err)
			}
			continue
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), value); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}
	return nil
}

// Repo represents a repository along with additional information.
type Repo struct {
	valid bool

	Path     string `json:",omitempty"`
	Relative string `json:",omitempty"`

	Score float64 `json:",omitempty"`

	Remote    string `json:",omitempty"`
	Canonical string `json:",omitempty"`
}

// getRepositories returns a list of repositories.
//
// the returned struct may only be partially populated, according to arguments.
// the struct may only return a limited number of repositories, according to arguments.
func (l *ls) getRepositories(cmd *cobra.Command, environment *env.Env) ([]Repo, error) {
	var canFile env.CanFile
	if l.Canonical {
		var err error
		if canFile, err = environment.LoadDefaultCANFILE(); err != nil {
			return nil, fmt.Errorf("%w: %w", errLsInvalidCanfile, err)
		}
	}

	// list all the repositories.
	repos, scores := environment.RepoScores(cmd.Context(), true)
	if l.Limit > 0 && len(repos) > l.Limit {
		repos = repos[:l.Limit]
		scores = scores[:l.Limit]
	}

	infos := make([]Repo, len(repos))

	// get information about each repository concurrently.
	var wg sync.WaitGroup
	for i, path := range repos {
		wg.Go(func() {
			infos[i] = l.getRepository(cmd, environment, path, scores[i], canFile)
		})
	}
	wg.Wait()

	return collection.KeepFunc(infos, func(repo Repo) bool {
		return repo.valid
	}), nil
}

// getRepository returns information about a single repository in accordance with flags.
func (ls *ls) getRepository(cmd *cobra.Command, environment *env.Env, path string, score float64, canFile env.CanFile) (r Repo) {
	r.Path = path
	r.Score = score

	if ls.Remote {
		var err error
		r.Remote, err = environment.Git.GetRemote(cmd.Context(), path, "")
		if err != nil {
			return Repo{valid: false}
		}
		if ls.Canonical {
			r.Canonical = env.ParseURL(r.Remote).CanonicalWith(canFile)
		}
	}

	if ls.Relative {
		var err error

		r.Relative, err = filepath.Rel(environment.Root, path)
		if err != nil {
			return Repo{valid: false}
		}
	}

	r.valid = true
	return r
}
