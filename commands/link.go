package commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/repos"
)

// LinkCommand is the entry point for the link command
func LinkCommand(parsed *GGArgs) (retval int, err string) {
	// 'link' takes no for
	retval, err = parsed.EnsureNoFor()
	if retval != 0 {
		return
	}

	// read the repo to link
	if len(parsed.Args) != 1 {
		err = constants.StringLinkTakesOneArgument
		retval = constants.ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := getRootOrPanic()
	if e != nil {
		err = constants.StringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
		return
	}

	return linkRepository(filepath.Clean(parsed.Args[0]), root)

	// figure out where it goes
	//targetPath := path.Join(append([]string{root}, remote.Components()...)...)
}

func linkRepository(from string, root string) (retval int, err string) {

	// open the source repository and get the remotre
	r, e := repos.GetRemote(from)
	if e != nil {
		err = constants.StringLinkDoesNotExist
		retval = constants.ErrorCodeCustom
		return
	}

	// get the remote url
	remote, e := repos.NewRepoURI(r)
	if e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

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
