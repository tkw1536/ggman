package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/repos"
)

// CompsCommand is the entry point for the compos command
func CompsCommand(parsed *GGArgs) (retval int, err string) {
	// 'comps' takes no for
	retval, err = parsed.EnsureNoFor()
	if retval != 0 {
		return
	}

	// we accept one argument
	if len(parsed.Args) != 1 {
		err = stringCompsTakesOneArgument
		retval = ErrorSpecificParseArgs
		return
	}

	// parse the repo uri
	uri, e := repos.NewRepoURI(parsed.Args[0])
	if e != nil {
		err = stringUnparsedRepoName
		retval = ErrorInvalidRepo
		return
	}

	// print each component on one line
	for _, comp := range uri.Components() {
		fmt.Println(comp)
	}

	// and finish
	return
}
