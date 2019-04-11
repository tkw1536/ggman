package gitwrap

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/src/constants"
	git "gopkg.in/src-d/go-git.v4"
)

// CloneRepository clones a repository into a given location
func CloneRepository(from string, to string) (retval int, err string) {
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

	// do the clone
	if _, e := git.PlainClone(to, false, &git.CloneOptions{URL: from, Progress: os.Stdout}); e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

	// and be done
	return
}
