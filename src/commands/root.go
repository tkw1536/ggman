package commands

import (
	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
)

// RootCommand is the entry point for the clone command
func RootCommand(parsed *program.GGArgs) (retval int, err string) {
	// 'root' takes no for
	retval, err = parsed.EnsureNoFor()
	if retval != 0 {
		return
	}

	// we accept no arguments
	if len(parsed.Args) != 0 {
		err = constants.StringRootTakesNoArguments
		retval = constants.ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := program.GetRootOrPanic()
	if e != nil {
		err = constants.StringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
		return
	}

	// and echo out the root directory
	println(root)

	// and exit
	return
}
