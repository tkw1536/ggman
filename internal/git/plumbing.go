package git

//spellchecker:words context errors exec path filepath runtime slices github plumbing pkglib exit stream
import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"slices"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
	"go.tkw01536.de/pkglib/stream"
)

//spellchecker:words worktree bref reflike gogit gitgit wrapf storer

// Plumbing is an interface that represents a working internal implementation of git.
// Plumbing is intended to be goroutine-safe, i.e. everything except the Init() method can be called from multiple goroutines at once.
//
// It is not intended to be called directly by an external caller, instead it is intended to be called by the Git interface only.
// The reason for this is that it requires initialization and places certain assumptions on the caller.
//
// For instance, to pull a repository, the following code is required:
//
//	plumbing.Init() // called exactly once
//	cache, isRepo := plumbing.IsRepository("/home/user/Projects/github.com/hello/world")
//	if !isRepo {
//	  // error, not a repository
//	}
//	err = plumbing.Pull(stream.NewEnvIOStream(), "/home/user/Projects/github.com/hello/world", cache)
//
// Such code is typically handled by a Git instance that wraps a Plumbing.
type Plumbing interface {

	// Init is used to initialize this Plumbing.
	// Init should only be called once per Plumbing instance.
	// When initialization fails, for example due to missing dependencies, returns a non-nil error.
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
	// This function suppresses all errors, and if something goes wrong assumed that isRepo is false.
	IsRepository(ctx context.Context, localPath string) (repoObject any, isRepo bool)

	// IsRepositoryUnsafe efficiently checks if the directly at localPath contains a repository.
	// It is like IsRepository, except that it may return false positives, but no false negatives.
	// This function is optimized to be called a lot of times.
	IsRepositoryUnsafe(ctx context.Context, localPath string) bool

	// GetHeadRef returns a reference to the current head of the repository cloned at clonePath.
	// The string ref should contain a git REFLIKE, that is a branch, a tag or a commit id.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetHeadRef(ctx context.Context, clonePath string, repoObject any) (ref string, err error)

	// GetRemotes returns the names and urls of the remotes of the repository cloned at clonePath.
	// If determining the remotes is not possible, and error is returned instead.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetRemotes(ctx context.Context, clonePath string, repoObject any) (remotes map[string][]string, err error)

	// GetCanonicalRemote gets the name of the canonical remote of the repository cloned at clonePath.
	// The Plumbing is free to decided what the canonical remote is, but it is typically the remote of the currently checked out branch or the 'origin' remote.
	// If no remote exists, an empty name is returned.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetCanonicalRemote(ctx context.Context, clonePath string, repoObject any) (name string, urls []string, err error)

	// SetRemoteURLs set the remote 'remote' of the repository at clonePath to urls.
	// The remote 'name' must exist.
	// Furthermore newURLs must be of the same length as the old URLs.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	SetRemoteURLs(ctx context.Context, clonePath string, repoObject any, name string, urls []string) (err error)

	// Clone tries to clone the repository at 'from' to the folder 'to'.
	// May attempt to read credentials from stream.Stdin.
	// Output is directed to stream.Stdout and stream.Stderr.
	//
	// remoteURI will be the uri of the remote repository.
	// clonePath will be the path to a local folder where the repository should be cloned to.
	// It's parent is guaranteed to exist.
	//
	// extraArgs will be additional arguments, in the form of arguments of a 'git clone' command.
	// When this implementation does not support arguments, it returns ErrArgumentsUnsupported whenever arguments is a list of length > 0.
	//
	// If the clone succeeds returns, err = nil.
	// If the underlying clone command returns a non-zero code, returns an error of type ExitError.
	// If something else goes wrong, may return any other error type.
	Clone(ctx context.Context, stream stream.IOStream, remoteURI, clonePath string, extraArgs ...string) error

	// Fetch should fetch new objects and refs from all remotes of the repository cloned at clonePath.
	// May attempt to read credentials from stream.Stdin.
	// Output is directed to stream.Stdout and stream.Stderr.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	Fetch(ctx context.Context, stream stream.IOStream, clonePath string, cache any) (err error)

	// Pull should fetch new objects and refs from all remotes of the repository cloned at clonePath.
	// It then merges them into the local branch wherever an upstream is set.
	// May attempt to read credentials from stream.Stdin.
	// Output is directed to stream.Stdout and stream.Stderr.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	Pull(ctx context.Context, stream stream.IOStream, clonePath string, cache any) (err error)

	// GetBranches gets the names of all branches contained in the repository at clonePath.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	GetBranches(ctx context.Context, clonePath string, cache any) (branches []string, err error)

	// ContainsBranch checks if the repository at clonePath contains a branch with the provided branch.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	ContainsBranch(ctx context.Context, clonePath string, cache any, branch string) (contains bool, err error)

	// IsDirty checks if the repository at clonePath contains uncommitted changes.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	IsDirty(ctx context.Context, clonePath string, cache any) (dirty bool, err error)

	// IsSync checks if the repository at clonePath does not have branches synced with their upstream.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	IsSync(ctx context.Context, clonePath string, cache any) (dirty bool, err error)
}

