package program

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
)

// TODO: Return a custom error type in these functions that isn't a string.

// SubCommand represents a command that can be run with the program
type SubCommand func(runtime *SubRuntime) (retval int, err string)

// SubCommandArgs represents the arguments passed to a gg command
// gg [for $pattern] $command [$args...]
type SubCommandArgs struct {
	Command string
	Pattern string
	Help    bool
	Version bool
	args    []string
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

// ParseArgs parses arguments from the command line
func ParseArgs(args []string) (parsed *SubCommandArgs, err string) {

	// if we have no arguments, that is an error
	count := len(args)
	if count == 0 {
		err = constants.StringNeedOneArgument
		return
	}

	// The ParseArguments method only needs to examine the first argument.
	// If this is a help, version or for argument, it gets treated accordingly.
	// Otherwise, we assume it is a subcommand to be run and run it

	switch args[0] {
	case helpLiteralForm, helpShortForm, helpLongForm:
		parsed = &SubCommandArgs{
			Help: true,
			args: args[1:],
		}
		return
	case versionLiteralForm, versionShortForm, versionLongForm:
		parsed = &SubCommandArgs{
			Version: true,
			args:    args[1:],
		}
		return
	case forLiteralForm, forShortForm, forLongForm:
		if count < 3 {
			err = constants.StringNeedTwoAfterFor
			return
		}

		parsed = &SubCommandArgs{
			Command: args[2],
			Pattern: args[1],
			args:    args[3:],
		}
		return
	}

	parsed = &SubCommandArgs{
		Command: args[0],
		args:    args[1:],
	}
	return
}

// ParseFlag parses a single flag
func (parsed *SubCommandArgs) ParseFlag(opt *SubOptions) (value bool, retval int, err string) {
	la := len(parsed.args)

	// check if the flag has been set
	// and if so remove the flag from the rest of the args
	if la > 0 && parsed.args[0] == opt.Flag {
		value = true
		parsed.args = parsed.args[1:]
	}

	// if we have exactly zero arguments, the flag is mandatory or to be omitted
	if opt.MinArgs == 0 && opt.MaxArgs == 0 {

		// when we got extra arguments, or we got an invalid flag value
		// show a dedicated error message
		if la > 1 || (la == 1 && !value) {
			err = fmt.Sprintf(constants.StringUnknownArgument, parsed.Command, opt.Flag)
			retval = constants.ErrorSpecificParseArgs
			return
		}

	}

	// and return
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
	if argc < min || ((max != -1) && (argc > max)) {
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

			// special case: maximal number of arguments == -1 => unlimited number of arguments allowed
		} else if max == -1 {
			err = fmt.Sprintf(constants.StringTakesMinArguments, parsed.Command, min)

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
