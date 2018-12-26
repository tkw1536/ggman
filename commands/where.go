package commands

import (
	"path"

	"github.com/tkw1536/ggman/utils"
)

// WhereCommand is the entry point for the where command
func WhereCommand(parsed *GGArgs) (retval int, err string) {
	// we have an error
	if parsed.Pattern != "" {
		err = stringWhereNoFor
		retval = ErrorSpecificParseArgs
		return
	}

	// we accept no arguments
	if len(parsed.Args) != 1 {
		err = stringWhereTakesOneArgument
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

	// get components of the argument
	components, e := utils.Components(parsed.Args[0])
	if e != nil {
		err = stringUnparsedRepoName
		retval = ErrorInvalidRepo
		return
	}

	// and join it into a path
	location := path.Join(append([]string{root}, components...)...)
	println(location)

	return
}
