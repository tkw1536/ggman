// Package git contains an implementation of git functionality.
//
// The implementation consists of the Git interface and the Plumbing interface.
//
// The Git interface (and it's default instance Default) provide a usable interface to Git Functionality.
// The Git interface will automatically choose between using a os.Exec() call to a native "git" wrapper, or using a pure golang git implementation.
// This should be used directly by callers.
//
// The Plumbing interface provides more direct control over which interface is used to interact with repositories.
// Calls to a Plumbing typically place assumptions on the caller and require some setup.
// For this reason, implementation of the Plumbing interface are not exported.
package git

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/tkw1536/goprogram/stream"
)

// Git represents a wrapper around a Plumbing instance.
// It is goroutine-safe and initialization free.
//
// As opposed to Plumbing, which poses certain requirements and assumptions on the caller, a Git does not.
// Using a Git can be as simple as:
//
//  err := git.Pull(stream.NewEnvIOStream(), "/home/user/Projects/github.com/hello/world")
//
type Git interface {
	// Plumbing returns the plumbing used by this git.
	Plumbing() Plumbing

	// IsRepository checks if the directory at localPath is the root of a git repository.
	IsRepository(localPath string) bool

	// IsRepositoryQuick efficiently checks if the directly at localPath contains a repository.
	// It is like IsRepository, except that it returns false more quickly than IsRepository.
	IsRepositoryQuick(localPath string) bool

	// Clone clones a remote repository from remoteURI to clonePath.
	// May attempt to read credentials from stream.Stdin.
	// Writes to stream.Stdout and stream.Stderr.
	//
	// remoteURI is the remote git uri to clone the repository from.
	// clonePath is the local path to clone the repository to.
	// extraargs are arguments as would be passed to a 'git clone' command.
	//
	// If there is already a repository at clonePath returns ErrCloneAlreadyExists.
	// If the underlying 'git' process exits abnormally, returns.
	// If extraargs is non-empty and extra arguments are not supported by this Wrapper, returns ErrArgumentsUnsupported.
	// May return other error types for other errors.
	Clone(stream stream.IOStream, remoteURI, clonePath string, extraargs ...string) error

	// GetHeadRef gets a resolved reference to head at the repository at clonePath.
	//
	// When getting the reference succeeded, returns err = nil.
	// If there is no repository at clonePath returns err = ErrNotARepository.
	// May return other error types for other errors.
	GetHeadRef(clonePath string) (ref string, err error)

	// Fetch fetches all remotes of the repository at clonePath.
	// May attempt to read credentials from stream.Stdin.
	// Writes to stream.Stdout and stream.Stderr.
	//
	// When fetching succeeded, returns nil.
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	Fetch(stream stream.IOStream, clonePath string) error

	// Pull fetches the repository at clonePath and merges in changes where appropriate.
	// May attempt to read credentials from stream.Stdin.
	// Writes to stream.Stdout and stream.Stderr.
	//
	// When pulling succeeded, returns nil.
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	Pull(stream stream.IOStream, clonePath string) error

	// GetRemote gets the url of the canonical remote at clonePath.
	// The semantics of 'canonical' are determined by the underlying git implementation.
	// Typically this function returns the url of the tracked remote of the currently checked out branch or the 'origin' remote.
	// If no remote exists, an empty url is returned.
	//
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	GetRemote(clonePath string) (url string, err error)

	// UpdateRemotes updates the urls of all remotes of the repository at clonePath.
	// updateFunc is a function that is called for each remote url to be updated.
	// It should return the new url corresponding to each old url.
	// If it returns a non-nil error, updating the current remote of the repository is instead aborted and error is returned.
	//
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	UpdateRemotes(clonePath string, updateFunc func(url, name string) (newURL string, err error)) error

	// GetBranches gets the names of all branches contained in the repository at clonePath.
	//
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	GetBranches(clonePath string) (branches []string, err error)

	// ContainsBranch checks if the repository at clonePath contains a branch with the provided name.
	//
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	ContainsBranch(clonePath, branch string) (exists bool, err error)

	// IsDirty checks if the repository at clonePath contains uncommitted changes.
	//
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	IsDirty(clonePath string) (dirty bool, err error)

	// IsSync checks if the repository at clonePath contains branches that are not yet synced with their upstream.
	//
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	IsSync(clonePath string) (sycned bool, err error)

	// GitPath returns the path to the git executable being used, if any.
	GitPath() string
}

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
	return impl.git
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

func (impl *dfltGitWrapper) Clone(stream stream.IOStream, remoteURI, clonePath string, extraargs ...string) error {
	impl.ensureInit()

	// check if the repository already exists
	if _, isRepo := impl.git.IsRepository(clonePath); isRepo {
		return ErrCloneAlreadyExists
	}

	// make the parent directory to clone the repository into
	if err := os.MkdirAll(filepath.Join(clonePath, ".."), os.ModePerm); err != nil {
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

func (impl *dfltGitWrapper) Fetch(stream stream.IOStream, clonePath string) error {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Fetch(stream, clonePath, repoObject)
}

func (impl *dfltGitWrapper) Pull(stream stream.IOStream, clonePath string) error {
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

func (impl *dfltGitWrapper) GetBranches(clonePath string) (branches []string, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return nil, ErrNotARepository
	}

	return impl.git.GetBranches(clonePath, repoObject)
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

func (impl *dfltGitWrapper) IsDirty(clonePath string) (dirty bool, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return false, ErrNotARepository
	}

	return impl.git.IsDirty(clonePath, repoObject)
}

func (impl *dfltGitWrapper) IsSync(clonePath string) (dirty bool, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return false, ErrNotARepository
	}

	return impl.git.IsSync(clonePath, repoObject)
}

func (impl *dfltGitWrapper) GitPath() string {
	impl.ensureInit()

	gitgit, isGitGit := impl.git.(*gitgit)
	if !isGitGit {
		return ""
	}
	return gitgit.gitPath
}
