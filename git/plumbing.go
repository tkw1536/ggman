package git

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	git "github.com/go-git/go-git/v5"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/internal/text"
)

// Plumbing is an interface that represents a working internal implementation of git.
// Plumbing is intended to be goroutine-safe, i.e. everything except the Init() method can be called from multiple goroutines at once.
//
// It is not intended to be called directly by an external caller, instead it is intended to be called by the Git interface only.
// The reason for this is that it requires initalization and places certain assumptions on the caller.
//
// For instance, to pull a repository, the following code is required:
//
//  plumbing.Init() // called exactly once
//  cache, isRepo := plumbing.IsRepository("/home/user/Projects/github.com/hello/world")
//  if !isRepo {
//    // error, not a repository
//  }
//  err = plumbining.Pull(ggman.NewEnvIOStream(), "/home/user/Projects/github.com/hello/world", cache)
//
// Such code is typically handled by a Git instance that wraps a Plumbing.
type Plumbing interface {

	// Init is used to initialize this Plumbing.
	// Init should only be called once per Plumbing instance.
	// When initialization fails, for example due to missing dependencies, returns a non-nil error,
	Init() error

	// IsRepository checks if the directory at localPath is the root of a git repository.
	// May assume that localPath exists and is a repository.
	//
	// This function returns a pair, a boolean isRepo that indicates if this object is a repository
	// and an optional repoObject value.
	// The repoObject value will only be taken into account when isRepo is true, and passed to other functions in this git implementation.
	// The semantics of the repoObject are determined by this Plumbing and should not be used outside of it.
	// Note that the repoObject may be used for more than one subsequent call.
	//
	// This function surpresses all errors, and if something goes wrong assumed that isRepo is false.
	IsRepository(localPath string) (repoObject interface{}, isRepo bool)

	// IsRepositoryUnsafe efficiently checks if the directly at localPath contains a repository.
	// It is like IsRepository, except that it may return false positives, but no false negatives.
	// This function is optimized to be called a lot of times.
	IsRepositoryUnsafe(localPath string) bool

	// GetHeadRef returns a reference to the current head of the repository cloned at clonePath.
	// The string ref should contain a git REFLIKE, that is a branch, a tag or a commit id.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetHeadRef(clonePath string, repoObject interface{}) (ref string, err error)

	// GetRemotes returns the names and urls of the remotes of the repository cloned at clonePath.
	// If determining the remotes is not possible, and error is returned instead.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetRemotes(clonePath string, repoObject interface{}) (remotes map[string][]string, err error)

	// GetCanonicalRemote gets the name of the canonical remote of the reposity cloned at clonePath.
	// The Plumbing is free to decided what the canonical remote is, but it is typically the remote of the currently checked out branch or the 'origin' remote.
	// If no remote exists, an empty name is returned.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetCanonicalRemote(clonePath string, repoObject interface{}) (name string, urls []string, err error)

	// SetRemoteURLs set the remote 'remote' of the repository at clonePath to urls.
	// The remote 'name' must exist.
	// Furthermore newURLs must be of the same length as the old URLs.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	SetRemoteURLs(clonePath string, repoObject interface{}, name string, urls []string) (err error)

	// Clone tries to clone the repository at 'from' to the folder 'to'.
	// May attempt to read credentials from stream.Stdin.
	// Output is directed to stream.Stdout and stream.Stderr.
	//
	// remoteURI will be the uri of the remote repository.
	// clonePath will be the path to a local folder where the repository should be cloned to.
	// It's parent is guaranteed to exist.
	//
	// extraargs will be additional arguments, in the form of arguments of a 'git clone' command.
	// When this implementation does not support arguments, it returns ErrArgumentsUnsupported whenever arguments is a list of length > 0.
	//
	// If the clone succeeds returns, err = nil.
	// If the underlying clone command returns a non-zero code, returns an error of type ExitError.
	// If something else goes wrong, may return any other error type.
	Clone(stream ggman.IOStream, remoteURI, clonePath string, extraargs ...string) error

	// Fetch should fetch new objects and refs from all remotes of the repository cloned at clonePath.
	// May attempt to read credentials from stream.Stdin.
	// Output is directed to stream.Stdout and stream.Stderr.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	Fetch(stream ggman.IOStream, clonePath string, cache interface{}) (err error)

	// Pull should fetch new objects and refs from all remotes of the repository cloned at clonePath.
	// It then merges them into the local branch wherever an upstream is set.
	// May attempt to read credentials from stream.Stdin.
	// Output is directed to stream.Stdout and stream.Stderr.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	Pull(stream ggman.IOStream, clonePath string, cache interface{}) (err error)

	// ContainsBranch checks if the repository at clonePath contains a branch with the provided branch.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	ContainsBranch(clonePath string, cache interface{}, branch string) (contains bool, err error)

	// IsDirty checks if the repository at clonePath contains uncommitted changes.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	IsDirty(clonePath string, cache interface{}) (dirty bool, err error)
}

