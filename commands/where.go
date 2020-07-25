package commands

import (
	"path"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"
)

// WhereCommand is the entry point for the where command
func WhereCommand(runtime *program.SubRuntime) (retval int, err string) {
	argv := runtime.Argv
	root := runtime.Root

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
