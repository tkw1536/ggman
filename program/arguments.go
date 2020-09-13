package program

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

// Arguments represent a set of partially parsed arguments for an invocation of the 'ggman' program.
//
// These should be further parsed into CommandArguments using the appropriate Parse() method.
type Arguments struct {
	Command string     // the command, if provided
	For     env.Filter // the 'for' argument, if provided

	Help    bool // the 'help' argument
	Version bool // the 'version' argument

	Argv []string // the rest of the arguments to be passed to the command
}

const helpLongForm = "--help"
const helpShortForm = "-h"
const helpLiteralForm = "help"

const versionLongForm = "--version"
const versionShortForm = "-v"
const versionLiteralForm = "version"

const forLongForm = "--for"
const forShortForm = "-f"
const forLiteralForm = "for"

var errParseArgsNeedOneArgument = ggman.Error{
	ExitCode: ggman.ExitGeneralArguments,
	Message:  "Unable to parse arguments: Need at least one argument. Use `ggman license` to view licensing information. ",
}

var errParseArgsNeedTwoAfterFor = ggman.Error{
	ExitCode: ggman.ExitGeneralArguments,
	Message:  "Unable to parse arguments: At least two arguments needed after 'for' keyword. ",
}

// Parse parses arguments
//
// When parsing fails, returns an error of type Error.
func (args *Arguments) Parse(argv []string) error {
	// if we have no arguments, that is an error
	count := len(argv)
	if count == 0 {
		return errParseArgsNeedOneArgument
	}

	args.Argv = argv[1:] // usually arguments are after the first argument

	// The Parse() method only needs to examine the first argument.
	// If this is a help, version or the '--for' argument, it gets treated accordingly.
	// Otherwise, we assume it is a subcommand to be run and run it

	switch argv[0] {
	case helpLiteralForm, helpShortForm, helpLongForm:
		args.Help = true
		return nil
	case versionLiteralForm, versionShortForm, versionLongForm:
		args.Version = true

		return nil
	case forLiteralForm, forShortForm, forLongForm:
		if count < 3 {
			args.Argv = nil
			return errParseArgsNeedTwoAfterFor
		}

		// parse the filter
		args.For = env.NewFilter(argv[1])
		args.Command = argv[2]
		args.Argv = argv[3:] // overwrite the existing arguments

		return nil
	}

	args.For = env.NoFilter
	args.Argv = argv[1:]
	args.Command = argv[0]

	return nil
}