// NewPlumbing returns an implementation of a plumbing that has no external dependencies.
// The plumbing is guaranteed to have been initialized.
//
// There is no guarantee as to what plumbing is returned.
func NewPlumbing() Plumbing {
	gg := gogit{}
	gg.Init()
	return gg
}

// ErrArgumentsUnsupported is an error that is returned when arguments are not supported by a Plumbing.
var ErrArgumentsUnsupported = errors.New("Plumbing does not support extra clone arguments")

//
// gitgit
//

type gitgit struct {
	gogit
	gitPath string
}

func (gg *gitgit) Init() (err error) {
	gg.gitPath, err = gg.findgit()
	return
}

func (gg gitgit) findgit() (git string, err error) {
	if runtime.GOOS == "windows" {
		return gg.findGitByExtension([]string{"exe"})
	}
	return gg.findGitByMode()
}

// findGitByMode finds git by finding a file named 'git' with executable flag set
func (gg gitgit) findGitByMode() (git string, err error) {
	// this code has been adapted from exec.LookPath in the standard library
	// it allows using a more generic path variables
	for _, git := range filepath.SplitList(gg.gitPath) {
		if git == "" { // unix shell behavior
			git = "."
		}
		git = filepath.Join(git, "git")
		d, err := os.Stat(git)
		if err != nil {
			continue
		}
		if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
			return git, nil
		}
	}
	return "", exec.ErrNotFound
}

// findGitByExtension finds the git executable by looking for a non-directory file named "git.extension" where extension is in ext
func (gg gitgit) findGitByExtension(exts []string) (git string, err error) {
	// this code has been adapted from exec.LookPath in the standard library
	// it allows using a more generic path variables
	for _, git := range filepath.SplitList(gg.gitPath) {
		if git == "" { // unix shell behavior
			git = "."
		}
		for _, ext := range exts {
			git = filepath.Join(git, "git."+ext) // TODO: Case insensitive extensions
			d, err := os.Stat(git)
			if err != nil {
				continue
			}
			if m := d.Mode(); !m.IsDir() {
				return git, nil
			}
		}
	}
	return "", exec.ErrNotFound
}

func (gg gitgit) Clone(stream ggman.IOStream, remoteURI, clonePath string, extraargs ...string) error {

	gargs := append([]string{"clone", remoteURI, clonePath}, extraargs...)

	cmd := exec.Command(gg.gitPath, gargs...)
	cmd.Stdin = stream.Stdin
	cmd.Stdout = stream.Stdout
	cmd.Stderr = stream.Stderr

	// run the underlying command, but treat ExitError specially by turning it into a ExitError
	err := cmd.Run()
	if exitError, isExitError := err.(*exec.ExitError); isExitError {
		err = ExitError{error: err, Code: exitError.ExitCode()}
	}
	return err
}

