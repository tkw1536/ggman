// Package gitwrap contains a wrapper for git functionality
package gitwrap

import (
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/repos"
)

// ErrNotARepository is an error that is returned when the clonePath parameter is not a repository
var ErrNotARepository = errors.New("not a repository")

// GitWrap is a git implementation
type GitWrap struct {
	git   GitImplementation
	mutex sync.Mutex
}

func (impl *GitWrap) ensureInit() {
	impl.mutex.Lock()
	defer impl.mutex.Unlock()

	// we try to initialize a git implementation
	// if there already is a git, we return immedialty because we are done.
	// else we first try to initialize a gitGitImpl, and then fallback to goGitImpl.

	if impl.git != nil {
		return
	}

	impl.git = &gitgit{}
	if impl.git.Init() == nil {
		return
	}

	impl.git = &gogit{}
	if err := impl.git.Init(); err != nil {
		panic(err)
	}
}

// Clone clones a repository from remoteURI to clonePath.
func (impl *GitWrap) Clone(remoteURI, clonePath string, extraargs ...string) (retval int, err string) {
	impl.ensureInit()

	fmt.Printf("Cloning %q into %q ...\n", remoteURI, clonePath)

	// check if the repository already exists
	if _, isRepo := impl.git.IsRepository(clonePath); isRepo {
		err = constants.StringRepoAlreadyExists
		retval = constants.ErrorCodeCustom
		return
	}

	// make the directory to clone the repository into
	if e := os.MkdirAll(clonePath, os.ModePerm); e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

	// run the clone code and return
	retval, e := impl.git.Clone(remoteURI, clonePath, extraargs...)
	if e != nil {
		err = e.Error()
	}
	return
}

// GetHeadRef gets a resolved reference to head at the repository at clonePath
func (impl *GitWrap) GetHeadRef(clonePath string) (name string, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return "", ErrNotARepository
	}

	// and return the reference to the head
	return impl.git.GetHeadRef(clonePath, repoObject)
}

// Fetch fetches all remotes of the repository at clonePath
func (impl *GitWrap) Fetch(clonePath string) (err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Fetch(clonePath, repoObject)
}

// Pull fetches and merges the main repository at clonePath
func (impl *GitWrap) Pull(clonePath string) (err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Pull(clonePath, repoObject)
}

// GetRemote gets the url of the canonical remote at clonePath
func (impl *GitWrap) GetRemote(clonePath string) (uri string, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		err = ErrNotARepository
		return
	}

	// get all the uris
	_, uris, err := impl.git.GetCanonicalRemote(clonePath, repoObject)
	if err != nil || len(uris) == 0 {
		return
	}

	// use the first uri
	uri = uris[0]
	return
}

// FixRemotes updates all remotes of a repository with a given CanLine array
func (impl *GitWrap) FixRemotes(clonePath string, simulate bool, initialLogLine string, lines []repos.CanLine) (err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	// get all the remotes listed in the repository
	remotes, err := impl.git.GetRemotes(clonePath, repoObject)
	if err != nil {
		return err
	}

	fmt.Println(initialLogLine)

	// iterate over all the remotes, and their URLs
	// then fix each url with the provided []repos.CanLine
	// and store them again if we're not simulating

	for remote, urls := range remotes {
		canonURLs := make([]string, len(urls))
		for i, url := range urls {
			current, err := repos.NewRepoURI(url)
			if err != nil {
				continue
			}
			canonURLs[i] = current.CanonicalWith(lines)
			if canonURLs[i] != url {
				fmt.Printf("Updating %s: %s -> %s\n", remote, url, canonURLs[i])
			}
		}

		if simulate {
			continue
		}

		err := impl.git.SetRemoteURLs(clonePath, repoObject, remote, canonURLs)
		if err != nil {
			return err
		}
	}

	return
}
