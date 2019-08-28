package commands

import (
	"fmt"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/gitwrap"
	"github.com/tkw1536/ggman/src/program"
	"github.com/tkw1536/ggman/src/repos"

	"github.com/skratchdot/open-golang/open"
)

// WebCommand is the entry point for the web command
func WebCommand(runtime *program.SubRuntime) (retval int, err string) {
	return webCommandInternal(runtime, true)
}

//URLCommand is the entry point for the url command
func URLCommand(runtime *program.SubRuntime) (retval int, err string) {
	return webCommandInternal(runtime, false)
}

func webCommandInternal(runtime *program.SubRuntime, openInstead bool) (retval int, err string) {
	root, relative := repos.Here(".", runtime.Root)
	if root == "" {
		return constants.ErrorInvalidRepo, constants.StringOutsideRepository
	}

	if relative == "." {
		relative = ""
	}

	// get the remote
	remote, e := gitwrap.GetRemote(root)
	if e != nil {
		return constants.ErrorInvalidRepo, constants.StringOutsideRepository
	}

	// parse it as a repo url
	uri, e := repos.NewRepoURI(remote)
	if e != nil {
		return constants.ErrorInvalidRepo, constants.StringUnparsedRepoName
	}

	// get the web url
	url := uri.Canonical("https://^/$")

	if runtime.Flag {
		ref, e := gitwrap.GetHeadRef(root)
		if e != nil {
			return constants.ErrorInvalidRepo, constants.StringUnparsedRepoName
		}
		// TODO: replace master with something useful
		url += "/tree/" + ref + "/" + relative
	}

	// open the url
	if openInstead {
		open.Start(url)

		// print the url
	} else {
		fmt.Println(url)
	}

	return
}
