package commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"
)

// LinkCommand is the entry point for the link command
func LinkCommand(runtime *program.SubRuntime) (retval int, err string) {
	argv := runtime.Argv
	root := runtime.Root

	return linkRepository(filepath.Clean(argv[0]), root)
}

func linkRepository(from string, root string) (retval int, err string) {
	// make sure that the path is absolute
	// to avoid relative symlinks
	from, e := filepath.Abs(from)
	if e != nil {
		err = constants.StringLinkDoesNotExist
		retval = constants.ErrorCodeCustom
		return
	}

	// open the source repository and get the remotre
	r, e := git.Default.GetRemote(from)
	if e != nil {
		err = constants.StringLinkDoesNotExist
		retval = constants.ErrorCodeCustom
		return
	}

	// get the remote url
	remote := repos.ParseRepoURL(r)

	// find the target path
	to := path.Join(append([]string{root}, remote.Components()...)...)
	parentTo := filepath.Dir(to)

	// if it's the same path, we throw an error
	if from == to {
		err = constants.StringLinkSamePath
		retval = constants.ErrorCodeCustom
		return
	}

	// make sure it doesn't exist
	if _, e := os.Stat(to); !os.IsNotExist(e) {
		err = constants.StringLinkAlreadyExists
		retval = constants.ErrorCodeCustom
		return
	}

	fmt.Printf("Linking %q -> %q\n", to, from)

	// make the parent folder
	if e := os.MkdirAll(parentTo, os.ModePerm); e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

	// and make the symlink
	if e := os.Symlink(from, to); e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

	// and be done
	return
}
