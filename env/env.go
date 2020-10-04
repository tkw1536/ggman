package env

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/util"
)

// Env represents an environment to be used by ggman.
//
// The environment defines which git repositories are managed by ggman and where these are stored.
// It furthremore determins how a URL is canonicalized using a CanFile.
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

	// Filter is an optional filter for the environment.
	// Repositories not matching the filter are considered to not be a part of the environment.
	// See the Repos() method.
	Filter Filter

	// CanFile is the CanFile used to canonicalize repositories.
	// See the Canonical() method.
	CanFile CanFile
}

// Requirement represents a set of requirements on the Environment.
type Requirement struct {
	// Does the environment require a root directory?
	NeedsRoot bool

	// Does the environment allow filtering?
	// AllowsFilter implies NeedsRoot.
	AllowsFilter bool

	// Does the environment require a CanFile?
	NeedsCanFile bool
}

// NewEnv returns a new Env that fullfills the requirement r.
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
func NewEnv(r Requirement, vars Variables, filter Filter) (Env, error) {
	env := &Env{
		Git:    git.NewGitFromPlumbing(nil, vars.PATH),
		Vars:   vars,
		Filter: filter,
	}

	if r.NeedsRoot || r.AllowsFilter { // AllowsFilter implies NeedsRoot
		if err := env.LoadDefaultRoot(); err != nil {
			return *env, err
		}
	}

	if r.NeedsCanFile {
		if err := env.LoadDefaultCANFILE(); err != nil {
			return *env, err
		}
	}

	return *env, nil
}

var errMissingRoot = ggman.Error{
	ExitCode: ggman.ExitInvalidEnvironment,
	Message:  "Unable to find GGROOT directory. ",
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
		err = errInvalidRoot.WithMessageF(err)
		return "", errMissingRoot
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

// LoadDefaultCANFILE sets env.CANFILE according to the environment variables in e.Vars.
// If the CANFILE is already set, immediatly returns nil.
//
// If the GGMAN_CANFILE variable is set, it will use it as a filepath to read the CanFile from.
// If it is not set it will attempt to load the file '.ggman' in the home directory.
// If neither is set, this function will load an in-memory default CanFile.
//
// If loading a CanFile fails, an error of type Error is returned.
// If loading succeeds, this function returns nil.
func (env *Env) LoadDefaultCANFILE() error {
	if env.CanFile != nil {
		return nil
	}

	files := make([]string, 0, 2)
	if env.Vars.CANFILE != "" {
		files = append(files, env.Vars.CANFILE)
	}
	if env.Vars.HOME != "" {
		files = append(files, filepath.Join(env.Vars.HOME, ".ggman"))
	}

	var cf CanFile

	// In order, if a file exists read it or fail.
	// If it doesn't exist continue to the next file.
	for _, file := range files {
		f, err := os.Open(file)
		switch {
		case err == nil: // do nothing
		case os.IsNotExist(err):
			continue
		default:
			return errors.Wrapf(err, "Unable to open CANFILE %q", file)
		}
		defer f.Close()
		if _, err := (&cf).ReadFrom(f); err != nil {
			return err
		}
		env.CanFile = cf
		return nil
	}

	(&cf).ReadDefault()
	env.CanFile = cf
	return nil
}

// Local is the localpath to the repository pointed to by URL
func (env Env) Local(url URL) string {
	root, err := env.absRoot()
	if err != nil {
		panic("Env.Local(): Root not resolved")
	}

	path := append([]string{root}, url.Components()...)
	return filepath.Join(path...)
}

var errInvalidRoot = ggman.Error{
	ExitCode: ggman.ExitInvalidEnvironment,
	Message:  "Unable to resolve root directory: %s",
}

var errNotResolved = ggman.Error{
	ExitCode: ggman.ExitInvalidRepo,
	Message:  "Unable to resolve repository %q",
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
	root, err := env.absRoot()
	if err != nil {
		panic("Env.At(): Root not resolved")
	}

	// find the absolute path that p points to
	// so that we can start searching
	path, err := filepath.Abs(p)
	if err != nil {
		return "", "", errNotResolved.WithMessageF(p)
	}

	// start recursively searching, starting at 'path' doing at most count iterations.
	// the regular exit condition is that repo should be the root of a repository.
	// we addtionally need to check that the path is inside of the root.
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

// Canonical returns the canonical version of the URL url.
// This requires that CanFile is not nil.
// See the CanonicalWith() method of URL.
//
// This function is untested.
func (env Env) Canonical(url URL) string {
	if env.CanFile == nil {
		panic("Env.Canonical(): CanFile is nil")
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

// Repos returns a list of local paths to all repositories in this Environment.
// This method silently ignores all errors.
//
// See the ScanRepos() method for more control.
//
// This function is untested, because ScanRepos() is tested.
func (env Env) Repos() []string {
	repos, _ := env.ScanRepos("")
	return repos
}

// ScanRepos scans for repositories in the provided folder that match the Filter of this environment.
// When an error occurs, this function may still return a list of (incomplete) repositories along with an error.
func (env Env) ScanRepos(folder string) ([]string, error) {
	if folder == "" {
		var err error
		folder, err = env.absRoot()
		if err != nil {
			panic("env.Repos(): Root not resolved")
		}
	}

	return util.Scan(util.ScanOptions{
		Root: folder,

		Filter: func(path string) (match, cont bool) {
			if env.Git.IsRepositoryQuick(path) {
				return env.Filter.Matches(env.Root, path), false // never continue even if a repository does not match
			}
			return false, true
		},

		BufferSize:  reposBufferSize,
		MaxParallel: reposMaxParallelScan,
		FollowLinks: false, // TODO: Make this configurable
	})
}
