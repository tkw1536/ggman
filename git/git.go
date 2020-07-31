// Package git contains a wrapper for git functionality
package git

// Git represents a wrapper around a Plumbing instance.
// As opposed to Plumbing, which poses certain requirements and assumptions on the caller, a Git does not.
// A Git is intended to be called by seperate package.
type Git interface {
	// Plumbing returns the plumbing used by this git.
	Plumbing() Plumbing

	// Clone clones a remote repository from remoteURI to clonePath.
	// Writes to os.Stdout and os.Stderr.
	//
	// remoteURI is the remote git uri to clone the repository from.
	// clonePath is the local path to clone the repository to.
	// extraargs are arguments as would be passed to a 'git clone' command.
	//
	// If there is already a repository at clonePath returns ErrCloneAlreadyExists.
	// If the underlying 'git' process exits abnormally, returns.
	// If extraargs is non-empty and extra arguments are not supported by this Wrapper, returns ErrArgumentsUnsupported.
	// May return other error types for other errors.
	Clone(remoteURI, clonePath string, extraargs ...string) error

	// GetHeadRef gets a resolved reference to head at the repository at clonePath.
	//
	// When getting the reference succeeded, returns err = nil.
	// If there is no repository at clonePath returns err = ErrNotARepository.
	// May return other error types for other errors.
	GetHeadRef(clonePath string) (ref string, err error)

	// Fetch fetches all remotes of the repository at clonePath.
	// May attempt to read credentials from os.Stdin.
	// Writes to os.Stdout and os.Stderr.
	//
	// When fetching succeeded, returns nil.
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	Fetch(clonePath string) error

	// Pull fetches the repository at clonePath and merges in changes where appropriate.
	// May attempt to read credentials from os.Stdin.
	// Writes to os.Stdout and os.Stderr.
	//
	// When pulling succeeded, returns nil.
	// If there is no repository at clonePath returns ErrNotARepository.
	// May return other error types for other errors.
	Pull(clonePath string) error

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
	UpdateRemotes(clonePath string, updateFunc func(url, remoteName string) (newRemote string, err error)) error
}

// Default is the default git
var Default = NewGitFromPlumbing(nil)