// NewPlumbing returns an implementation of a plumbing that has no external dependencies.
// The plumbing is guaranteed to have been initialized.
//
// There is no guarantee as to what plumbing is returned.
func NewPlumbing() Plumbing {
	// NOTE: We cast here to avoid a warning that the Init method is a noop.
	// We want to keep it in case it does something in the future.
	gg := Plumbing(gogit{})
	_ = gg.Init() // ignore cause gogit always returns a nil error
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
	gg.gitPath, err = gg.findGit()
	return
}

func (gg *gitgit) findGit() (git string, err error) {
	if runtime.GOOS == "windows" {
		return gg.findGitByExtension([]string{"exe"})
	}
	return gg.findGitByMode()
}

// findGitByMode finds git by finding a file named 'git' with executable flag set.
func (gg *gitgit) findGitByMode() (git string, err error) {
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

// findGitByExtension finds the git executable by looking for a non-directory file named "git.extension" where extension is in ext.
func (gg *gitgit) findGitByExtension(exts []string) (git string, err error) {
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

func (gg *gitgit) Clone(ctx context.Context, stream stream.IOStream, remoteURI, clonePath string, extraArgs ...string) error {
	gitArgs := append([]string{"clone", remoteURI, clonePath}, extraArgs...)

	cmd := exec.CommandContext(ctx, gg.gitPath, gitArgs...) /* #nosec G204  -- user-controlled by design */
	cmd.Stdin = stream.Stdin
	cmd.Stdout = stream.Stdout
	cmd.Stderr = stream.Stderr

	// run the underlying command, but treat ExitError specially by turning it into a ExitError
	err := cmd.Run()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		err = exit.FromExitError(exitError)
	}

	return err
}

func (gg *gitgit) Fetch(ctx context.Context, stream stream.IOStream, clonePath string, cache any) error {
	cmd := exec.CommandContext(ctx, gg.gitPath, "fetch", "--all") /* #nosec G204  -- gitPath user-controlled by design */
	cmd.Dir = clonePath
	cmd.Stdin = stream.Stdin
	cmd.Stdout = stream.Stdout
	cmd.Stderr = stream.Stderr

	// run the underlying command, but treat ExitError specially by turning it into a ExitError
	err := cmd.Run()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		err = exit.FromExitError(exitError)
	}
	return err
}

func (gg *gitgit) Pull(ctx context.Context, stream stream.IOStream, clonePath string, cache any) error {
	cmd := exec.CommandContext(ctx, gg.gitPath, "pull") /* #nosec G204  -- gitPath user-controlled by design */
	cmd.Dir = clonePath
	cmd.Stdin = stream.Stdin
	cmd.Stdout = stream.Stdout
	cmd.Stderr = stream.Stderr

	// run the underlying command, but treat ExitError specially by turning it into a ExitError
	err := cmd.Run()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		err = exit.FromExitError(exitError)
	}
	return err
}

