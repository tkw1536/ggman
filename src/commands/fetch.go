package commands

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"
)

// FetchCommand is the entry point for the fetch command
func FetchCommand(runtime *program.SubRuntime) (retval int, err string) {
	root := runtime.Root

	// find all the repos
	rs := repos.Repos(root, runtime.For)
	hasError := false

	// and fetch them
	for _, repo := range rs {
		fmt.Printf("Fetching %q\n", repo)
		if e := repos.FetchRepo(repo); e != nil {
			fmt.Fprintln(os.Stderr, e.Error())
			hasError = true
		}
	}

	// if we had an error, indicate that to the user
	if hasError {
		retval = constants.ErrorCodeCustom
	}

	// and finish
	return
}
