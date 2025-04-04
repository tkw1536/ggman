package cmd

//spellchecker:words path filepath slices github browser ggman internal goprogram exit
import (
	"path/filepath"
	"slices"

	"github.com/pkg/browser"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/path"
	"github.com/tkw1536/goprogram/exit"
)

//spellchecker:words CANSPEC godoc localgodoc reclone urlweb positionals GGROOT worktree weburl workdir

// Optional name of git remote to show url for.
var Web ggman.Command = urlweb{
	isWebCommand: true,
}

// URL is the 'ggman url' command.
//
// The ggman url command behaves exactly like the ggman web command, except that instead of opening the URL in a web browser it prints it to standard output.
var URL ggman.Command = urlweb{
	isWebCommand: false,
}

type urlweb struct {
	isWebCommand bool // if true, execute the web command; else the url command

	Positionals struct {
		Base string `description:"if provided, replace the first component with the provided base url. alternatively you can use one of the predefined base URLs. use \"--list-bases\" to see a list of predefined base URLs" positional-arg-name:"BASE"`
	} `positional-args:"true"`
	List bool `description:"print a list of all predefined base URLs" long:"list-bases" short:"l"`

	ForceRepoHere bool   `description:"pretend there is a repository in the current path and use the path relative to the GGROOT directory as the remote url" long:"force-repo-here" short:"f"`
	Branch        bool   `description:"if provided, include the HEAD reference in the resolved URL"                                                           long:"branch"          short:"b"`
	Tree          bool   `description:"if provided, additionally use the HEAD reference and relative path to the root of the git worktree"                    long:"tree"            short:"t"`
	BaseAsPrefix  bool   `description:"treat the base argument as a prefix, instead of the hostname"                                                          long:"prefix"          short:"p"`
	Clone         bool   `description:"if provided to the url command, print a \"git clone\" command that can be used to clone the current repository"        long:"clone"           short:"c"`
	ReClone       bool   `description:"like clone, but uses the current remote url as opposed to the https one"                                               long:"reclone"         short:"r"`
	Remote        string `description:"optional name of git remote to show url for"                                                                           long:"remote"          short:"R"`
}

// WebBuiltInBases is a map of built-in bases for the url and web commands.
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
	Message:  "not inside a ggman-controlled repository",
}

var errWebNoRemote = exit.Error{
	ExitCode: env.ExitInvalidRepo,
	Message:  "repository does not have a remote",
}

var errNoRelativeRepository = exit.Error{
	ExitCode: env.ExitInvalidRepo,
	Message:  "unable to use `--relative`: not inside GGROOT",
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
		return uw.listBases(context)
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

			// lookup in builtin
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

		if !uw.Clone && !uw.ReClone {
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
		return browser.OpenURL(weburl)
	} else {
		_, err := context.Println(weburl)
		return ggman.ErrGenericOutput.WrapError(err)
	}
}

func (uw urlweb) listBases(context ggman.Context) error {
	bases := make([]string, 0, len(WebBuiltInBases))
	for key := range WebBuiltInBases {
		bases = append(bases, key)
	}
	slices.Sort(bases)

	for _, name := range bases {
		base := WebBuiltInBases[name]
		if _, err := context.Printf("%s: %s\n", name, base.URL); err != nil {
			return ggman.ErrGenericOutput.WrapError(err)
		}
	}
	return nil
}

// getRemoteURL gets the remote url of current repository in the context.
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
	remote, err = context.Environment.Git.GetRemote(root, uw.Remote)
	if err != nil {
		return "", "", "", errOutsideRepository
	}

	return
}

func (uw urlweb) getRemoteURLFake(context ggman.Context) (root string, remote string, relative string, err error) {
	// get the absolute path to the current working directory
	workdir, err := context.Environment.Abs("")
	if err != nil {
		return "", "", "", err
	}

	// check that the
	if !path.HasChild(context.Environment.Root, workdir) {
		return "", "", "", errNoRelativeRepository
	}

	// determine the relative path to the root directory
	relPath, err := filepath.Rel(context.Environment.Root, workdir)
	if err != nil {
		return "", "", "", errNoRelativeRepository
	}

	// turn it into a fake url by prepending a protocol
	return "", "file://" + relPath, "", nil
}
