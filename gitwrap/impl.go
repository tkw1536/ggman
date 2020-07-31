package gitwrap

import (
	"errors"
)

// GitImplementation is an interface that represents a working implementation of git.
type GitImplementation interface {

	// Init is used to initialize this git implementation
	// When initialization fails, for example due to missing dependencies, returns a non-nil error
	Init() error

	// IsRepository checks if the directory at localPath is the root of a git repository.
	// May assume that localPath exists and is a repository.
	//
	// This function returns a pair, a boolean isRepo that indicates if this object is a repository
	// and an optional repoObject value.
	// The repoObject value will only be taken into account when isRepo is true, and passed to other functions in this git implementation.
	// The semantics of the repoObject are determined by this GitImplementation and should not be used outside of it.
	// Note that the repoObject may be used for more than one subsequent call.
	//
	// This function surpresses all errors, and if something goes wrong assumed that isRepo is false.
	IsRepository(localPath string) (repoObject interface{}, isRepo bool)

	// IsRepositoryUnsafe efficiently checks if the directly at localPath contains a repository.
	// It is like IsRepository, except that it may return false positives, but no false negatives.
	// This function is optimized to be called a lot of times.
	IsRepositoryUnsafe(localPath string) bool

	// GetHeadRef returns a reference to the current head of the repository cloned at clonePath.
	// The string ref should contain a git REFLIKE, that is a branch, a tag or a commit id.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetHeadRef(clonePath string, repoObject interface{}) (ref string, err error)

	// GetRemotes returns the names and urls of the remotes of the repository cloned at clonePath.
	// If determining the remotes is not possible, and error is returned instead.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetRemotes(clonePath string, repoObject interface{}) (remotes map[string][]string, err error)

	// GetCanonicalRemote gets the name of the canonical remote of the reposity cloned at clonePath.
	// The GitImplementation is free to decided what the canonical remote is, but it is typically the remote of the currently checked out branch or the 'origin' remote.
	// If no remote exists, an empty name is returned.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	GetCanonicalRemote(clonePath string, repoObject interface{}) (remoteName string, remoteURLs []string, err error)

	// SetRemoteURLs set the remote 'remoteName' of the repository at clonePath to newURLs.
	// The remote remoteName must exist.
	// Furthermore newURLs must be of the same length as the old URLs.
	//
	// This function should only be called if IsRepository(clonePath) returns true.
	// The second parameter must be the returned value from IsRepository().
	SetRemoteURLs(clonePath string, repoObject interface{}, remoteName string, newURLs []string) (err error)

	// Clone tries to clone the repository at 'from' to the folder 'to'.
	// Output should be directed to os.Stdout and os.Stderr.
	//
	// remoteURI will be the uri of the remote repository.
	// clonePath will be the path to a local folder where the repository should be cloned to.
	// It is guaranteed to exist, and be empty.
	//
	// extraargs will be additional arguments, in the form of arguments of a 'git clone' command.
	// When this implementation does not support arguments, it should return ErrArgumentsUnsupported whenever arguments is a list of length > 0.
	//
	// code is the return code that a normal git command would return
	Clone(remoteURI, clonePath string, extraargs ...string) (code int, err error)

	// Fetch should fetch new objects and refs from all remotes of the repository cloned at clonePath.
	// Output should be directed to os.Stdout and os.Stderr.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	Fetch(clonePath string, cache interface{}) (err error)

	// Pull should fetch new objects and refs from all remotes of the repository cloned at clonePath.
	// It then merges them into the local branch wherever an upstream is set.
	// Output should be directed to os.Stdout and os.Stderr.
	//
	// This function will only be called if IsRepository(clonePath) returns true.
	// The second parameter passed will be the returned value from IsRepository().
	Pull(clonePath string, cache interface{}) (err error)
}

// ErrArgumentsUnsupported is an error that is returned when arguments are not supported by a GitImplementation.
var ErrArgumentsUnsupported = errors.New("GitImplementation does not support extra clone arguments")

// Implementation is *the* GitImpl
var Implementation = &GitWrap{}
