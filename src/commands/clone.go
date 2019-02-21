package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"
	"gopkg.in/src-d/go-git.v4"
)

// CloneCommand is the entry point for the clone command
func CloneCommand(parsed *program.SubCommandArgs) (retval int, err string) {
	// 'clone' takes no for
	retval, err = parsed.EnsureNoFor()
	if retval != 0 {
		return
	}

	// read the repo to clone
	if len(parsed.Args) != 1 {
		err = constants.StringCloneTakesOneArgument
		retval = constants.ErrorSpecificParseArgs
		return
	}

	// parse the repo uri
	remote, e := repos.NewRepoURI(parsed.Args[0])
	if e != nil {
		err = constants.StringUnparsedRepoName
		retval = constants.ErrorInvalidRepo
		return
	}

	// get the canfile
	lines, e := program.GetCanonOrPanic()
	if e != nil {
		err = constants.StringInvalidCanfile
		retval = constants.ErrorMissingConfig
		return
	}

	// get the canonical url
	cloneURI := remote.CanonicalWith(lines)

	// get the root directory or panic
	root, e := program.GetRootOrPanic()
	if e != nil {
		err = constants.StringUnableParseRootDirectory
		retval = constants.ErrorMissingConfig
		return
	}

	// figure out where it goes
	targetPath := path.Join(append([]string{root}, remote.Components()...)...)

	// and finish
	return cloneRepository(cloneURI, targetPath)
}

func cloneRepository(from string, to string) (retval int, err string) {
	// tell the user what we are doing
	fmt.Printf("Cloning %q into %q ...\n", from, to)

	// if we can open a repository in 'to', it already exists
	if _, e := git.PlainOpen(to); e == nil {
		err = constants.StringRepoAlreadyExists
		retval = constants.ErrorCodeCustom
		return
	}

	// make the folder to clone into
	if e := os.MkdirAll(to, os.ModePerm); e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

	// do the clone
	if _, e := git.PlainClone(to, false, &git.CloneOptions{URL: from, Progress: os.Stdout}); e != nil {
		err = e.Error()
		retval = constants.ErrorCodeCustom
		return
	}

	// and be done
	return
}
