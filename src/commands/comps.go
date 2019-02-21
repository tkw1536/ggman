package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"
)

// CompsCommand is the entry point for the compos command
func CompsCommand(parsed *program.SubCommandArgs) (retval int, err string) {
	// 'comps' takes no for
	retval, err = parsed.EnsureNoFor()
	if retval != 0 {
		return
	}

	// comps takes exactly 1 argument
	_, argv, retval, err := parsed.EnsureArguments(1, 1)
	if retval != 0 {
		return
	}

	// parse the repo uri
	uri, e := repos.NewRepoURI(argv[0])
	if e != nil {
		err = constants.StringUnparsedRepoName
		retval = constants.ErrorInvalidRepo
		return
	}

	// print each component on one line
	for _, comp := range uri.Components() {
		fmt.Println(comp)
	}

	// and finish
	return
}
