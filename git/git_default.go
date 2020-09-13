package git

import (
	"os"
	"sync"

	"github.com/tkw1536/ggman"
)

// NewGitFromPlumbing creates a new Git wrapping a specific Plumbing.
//
// When git is nil, attempts to automatically select a Plumbing automatically.
// The second parameter is only used when plumbing is nil; it should contain the value of 'PATH' environment variable.
//
// The implementation of this function relies on the underlying Plumbing (be it a default one or a caller provided one) to conform according to the specification.
// In particular, this function does not checks on the error values returned and passes them directly from the implementation to the caller
func NewGitFromPlumbing(plumbing Plumbing, path string) Git {
	return &dfltGitWrapper{git: plumbing, path: path}
}

type dfltGitWrapper struct {
	once sync.Once

	git  Plumbing
	path string // the path to lookup 'git' in, if needed.
}

func (impl *dfltGitWrapper) Plumbing() Plumbing {
	impl.ensureInit()
	return impl.Plumbing()
}

func (impl *dfltGitWrapper) ensureInit() {
	impl.once.Do(func() {
		if impl.git != nil {
			return
		}

		// first try to use a gitgit
		impl.git = &gitgit{gitPath: impl.path}
		if impl.git.Init() == nil {
			return
		}

		// then fallback to a gogit
		impl.git = &gogit{}
		if err := impl.git.Init(); err != nil {
			panic(err)
		}
	})
}

func (impl *dfltGitWrapper) IsRepository(localPath string) bool {
	impl.ensureInit()

	_, isRepo := impl.git.IsRepository(localPath)
	return isRepo
}

func (impl *dfltGitWrapper) IsRepositoryQuick(localPath string) bool {
	impl.ensureInit()

	if !impl.git.IsRepositoryUnsafe(localPath) { // IsRepositoryUnsafe may not return false negatives
		return false
	}

	return impl.IsRepository(localPath)
}

func (impl *dfltGitWrapper) Clone(stream ggman.IOStream, remoteURI, clonePath string, extraargs ...string) error {
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
	return impl.git.Clone(stream, remoteURI, clonePath, extraargs...)
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

func (impl *dfltGitWrapper) Fetch(stream ggman.IOStream, clonePath string) error {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Fetch(stream, clonePath, repoObject)
}

func (impl *dfltGitWrapper) Pull(stream ggman.IOStream, clonePath string) error {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Pull(stream, clonePath, repoObject)
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
	// then fix each url with the provided []env.CanLine
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

func (impl *dfltGitWrapper) ContainsBranch(clonePath, branch string) (exists bool, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return false, ErrNotARepository
	}

	return impl.git.ContainsBranch(clonePath, repoObject, branch)
}
