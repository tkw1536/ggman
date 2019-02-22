package program

import (
	"fmt"

	"github.com/tkw1536/ggman/src/constants"
)

// SubCommand represents a command that can be run with the program
type SubCommand func(runtime *SubRuntime) (retval int, err string)

// SubCommandArgs represents the arguments passed to a gg command
// gg [for $pattern] $command [$args...]
type SubCommandArgs struct {
	Command string
	Pattern string
	Help    bool
	args    []string
}

const helpLongForm = "--help"
const helpShortForm = "-h"
const helpLiteralForm = "help"

const forLongForm = "--for"
const forShortForm = "-f"
const forLiteralForm = "for"

// ParseArgs parses arguments from the command line
func ParseArgs(args []string) (parsed *SubCommandArgs, err string) {

	// if we have no arguments, that is an error
	count := len(args)
	if count == 0 {
		err = constants.StringNeedOneArgument
		return
	}

	// if the first argument is help, everything else is ignroed
	head := args[0]

	// if the first argument is help, return help
	if head == helpLiteralForm || head == helpShortForm || head == helpLongForm {
		parsed = &SubCommandArgs{"", "", true, args[1:]}
		return
	}

	if head == forLiteralForm || head == forShortForm || head == forLongForm {
		// gg for $pattern $command
		if count < 3 {
			err = constants.StringNeedTwoAfterFor
			return
		}

		parsed = &SubCommandArgs{args[2], args[1], false, args[3:]}
		return
	}

	// gg $pattern $command
	parsed = &SubCommandArgs{args[0], "", false, args[1:]}
	return
}

// ParseSingleFlag parses a single optional flag
func (parsed *SubCommandArgs) ParseSingleFlag(flag string) (value bool, retval int, err string) {
	la := len(parsed.args)

	// if we have too many arguments throw an error
	if la > 1 || (la == 1 && parsed.args[0] != flag) {
		err = fmt.Sprintf(constants.StringUnknownArgument, parsed.Command, flag)
		retval = constants.ErrorSpecificParseArgs
		return
	}

	// return the value
	value = la == 1
	return
}

// EnsureNoFor markes a command as taking no for arguments
func (parsed *SubCommandArgs) EnsureNoFor() (retval int, err string) {
	if parsed.Pattern != "" {
		err = fmt.Sprintf(constants.StringCmdNoFor, parsed.Command)
		retval = constants.ErrorSpecificParseArgs
	}

	return
}

// EnsureNoArguments marks a command as taking no arguments
func (parsed *SubCommandArgs) EnsureNoArguments() (retval int, err string) {
	_, _, retval, err = parsed.EnsureArguments(0, 0)
	return
}

// EnsureArguments ensures that between min and max (both inclusive) arguments are given
func (parsed *SubCommandArgs) EnsureArguments(min int, max int) (argc int, argv []string, retval int, err string) {
	argc = len(parsed.args)

	// if we are outside of the range
	if argc < min || argc > max {
		// reset argc and argv
		argc = 0

		// error is specific
		retval = constants.ErrorSpecificParseArgs

		// if we have min == max we take an exact number of arguments
		if min == max {
			if min != 0 {
				err = fmt.Sprintf(constants.StringTakesExactlyArguments, parsed.Command, min)

				// special case: no arguments
			} else {
				err = fmt.Sprintf(constants.StringTakesNoArguments, parsed.Command)
			}

			// if we do not, we have a range
		} else {
			err = fmt.Sprintf(constants.StringTakesBetweenArguments, parsed.Command, min, max)
		}

		return
	}

	// return the arguments too
	argv = parsed.args
	return
}
