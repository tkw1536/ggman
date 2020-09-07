package program

import (
	flag "github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/util"
)

// CommandArguments represent a parsed set of options for a specific subcommand
// The zero value is ready to use, see the "Parse" method.
type CommandArguments struct {
	Options Options       // options for this command
	Flagset *flag.FlagSet // flagset for custom options

	Arguments // Arguments, first pass

	// TODO: Move all these into FlagSet instead
	Flag bool // Flag indicates if the flag was set
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

	if err := args.parseFlag(); err != nil {
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

var errParseUnknownFlag = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Unknown argument: '%s' must be called with either '%s' or no arguments. ",
}

// parseFlag parses the flag (if set in options) passed to this command.
// This function implicitly assumes that Options and Arguments are set appropriatly.
func (args *CommandArguments) parseFlag() error {
	if args.Options.FlagValue == "" {
		return nil
	}

	la := len(args.Argv)

	// check if the flag has been set
	// and if so remove the flag from the rest of the args
	if la > 0 && args.Argv[0] == args.Options.FlagValue {
		args.Flag = true
		args.Argv = args.Argv[1:]
	}

	// if we have exactly zero arguments, the flag is mandatory or to be omitted
	if args.Options.MinArgs == 0 && args.Options.MaxArgs == 0 {

		// when we got extra arguments, or we got an invalid flag value
		// show a dedicated error message
		if la > 1 || (la == 1 && !args.Flag) {
			return errParseUnknownFlag.WithMessageF(args.Command, args.Options.FlagValue)
		}

	}

	return nil
}

var errParseFlagSet = ggman.Error{
	ExitCode: ggman.ExitCommandArguments,
	Message:  "Error parsing flags: %s",
}

// parseFlagset calls Parse() on the flagset, iff it is not nil and at least one flg is defined
//
// When an error occurs, returns an error of type Error.
func (args *CommandArguments) parseFlagset() (err error) {
	// if the flagset is nil, do nothing
	if args.Flagset == nil {
		return nil
	}

	// the only way to check if a flagset has no flags is to call VisitAll
	hasFlag := false
	args.Flagset.VisitAll(func(_ *flag.Flag) { hasFlag = true })
	if !hasFlag {
		return nil
	}

	args.Flagset.Usage = func() {} // don't print any usage messages please

	err = args.Flagset.Parse(args.Argv)
	switch err {
	case nil: /* do nothing */
	case flag.ErrHelp: /* help error, set the help flag but nothing else */
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

	if args.For != "" {
		return errParseNoFor.WithMessageF(args.Command)
	}

	return nil
}
