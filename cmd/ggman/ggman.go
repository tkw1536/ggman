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
//  --for filter, -f filter, --here, -H, --dirty, -d, --clean, -c,
//  --synced, -s, --unsynced, -u, --pristine, -p, --tarnished, -t
//
// Apply FILTER to list of repositories. See Environment section below.
//
//  --no-fuzzy-filter, -n
//
// Disable fuzzy filter matching. See Environment section below.
//
// See the Arguments type of the github.com/tkw1536/ggman/program package for more details of package argument parsing.
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
// Furthermore the '--here' flag can also be used to match the repository in the current working directory.
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
	"github.com/tkw1536/ggman/program"
)

// the main ggman program, created from the environment
var ggmanExe *program.Program = &program.Program{IOStream: ggman.NewEnvIOStream()}

// register all the commands to the ggman program!
func init() {
	for _, c := range []program.Command{
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
		cmd.Root,
		cmd.Sweep,
		cmd.URL,
		cmd.Web,
		cmd.Where,
	} {
		ggmanExe.Register(c)
	}

	for _, a := range []program.Alias{
		{Name: "git", Command: "exec", Args: []string{"git"}, Description: "Execute a git command using a native 'git' executable. "},
	} {
		ggmanExe.RegisterAlias(a)
	}
}

func main() {

	// recover from calls to panic(), and exit the program appropriatly.
	// This has to be in the main() function because any of the libary functions might be broken.
	// For this reason, no ggman functions are used here; just stuff from the main package.
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, fatalPanicMessage, err)
			debug.PrintStack()
			ggman.ExitPanic.Return()
		}
	}()

	// Create the 'ggman' program and register all the subcommands
	// Then execute the program and handle the exit code.
	err := ggman.AsError(ggmanExe.Main(env.EnvironmentParameters{
		Variables: env.ReadVariables(),
		Plumbing:  nil,
		Workdir:   "",
	}, os.Args[1:]))
	err.Return()
}

const fatalPanicMessage = `Fatal Error: Panic

The ggman program panicked and had to abort execution. This is usually
indicative of a bug. If this occurs repeatedly you might want to consider
filing an issue in the issue tracker at
https://github.com/tkw1536/ggman/issues. Below is debug information that might
help the developer to track down what happened. 

panic: %v
`
