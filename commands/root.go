package commands

import "github.com/tkw1536/ggman/constants"

// RootCommand is the entry point for the clone command
func RootCommand(parsed *GGArgs) (retval int, err string) {
	// 'root' takes no for
	retval, err = parsed.EnsureNoFor()
	if retval != 0 {
		return
	}

	// we accept no arguments
	if len(parsed.Args) != 0 {
		err = stringRootTakesNoArguments
		retval = constants.ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := getRootOrPanic()
	if e != nil {
		err = stringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
		return
	}

	// and echo out the root directory
	println(root)

	// and exit
	return
}
