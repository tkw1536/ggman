package git

//spellchecker:words errors
import "errors"

// ErrNotARepository is an error that is returned when the clonePath parameter is not a repository
var ErrNotARepository = errors.New("not a repository")

// ErrCloneAlreadyExists is an error that is returned when an operation can not be completed because a clone at the provided path already exists.
var ErrCloneAlreadyExists = errors.New("repository already exists")

// ExitError is an error that indicates the 'git' process exited abnormally
// This type is compatible with https://golang.org/pkg/errors/
type ExitError struct {
	// underlying error message
	error

	// Code that the git process exited with
	Code int
}

// Cause returns the cause of this error
func (err ExitError) Cause() error {
	return err.error
}
