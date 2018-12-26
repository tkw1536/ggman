package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/repos"
)

// LSCommand is the entry point for the ls command
func LSCommand(parsed *GGArgs) (retval int, err string) {
	// we accept no arguments
	if len(parsed.Args) != 0 {
		err = stringLSTakesNoArguments
		retval = ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := getRootOrPanic()
	if e != nil {
		err = stringUnableParseRootDirectory
		retval = ErrorNoRoot
		return
	}

	// find all the repos
	repos := repos.Repos(root, parsed.Pattern)

	// and print them
	for _, repo := range repos {
		fmt.Println(repo)
	}

	return
}
