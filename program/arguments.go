package program

import (
	"github.com/spf13/pflag"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/util"
)

// Arguments represent a set of partially parsed arguments for an invocation of the 'ggman' program.
//
// These should be further parsed into CommandArguments using the appropriate Parse() method.
type Arguments struct {
	Command string

	For env.Filter

	Help    bool
	Version bool

	Args []string // non-flag arguments

	flagsetGlobal *pflag.FlagSet // flagset used for global argument parsing
}

const helpLiteralForm = "help"
const versionLiteralForm = "version"

var errParseArgsNeedOneArgument = ggman.Error{
	ExitCode: ggman.ExitGeneralArguments,
	Message:  "Unable to parse arguments: Need at least one argument. Use `ggman license` to view licensing information. ",
}

var errParseArgsUnknownError = ggman.Error{
	ExitCode: ggman.ExitGeneralArguments,
	Message:  "Unable to parse arguments: %s",
}

var errParseArgsNeedTwoAfterFor = ggman.Error{
	ExitCode: ggman.ExitGeneralArguments,
	Message:  "Unable to parse arguments: At least two arguments needed after 'for' keyword. ",
}

const errForNeedsArgument = "flag needs an argument: --for"
const errFNeedsArgument = "flag needs an argument: 'f' in -f"

// Parse parses arguments
//
// When parsing fails, returns an error of type Error.
func (args *Arguments) Parse(argv []string) error {

	// first parse arguments using the flagset
	// and intercept special 'for' error messages.
	fs := args.setflagsetGlobal()
	if err := fs.Parse(argv); err != nil {
		msg := err.Error()
		switch msg {
		case errForNeedsArgument, errFNeedsArgument:
			return errParseArgsNeedTwoAfterFor
		default:
			return errParseArgsUnknownError.WithMessageF(msg)
		}
	}

	// store the arguments we got and complain if there are none.
	// If we had a 'for' argument though, we should raise an error.
	args.Args = fs.Args()
	if len(args.Args) == 0 {
		switch {
		case args.Help || args.Version:
			return nil
		case !args.For.IsEmpty():
			return errParseArgsNeedTwoAfterFor
		default:
			return errParseArgsNeedOneArgument
		}
	}

	// if we had help or version arguments we don't need to do
	// any more parsing and can bail out.
	if args.Help || args.Version {
		return nil
	}

	// setup command and arguments
	args.Command = args.Args[0]
	args.Args = args.Args[1:]

	// catch special undocumented legacy flags
	// these can be provided with '--'s in front of their arguments
	switch args.Command {
	// ggman help
	case "help":
		args.Command = ""
		args.Help = true
	// ggman version
	case "version":
		args.Command = ""
		args.Version = true

	// ggman for FILTER command args...
	case "for":
		if len(args.Args) < 2 {
			return errParseArgsNeedTwoAfterFor
		}
		args.For.Set(args.Args[0])
		args.Command = args.Args[1]
		args.Args = args.Args[2:]
	}

	return nil
}

// setflagsetGlobal sets flagset to a new flagset for argument parsing
func (args *Arguments) setflagsetGlobal() (fs *pflag.FlagSet) {
	if args.flagsetGlobal != nil {
		return args.flagsetGlobal
	}
	defer func() { args.flagsetGlobal = fs }()

	fs = pflag.NewFlagSet("ggman", pflag.ContinueOnError)
	fs.Usage = func() {}      // don't print a usage message on error
	fs.SetInterspersed(false) // stop at the first regular argument
	fs.SortFlags = false      // flag's shouldn't be sorted

	fs.BoolVarP(&args.Help, "help", "h", false, "Print this usage dialog and exit.")
	fs.BoolVarP(&args.Version, "version", "v", false, "Print version message and exit.")

	fs.VarP(&args.For, "for", "f", "Filter the list of repositories to apply command to by `filter`.")

	return fs
}

// CommandArguments represent a parsed set of options for a specific subcommand
// The zero value is ready to use, see the "Parse" method.
type CommandArguments struct {
	Arguments // Arguments that were passed to the command globally

	options        Options
	flagsetCommand *pflag.FlagSet
}

// Parse parses arguments from a set of parsed command arguments.
// It also calls .Parse() with the provided arguments on the flagset
//
// It expects that neither the Help nor Version flag of Arguments are true.
//
// When parsing fails, returns an error of type Error.
func (args *CommandArguments) Parse(command Command, arguments Arguments) error {
	args.prepare(command, arguments)

	// We first have to check the following (in order):
	// - a help flag
	// - the 'for' flag
	// - the custom flag(s)
	// - the right number of arguments

	if util.SliceContainsAny(args.Args, "--help", "-h", "help") {
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

// prepare prepares this CommandArguments for parsing arguments for command
func (args *CommandArguments) prepare(command Command, arguments Arguments) {
	args.flagsetCommand = pflag.NewFlagSet("ggman "+args.Command, pflag.ContinueOnError)

	args.options = command.Options(args.flagsetCommand)
	args.Arguments = arguments
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
	if args.flagsetCommand == nil || !args.flagsetCommand.HasFlags() {
		return nil
	}

	args.flagsetCommand.Usage = func() {} // don't print any usage messages please

	err = args.flagsetCommand.Parse(args.Args)
	switch err {
	case nil: /* do nothing */
	case pflag.ErrHelp: /* help error, set the help flag but nothing else */
		args.Help = true
		err = nil
	default:
		err = errParseFlagSet.WithMessageF(err.Error())
	}

	// store back the parsed arguments
	args.Args = args.flagsetCommand.Args()
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

	min := args.options.MinArgs
	max := args.options.MaxArgs

	argc := len(args.Args)

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
	if args.options.Environment.AllowsFilter {
		return nil
	}

	if !args.For.IsEmpty() {
		return errParseNoFor.WithMessageF(args.Command)
	}

	return nil
}
