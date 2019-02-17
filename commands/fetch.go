package commands

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/repos"
)

// FetchCommand is the entry point for the fetch command
func FetchCommand(parsed *GGArgs) (retval int, err string) {
	la := len(parsed.Args)
	// we accept no arguments
	if la != 0 {
		err = stringFetchTakesNoArguments
		retval = ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := getRootOrPanic()
	if e != nil {
		err = stringUnableParseRootDirectory
		retval = ErrorMissingConfig
		return
	}

	// find all the repos
	rs := repos.Repos(root, parsed.Pattern)
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
		retval = ErrorCodeCustom
	}

	// and finish
	return
}
