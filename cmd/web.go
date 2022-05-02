package cmd

import (
	"path/filepath"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/path"
	"github.com/tkw1536/goprogram/exit"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

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
//  --force-repo-here
// Instead of looking for a repository, always pretend there is one in the current directory.
// If the current directory is outside of the ggman root, this will cause an error.
// Use the path relative to the 'ggman root' as the URL to the repository.
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
//
// --clone
// When provided, instead of printing only the URL, print a 'git clone' command invocation that can be used by an anonymous user to clone the current repository.
// The clone url will always append '.git' to the web url, which may not work with every server.
// Only compatible with the '--branch' flag, but not '--tree', '--prefix', and 'BASE'.
//
// --reclone
// Like the --clone flag, but instead of using a normalized url, use the exact one found in the current repository.
var Web ggman.Command = &urlweb{
	isWebCommand: true,
}

// URL is the 'ggman url' command.
//
// The ggman url command behaves exactly like the ggman web command, except that instead of opening the URL in a webbrowser it prints it to standard output.
var URL ggman.Command = &urlweb{
	isWebCommand: false,
}

type urlweb struct {
	isWebCommand bool // if true, execute the web command; else the url command

	Positionals struct {
		Base string `positional-arg-name:"BASE" description:"if provided, replace the first component with the provided base url. alternatively you can use one of the predefined base URLs. use \"--list-bases\" to see a list of predefined base URLs"`
	} `positional-args:"true"`
	List bool `short:"l" long:"list-bases" descrioption:"print a list of all predefined base URLs"`

	ForceRepoHere bool `short:"f" long:"force-repo-here" description:"pretend there is a repository in the current path and use the path relative to the GGROOT directory as the remote url"`
	Branch        bool `short:"b" long:"branch" description:"if provided, include the HEAD reference in the resolved URL"`
	Tree          bool `short:"t" long:"tree" description:"if provided, additionally use the HEAD reference and relative path to the root of the git worktree"`
	BaseAsPrefix  bool `short:"p" long:"prefix" description:"treat the base argument as a prefix, instead of the hostname"`
	Clone         bool `short:"c" long:"clone" description:"if provided to the url command, print a \"git clone\" command that can be used to clone the current repository"`
	ReClone       bool `short:"r" long:"reclone" description:"like clone, but uses the current remote url as opposed to the https one"`
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

func (uw urlweb) Description() ggman.Description {
	var Name string
	if uw.isWebCommand {
		Name = "web"
	} else {
		Name = "url"
	}

	var Description string
	if uw.isWebCommand {
		Description = "open the URL of this repository in a web browser"
	} else {
		Description = "print the URL to this repository for opening a web browser"
	}

	return ggman.Description{
		Command:     Name,
		Description: Description,

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (uw urlweb) AfterParse() error {

	var cloneFlag string
	if uw.Clone {
		cloneFlag = "clone"
	} else {
		cloneFlag = "reclone"
	}

	isClone := uw.Clone || uw.ReClone

	if uw.isWebCommand && isClone {
		return errWebFlagUnsupported.WithMessageF(cloneFlag)
	}

	if isClone && uw.Tree {
		return errURLCloneAndUnsupported.WithMessageF(cloneFlag, "tree")
	}

	if isClone && uw.BaseAsPrefix {
		return errURLCloneAndUnsupported.WithMessageF(cloneFlag, "prefix")
	}

	return nil
}

var errOutsideRepository = exit.Error{
	ExitCode: env.ExitInvalidRepo,
	Message:  "Not inside a ggman-controlled repository",
}

var errWebNoRemote = exit.Error{
	ExitCode: env.ExitInvalidRepo,
	Message:  "Repository does not have a remote",
}

var errNoRelativeRepository = exit.Error{
	ExitCode: env.ExitInvalidRepo,
	Message:  "Unable to use '--relative': Not inside GGROOT",
}

var errWebFlagUnsupported = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "`ggman web` does not support the %s flag",
}

var errURLCloneAndUnsupported = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "`ggman url` does not support %s and %s arguments at the same time",
}

func (uw urlweb) Run(context ggman.Context) error {
	if uw.List {
		uw.listBases(context)
		return nil
	}

	// get the remote url of the current repository
	root, remote, relative, err := uw.getRemoteURL(context)
	if err != nil {
		return err
	}
	if remote == "" {
		return errWebNoRemote
	}

	var weburl string
	if !uw.ReClone {
		// parse it as a repo url
		url := env.ParseURL(remote)

		// set the base host
		base := "https://" + url.HostName

		// if we have a base argument, we need to use it
		if len(uw.Positionals.Base) > 0 {
			base = uw.Positionals.Base

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
		canspec := "^/$"
		if uw.Clone {
			canspec = "git clone ^/$.git"
		}
		weburl = url.Canonical(canspec)
	} else {
		weburl = "git clone " + remote
	}

	if root != "" && (uw.Tree || uw.Branch) {
		ref, err := context.Environment.Git.GetHeadRef(root)
		if err != nil {
			return errOutsideRepository
		}

		if !(uw.Clone || uw.ReClone) {
			weburl += "/tree/" + ref
			if uw.Tree {
				weburl += "/" + relative
			}
		} else {
			weburl += " --branch " + ref
		}
	}

	// print or open the url
	if uw.isWebCommand {
		browser.OpenURL(weburl)
	} else {
		context.Println(weburl)
	}

	return nil
}

func (uw urlweb) listBases(context ggman.Context) {
	bases := maps.Keys(WebBuiltInBases)
	slices.Sort(bases)

	for _, name := range bases {
		base := WebBuiltInBases[name]
		context.Printf("%s: %s\n", name, base.URL)
	}
}

// getRemoteURL gets the remote url of current repository in the context
func (uw urlweb) getRemoteURL(context ggman.Context) (root string, remote string, relative string, err error) {

	if uw.ForceRepoHere { // don't use a repository, instead fake one!
		return uw.getRemoteURLFake(context)
	}

	return uw.getRemoteURLReal(context)
}

func (uw urlweb) getRemoteURLReal(context ggman.Context) (root string, remote string, relative string, err error) {
	// find the repository at the current location
	root, relative, err = context.Environment.At(".")
	if err != nil {
		return "", "", "", err
	}

	if relative == "." {
		relative = ""
	}

	// get the remote
	remote, err = context.Environment.Git.GetRemote(root)
	if err != nil {
		return "", "", "", errOutsideRepository
	}

	return
}

func (uw urlweb) getRemoteURLFake(context ggman.Context) (root string, remote string, relative string, err error) {
	// get the absolute path to the current workdir
	workdir, err := context.Environment.Abs("")
	if err != nil {
		return "", "", "", err
	}

	// determine the relative path to the root directory
	relpath, err := filepath.Rel(context.Environment.Root, workdir)
	if err != nil || path.GoesUp(relpath) {
		return "", "", "", errNoRelativeRepository
	}

	// turn it into a fake url by prepending a protocol
	return "", "file://" + relpath, "", nil
}
