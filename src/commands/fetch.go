package commands

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"
)

// FetchCommand is the entry point for the fetch command
func FetchCommand(parsed *program.SubCommandArgs) (retval int, err string) {
	la := len(parsed.Args)
	// we accept no arguments
	if la != 0 {
		err = constants.StringFetchTakesNoArguments
		retval = constants.ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := program.GetRootOrPanic()
	if e != nil {
		err = constants.StringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
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
		retval = constants.ErrorCodeCustom
	}

	// and finish
	return
}
