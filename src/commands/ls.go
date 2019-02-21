package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/src/args"
	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/repos"
)

// LSCommand is the entry point for the ls command
func LSCommand(parsed *args.GGArgs) (retval int, err string) {

	// read the --exit-code flag
	exitCodeFlag, ie := parsed.ParseSingleFlag("--exit-code")
	if ie {
		err = constants.StringLSArguments
		retval = constants.ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := args.GetRootOrPanic()
	if e != nil {
		err = constants.StringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
		return
	}

	// find all the repos
	repos := repos.Repos(root, parsed.Pattern)

	// and print them
	for _, repo := range repos {
		fmt.Println(repo)
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if exitCodeFlag && len(repos) == 0 {
		retval = constants.ErrorCodeCustom
	}

	return
}
