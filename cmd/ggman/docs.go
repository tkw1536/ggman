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
//	help, --help, -h
//
// Instead of running a subcommand, print a short usage dialog to STDOUT and exit.
//
//  version|--version|-v
//
// Instead of running a subcommand, print version information to STDOUT and exit.
//
//  for FILTER, --for FILTER, -f FILTER
//
// Apply FILTER to list of repositories. See Environment section below.
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
// Any subcommand that iterates over local repositories will recursively find all repositories inside the 'GGROOT' directory.
// In some scenarios it is desired to filter the local list of repositories, e.g. applying only to those inside a specific namespace.
// This can be achieved using the '--for' flag, which will match to any component of the url.
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
