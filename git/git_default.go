package git

import (
	"os"
	"sync"
)

// NewGitFromPlumbing creates a new Git wrapping a specific Plumbing.
//
// When git is nil, attempts to automatically select a Plumbing automatically.
//
// The implementation of this function relies on the underlying Plumbing (be it a default one or a caller provided one) to conform according to the specification.
// In particular, this function does not checks on the error values returned and passes them directly from the implementation to the caller
func NewGitFromPlumbing(plumbing Plumbing) Git {
	return &dfltGitWrapper{git: plumbing}
}

type dfltGitWrapper struct {
	m   sync.Mutex
	git Plumbing
}

func (impl *dfltGitWrapper) Plumbing() Plumbing {
	impl.ensureInit()
	return impl.Plumbing()
}

func (impl *dfltGitWrapper) ensureInit() {
	impl.m.Lock()
	defer impl.m.Unlock()

	// We try to initialize a dflt
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

func (impl *dfltGitWrapper) Clone(remoteURI, clonePath string, extraargs ...string) error {
	impl.ensureInit()

	// check if the repository already exists
	if _, isRepo := impl.git.IsRepository(clonePath); isRepo {
		return ErrCloneAlreadyExists
	}

	// make the directory to clone the repository into
	if err := os.MkdirAll(clonePath, os.ModePerm); err != nil {
		return err
	}

	// run the clone code and return
	return impl.git.Clone(remoteURI, clonePath, extraargs...)
}

func (impl *dfltGitWrapper) GetHeadRef(clonePath string) (ref string, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return "", ErrNotARepository
	}

	// and return the reference to the head
	return impl.git.GetHeadRef(clonePath, repoObject)
}

func (impl *dfltGitWrapper) Fetch(clonePath string) error {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Fetch(clonePath, repoObject)
}

func (impl *dfltGitWrapper) Pull(clonePath string) error {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Pull(clonePath, repoObject)
}

func (impl *dfltGitWrapper) GetRemote(clonePath string) (uri string, err error) {
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

func (impl *dfltGitWrapper) UpdateRemotes(clonePath string, updateFunc func(url, name string) (string, error)) (err error) {
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

	// iterate over all the remotes, and their URLs
	// then fix each url with the provided []repos.CanLine
	// and store them again if we're not simulating

	for remoteName, urls := range remotes {
		canonURLs := make([]string, len(urls))
		for i, url := range urls {
			if canonURLs[i], err = updateFunc(url, remoteName); err != nil {
				return err
			}
		}

		err := impl.git.SetRemoteURLs(clonePath, repoObject, remoteName, canonURLs)
		if err != nil {
			return err
		}
	}

	return
}
