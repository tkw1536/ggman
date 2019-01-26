package commands

// RootCommand is the entry point for the clone command
func RootCommand(parsed *GGArgs) (retval int, err string) {
	// 'root' takes not for
	if parsed.Pattern != "" {
		err = stringRootNoFor
		retval = ErrorSpecificParseArgs
		return
	}

	// we accept no arguments
	if len(parsed.Args) != 0 {
		err = stringRootTakesNoArguments
		retval = ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := getRootOrPanic()
	if e != nil {
		err = stringUnableParseRootDirectory
		retval = ErrorNoRoot
		return
	}

	// and echo out the root directory
	println(root)

	// and exit
	return
}
