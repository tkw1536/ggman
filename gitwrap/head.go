package gitwrap

import (
	git "gopkg.in/src-d/go-git.v4"
)

// GetHeadRef returns the reference to the current head
func GetHeadRef(filepath string) (ref string, err error) {
	// open filepath or error out
	repo, err := git.PlainOpen(filepath)
	if err != nil {
		return
	}

	// get the current head
	head, err := repo.Head()
	if err != nil {
		return
	}

	name := head.Name()

	// if we are on a branch or a tag
	// we can return the appropriate short version
	if name.IsBranch() || name.IsTag() {
		ref = name.Short()

		// else we need to resolve it
		// because we probably have a detached HEAD
	} else {
		ref = head.Hash().String()
	}
	return
}
