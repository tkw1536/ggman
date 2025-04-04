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

//spellchecker:words path filepath sync github ggman internal dirs pkglib stream
import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/tkw1536/ggman/internal/dirs"
	"github.com/tkw1536/pkglib/stream"
)

//spellchecker:words gitgit gogit

// Git represents a wrapper around a Plumbing instance.
// It is goroutine-safe and initialization free.
//
// As opposed to Plumbing, which poses certain requirements and assumptions on the caller, a Git does not.
// Using a Git can be as simple as:
//
//	err := git.Pull(stream.NewEnvIOStream(), "/home/user/Projects/github.com/hello/world")
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
	// extraArgs are arguments as would be passed to a 'git clone' command.
	//
	// If there is already a repository at clonePath returns ErrCloneAlreadyExists.
	// If the underlying 'git' process exits abnormally, returns.
	// If extraArgs is non-empty and extra arguments are not supported by this Wrapper, returns ErrArgumentsUnsupported.
	// May return other error types for other errors.
	Clone(stream stream.IOStream, remoteURI, clonePath string, extraArgs ...string) error

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

	// GetRemote gets the url of the remote at clonePath.
	// Name is the name of the remote.
	// When empty, picks the primary remote, as determined by the underlying git implementation.
	// Typically this function returns the url of the tracked remote of the currently checked out branch or the 'origin' remote.
	// If no remote exists, an empty url is returned.
	//
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	GetRemote(clonePath string, name string) (url string, err error)

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
	IsSync(clonePath string) (synced bool, err error)

	// GitPath returns the path to the git executable being used, if any.
	GitPath() string
}

// In particular, this function does not checks on the error values returned and passes them directly from the implementation to the caller.
func NewGitFromPlumbing(plumbing Plumbing, path string) Git {
	return &defaultGitWrapper{git: plumbing, path: path}
}

type defaultGitWrapper struct {
	once sync.Once

	git  Plumbing
	path string // the path to lookup 'git' in, if needed.
}

func (impl *defaultGitWrapper) Plumbing() Plumbing {
	impl.ensureInit()
	return impl.git
}

func (impl *defaultGitWrapper) ensureInit() {
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

func (impl *defaultGitWrapper) IsRepository(localPath string) bool {
	impl.ensureInit()

	_, isRepo := impl.git.IsRepository(localPath)
	return isRepo
}

func (impl *defaultGitWrapper) IsRepositoryQuick(localPath string) bool {
	impl.ensureInit()

	if !impl.git.IsRepositoryUnsafe(localPath) { // IsRepositoryUnsafe may not return false negatives
		return false
	}

	return impl.IsRepository(localPath)
}

func (impl *defaultGitWrapper) Clone(stream stream.IOStream, remoteURI, clonePath string, extraArgs ...string) error {
	impl.ensureInit()

	// check if the repository already exists
	if _, isRepo := impl.git.IsRepository(clonePath); isRepo {
		return ErrCloneAlreadyExists
	}

	// make the parent directory to clone the repository into
	if err := os.MkdirAll(filepath.Join(clonePath, ".."), dirs.NewModBits); err != nil {
		return err
	}

	// run the clone code and return
	return impl.git.Clone(stream, remoteURI, clonePath, extraArgs...)
}

func (impl *defaultGitWrapper) GetHeadRef(clonePath string) (ref string, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return "", ErrNotARepository
	}

	// and return the reference to the head
	return impl.git.GetHeadRef(clonePath, repoObject)
}

func (impl *defaultGitWrapper) Fetch(stream stream.IOStream, clonePath string) error {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Fetch(stream, clonePath, repoObject)
}

func (impl *defaultGitWrapper) Pull(stream stream.IOStream, clonePath string) error {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return ErrNotARepository
	}

	return impl.git.Pull(stream, clonePath, repoObject)
}

var errNoRemoteURL = errors.New("no remote URL found")

func (impl *defaultGitWrapper) GetRemote(clonePath string, name string) (uri string, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		err = ErrNotARepository
		return
	}

	// if no name is provided, use the canonical remote!
	if name == "" {
		_, uris, err := impl.git.GetCanonicalRemote(clonePath, repoObject)
		if err != nil {
			return "", err
		}
		if len(uris) == 0 {
			return "", errNoRemoteURL
		}

		// use the first uri
		return uris[0], nil
	}

	// get all the remotes
	remotes, err := impl.git.GetRemotes(clonePath, repoObject)
	if err != nil {
		return "", err
	}

	// pick the canonical one!
	urls, ok := remotes[name]
	if !ok {
		return "", fmt.Errorf("remote %q not found", name)
	}
	if len(urls) == 0 {
		return "", fmt.Errorf("remote %q: %w", name, errNoRemoteURL)
	}

	return urls[0], nil
}

func (impl *defaultGitWrapper) UpdateRemotes(clonePath string, updateFunc func(url, name string) (string, error)) (err error) {
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

func (impl *defaultGitWrapper) GetBranches(clonePath string) (branches []string, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return nil, ErrNotARepository
	}

	return impl.git.GetBranches(clonePath, repoObject)
}

func (impl *defaultGitWrapper) ContainsBranch(clonePath, branch string) (exists bool, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return false, ErrNotARepository
	}

	return impl.git.ContainsBranch(clonePath, repoObject, branch)
}

func (impl *defaultGitWrapper) IsDirty(clonePath string) (dirty bool, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return false, ErrNotARepository
	}

	return impl.git.IsDirty(clonePath, repoObject)
}

func (impl *defaultGitWrapper) IsSync(clonePath string) (dirty bool, err error) {
	impl.ensureInit()

	// check that the given folder is actually a repository
	repoObject, isRepo := impl.git.IsRepository(clonePath)
	if !isRepo {
		return false, ErrNotARepository
	}

	return impl.git.IsSync(clonePath, repoObject)
}

func (impl *defaultGitWrapper) GitPath() string {
	impl.ensureInit()

	gitgit, isGitGit := impl.git.(*gitgit)
	if !isGitGit {
		return ""
	}
	return gitgit.gitPath
}
