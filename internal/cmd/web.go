package cmd

//spellchecker:words context path filepath slices github browser cobra ggman internal pkglib exit
import (
	"context"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/ggman/internal/path"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words godoc localgodoc reclone urlweb positionals GGROOT worktree weburl workdir wrapcheck

// NewWebCommand creates the 'ggman web' command.
func NewWebCommand() *cobra.Command {
	impl := &urlweb{isWebCommand: true}
	cmd := &cobra.Command{
		Use:   "web [BASE]",
		Short: "Open the URL of this repository in a web browser",
		Long:  "Web opens the URL of this repository in a web browser.",
		Args:  cobra.MaximumNArgs(1),

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	addURLWebFlags(cmd, impl)

	return cmd
}

func NewURLCommand() *cobra.Command {
	impl := &urlweb{isWebCommand: false}

	cmd := &cobra.Command{
		Use:   "url [BASE]",
		Short: "Print the URL to this repository for opening a web browser",
		Long:  "Url behaves exactly like the 'ggman web' command, except that instead of opening the URL in a web browser it prints it to standard output.",
		Args:  cobra.MaximumNArgs(1),

		PreRunE: impl.ParseArgs,
		RunE:    impl.Exec,
	}

	addURLWebFlags(cmd, impl)

	return cmd
}

func addURLWebFlags(cmd *cobra.Command, impl *urlweb) {
	flags := cmd.Flags()
	flags.BoolVarP(&impl.List, "list-bases", "l", false, "print a list of all predefined base URLs")
	flags.BoolVarP(&impl.ForceRepoHere, "force-repo-here", "f", false, "pretend there is a repository in the current path and use the path relative to the GGROOT directory as the remote url")
	flags.BoolVarP(&impl.Branch, "branch", "b", false, "if provided, include the HEAD reference in the resolved URL")
	flags.BoolVarP(&impl.Tree, "tree", "t", false, "if provided, additionally use the HEAD reference and relative path to the root of the git worktree")
	flags.BoolVarP(&impl.BaseAsPrefix, "prefix", "p", false, "treat the base argument as a prefix, instead of the hostname")
	flags.BoolVarP(&impl.Clone, "clone", "c", false, "if provided to the url command, print a \"git clone\" command that can be used to clone the current repository")
	flags.BoolVarP(&impl.ReClone, "reclone", "r", false, "like clone, but uses the current remote url as opposed to the https one")
	flags.StringVarP(&impl.Remote, "remote", "g", "", "optional name of git remote to show url for")
}

type urlweb struct {
	isWebCommand bool // if true, execute the web command; else the url command

	Positionals struct {
		Base string
	}

	List          bool
	ForceRepoHere bool
	Branch        bool
	Tree          bool
	BaseAsPrefix  bool
	Clone         bool
	ReClone       bool
	Remote        string
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

func (uw *urlweb) ParseArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		uw.Positionals.Base = args[0]
	}

	var cloneFlag string
	if uw.Clone {
		cloneFlag = "clone"
	} else if uw.ReClone {
		cloneFlag = "reclone"
	}

	isClone := uw.Clone || uw.ReClone

	if uw.isWebCommand && isClone {
		return fmt.Errorf("%w: %q", errWebFlagUnsupported, cloneFlag)
	}

	if isClone && uw.Tree {
		return fmt.Errorf("%w: %q and %q", errURLFlagsUnsupported, cloneFlag, "tree")
	}

	if isClone && uw.BaseAsPrefix {
		return fmt.Errorf("%w: %q and %q", errURLFlagsUnsupported, cloneFlag, "prefix")
	}

	return nil
}

var (
	errURLWebNoRelative        = exit.NewErrorWithCode("unable to use `--relative`: not inside GGROOT", env.ExitInvalidRepo)
	errURLWebNoRemote          = exit.NewErrorWithCode("repository does not have a remote", env.ExitInvalidRepo)
	errURLWebOutsideRepository = exit.NewErrorWithCode("not inside a ggman-controlled repository", env.ExitInvalidRepo)

	errWebFlagUnsupported = exit.NewErrorWithCode("flag unsupported by `ggman web`", env.ExitCommandArguments)

	errURLFlagsUnsupported = exit.NewErrorWithCode("flag combination unsupported by `ggman url`", env.ExitCommandArguments)
	errURLFailedBrowser    = exit.NewErrorWithCode("failed to open browser", env.ExitGeneric)
)

func (uw *urlweb) Exec(cmd *cobra.Command, args []string) error {
	if uw.List {
		return uw.listBases(cmd)
	}

	environment, err := env.GetEnv(cmd, env.Requirement{
		NeedsRoot: true,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericEnvironment, err)
	}

	// get the remote url of the current repository
	root, remote, relative, err := uw.getRemoteURL(cmd.Context(), environment)
	if err != nil {
		return err
	}
	if remote == "" {
		return errURLWebNoRemote
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
		ref, err := environment.Git.GetHeadRef(cmd.Context(), root)
		if err != nil {
			return errURLWebOutsideRepository
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
		err := browser.OpenURL(weburl)
		if err != nil {
			return fmt.Errorf("%w: %w", errURLFailedBrowser, err)
		}
		return nil
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), weburl)
	if err != nil {
		return fmt.Errorf("%w: %w", errGenericOutput, err)
	}
	return nil
}

func (uw *urlweb) listBases(cmd *cobra.Command) error {
	bases := make([]string, 0, len(WebBuiltInBases))
	for key := range WebBuiltInBases {
		bases = append(bases, key)
	}
	slices.Sort(bases)

	for _, name := range bases {
		base := WebBuiltInBases[name]
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", name, base.URL); err != nil {
			return fmt.Errorf("%w: %w", errGenericOutput, err)
		}
	}
	return nil
}

// getRemoteURL gets the remote url of current repository in the context.
func (uw *urlweb) getRemoteURL(ctx context.Context, environment *env.Env) (root string, remote string, relative string, err error) {
	if uw.ForceRepoHere { // don't use a repository, instead fake one!
		return uw.getRemoteURLFake(environment)
	}

	return uw.getRemoteURLReal(ctx, environment)
}

func (uw *urlweb) getRemoteURLReal(ctx context.Context, environment *env.Env) (root string, remote string, relative string, err error) {
	// find the repository at the current location
	root, relative, err = environment.At(ctx, ".")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get local path to repository: %w", err)
	}

	if relative == "." {
		relative = ""
	}

	// get the remote
	remote, err = environment.Git.GetRemote(ctx, root, uw.Remote)
	if err != nil {
		return "", "", "", errURLWebOutsideRepository
	}

	return
}

func (uw *urlweb) getRemoteURLFake(environment *env.Env) (root string, remote string, relative string, err error) {
	// get the absolute path to the current working directory
	workdir, err := environment.Abs("")
	if err != nil {
		return "", "", "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// check that the
	if !path.HasChild(environment.Root, workdir) {
		return "", "", "", errURLWebNoRelative
	}

	// determine the relative path to the root directory
	relPath, err := filepath.Rel(environment.Root, workdir)
	if err != nil {
		return "", "", "", errURLWebNoRelative
	}

	// turn it into a fake url by prepending a protocol
	return "", "file://" + relPath, "", nil
}
