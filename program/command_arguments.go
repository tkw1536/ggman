package program

import (
	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/util"
)

// CommandArguments represent a parsed set of options for a specific subcommand
// The zero value is ready to use, see the "Parse" method.
type CommandArguments struct {
	Options Options        // options for this command
	Flagset *pflag.FlagSet // flagset for custom options

	Arguments // Arguments, first pass
}

// Parse parses arguments from a set of parsed command arguments.
// It also calls .Parse() with the provided arguments on the flagset
//
// Parse expects CommandOptions and Arguments to be set.
// It expects that neither the Help nor Version flag of Arguments are true.
//
// When parsing fails, returns an error of type Error.
func (args *CommandArguments) Parse() error {

	// We first have to check the following (in order):
	// - a help flag
	// - the 'for' flag
	// - the custom flag(s)
	// - the right number of arguments

	if util.SliceContainsAny(args.Argv, helpLongForm, helpShortForm, helpLiteralForm) {
		args.Help = true
		return nil
	}

	if err := args.checkForArgument(); err != nil {
		return err
	}

	if err := args.parseFlagset(); err != nil {
		return err
	}

	if err := args.checkArgumentCount(); err != nil {
		return err
	}

	return nil
}

var errParseFlagSet = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Error parsing flags: %s",
}

// parseFlagset calls Parse() on the flagset.
// If the flagset has no defined flags (or is nil), immediatly returns nil
//
// When an error occurs, returns an error of type Error.
func (args *CommandArguments) parseFlagset() (err error) {
	if args.Flagset == nil || !args.Flagset.HasFlags() {
		return nil
	}

	args.Flagset.Usage = func() {} // don't print any usage messages please

	err = args.Flagset.Parse(args.Argv)
	switch err {
	case nil: /* do nothing */
	case pflag.ErrHelp: /* help error, set the help flag but nothing else */
		args.Help = true
		err = nil
	default:
		err = errParseFlagSet.WithMessageF(err.Error())
	}

	// store back the parsed arguments
	args.Argv = args.Flagset.Args()
	return err
}

var errParseTakesExactlyArguments = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes exactly %d argument(s). ",
}

var errParseTakesNoArguments = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes no arguments. ",
}

var errParseTakesMinArguments = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes at least %d argument(s). ",
}

var errParseTakesBetweenArguments = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes between %d and %d arguments. ",
}

// checkArgumentCount checks that the correct number of arguments was passed to this command.
// This function implicitly assumes that Options, Arguments and Argv are set appropriatly.
// When the wrong number of arguments is passed, returns an error of type Error.
func (args CommandArguments) checkArgumentCount() error {

	min := args.Options.MinArgs
	max := args.Options.MaxArgs

	argc := len(args.Argv)

	// If we are outside the range for the arguments, we reset the counter to 0
	// and return the appropriate error message.
	//
	// - we always need to be more than the minimum
	// - we need to be below the max if the maximum is not unlimited
	if argc < min || ((max != -1) && (argc > max)) {
		switch {
		case min == max && min == 0: // 0 arguments, but some given
			return errParseTakesNoArguments.WithMessageF(args.Command)
		case min == max: // exact number of arguments is wrong
			return errParseTakesExactlyArguments.WithMessageF(args.Command, min)
		case max == -1: // less than min arguments
			return errParseTakesMinArguments.WithMessageF(args.Command, min)
		default: // between set number of arguments
			return errParseTakesBetweenArguments.WithMessageF(args.Command, min, max)
		}
	}

	return nil
}

var errParseNoFor = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes no 'for' argument. ",
}

// checkForArgument checks that if a 'for' argument is not allowed it is not passed.
// It expects args.For to be set appropriatly
//
// If the check fails, returns an error of type Error.
func (args CommandArguments) checkForArgument() error {
	if args.Options.Environment.AllowsFilter {
		return nil
	}

	if !args.For.IsEmpty() {
		return errParseNoFor.WithMessageF(args.Command)
	}

	return nil
}
