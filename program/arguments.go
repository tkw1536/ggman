package program

import (
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/internal/text"
	"github.com/tkw1536/ggman/program/exit"
)

// Arguments represent a set of partially parsed arguments for an invocation of the 'ggman' program.
// These should be further parsed into CommandArguments using the appropriate Parse() method.
//
// Command line argument are annotated using syntax provided by "github.com/jessevdk/go-flags".
type Arguments[Flags any] struct {
	Universals Universals
	Flags      Flags

	Command string   // command to run
	Pos     []string // positional arguments
}

// Universals holds flags added to every executable.
//
// Command line arguments are annotated using syntax provided by "github.com/jessevdk/go-flags".
type Universals struct {
	Help    bool `short:"h" long:"help" description:"Print a help message and exit"`
	Version bool `short:"v" long:"version" description:"Print a version message and exit"`
}

var ErrParseArgsNeedOneArgument = exit.Error{ // TODO: Public because test
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: Need at least one argument. Use `ggman license` to view licensing information.",
}

var errParseArgsUnknownError = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: %s",
}

// parser returns a new parser for the arguments
func (args *Arguments[Flags]) parser() *flags.Parser {
	return makeFlagsParser(args, flags.PassAfterNonOption|flags.PassDoubleDash)
}

// Parse parses arguments.
//
// When parsing fails, returns an error of type Error.
func (args *Arguments[Flags]) Parse(argv []string) error {
	// create a parser and parse the arguments
	var err error
	args.Pos, err = args.parser().ParseArgs(argv)

	// intercept unknonw flags
	if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrUnknownFlag {
		err = errParseArgsUnknownError.WithMessageF(e.Message)
	}

	// store the arguments we got and complain if there are none.
	// If we had a 'for' argument though, we should raise an error.
	if len(args.Pos) == 0 {
		switch {
		case args.Universals.Help || args.Universals.Version:
			return nil
		default:
			return ErrParseArgsNeedOneArgument
		}
	}

	// if we had help or version arguments we don't need to do
	// any more parsing and can bail out.
	if args.Universals.Help || args.Universals.Version {
		return nil
	}

	// setup command and arguments
	args.Command = args.Pos[0]
	args.Pos = args.Pos[1:]

	return err
}

var flagNameCutset = "/-"

// parseFlagNames parses flag names between `' from a flags.Error
func parseFlagNames(err *flags.Error) (names []string, ok bool) {

	// find the `' delimiters
	start := strings.IndexRune(err.Message, '`')
	end := strings.IndexRune(err.Message, '\'')

	// if they can't be found (or aren't in the right order)
	if start == -1 || end == -1 || start >= end-1 {
		return
	}

	// extract the description of the flags
	description := err.Message[start+1 : end]
	ok = true

	// trim off the names
	names = strings.Split(description, ", ")
	for i, name := range names {
		names[i] = strings.TrimLeft(name, flagNameCutset)
	}

	return
}

// CommandArguments represent a parsed set of options for a specific subcommand
// The zero value is ready to use, see the "Parse" method.
type CommandArguments[Runtime any, Parameters any, Flags any, Requirements Requirement[Flags]] struct {
	Arguments Arguments[Flags] // Arguments that were passed to the command globally

	Parser      *flags.Parser                    // TODO: Public because test
	Description Description[Flags, Requirements] // TODO: Public because test
}

// Parse parses arguments from a set of parsed command arguments.
// It also calls .Parse() with the provided arguments on the flagset
//
// It expects that neither the Help nor Version flag of Arguments are true.
//
// When parsing fails, returns an error of type Error.
func (args *CommandArguments[Runtime, Parameters, Flags, Requirements]) Parse(command Command[Runtime, Parameters, Flags, Requirements], arguments Arguments[Flags]) error {
	args.prepare(command, arguments)

	// We first have to check the following (in order):
	// - a help flag
	// - the 'for' flag
	// - the custom flag(s)
	// - the right number of arguments

	if text.SliceContainsAny(args.Arguments.Pos, "--help", "-h", "help") {
		args.Arguments.Universals.Help = true
		return nil
	}

	if err := args.CheckFilterArgument(); err != nil {
		return err
	}

	if err := args.parseFlags(); err != nil {
		return err
	}

	if err := args.CheckPositionalCount(); err != nil {
		return err
	}

	if err := command.AfterParse(); err != nil {
		return err
	}

	return nil
}

