package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"
)

// CompsCommand is the entry point for the compos command
func CompsCommand(runtime *program.SubRuntime) (retval int, err string) {
	argv := runtime.Argv

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
