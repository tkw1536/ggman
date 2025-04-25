package env

//spellchecker:words errors path filepath strings github ggman internal walker goprogram exit pkglib
import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/internal/path"
	"github.com/tkw1536/ggman/internal/walker"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/pkglib/fsx"
)

//spellchecker:words worktree canonicalized canonicalize CANFILE workdir GGNORM GGROOT Wrapf nolint wrapcheck

// Env represents an environment to be used by ggman.
//
// The environment defines which git repositories are managed by ggman and where these are stored.
// It furthermore determines how a URL is canonicalized using a CanFile.
//
// An environment consists of four parts, each are defined as a part of this struct.
// See NewEnv on the defaults used by ggman.
type Env struct {
	// Git is a method of interacting with on-disk git repositories.
	Git git.Git

	// Vars are the values of environment variables.
	// These are used to conditionally initialize the root and CanFile values.
	Vars Variables

	// Root is the Root folder of the environment.
	// Repositories managed by ggman should be stored in this folder.
	// See the Local() method.
	Root string

	// Workdir is the working directory of this environment.
	Workdir string

	// Filter is an optional filter for the environment.
	// Repositories not matching the filter are considered to not be a part of the environment.
	// See the Repos() method.
	Filter Filter

	// CanFile is the CanFile used to canonicalize repositories.
	// See the Canonical() method.
	CanFile CanFile
}

// Normalization returns the path Normalization used by this environment.
func (env Env) Normalization() path.Normalization {
	switch strings.ToLower(env.Vars.GGNORM) {
	case "exact":
		return path.NoNorm
	case "fold":
		return path.FoldNorm
	default:
		return path.FoldPreferExactNorm
	}
}

// Parameters represent additional parameters to create a new environment.
type Parameters struct {
	Variables
	Workdir  string
	Plumbing git.Plumbing
}

// NewEnv returns a new Env that fulfills the requirement r.
//
// See methods LoadDefaultRoot() and LoadDefaultCanFile() for a description of default values.
//
// When a CanFile is requested, it will try to load the file pointed to by the GGMAN_CANFILE environment variable.
// If the variable does not exist, it will attempt to load the file ".ggman" in the users HOME directory.
// Failing to open a CanFile, e.g. because of invalid syntax, results in an error of type Error.
//
// If r.AllowsFilter is false, NoFilter should be passed for the filter argument.
// If r.AllowsFilter is true, a filter may be passed via the filter argument.
//
// This function is untested.
func NewEnv(r Requirement, params Parameters) (Env, error) {
	env := Env{
		Git:     git.NewGitFromPlumbing(params.Plumbing, params.PATH),
		Vars:    params.Variables,
		Filter:  NoFilter,
		Workdir: params.Workdir,
	}

	if r.NeedsRoot || r.AllowsFilter { // AllowsFilter implies NeedsRoot
		if err := env.LoadDefaultRoot(); err != nil {
			return Env{}, err
		}
	}

	if r.NeedsCanFile {
		if _, err := env.LoadDefaultCANFILE(); err != nil {
			return Env{}, err
		}
	}

	return env, nil
}

var errMissingRoot = exit.Error{
	ExitCode: ExitInvalidEnvironment,
	Message:  "unable to find GGROOT directory",
}

// absRoot returns the absolute path to the root directory.
// If the root directory is not set, returns an error of type Error.
//
// This function is untested.
func (env Env) absRoot() (string, error) {
	if env.Root == "" {
		return "", errMissingRoot
	}
	root, err := filepath.Abs(env.Root)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errInvalidRoot, err)
	}
	return root, nil
}

// LoadDefaultRoot sets env.Root according to the environment variables in e.Vars.
// If e.Root is already set, does nothing and returns nil.
//
// If the GGROOT variable is set, it is used as the root directory.
// If it is not set, the subdirectory 'Projects' of the home directory is used.
//
// The root directory does not have to exist for this function to return nil.
// However if both GGROOT and Home are unset, this function returns an error of type Error.
func (env *Env) LoadDefaultRoot() error {
	if env.Root != "" {
		return nil
	}

	env.Root = env.Vars.GGROOT
	if env.Root != "" {
		return nil
	}

	if env.Vars.HOME == "" {
		return errMissingRoot
	}

	env.Root = filepath.Join(env.Vars.HOME, "Projects")
	return nil
}

