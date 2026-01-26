package git

//spellchecker:words errors
import "errors"

// ErrNotARepository is an error that is returned when the clonePath parameter is not a repository.
var ErrNotARepository = errors.New("failed to resolve path: not a repository")

// ErrCloneAlreadyExists is an error that is returned when an operation can not be completed because a clone at the provided path already exists.
var ErrCloneAlreadyExists = errors.New("repository already exists")
