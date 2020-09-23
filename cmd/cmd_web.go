package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"

	"github.com/pkg/browser"
)

// Web is the 'ggman web' command
var Web program.Command = web{&urlweb{openInstead: true}}

type web struct{ *urlweb }

func (web) Name() string {
	return "web"
}

// URL is the 'ggman url' command
var URL program.Command = url{&urlweb{openInstead: false}}

type url struct{ *urlweb }

func (url) Name() string {
	return "url"
}

type urlweb struct {
	openInstead bool

	Branch bool
	Tree   bool
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

var urlwebBuiltinBaseNames string = FmtWebBuiltInBaseNames()

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

var stringWebBaseUsage = "If provided, replace the first component with the provided base url. Alternatively you can use one of the predefined urls %s. "

func (uw *urlweb) Options(flagset *pflag.FlagSet) program.Options {
	flagset.BoolVarP(&uw.Tree, "tree", "t", uw.Tree, "If provided, additionally use the HEAD reference and relative path to the root of the git worktree. ")
	flagset.BoolVarP(&uw.Branch, "branch", "b", uw.Branch, "If provided, include the HEAD reference in the resolved URL")
	return program.Options{
		MinArgs: 0,
		MaxArgs: 1,
		Metavar: "BASE",

		UsageDescription: fmt.Sprintf(stringWebBaseUsage, urlwebBuiltinBaseNames),

		Environment: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (urlweb) AfterParse() error {
	return nil
}

var errOutsideRepository = ggman.Error{
	ExitCode: ggman.ExitInvalidRepo,
	Message:  "Not inside a ggman-controlled repository. ",
}

func (uw urlweb) Run(context program.Context) error {
	root, relative, err := context.At(".")
	if err != nil {
		return err
	}

	if relative == "." {
		relative = ""
	}

	// get the remote
	remote, e := context.Git.GetRemote(root)
	if e != nil {
		return errOutsideRepository
	}

	// parse it as a repo url
	url := context.ParseURL(remote)

	// set the base host
	base := "https://" + url.HostName
	if len(context.Args) > 0 {
		base = context.Args[0]
	}

	// lookup in the builtins
	// we can do this safely because none of them start with https://
	if builtIn, ok := WebBuiltInBases[base]; ok {
		base = builtIn.URL
		if builtIn.IncludeHost {
			base += url.HostName
		}
	}

	// set the hostname to the base
	url.HostName = base

	// get the web url
	weburl := url.Canonical("^/$")

	if uw.Tree || uw.Branch {
		ref, e := context.Git.GetHeadRef(root)
		if e != nil {
			return errOutsideRepository
		}

		// TODO: do we want to replace the HEAD branch with something more useful?
		weburl += "/tree/" + ref
		if uw.Tree {
			weburl += "/" + relative
		}
	}

	// print or open the url
	if uw.openInstead {
		browser.OpenURL(weburl)
	} else {
		context.Println(weburl)
	}

	return nil
}
