package commands

import (
	"path"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/gitwrap"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"
)

// CloneCommand is the entry point for the clone command
func CloneCommand(runtime *program.SubRuntime) (retval int, err string) {
	argv := runtime.Argv
	lines := runtime.Canfile
	root := runtime.Root

	// parse the repo uri
	remote, e := repos.NewRepoURI(argv[0])
	if e != nil {
		err = constants.StringUnparsedRepoName
		retval = constants.ErrorInvalidRepo
		return
	}

	// get the canonical url
	cloneURI := remote.CanonicalWith(lines)

	// figure out where it goes
	targetPath := path.Join(append([]string{root}, remote.Components()...)...)

	// and finish
	return gitwrap.CloneRepository(cloneURI, targetPath, argv[1:]...)
}
