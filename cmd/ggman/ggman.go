// Command ggman is the main entry point into ggman.
// ggman is a golang tool that can manage all your git repositories.
// An invocation of ggman is as follows:
//
//  ggman [general arguments...] SUBCOMMAND [command arguments...]
//
// Each invocation of ggman calls one subcommand.
// Arguments passed to ggman are split into general arguments and command arguments.
//
// General Arguments
//
// General Arguments are supported by every call to 'ggman'. The following arguments are supported:
//
//	--help, -h
//
// Instead of running a subcommand, print a short usage dialog to STDOUT and exit.
//
//  --version, -v
//
// Instead of running a subcommand, print version information to STDOUT and exit.
//
//  --for filter, -f filter, -i file, --from-file file, --here, -H, --path, -P, --dirty, -d, --clean, -c,
//  --synced, -s, --unsynced, -u, --pristine, -p, --tarnished, -t
//
// Apply FILTER to list of repositories. See Environment section below.
//
//  --no-fuzzy-filter, -n
//
// Disable fuzzy filter matching. See Environment section below.
//
// See the Arguments type of the github.com/tkw1536/goprogram package for more details of argument parsing.
//
// Subcommands and their Arguments
//
// Each subcommand is defined as a single variable (and private associated struct) in the github.com/tkw1536/ggman/cmd package.
//
// As the name implies, ggman supports command specific arguments to be passed to each subcommand.
// These are documented in the cmd package on a per-subcommand package.
//
// In addition to subcommand specific commands, one can also use the 'help' argument safely with each subcommand.
// Using this has the effect of printing a short usage message to the command line, instead of running the command.
//
// Environment
//
// ggman manages all git repositories inside a given root directory, and automatically sets up new repositories relative to the URLs they are cloned from.
// This root folder defaults to '$HOME/Projects' but can be customized using the 'GGROOT' environment variable.
//
// For example, when 'ggman' clones a repository 'https://github.com/hello/world.git', this would automatically end up in 'GGROOT/github.com/hello/world'.
// This behavior is not limited to 'github.com' urls, but instead works for any git remote url.
//
// As of ggman 1.12, this translation of URLs into paths takes existing paths into account.
// In particular, it re-uses existing sub-paths if they differ from the requested path only by casing.
// By default, the first matching directory is used as opposed to creating a new one.
// If a directory with the exact name exists, this is prefered over a case-insensitive match.
//
// This normalization behavior can be controlled using the 'GGNORM' environment variable.
// It has three values: 'smart' (default behavior), 'fold' (fold paths, but do not prefer exact matches) and 'none' (always use exact paths, legacy behavior).
//
// Any subcommand that iterates over local repositories will recursively find all repositories inside the 'GGROOT' directory.
// In some scenarios it is desired to filter the local list of repositories, e.g. applying only to those inside a specific namespace.
// This can be achieved using the '--for' flag, which will match to any component of the url.
// This matching is fuzzy by default, by the fuzzyness can be disabled by passing the '--no-fuzzy-filter' flag.
// The '--for' flag also matches (relative or absolute) filesystem paths, as well as full clone URLs.
//
// A list of '--for' patterns or paths can also be read from a file name by using the '--from-file' argument.
// This functions exactly as if one would provide a path via a --for argument except that it is read from a file.
// Each filter should be on it's own line.
// Lines are trimmed for whitespace, and blank lines or lines starting with ';', '//' or '#' are ignored.
//
// Furthermore the '--path' flag can be used to match the repository inside of or contained inside the provided directories.
// '--here' is an alias for '--path .'. The '--path' flag can be provided multiple times.
//
// The '--dirty' and '--clean' flags can be used to only match repositories that have a dirty or clean working directory.
// The '--synced' and '--unsynced' flags can be used to only match repositories that are or are not synced with the remote.
//
// The '--pristine' filter can be used to only match pristine repositories: Those are clean and have all changes synced.
// '--tarnished' can be used to match all non-pristine repositories.
//
// On 'github.com' and multiple other providers, it is usually possible to clone repositories via multiple urls.
// For example, the repository at https://github.com/hello/world can be cloned using both
//  git clone https://github.com/hello/world.git
// and
//  git clone git@github.com:hello/world.git
//
// Usually the latter url is prefered to make use of SSH authentication.
// This avoids having to repeatedly type a password.
// For this purpose, ggman implements the concept of 'canonical urls'.
// This causes it to treat the latter url as the main one and uses it to clone the repository.
// The exact canonical URLs being used can be configured by the user using a so-called 'CANFILE'.
//
// See Package github.com/tkw1536/ggman/env for more information about repository urls and the environment.
//
// Exit Code
//
// When a subcommand succeeds, ggman exits with code 0.
// When something goes wrong, it instead exits with a non-zero exit code.
// Exit codes are defined by the ExitCode type in the github.com/tkw1536/ggman package.
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/cmd"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/stream"
)

// the main ggman program that will contain everything
var ggmanExe = ggman.NewProgram()

func init() {
	// register all the known commands to the ggman program!
	for _, c := range []ggman.Command{
		cmd.Canon,
		cmd.Clone,
		cmd.Comps,
		cmd.Exec,
		cmd.Fetch,
		cmd.FindBranch,
		cmd.Fix,
		cmd.Here,
		cmd.License,
		cmd.Link,
		cmd.Ls,
		cmd.Lsr,
		cmd.Pull,
		cmd.Relocate,
		cmd.Env,
		cmd.Sweep,
		cmd.URL,
		cmd.Web,
		cmd.Where,
	} {
		ggmanExe.Register(c)
	}

	// register all the aliases to the program
	for _, a := range []goprogram.Alias{
		{Name: "git", Command: "exec", Args: []string{"git"}, Description: "Execute a git command using a native 'git' executable. "},
		{Name: "root", Command: "env", Args: []string{"GGROOT"}, Description: "Print the ggman root folder. "},
		{Name: "require", Command: "clone", Args: []string{"--force"}, Description: "Require a remote git repository to be installed locally. "},
	} {
		ggmanExe.RegisterAlias(a)
	}
}

// an error when nor arguments are provided.
var errNoArgumentsProvided = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Need at least one argument. Use `ggman license` to view licensing information",
}

func main() {
	// recover from calls to panic(), and exit the program appropriatly.
	// This has to be in the main() function because any of the library functions might be broken.
	// For this reason, as few ggman functions as possible are used here; just stuff from the top-level ggman package.
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, fatalPanicMessage, err)
			debug.PrintStack()
			exit.ExitPanic.Return()
		}
	}()

	streams := stream.FromEnv()

	// when there are no arguments then parsing argument *will* fail
	//
	// we don't need to even bother with the rest of the program
	// just immediatly return a custom error message.
	if len(os.Args) == 1 {
		streams.Die(errNoArgumentsProvided)
		errNoArgumentsProvided.Return()
		return
	}

	// execute the main program with the real environment!
	err := exit.AsError(ggmanExe.Main(streams, env.Parameters{
		Variables: env.ReadVariables(),
		Plumbing:  nil,
		Workdir:   "",
	}, os.Args[1:]))

	err.Return()
}

const fatalPanicMessage = `Fatal Error: Panic

The ggman program panicked and had to abort execution. This is usually
indicative of a bug. If this occurs repeatedly you might want to consider
filing an issue in the issue tracker at:

https://github.com/tkw1536/ggman/issues

Below is debug information that might help the developer track down what
happened.

panic: %v
`