// LoadDefaultCANFILE sets and returns env.CANFILE according to the environment variables in e.Vars.
// If the CANFILE is already set, immediately returns nil.
//
// If the GGMAN_CANFILE variable is set, it will use it as a filepath to read the CanFile from.
// If it is not set it will attempt to load the file '.ggman' in the home directory.
// If neither is set, this function will load an in-memory default CanFile.
//
// If loading a CanFile fails, an error of type Error is returned.
// If loading succeeds, this function returns nil.
func (env *Env) LoadDefaultCANFILE() (cf CanFile, err error) {
	if env.CanFile != nil {
		return nil, nil
	}

	files := make([]string, 0, 2)
	if env.Vars.CANFILE != "" {
		files = append(files, env.Vars.CANFILE)
	}
	if env.Vars.HOME != "" {
		files = append(files, filepath.Join(env.Vars.HOME, ".ggman"))
	}

	// In order, if a file exists read it or fail.
	// If it doesn't exist continue to the next file.
	for _, file := range files {
		f, oErr := os.Open(file) /* #nosec G304 -- different files used by design */
		switch {
		case oErr == nil: // do nothing
		case errors.Is(oErr, fs.ErrNotExist):
			continue
		default:
			return nil, fmt.Errorf("unable to open CANFILE %q: %w", file, oErr)
		}
		defer func() {
			errClose := f.Close()
			if errClose == nil {
				return
			}
			if err == nil {
				err = errClose
			}
		}()
		if _, err := cf.ReadFrom(f); err != nil {
			return nil, err
		}
		env.CanFile = cf
		return cf, nil
	}

	cf.ReadDefault()
	env.CanFile = cf
	return cf, nil
}

var errUnableToReadDirectory = errors.New("unable to read directory")

// Local returns the path that a repository named URL should be cloned to.
// Normalization of paths is controlled by the norm parameter
func (env Env) Local(url URL) (string, error) {
	root, err := env.absRoot()
	if err != nil {
		panic("Env.Local: Root not resolved")
	}

	path, err := path.JoinNormalized(env.Normalization(), root, url.Components()...)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errUnableToReadDirectory, err)
	}
	return path, nil
}

const (
	// ExitInvalidEnvironment indicates that the environment for the ggman command is setup incorrectly.
	// This typically means that the CANFILE or GGROOT is configured incorrectly, but could also indicate a different error.
	ExitInvalidEnvironment exit.ExitCode = 5

	// ExitInvalidRepo indicates that the user attempted to perform an operation on an invalid repository.
	// This typically means that the current directory is not inside GGROOT.
	ExitInvalidRepo exit.ExitCode = 6
)

var errInvalidRoot = exit.Error{
	ExitCode: ExitInvalidEnvironment,
	Message:  "unable to resolve root directory",
}

var errNotResolved = exit.Error{
	ExitCode: ExitInvalidRepo,
	Message:  "unable to resolve repository %q",
}

// Abs returns the absolute path to path, unless it is already absolute.
// path is resolved relative to the working directory of this environment.
func (env Env) Abs(path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	abs, err := filepath.Abs(filepath.Join(env.Workdir, path))
	if err != nil {
		return "", fmt.Errorf("failed to make absolute: %w", err)
	}
	return abs, nil
}

// atMaxIterCount is the maximum number of recursions for the At function.
// This prevents infinite loops in a symlinked filesystem.
const atMaxIterCount = 1000

// At returns the local path to a repository at the provided path, as well as the relative path within the repository.
//
// The algorithm proceeds as follows:
//
// First check if there is a repository at the provided path.
// If there is a repository, returns it.
// If there is not, recursively try parent directories until outside of the root directory.
//
// Assumes that the root directory is set.
// If that is not the case, calls panic().
// If no repository is found, returns an error of type Error.
func (env Env) At(p string) (repo, worktree string, err error) {
	// Changes here should be reflected in AtRoot().
	root, err := env.absRoot()
	if err != nil {
		panic("Env.At: Root not resolved")
	}

	// find the absolute path that p points to
	// so that we can start searching
	path, err := env.Abs(p)
	if err != nil {
		return "", "", errNotResolved.WithMessageF(p)
	}

	// start recursively searching, starting at 'path' doing at most count iterations.
	// the regular exit condition is that repo should be the root of a repository.
	// we additionally need to check that the path is inside of the root.
	repo = path
	count := atMaxIterCount
	for !env.Git.IsRepository(repo) {
		count--
		repo = filepath.Join(repo, "..")
		if !strings.HasPrefix(repo, root) || root == "" || root == "/" || count == 0 {
			return "", "", errNotResolved.WithMessageF(p)
		}
	}

	// we have found the worktree path and the repository.
	worktree, err = filepath.Rel(repo, path)
	if err != nil {
		return "", "", errNotResolved.WithMessageF(root)
	}

	return
}

