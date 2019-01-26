package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/repos"
)

// LSRCommand is the entry point for the lsr command
func LSRCommand(parsed *GGArgs) (retval int, err string) {
	la := len(parsed.Args)
	// we accept no arguments
	if la > 1 {
		err = stringLSRArguments
		retval = ErrorSpecificParseArgs
		return
	} else if la == 1 && parsed.Args[0] != "--canonical" {
		err = stringLSRArguments
		retval = ErrorSpecificParseArgs
		return
	}

	// should we show the canonical url?
	shouldCanon := la == 1

	var lines []repos.CanLine
	var e error
	if shouldCanon {
		lines, e = getCanonOrPanic()
		if e != nil {
			err = stringInvalidCanfile
			retval = ErrorMissingConfig
			return
		}
	}

	// get the root directory or panic
	root, e := getRootOrPanic()
	if e != nil {
		err = stringUnableParseRootDirectory
		retval = ErrorMissingConfig
		return
	}

	// find all the repos
	rs := repos.Repos(root, parsed.Pattern)

	// and print them
	for _, repo := range rs {
		remote, err := repos.GetRemote(repo)
		if err == nil {
			if shouldCanon {
				printCanonOrError(lines, remote)
			} else {
				fmt.Println(remote)

			}
		}
	}

	// if we have --exit-code set and no results
	// we need to exit with an error code
	if la == 1 && len(rs) == 0 {
		retval = ErrorCodeCustom
	}

	return
}
