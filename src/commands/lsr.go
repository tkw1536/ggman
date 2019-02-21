package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"
)

// LSRCommand is the entry point for the lsr command
func LSRCommand(parsed *program.SubCommandArgs) (retval int, err string) {

	// read the --canonical flag
	shouldCanon, retval, err := parsed.ParseSingleFlag("--canonical")
	if retval != 0 {
		return
	}

	var lines []repos.CanLine
	var e error
	if shouldCanon {
		lines, e = program.GetCanonOrPanic()
		if e != nil {
			err = constants.StringInvalidCanfile
			retval = constants.ErrorMissingConfig
			return
		}
	}

	// get the root directory or panic
	root, e := program.GetRootOrPanic()
	if e != nil {
		err = constants.StringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
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
	if shouldCanon && len(rs) == 0 {
		retval = constants.ErrorCodeCustom
	}

	return
}
