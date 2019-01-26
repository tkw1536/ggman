package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/repos"
)

// LSCommand is the entry point for the ls command
func LSCommand(parsed *GGArgs) (retval int, err string) {
	la := len(parsed.Args)
	// we accept no arguments
	if la > 1 {
		err = stringLSArguments
		retval = ErrorSpecificParseArgs
		return
	} else if la == 1 && parsed.Args[0] != "--exit-code" {
		err = stringLSArguments
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
	repos := repos.Repos(root, parsed.Pattern)

	// and print them
	for _, repo := range repos {
		fmt.Println(repo)
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if la == 1 && len(repos) == 0 {
		retval = ErrorCodeCustom
	}

	return
}
