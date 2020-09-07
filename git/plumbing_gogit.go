package git

import (
	"os"
	"path"

	git "github.com/go-git/go-git/v5"
	"github.com/pkg/errors"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/util"
)

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
	return !os.IsNotExist(err) && s.Mode().IsDir()
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
	if util.SliceEquals(remotes.URLs, urls) {
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
		err = ignoreErrUpToDate(err)

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
	err = ignoreErrUpToDate(err)
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

func ignoreErrUpToDate(err error) error {
	if err == git.NoErrAlreadyUpToDate {
		os.Stdout.WriteString(err.Error() + "\n")
		err = nil
	}
	return err
}

func init() {
	// check that goGitImpl is a git implementation
	var _ Plumbing = (*gogit)(nil)
}
