package git

import (
	"os"
	"path"

	"github.com/pkg/errors"
	git "gopkg.in/src-d/go-git.v4"
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

func (gogit) GetHeadRef(clonePath string, repoObject interface{}) (ref string, err error) {
	// get the repository
	repo := repoObject.(*git.Repository)

	// get the current head
	head, err := repo.Head()
	if err != nil {
		err = errors.Wrap(err, "Cannot resolve HEAD")
		return
	}

	name := head.Name()

	// if we are on a branch or a tag
	// we can return the appropriate short version
	if name.IsBranch() || name.IsTag() {
		ref = name.Short()

		// else we need to resolve it
		// because we probably have a detached HEAD
	} else {
		ref = head.Hash().String()
	}
	return
}

func (gogit) GetRemotes(clonePath string, repoObject interface{}) (remoteMap map[string][]string, err error) {
	// get the repository
	r := repoObject.(*git.Repository)

	// get all the remotes for the repository
	remotes, err := r.Remotes()
	if err != nil {
		err = errors.Wrap(err, "Unable to get remotes")
		return
	}

	// make a map for remotes
	remoteMap = make(map[string][]string, len(remotes))
	for _, r := range remotes {
		cfg := r.Config()
		remoteMap[cfg.Name] = cfg.URLs
	}

	return
}

// originRemoteName is the name of the canonical remote
const originRemoteName = "origin"

func (gg gogit) GetCanonicalRemote(clonePath string, repoObject interface{}) (remoteName string, remoteURLs []string, err error) {
	// get a map of remotes
	remotes, err := gg.GetRemotes(clonePath, repoObject)
	if err != nil {
		err = errors.Wrap(err, "Unabel to get remotes")
		return
	}

	// if we don't have any remotes we're done
	if len(remotes) == 0 {
		return
	}

	// if the current branch has a remote, use it
	r := repoObject.(*git.Repository)
	remoteName, _ = gg.getCurrentBranchRemote(r)
	if remoteName != "" {
		remoteURLs = remotes[remoteName]
		return
	}

	// else if we have an 'origin' remote we use that
	if originRemote, originRemoteExists := remotes[originRemoteName]; originRemoteExists {
		remoteURLs = originRemote
		remoteName = originRemoteName
		return
	}

	// else randomly use the first remote that we have
	for rn, ru := range remotes {
		remoteURLs = ru
		remoteName = rn
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

func (gogit) SetRemoteURLs(clonePath string, repoObject interface{}, remoteName string, newURLs []string) (err error) {
	// get the repository
	r := repoObject.(*git.Repository)

	// get the desired remote
	remote, err := r.Remote(remoteName)
	if err != nil {
		err = errors.Wrapf(err, "Unable to find remote %s", remoteName)
		return
	}

	// fetch the current configuration
	cfg, err := r.Storer.Config()
	if err != nil {
		return
	}

	// update the urls
	if len(cfg.Remotes[remote.Config().Name].URLs) != len(newURLs) {
		return errors.New("Cannot set remoteURL: Length of old and new urls must be identical")
	}
	cfg.Remotes[remote.Config().Name].URLs = newURLs

	return
}

func (gogit) Clone(remoteURI, clonePath string, extraargs ...string) error {
	// doesn't support extra arguments
	if len(extraargs) > 0 {
		return ErrArgumentsUnsupported
	}

	// run a plain git clone but intercept all errors
	_, err := git.PlainClone(clonePath, false, &git.CloneOptions{URL: remoteURI, Progress: os.Stdout})
	if err != nil {
		err = ExitError{error: errors.Wrap(err, "Unable clone repository"), Code: 1}
	}

	return err
}

func (gogit) Fetch(clonePath string, cache interface{}) (err error) {
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
		err = remote.Fetch(&git.FetchOptions{Progress: os.Stdout})
		err = ignoreErrUpToDate(err)

		// fail on other errors
		if err != nil {
			err = errors.Wrapf(err, "Unable to fetch remote %s", remote.Config().Name)
			return
		}
	}

	return
}

func (gogit) Pull(clonePath string, cache interface{}) (err error) {
	// get the repository
	r := cache.(*git.Repository)

	// get the worktree
	w, err := r.Worktree()
	if err != nil {
		err = errors.Wrap(err, "Unable to find worktree")
		return
	}

	// do a git pull, and ignore error already up-to-date
	err = w.Pull(&git.PullOptions{Progress: os.Stdout})
	err = ignoreErrUpToDate(err)
	if err != nil {
		err = errors.Wrap(err, "Unable to pull")
	}

	return
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