// prepare prepares this CommandArguments for parsing arguments for command
func (args *CommandArguments[Runtime, Parameters, Flags, Requirements]) prepare(command Command[Runtime, Parameters, Flags, Requirements], arguments Arguments[Flags]) {
	// setup options and arguments!
	args.Description = command.Description()
	args.Arguments = arguments

	// make a flag parser
	var options flags.Options = flags.PassDoubleDash | flags.HelpFlag
	if args.Description.SkipUnknownOptions {
		options |= flags.IgnoreUnknown
	}
	args.Parser = makeFlagsParser(command, options)
}

var errParseFlagSet = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Error parsing flags: %s",
}

// parseFlagset calls Parse() on the flagset.
// If the flagset has no defined flags (or is nil), immediatly returns nil
//
// When an error occurs, returns an error of type Error.
func (args *CommandArguments[Runtime, Parameters, Flags, Requirements]) parseFlags() (err error) {
	args.Arguments.Pos, err = args.Parser.ParseArgs(args.Arguments.Pos)

	// catch the help error
	if flagErr, ok := err.(*flags.Error); ok && flagErr.Type == flags.ErrHelp {
		args.Arguments.Universals.Help = true
		err = nil
	}

	// if an error occured, return it!
	if err != nil {
		err = errParseFlagSet.WithMessageF(err.Error())
	}

	return err
}

var errParseTakesExactlyArguments = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes exactly %d argument(s). ",
}

var errParseTakesNoArguments = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes no arguments. ",
}

var errParseTakesMinArguments = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes at least %d argument(s). ",
}

var errParseTakesBetweenArguments = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes between %d and %d arguments. ",
}

// checkPositionalCount checks that the correct number of arguments was passed to this command.
// This function implicitly assumes that Options, Arguments and Argv are set appropriatly.
// When the wrong number of arguments is passed, returns an error of type Error.
func (args CommandArguments[Runtime, Parameters, Flags, Requirements]) CheckPositionalCount() error {
	// TODO: Public because test!

	min := args.Description.PosArgsMin
	max := args.Description.PosArgsMax

	argc := len(args.Arguments.Pos)

	// If we are outside the range for the arguments, we reset the counter to 0
	// and return the appropriate error message.
	//
	// - we always need to be more than the minimum
	// - we need to be below the max if the maximum is not unlimited
	if argc < min || ((max != -1) && (argc > max)) {
		switch {
		case min == max && min == 0: // 0 arguments, but some given
			return errParseTakesNoArguments.WithMessageF(args.Arguments.Command)
		case min == max: // exact number of arguments is wrong
			return errParseTakesExactlyArguments.WithMessageF(args.Arguments.Command, min)
		case max == -1: // less than min arguments
			return errParseTakesMinArguments.WithMessageF(args.Arguments.Command, min)
		default: // between set number of arguments
			return errParseTakesBetweenArguments.WithMessageF(args.Arguments.Command, min, max)
		}
	}

	return nil
}

// checkFilterArgument checks that any filter argument (like --for) which is not allowed is not passed.
// It expects argument passing to have occured.
//
// When filter arguments are allowed, immediatly returns nil.
// When filter arguments are not allowed returns an error of type Error iff the check fails.
func (args CommandArguments[Runtime, Parameters, Flags, Requirements]) CheckFilterArgument() error {
	// TODO: public because test!
	return args.Description.Requirements.Validate(args.Arguments)
}
