package commands

import (
	"path"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"
)

// WhereCommand is the entry point for the where command
func WhereCommand(parsed *program.SubCommandArgs) (retval int, err string) {
	// 'where' takes no for
	retval, err = parsed.EnsureNoFor()
	if retval != 0 {
		return
	}

	// where takes exactly 1 argument
	_, argv, retval, err := parsed.EnsureArguments(1, 1)
	if retval != 0 {
		return
	}

	// get the root directory or panic
	root, e := program.GetRootOrPanic()
	if e != nil {
		err = constants.StringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
		return
	}

	// parse the repository in questions
	r, e := repos.NewRepoURI(argv[0])
	if e != nil {
		err = constants.StringUnparsedRepoName
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
