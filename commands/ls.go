package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"
)

// LSCommand is the entry point for the ls command
func LSCommand(runtime *program.SubRuntime) (retval int, err string) {
	exitCodeFlag := runtime.Flag
	root := runtime.Root

	// find all the repos
	repos := repos.Repos(root, runtime.For)

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
