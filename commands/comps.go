package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"
)

// CompsCommand is the entry point for the compos command
func CompsCommand(runtime *program.SubRuntime) (retval int, err string) {
	argv := runtime.Argv

	// parse the repo uri
	url := repos.ParseRepoURL(argv[0])

	// print each component on one line
	for _, comp := range url.Components() {
		fmt.Println(comp)
	}

	// and finish
	return
}
