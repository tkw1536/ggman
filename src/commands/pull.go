package commands

import (
	"fmt"
	"os"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"
)

// PullCommand is the entry point for the fetch command
func PullCommand(parsed *program.SubCommandArgs) (retval int, err string) {
	retval, err = parsed.EnsureNoArguments()
	if retval != 0 {
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

	// and pull them
	for _, repo := range rs {
		fmt.Printf("Pulling %q\n", repo)
		if e := repos.PullRepo(repo); e != nil {
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
