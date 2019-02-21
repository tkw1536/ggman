package commands

import (
	"path"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/repos"
)

// WhereCommand is the entry point for the where command
func WhereCommand(parsed *GGArgs) (retval int, err string) {
	// 'where' takes no for
	retval, err = parsed.EnsureNoFor()
	if retval != 0 {
		return
	}

	// we accept no arguments
	if len(parsed.Args) != 1 {
		err = stringWhereTakesOneArgument
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

	// parse the repository in questions
	r, e := repos.NewRepoURI(parsed.Args[0])
	if e != nil {
		err = stringUnparsedRepoName
		retval = constants.ErrorInvalidRepo
		return
	}

	// and get it's components
	components := r.Components()

	// and join it into a path
	location := path.Join(append([]string{root}, components...)...)
	println(location)

	return
}
