package repos

import (
	"errors"

	"gopkg.in/src-d/go-git.v4"
)

// GetRemote gets the remote of a repository or nu
func GetRemote(filepath string) (uri string, err error) {
	// open filepath or error out
	repo, err := git.PlainOpen(filepath)
	if err != nil {
		return
	}

	// get the current head
	remotes, err := repo.Remotes()
	if err != nil {
		return
	}

	// get the canonical remote (commonly named 'origin'
	canon, err := getCanonicalRemote(remotes, "origin")
	if err != nil {
		return
	}

	// and return it's fetching url
	uri = canon.Config().URLs[0]
	return
}

// gets the canonical remote of a repository
func getCanonicalRemote(remotes []*git.Remote, name string) (canonical *git.Remote, err error) {
	// find the main remote
	remoteCount := len(remotes)

	// no remotes => error
	if remoteCount == 0 {
		err = errors.New("No remotes found")
		return
	}

	// if we have a remote named 'origin' use it
	for _, r := range remotes {
		if r.Config().Name == name {
			canonical = r
		}
	}

	// else pick the first remote
	if canonical == nil {
		canonical = remotes[0]
	}

	return
}
