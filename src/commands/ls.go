package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"
)

// LSCommand is the entry point for the ls command
func LSCommand(parsed *program.SubCommandArgs) (retval int, err string) {

	// read the --exit-code flag
	exitCodeFlag, retval, err := parsed.ParseSingleFlag("--exit-code")
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
