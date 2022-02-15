package program

import (
	"reflect"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/internal/slice"
	"github.com/tkw1536/ggman/program/exit"
)

// TODO: fix error message
var errParseArgsNeedOneArgument = exit.Error{
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

// makeFlagsParser creates a new flags parser for data.
// When data is nil or not a pointer to a struct, returns an empty parser.
//
// This function is untested.
func makeFlagsParser(data interface{}, options flags.Options) *flags.Parser {
	var actual interface{} = data
	if ptrval := reflect.ValueOf(actual); data == nil || ptrval.Type().Kind() != reflect.Ptr {
		// not a pointer to struct
		actual = &struct{}{}
	}

	return flags.NewParser(actual, options)
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
			return errParseArgsNeedOneArgument
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

// Parse parses arguments from a set of parsed command arguments.
// It also calls .Parse() with the provided arguments on the flagset
//
// It expects that neither the Help nor Version flag of Arguments are true.
//
// When parsing fails, returns an error of type Error.
func (context *Context[E, P, F, R]) Parse(command Command[E, P, F, R], arguments Arguments[F]) error {
	context.prepare(command, arguments)

	// We first have to check the following (in order):
	// - a help flag
	// - the 'for' flag
	// - the custom flag(s)
	// - the right number of arguments

	if slice.ContainsAny(context.Args.Pos, "--help", "-h", "help") {
		context.Args.Universals.Help = true
		return nil
	}

	if err := context.checkFilterArgument(); err != nil {
		return err
	}

	if err := context.parseFlags(); err != nil {
		return err
	}

	if err := context.checkPositionalCount(); err != nil {
		return err
	}

	if err := command.AfterParse(); err != nil {
		return err
	}

	return nil
}

// prepare prepares this CommandArguments for parsing arguments for command
func (context *Context[E, P, F, R]) prepare(command Command[E, P, F, R], arguments Arguments[F]) {
	// setup options and arguments!
	context.Description = command.Description()
	context.Args = arguments

	// make a flag parser
	var options flags.Options = flags.PassDoubleDash | flags.HelpFlag
	if context.Description.SkipUnknownOptions {
		options |= flags.IgnoreUnknown
	}
	context.parser = makeFlagsParser(command, options)
}

var errParseFlagSet = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Error parsing flags: %s",
}

// parseFlagset calls Parse() on the flagset.
// If the flagset has no defined flags (or is nil), immediatly returns nil
//
// When an error occurs, returns an error of type Error.
func (context *Context[E, P, F, R]) parseFlags() (err error) {
	context.Args.Pos, err = context.parser.ParseArgs(context.Args.Pos)

	// catch the help error
	if flagErr, ok := err.(*flags.Error); ok && flagErr.Type == flags.ErrHelp {
		context.Args.Universals.Help = true
		err = nil
	}

	// if an error occured, return it!
	if err != nil {
		err = errParseFlagSet.WithMessageF(err.Error())
	}

	return err
}

// checkFilterArgument checks that any filter argument (like --for) which is not allowed is not passed.
// It expects argument passing to have occured.
//
// When filter arguments are allowed, immediatly returns nil.
// When filter arguments are not allowed returns an error of type Error iff the check fails.
func (context Context[E, P, F, R]) checkFilterArgument() error {
	return context.Description.Requirements.Validate(context.Args)
}

var errParseArgCount = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of positional arguments for %s: %s",
}

func (context Context[E, P, F, R]) checkPositionalCount() error {
	err := context.Description.Positional.Validate(len(context.Args.Pos))
	if err != nil {
		return errParseArgCount.WithMessageF(context.Args.Command, err)
	}
	return nil
}