func (gg gitgit) IsDirty(clonePath string, cache interface{}) (dirty bool, err error) {
	cmd := exec.Command(gg.gitPath, "diff", "--quiet")
	cmd.Dir = clonePath

	// run the underlying command
	err = cmd.Run()
	if exitError, isExitError := err.(*exec.ExitError); isExitError {
		// code 1: it is dirty
		if exitError.ExitCode() == 1 {
			return true, nil
		}
		err = ExitError{error: err, Code: exitError.ExitCode()}
	}
	return false, err
}

//
// gogit
//

type gogit struct{}

func (gogit) Init() error {
	// this git implementation has no data to be initialized
	return nil
}

func (gogit) IsRepository(localPath string) (repoObject interface{}, isRepo bool) {
	repoObject, err := git.PlainOpen(localPath)
	return repoObject, err == nil
}

func (gogit) IsRepositoryUnsafe(localPath string) bool {
	// path to .git
	gitPath := path.Join(localPath, ".git")

	// check that it exists and is a folder
	s, err := os.Stat(gitPath)
	isNotNotExist := !os.IsNotExist(err)
	if err != nil && isNotNotExist {
		return false
	}
	return isNotNotExist && s.Mode().IsDir()
}

func (gogit) GetHeadRef(clonePath string, repoObject interface{}) (string, error) {
	repo := repoObject.(*git.Repository)

	// first get the name of the current HEAD
	// or fail if that isn't possible
	head, err := repo.Head()
	if err != nil {
		err = errors.Wrap(err, "Cannot resolve HEAD")
		return "", err
	}
	name := head.Name()

	// if we have a branch or a tag, return the reference to it
	if name.IsBranch() {
		return name.Short(), nil
	}

	// else just return the plain old hash
	return head.Hash().String(), nil
}

func (gogit) GetRemotes(clonePath string, repoObject interface{}) (remotes map[string][]string, err error) {
	// get the repository
	r := repoObject.(*git.Repository)

	// get all the remotes for the repository
	gitRemotes, err := r.Remotes()
	if err != nil {
		err = errors.Wrap(err, "Unable to get remotes")
		return
	}

	// make a map for remotes
	remotes = make(map[string][]string, len(gitRemotes))
	for _, r := range gitRemotes {
		cfg := r.Config()
		remotes[cfg.Name] = cfg.URLs
	}

	return
}

// originRemoteName is the name of the canonical remote
const originRemoteName = "origin"

func (gg gogit) GetCanonicalRemote(clonePath string, repoObject interface{}) (name string, urls []string, err error) {
	// get a map of remotes
	remotes, err := gg.GetRemotes(clonePath, repoObject)
	if err != nil {
		err = errors.Wrap(err, "Unable to get remotes")
		return
	}

	// if we don't have any remotes we're done
	if len(remotes) == 0 {
		return
	}

	// if the current branch has a remote, use it
	r := repoObject.(*git.Repository)
	name, _ = gg.getCurrentBranchRemote(r)
	if name != "" {
		urls = remotes[name]
		return
	}

	// else if we have an 'origin' remote we use that
	if originRemote, originRemoteExists := remotes[originRemoteName]; originRemoteExists {
		urls = originRemote
		name = originRemoteName
		return
	}

	// else randomly use the first remote that we have
	for rn, ru := range remotes {
		urls = ru
		name = rn
		return
	}

	panic("never reached")
}

func (gogit) getCurrentBranch(r *git.Repository) (name string, err error) {

	// determine the current head and name of it
	head, err := r.Head()
	if err != nil {
		err = errors.Wrap(err, "Cannot resolve HEAD")
		return
	}

	// ensure that it's a branch
	headName := head.Name()
	if !headName.IsBranch() {
		err = errors.New("Not on a branch")
		return
	}

	// return the name
	name = headName.String()
	return
}

func (gg gogit) getCurrentBranchRemote(r *git.Repository) (name string, err error) {
	// get the current branch
	branchName, err := gg.getCurrentBranch(r)
	if err != nil {
		err = errors.Wrap(err, "Unable to get current branch")
		return
	}

	// get its' configuration
	branch, err := r.Branch(branchName)
	if err != nil {
		err = errors.Wrap(err, "Cannot find branch config")
		return
	}

	// and check that the remote is non-empty
	name = branch.Remote
	if name == "" {
		err = errors.New("Branch does not have an associated remote")
		return
	}

	return
}