func (gg *gitgit) IsDirty(ctx context.Context, clonePath string, cache any) (dirty bool, err error) {
	cmd := exec.CommandContext(ctx, gg.gitPath, "diff", "--quiet") /* #nosec G204 -- gitPath user-controlled by design */
	cmd.Dir = clonePath

	// run the underlying command
	err = cmd.Run()

	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		// code 1: it is dirty
		if exitError.ExitCode() == 1 {
			return true, nil
		}
		err = exit.FromExitError(exitError)
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

func (gogit) IsRepository(ctx context.Context, localPath string) (repoObject any, isRepo bool) {
	defer func() {
		if recover() != nil {
			repoObject = nil
			isRepo = false
		}
	}()
	repoObject, err := git.PlainOpen(localPath)
	return repoObject, err == nil
}

func (gogit) IsRepositoryUnsafe(ctx context.Context, localPath string) bool {
	// path to .git
	gitPath := path.Join(localPath, ".git")

	// check that it exists and is a directory
	ok, _ := fsx.IsDirectory(gitPath, false)
	return ok
}

func (gogit) GetHeadRef(ctx context.Context, clonePath string, repoObject any) (string, error) {
	repo := repoObject.(*git.Repository)

	// first get the name of the current HEAD
	// or fail if that isn't possible
	head, err := repo.Head()
	if err != nil {
		err = fmt.Errorf("%q: cannot resolve HEAD: %w", clonePath, err)
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

func (gogit) GetRemotes(ctx context.Context, clonePath string, repoObject any) (remotes map[string][]string, err error) {
	// get the repository
	r := repoObject.(*git.Repository)

	// get all the remotes for the repository
	gitRemotes, err := r.Remotes()
	if err != nil {
		err = fmt.Errorf("%q: unable to get remotes: %w", clonePath, err)
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

// originRemoteName is the name of the canonical remote.
const originRemoteName = "origin"

func (gg gogit) GetCanonicalRemote(ctx context.Context, clonePath string, repoObject any) (name string, urls []string, err error) {
	// get a map of remotes
	remotes, err := gg.GetRemotes(ctx, clonePath, repoObject)
	if err != nil {
		err = fmt.Errorf("%q: unable to get remotes: %w", clonePath, err)
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

var errNotOnABranch = errors.New("not on a branch")

func (gogit) getCurrentBranch(r *git.Repository) (name string, err error) {
	// determine the current head and name of it
	head, err := r.Head()
	if err != nil {
		err = fmt.Errorf("cannot resolve HEAD: %w", err)
		return
	}

	// ensure that it's a branch
	headName := head.Name()
	if !headName.IsBranch() {
		err = errNotOnABranch
		return
	}

	// return the name
	name = headName.String()
	return
}

var errBranchNoRemote = errors.New("branch does not have an associated remote")

func (gg gogit) getCurrentBranchRemote(r *git.Repository) (name string, err error) {
	// get the current branch
	branchName, err := gg.getCurrentBranch(r)
	if err != nil {
		err = fmt.Errorf("unable to get current branch: %w", err)
		return
	}

	// get its' configuration
	branch, err := r.Branch(branchName)
	if err != nil {
		err = fmt.Errorf("cannot find branch config: %w", err)
		return
	}

	// and check that the remote is non-empty
	name = branch.Remote
	if name == "" {
		err = errBranchNoRemote
		return
	}

	return
}

var errLengthMustBeEqual = errors.New("cannot set remoteURL: Length of old and new urls must be identical")

func (gogit) SetRemoteURLs(ctx context.Context, clonePath string, repoObject any, name string, urls []string) (err error) {
	// get the repository
	r := repoObject.(*git.Repository)

	// get the desired remote
	remote, err := r.Remote(name)
	if err != nil {
		err = fmt.Errorf("%q: unable to find remote %s: %w", clonePath, name, err)
		return
	}

	// fetch the current configuration
	cfg, err := r.Storer.Config()
	if err != nil {
		return
	}

	// get the current remotes
	remotes := cfg.Remotes[remote.Config().Name]

	// if they haven't changed, we can return immediately
	if slices.Equal(remotes.URLs, urls) {
		return nil
	}

	// check that they are of the new length
	if len(remotes.URLs) != len(urls) {
		return errLengthMustBeEqual
	}

	// Write back the URLs
	remotes.URLs = urls
	if err = r.SetConfig(cfg); err != nil {
		err = fmt.Errorf("%q: unable to store config: %w", clonePath, err)
		return
	}

	return
}

func (gogit) Clone(ctx context.Context, stream stream.IOStream, remoteURI, clonePath string, extraArgs ...string) error {
	// doesn't support extra arguments
	if len(extraArgs) > 0 {
		return ErrArgumentsUnsupported
	}

	// run a plain git clone but intercept all errors
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{URL: remoteURI, Progress: stream.Stderr})
	if err != nil {
		err = fmt.Errorf("%w: %w", exit.NewErrorWithCode(fmt.Sprintf("%q: unable clone repository", remoteURI), 1), err)
	}

	return err
}

func (gogit) Fetch(ctx context.Context, stream stream.IOStream, clonePath string, cache any) (err error) {
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
		err = remote.Fetch(&git.FetchOptions{Progress: stream.Stderr})
		err = ignoreErrUpToDate(stream, err)

		// fail on other errors
		if err != nil {
			err = fmt.Errorf("%q: unable to fetch remote %q: %w", clonePath, remote.Config().Name, err)
			return
		}
	}

	return
}

func (gogit) Pull(ctx context.Context, stream stream.IOStream, clonePath string, cache any) (err error) {
	// get the repository
	r := cache.(*git.Repository)

	// get the worktree
	w, err := r.Worktree()
	if err != nil {
		err = fmt.Errorf("%q: unable to find worktree: %w", clonePath, err)
		return
	}

	// do a git pull, and ignore error already up-to-date
	err = w.Pull(&git.PullOptions{Progress: stream.Stderr})
	err = ignoreErrUpToDate(stream, err)
	if err != nil {
		err = fmt.Errorf("%q: unable to pull: %w", clonePath, err)
	}

	return
}

func ignoreErrUpToDate(stream stream.IOStream, err error) error {
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		_, err = stream.Println(err.Error())
		if err != nil {
			err = fmt.Errorf("failed to print %w message: %w", git.NoErrAlreadyUpToDate, err)
		}
	}
	return err
}

func (gogit) GetBranches(ctx context.Context, clonePath string, cache any) (branches []string, err error) {
	// get the repository
	r := cache.(*git.Repository)

	// list the branches
	branchRefs, err := r.Branches()
	if err != nil {
		return nil, fmt.Errorf("%q: unable to get branches: %w", clonePath, err)
	}
	defer branchRefs.Close()

	// get their names
	if err := branchRefs.ForEach(func(bref *plumbing.Reference) error {
		branches = append(branches, bref.Name().Short())
		return nil
	}); err != nil {
		return nil, fmt.Errorf("%q: failed iterate branch refs: %w", clonePath, err)
	}

	return
}

func (gogit) ContainsBranch(ctx context.Context, clonePath string, cache any, branch string) (contains bool, err error) {
	// get the repository
	r := cache.(*git.Repository)

	// try to open the branch

	switch _, err := r.Branch(branch); {
	case errors.Is(err, git.ErrBranchNotFound):
		return false, nil
	case err == nil:
		return true, nil
	default:
		return false, fmt.Errorf("%q: unable to read branch %q: %w", clonePath, branch, err)
	}
}

func (gogit) IsDirty(ctx context.Context, clonePath string, cache any) (dirty bool, err error) {
	// get the repository
	r := cache.(*git.Repository)

	// get the worktree
	wt, err := r.Worktree()
	if err != nil {
		return false, fmt.Errorf("%q: unable to get worktree: %w", clonePath, err)
	}

	// check the status
	status, err := wt.Status()
	if err != nil {
		return false, fmt.Errorf("%q: unable to get status: %w", clonePath, err)
	}

	// return if it is dirty!
	return !status.IsClean(), nil
}

func (gg gogit) IsSync(ctx context.Context, clonePath string, cache any) (sync bool, err error) {
	r := cache.(*git.Repository)

	// get all the branches
	branches, err := gg.GetBranches(ctx, clonePath, cache)
	if err != nil {
		return false, fmt.Errorf("%q: unable to get branch names: %w", clonePath, err)
	}

	// check that all the upstream branches are synced!
	for _, b := range branches {
		src, dst, err := getTrackingRefs(r, b)
		if errors.Is(err, errNoUpstream) {
			continue // there is no upstream, that is ok!
		}
		if err != nil {
			return false, fmt.Errorf("%q: unable to get tracking refs: %w", clonePath, err)
		}
		srcRef, err := r.ResolveRevision(plumbing.Revision(src))
		if err != nil {
			return false, fmt.Errorf("%q: unable to resolve src revision: %w", clonePath, err)
		}
		dstRef, err := r.ResolveRevision(plumbing.Revision(dst))
		if err != nil {
			return false, fmt.Errorf("%q: unable to resolve destination revision: %w", clonePath, err)
		}
		if srcRef.String() != dstRef.String() {
			return false, nil
		}
	}
	return true, nil
}

var errNoUpstream = errors.New("no corresponding upstream to track")

// getTrackingRefs returns the src and dst upstream tracking refs for the provided branch.
// When the branch, or the upstream tracking refs do not exist, returns ErrNoUpstream.
func getTrackingRefs(repo *git.Repository, branch string) (src, dst plumbing.ReferenceName, err error) {
	br, err := repo.Branch(branch)
	if errors.Is(err, git.ErrBranchNotFound) {
		return "", "", errNoUpstream
	}
	if err != nil {
		return "", "", fmt.Errorf("unable to resolve branch %q: %w", branch, err)
	}
	if br.Remote == "" {
		return "", "", errNoUpstream
	}
	remote, err := repo.Remote(br.Remote)
	if err != nil {
		return "", "", fmt.Errorf("unable to resolve remote %q: %w", br.Remote, err)
	}
	for _, spec := range remote.Config().Fetch {
		if spec.Match(br.Merge) {
			return br.Merge, spec.Dst(br.Merge), nil
		}
	}
	return "", "", errNoUpstream
}

//spellchecker:words nosec
