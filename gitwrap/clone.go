package gitwrap

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/tkw1536/ggman/constants"
	git "gopkg.in/src-d/go-git.v4"
)

// CloneRepository clones a repository into a given location
func CloneRepository(from string, to string, args ...string) (retval int, err string) {
	// tell the user what we are doing
	fmt.Printf("Cloning %q into %q ...\n", from, to)

	// if we can open a repository in 'to', it already exists
	if _, e := git.PlainOpen(to); e == nil {
		err = constants.StringRepoAlreadyExists
		retval = constants.ErrorCodeCustom
		return
	}

	// make the folder to clone into
	if e := os.MkdirAll(to, os.ModePerm); e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

	if e := cloneRepositorySmart(from, to, args...); e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

	// and be done
	return
}

// clones a repository either internally or externally
func cloneRepositorySmart(from string, to string, args ...string) error {
	if hasExternalGit() {
		return cloneRepositoryExternal(from, to, args...)
	} else if len(args) != 0 {
		return errors.New("External 'git' not found, can not pass any additional arguments to 'git clone'. ")
	} else {
		return cloneRepositoryInternal(from, to)
	}
}

// checks if we have an external git in $PATH
func hasExternalGit() bool {
	_, err := exec.LookPath("git")
	return err == nil
}

// clones a repository using the built-in git
func cloneRepositoryInternal(from string, to string) error {
	_, e := git.PlainClone(to, false, &git.CloneOptions{URL: from, Progress: os.Stdout})
	return e
}

// clones a repository using an external 'git' command
func cloneRepositoryExternal(from string, to string, args ...string) error {
	gargs := append([]string{"clone", from, to}, args...)
	cmd := exec.Command("git", gargs...)

	cmd.Dir = to

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
