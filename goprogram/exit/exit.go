// Package exit defines exit and error behavior of programs and commands.
package exit

import (
	"os"
)

// ExitCode determines the exit behavior of a program.
// These are returned as an exit code to the operating system.
// See ExitCode.Return().
type ExitCode uint8

// Return returns this ExitCode to the operating system by invoking os.Exit().
//
// This function is untested.
func (code ExitCode) Return() {
	os.Exit(int(code))
}

const (
	// ExitZero indicates that no error occured.
	// It is the zero value of type ExitCode.
	ExitZero ExitCode = 0

	// ExitGeneric indicates a generic error occured within this invocation.
	// This typically implies a subcommand-specific behavior wants to return failure to the caller.
	ExitGeneric ExitCode = 1

	// ExitUnknownCommand indicates that the user attempted to call a subcommand that is not defined.
	ExitUnknownCommand ExitCode = 2

	// ExitGeneralArguments indiciates that the user attempted to pass invalid general arguments to the ggman command.
	ExitGeneralArguments ExitCode = 3
	// ExitCommandArguments indicates that the user attempted to pass invalid command-specific arguments to a ggman subcommand.
	ExitCommandArguments ExitCode = 4

	// ExitPanic indicates that the go code called panic() inside the executation of the ggman program.
	// This typically implies a bug inside a program.
	ExitPanic ExitCode = 255
)
