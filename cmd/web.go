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

// Web is the 'ggman web' command.
//
// It attempts to open the url of the current repository in a webrowser.
// The browser being used is determined by the underlying operating system.
//
// To determine the url it uses the CANSPEC `https://^/$`, which may not work with all git hosts.
// For instance when the url of the repository git@github.com:tkw1536/ggman.git, this will open the url https://github.com/tkw1536/ggman in a browser.
//
//  --tree
// This optional argument appends the string `/tree/$BRANCH/$PATH` to the url being opened, where $BRANCH is currently checked out branch and $PATH is the relative path from the repository root to the current folder.
// On common Git Hosts, such as GitHub and GitLab, this shows a page of the current folder on the current branch.
//  --branch
// This argument works like '--tree', except that it does not append the local path to the url.
//  BASE
// An optional argument of the form 'BASE'.
// If it is provided, the first component of the url is replace with the given base.
// For instance, using the base 'https://pkg.go.dev' would open the current repository on the golang documentation homepage.
// In addition to using a custom BASE, the following pre-defined bases 'travis' (TravisCI), 'circle' (CirclCI), 'godoc' (GoDoc) and 'localgodoc' (GoDoc when run on the local machine) can be used.
//  --prefix
// When provided, instead of replacing the hostname with the base, prefix it with the base instead.
// This flag is ignored when no base is provided, or a built-in base is used.
var Web program.Command = &web{}

type web struct{ urlweb }

func (w web) Run(context program.Context) error {
	w.urlweb.openInstead = true
	return w.urlweb.Run(context)
}

func (web) Name() string {
	return "web"
}

// URL is the 'ggman url' command.
//
// The ggman url command behaves exactly like the ggman web command, except that instead of opening the URL in a webbrowser it prints it to standard output.
var URL program.Command = &url{}

func (u url) Run(context program.Context) error {
	u.urlweb.openInstead = false
	return u.urlweb.Run(context)
}

type url struct{ urlweb }

func (url) Name() string {
	return "url"
}

type urlweb struct {
	openInstead bool

	Branch       bool
	Tree         bool
	BaseAsPrefix bool
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
	flagset.BoolVarP(&uw.Branch, "branch", "b", uw.Branch, "If provided, include the HEAD reference in the resolved URL. ")
	flagset.BoolVarP(&uw.BaseAsPrefix, "prefix", "p", uw.BaseAsPrefix, "Treat the base argument as a prefix, instead of the hostname. ")
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
	url := env.ParseURL(remote)

	// set the base host
	base := "https://" + url.HostName

	// if we have a base argument, we need to use it
	if len(context.Args) > 0 {
		base = context.Args[0]

		// lookup in builtins
		if builtIn, ok := WebBuiltInBases[base]; ok {
			base = builtIn.URL
			uw.BaseAsPrefix = builtIn.IncludeHost
		}

		// if we want to use the base as a prefix, add back the hostname
		if uw.BaseAsPrefix {
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
		browser.OpenURL(weburl) // TODO: This breaks test isolation and is very hard to test.
	} else {
		context.Println(weburl)
	}

	return nil
}