// AtRoot checks if the path p represents the root of a repository.
// If p is a relative path, it will be resolved relative to the current directory.
//
// When true it returns the absolute path to p, and no error.
// When false, returns the empty string and no error.
// When something goes wrong, returns an error.
func (env Env) AtRoot(p string) (repo string, err error) {
	// This function could check if At(p) returns worktree = "."
	// but that would create additional disk I/O!

	path, err := env.Abs(p)
	if err != nil {
		return "", errNotResolved.WithMessageF(p)
	}

	if !env.Git.IsRepository(path) {
		return "", nil
	}

	return path, nil
}

// Canonical returns the canonical version of the URL url.
// This requires that CanFile is not nil.
// See the CanonicalWith() method of URL.
//
// This function is untested.
func (env Env) Canonical(url URL) string {
	if env.CanFile == nil {
		panic("Env.Canonical: CanFile is nil")
	}
	return url.CanonicalWith(env.CanFile)
}

// reposBufferSize is the (currently hard-coded) size for the cache of the Repos function.
// 200 should be larger than the largest number of repositories expected.
// Note that this is only an optimization, the algorithm should perform even for a non-buffered channel.
const reposBufferSize = 200

// reposMaxParallelScan is the maximum number of folders to scan concurrently.
// Set to 0 for unlimited.
const reposMaxParallelScan = 0

// RepoScores returns a list of local paths to all repositories in this Environment.
// It also returns their scores along with the repositories.
// resolved indicates if their final path should be resolved
//
// This method silently ignores all errors.
//
// See the ScanReposScores() method for more control.
func (env Env) RepoScores(resolved bool) ([]string, []float64) {
	// NOTE: This function is untested, because only the score-less variant is tested.
	repos, scores, _ := env.ScanReposScores("", resolved)
	return repos, scores
}

// Repos returns a list of local paths to all repositories in this Environment.
// Resolved indicates if the final repository paths should be resolved.
// This method silently ignores all errors.
//
// See the ScanRepos() method for more control.
func (env Env) Repos(resolved bool) []string {
	// NOTE: This function is untested, because ScanRepos() is tested.
	repos, _ := env.RepoScores(resolved)
	return repos
}

// ScanRepoScores scans for repositories in the provided folder that match the Filter of this environment.
// Resolved indicates if the paths returned should resolve the final path of repositories.
// Repositories are returned in order of their scores, which are returned in the second argument.
//
// When an error occurs, this function may still return a list of (incomplete) repositories along with an error.
func (env Env) ScanReposScores(folder string, resolved bool) ([]string, []float64, error) {
	// NOTE: This function is untested, only ScanRepos() itself is tested
	if folder == "" {
		var err error
		folder, err = env.absRoot()
		if err != nil {
			panic("Env.Repos: Root not resolved")
		}
	}

	// grab extra candidates from the filter
	extraRoots := Candidates(env.Filter)
	n := 0
	for _, path := range extraRoots {
		if ok, err := fsx.IsDirectory(path, true); err == nil && ok {
			extraRoots[n] = path
			n++
		}
	}
	extraRoots = extraRoots[:n]

	extraFS := make([]walker.FS, len(extraRoots))
	for i, root := range extraRoots {
		extraFS[i] = walker.NewRealFS(root, true)
	}

	scanner := &walker.Walker[struct{}]{
		Process: walker.ScanProcess(func(path string, _ walker.FS, _ int) (score float64, cont bool, err error) {
			if env.Git.IsRepositoryQuick(path) {
				return env.Filter.Score(env, path), false, nil // never continue even if a repository does not match
			}
			return walker.ScanMatch(false), true, nil
		}),
		Params: walker.Params{
			Root: walker.NewRealFS(folder, true),

			ExtraRoots: extraFS,

			BufferSize:  reposBufferSize,
			MaxParallel: reposMaxParallelScan,
		},
	}

	err := scanner.Walk()
	return scanner.Paths(resolved), scanner.Scores(), err
}

// ScanRepos is like ScanReposScores, but returns only the first and last return value.
func (env Env) ScanRepos(folder string, resolved bool) ([]string, error) {
	results, _, err := env.ScanReposScores(folder, resolved)
	return results, err
}

//spellchecker:words nosec
