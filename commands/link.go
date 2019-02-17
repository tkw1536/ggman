package commands

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/tkw1536/ggman/repos"
)

// LinkCommand is the entry point for the link command
func LinkCommand(parsed *GGArgs) (retval int, err string) {

	// 'link' takes not for
	if parsed.Pattern != "" {
		err = stringLinkNoFor
		retval = ErrorSpecificParseArgs
		return
	}

	// read the repo to link
	if len(parsed.Args) != 1 {
		err = stringLinkTakesOneArgument
		retval = ErrorSpecificParseArgs
		return
	}

	// get the root directory or panic
	root, e := getRootOrPanic()
	if e != nil {
		err = stringUnableParseRootDirectory
		retval = ErrorMissingConfig
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
		err = stringLinkDoesNotExist
		retval = ErrorCodeCustom
		return
	}

	// get the remote url
	remote, e := repos.NewRepoURI(r)
	if e != nil {
		err = e.Error()
		retval = ErrorCodeCustom
		return
	}

	// find the target path
	to := path.Join(append([]string{root}, remote.Components()...)...)
	parentTo := filepath.Dir(to)

	// if it's the same path, we throw an error
	if from == to {
		err = stringLinkSamePath
		retval = ErrorCodeCustom
		return
	}

	// make sure it doesn't exist
	if _, e := os.Stat(to); !os.IsNotExist(e) {
		err = stringLinkAlreadyExists
		retval = ErrorCodeCustom
		return
	}

	fmt.Printf("Linking %q -> %q\n", to, from)

	// make the parent folder
	if e := os.MkdirAll(parentTo, os.ModePerm); e != nil {
		err = e.Error()
		retval = ErrorCodeCustom
		return
	}

	// and make the symlink
	if e := os.Symlink(from, to); e != nil {
		err = e.Error()
		retval = ErrorCodeCustom
		return
	}

	// and be done
	return
}
