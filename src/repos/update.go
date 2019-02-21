package repos

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
)

// FetchRepo fetches a repository
func FetchRepo(root string) (err error) {
	// open the repository
	r, err := git.PlainOpen(root)
	if err != nil {
		return
	}

	// list all of the remotes
	remotes, err := r.Remotes()
	if err != nil {
		return
	}

	// fetch all of the remotes for this repository
	for _, remote := range remotes {
		localError := checkUpdateError(remote.Fetch(&git.FetchOptions{Progress: os.Stdout}))
		if localError != nil && err == nil {
			err = localError
		}
	}

	return
}

// PullRepo pulls a repository
func PullRepo(root string) (err error) {
	// open the repository
	r, err := git.PlainOpen(root)
	if err != nil {
		return
	}

	// list all of the remotes
	w, err := r.Worktree()
	if err != nil {
		return
	}

	// and run a pull
	return checkUpdateError(w.Pull(&git.PullOptions{Progress: os.Stdout}))
}

func checkUpdateError(errIn error) (err error) {
	err = errIn
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println(errIn.Error())
		err = nil
	}
	return
}
