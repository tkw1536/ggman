package commands

import (
	"fmt"
	"path"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/git"
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

	// get the canonical uri
	remoteURI := remote.CanonicalWith(lines)
	clonePath := path.Join(append([]string{root}, remote.Components()...)...)

	// and do the actual command
	fmt.Printf("Cloning %q into %q ...\n", remoteURI, clonePath)
	// catch special error types, and set the appropriate error messages and return values
	switch cloneErr := git.Default.Clone(remoteURI, clonePath, argv[1:]...); cloneErr {
	case nil:
	case git.ErrCloneAlreadyExists:
		err = constants.StringRepoAlreadyExists
		retval = constants.ErrorCodeCustom
	case git.ErrArgumentsUnsupported:
		err = constants.StringNoExternalGitnoArguments
		retval = constants.ErrorCodeCustom
	default:
		err = cloneErr.Error()
		retval = constants.ErrorCodeCustom
	}

	return
}
