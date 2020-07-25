package commands

import (
	"fmt"
	"sort"
	"strings"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/gitwrap"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/repos"

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

// WebBuiltInBases is a map of built-in bases for the url and web commands
var WebBuiltInBases = map[string]struct {
	URL         string
	IncludeHost bool
}{
	"travis":     {"https://travis-ci.com", false},
	"circle":     {"https://app.circleci.com/pipelines/github", false},
	"godoc":      {"https://pkg.go.dev/", true},
	"localgodoc": {"http://localhost:6060/pkg/", true},
}

// FmtWebBuiltInBaseNames returns a formatted string with all builtin bases
func FmtWebBuiltInBaseNames() string {
	// create a list of bases
	bases := make([]string, len(WebBuiltInBases))
	var count int
	for name := range WebBuiltInBases {
		bases[count] = "'" + name + "'"
		count++
	}

	// sort them, we don't care about stability
	sort.Slice(bases, func(i, j int) bool {
		return bases[i] < bases[j]
	})

	// and return
	return strings.Join(bases, ", ")
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

	// set the base host
	base := "https://" + uri.HostName
	if runtime.Argc > 0 {
		base = runtime.Argv[0]
	}

	// lookup in the builtins
	// we can do this safely because none of them start with https://
	if builtIn, ok := WebBuiltInBases[base]; ok {
		base = builtIn.URL
		if builtIn.IncludeHost {
			base += uri.HostName
		}
	}

	// set the hostname to the base
	uri.HostName = base

	// get the web url
	url := uri.Canonical("^/$")

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