func (gogit) SetRemoteURLs(clonePath string, repoObject interface{}, name string, urls []string) (err error) {
	// get the repository
	r := repoObject.(*git.Repository)

	// get the desired remote
	remote, err := r.Remote(name)
	if err != nil {
		err = errors.Wrapf(err, "Unable to find remote %s", name)
		return
	}

	// fetch the current configuration
	cfg, err := r.Storer.Config()
	if err != nil {
		return
	}

	// get the current remotes
	remotes := cfg.Remotes[remote.Config().Name]

	// if they haven't changed, we can return immediatly
	if text.SliceEquals(remotes.URLs, urls) {
		return nil
	}

	// check that they are of the new length
	if len(remotes.URLs) != len(urls) {
		return errors.New("Cannot set remoteURL: Length of old and new urls must be identical")
	}

	// Write back the URLs
	remotes.URLs = urls
	if err = r.SetConfig(cfg); err != nil {
		err = errors.Wrap(err, "Unable to store config")
		return
	}

	return
}

func (gogit) Clone(stream ggman.IOStream, remoteURI, clonePath string, extraargs ...string) error {
	// doesn't support extra arguments
	if len(extraargs) > 0 {
		return ErrArgumentsUnsupported
	}

	// run a plain git clone but intercept all errors
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{URL: remoteURI, Progress: stream.Stdout})
	if err != nil {
		err = ExitError{error: errors.Wrap(err, "Unable clone repository"), Code: 1}
	}

	return err
}

func (gogit) Fetch(stream ggman.IOStream, clonePath string, cache interface{}) (err error) {
	// get the repository
	r := cache.(*git.Repository)

	// list all of the remotes
	remotes, err := r.Remotes()
	if err != nil {
		return
	}

	// fetch all of the remotes for this repository
	for _, remote := range remotes {
		// fetch and write out an 'already up-to-date'
		err = remote.Fetch(&git.FetchOptions{Progress: stream.Stdout})
		err = ignoreErrUpToDate(stream, err)

		// fail on other errors
		if err != nil {
			err = errors.Wrapf(err, "Unable to fetch remote %s", remote.Config().Name)
			return
		}
	}

	return
}

func (gogit) Pull(stream ggman.IOStream, clonePath string, cache interface{}) (err error) {
	// get the repository
	r := cache.(*git.Repository)

	// get the worktree
	w, err := r.Worktree()
	if err != nil {
		err = errors.Wrap(err, "Unable to find worktree")
		return
	}

	// do a git pull, and ignore error already up-to-date
	err = w.Pull(&git.PullOptions{Progress: stream.Stdout})
	err = ignoreErrUpToDate(stream, err)
	if err != nil {
		err = errors.Wrap(err, "Unable to pull")
	}

	return
}

func (gogit) ContainsBranch(clonePath string, cache interface{}, branch string) (contains bool, err error) {
	// get the repository
	r := cache.(*git.Repository)

	// try to open the branch
	switch _, err := r.Branch(branch); err {
	case git.ErrBranchNotFound:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, errors.Wrap(err, "Unable to read branch")
	}
}

func (gogit) IsDirty(clonePath string, cache interface{}) (dirty bool, err error) {
	// get the repository
	r := cache.(*git.Repository)

	// get the worktree
	wt, err := r.Worktree()
	if err != nil {
		return false, errors.Wrap(err, "Unable to get worktree")
	}

	// check the status
	status, err := wt.Status()
	if err != nil {
		return false, errors.Wrap(err, "Unable to get status")
	}

	// return if it is dirty!
	return !status.IsClean(), nil
}

func ignoreErrUpToDate(stream ggman.IOStream, err error) error {
	if err == git.NoErrAlreadyUpToDate {
		stream.StdoutWriteWrap(err.Error())
		err = nil
	}
	return err
}
