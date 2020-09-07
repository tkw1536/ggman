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
	"github.com/tkw1536/ggman"
)

// Git represents a wrapper around a Plumbing instance.
// It is goroutine-safe and initialization free.
//
// As opposed to Plumbing, which poses certain requirements and assumptions on the caller, a Git does not.
// Using a Git can be as simple as:
//
//  err := git.Pull(ggman.NewEnvIOStream(), "/home/user/Projects/github.com/hello/world")
//
type Git interface {
	// Plumbing returns the plumbing used by this git.
	Plumbing() Plumbing

	// IsRepository checks if the directory at localPath is the root of a git repository.
	IsRepository(localPath string) bool

	// IsRepositoryQuick efficiently checks if the directly at localPath contains a repository.
	// It is like IsRepository, except that it more quickly returns false than IsRepository.
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
	Clone(stream ggman.IOStream, remoteURI, clonePath string, extraargs ...string) error

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
	Fetch(stream ggman.IOStream, clonePath string) error

	// Pull fetches the repository at clonePath and merges in changes where appropriate.
	// May attempt to read credentials from stream.Stdin.
	// Writes to stream.Stdout and stream.Stderr.
	//
	// When pulling succeeded, returns nil.
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	Pull(stream ggman.IOStream, clonePath string) error

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

	// ContainsBranch checks if the repository at clonePath contains a branch with the provided branch.
	//
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	ContainsBranch(clonePath, branch string) (exists bool, err error)
}
