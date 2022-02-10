package program

import (
	"reflect"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/internal/text"
	"github.com/tkw1536/ggman/program/exit"
)

// Arguments represent a set of partially parsed arguments for an invocation of the 'ggman' program.
// These should be further parsed into CommandArguments using the appropriate Parse() method.
//
// Command line argument are annotated using syntax provided by "github.com/jessevdk/go-flags".
// Additionally flags that do not imply any filter semantics are tagged as `nofilter:"true"`.
//
// See Also GetMainOpts and Parse.
type Arguments struct {
	Help    bool `short:"h" long:"help" description:"Print a help message and exit" nofilter:"true"`
	Version bool `short:"v" long:"version" description:"Print a version message and exit" nofilter:"true"`

	Filters       []string `short:"f" long:"for" value-name:"filter" description:"Filter list of repositories to apply COMMAND to by filter. Filter can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url"`
	NoFuzzyFilter bool     `short:"n" long:"no-fuzzy-filter" description:"Disable fuzzy matching for filters"`

	Here bool     `short:"H" long:"here" description:"Filter the list of repositories to apply COMMAND to only contain repository in the current directory or subtree. Alias for '-p .'"`
	Path []string `short:"P" long:"path" description:"Filter the list of repositories to apply COMMAND to only contain repositories in or under the specified path. May be used multiple times"`

	Dirty bool `short:"d" long:"dirty" description:"List only repositories with uncommited changes"`
	Clean bool `short:"c" long:"clean" description:"List only repositories without uncommited changes"`

	Synced   bool `short:"s" long:"synced" description:"List only repositories which are up-to-date with remote"`
	UnSynced bool `short:"u" long:"unsynced" description:"List only repositories not up-to-date with remote"`

	Tarnished bool `short:"t" long:"tarnished" description:"List only repositories which are dirty or unsynced"`
	Pristine  bool `short:"p" long:"pristine" description:"List only repositories which are clean and synced"`

	Command string   // command to run
	Args    []string // remaining arguments
}

// reflect access to the arguments type
var argumentsType reflect.Type = reflect.TypeOf((*Arguments)(nil)).Elem() // TypeOf[Arguments]

var argumentsGeneralOptions []string // names of options that are considered non-filter
var argumentsFilterIndexes [][]int   // indexes of filter options

func init() {
	// iterate over the fields of the type
	fieldCount := argumentsType.NumField()
	for i := 0; i < fieldCount; i++ {
		field := argumentsType.Field(i)

		// skip over options that do not have a 'long' name
		longName, hasLongName := field.Tag.Lookup("long")
		if !hasLongName {
			continue
		}

		// argument is a nonfilter argument!
		if field.Tag.Get("nofilter") == "true" {
			argumentsGeneralOptions = append(argumentsGeneralOptions, longName)
			continue
		}

		// it's a long filter name
		argumentsFilterIndexes = append(argumentsFilterIndexes, field.Index)
	}
}

var errParseArgsNeedOneArgument = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: Need at least one argument. Use `ggman license` to view licensing information.",
}

var errParseArgsUnknownError = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: %s",
}

var errParseArgsNeedTwoAfterFor = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: At least two arguments needed after 'for' keyword. ",
}

// parser returns a new parser for the arguments
func (args *Arguments) parser() *flags.Parser {
	return makeFlagsParser(args, flags.PassAfterNonOption|flags.PassDoubleDash)
}

// Parse parses arguments.
//
// When parsing fails, returns an error of type Error.
func (args *Arguments) Parse(argv []string) error {
	// create a parser and parse the arguments
	var err error
	args.Args, err = args.parser().ParseArgs(argv)

	if e, ok := err.(*flags.Error); ok {
		switch e.Type {

		// --for, -f was passed without an argument!
		case flags.ErrExpectedArgument:
			if names, ok := parseFlagNames(e); ok && text.SliceContainsAny(names, "f", "for") {
				err = errParseArgsNeedTwoAfterFor
			}

		// encounted an unknown flag
		case flags.ErrUnknownFlag:
			err = errParseArgsUnknownError.WithMessageF(e.Message)
		}
	}

	// store the arguments we got and complain if there are none.
	// If we had a 'for' argument though, we should raise an error.
	if len(args.Args) == 0 {
		switch {
		case args.Help || args.Version:
			return nil
		case len(args.Filters) > 0:
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
		args.Filters = append(args.Filters, args.Args[0])
		args.Command = args.Args[1]
		args.Args = args.Args[2:]
	}

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
type CommandArguments[Runtime any, Parameters any, Requirements any] struct {
	Arguments // Arguments that were passed to the command globally

	parser      *flags.Parser
	description Description
}

// Parse parses arguments from a set of parsed command arguments.
// It also calls .Parse() with the provided arguments on the flagset
//
// It expects that neither the Help nor Version flag of Arguments are true.
//
// When parsing fails, returns an error of type Error.
func (args *CommandArguments[Runtime, Parameters, Requirements]) Parse(command Command[Runtime, Parameters, Requirements], arguments Arguments) error {
	args.prepare(command, arguments)

	// We first have to check the following (in order):
	// - a help flag
	// - the 'for' flag
	// - the custom flag(s)
	// - the right number of arguments

	if text.SliceContainsAny(args.Args, "--help", "-h", "help") {
		args.Help = true
		return nil
	}

	if err := args.checkFilterArgument(); err != nil {
		return err
	}

	if err := args.parseFlags(); err != nil {
		return err
	}

	if err := args.checkPositionalCount(); err != nil {
		return err
	}

	if err := command.AfterParse(); err != nil {
		return err
	}

	return nil
}

// prepare prepares this CommandArguments for parsing arguments for command
func (args *CommandArguments[Runtime, Parameters, Requirements]) prepare(command Command[Runtime, Parameters, Requirements], arguments Arguments) {
	// setup options and arguments!
	args.description = command.Description()
	args.Arguments = arguments

	// make a flag parser
	var options flags.Options = flags.PassDoubleDash | flags.HelpFlag
	if args.description.SkipUnknownOptions {
		options |= flags.IgnoreUnknown
	}
	args.parser = makeFlagsParser(command, options)
}

var errParseFlagSet = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Error parsing flags: %s",
}

// parseFlagset calls Parse() on the flagset.
// If the flagset has no defined flags (or is nil), immediatly returns nil
//
// When an error occurs, returns an error of type Error.
func (args *CommandArguments[Runtime, Parameters, Requirements]) parseFlags() (err error) {
	args.Args, err = args.parser.ParseArgs(args.Args)

	// catch the help error
	if flagErr, ok := err.(*flags.Error); ok && flagErr.Type == flags.ErrHelp {
		args.Help = true
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
func (args CommandArguments[Runtime, Parameters, Requirements]) checkPositionalCount() error {

	min := args.description.PosArgsMin
	max := args.description.PosArgsMax

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

var errTakesNoArgument = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes no '%s' argument. ",
}

// checkFilterArgument checks that any filter argument (like --for) which is not allowed is not passed.
// It expects argument passing to have occured.
//
// When filter arguments are allowed, immediatly returns nil.
// When filter arguments are not allowed returns an error of type Error iff the check fails.
func (args CommandArguments[Runtime, Parameters, Requirements]) checkFilterArgument() error {
	// we don't have to do any checking!
	if args.description.Environment.AllowsFilter {
		return nil
	}

	// check the value of the arguments struct
	aVal := reflect.ValueOf(args.Arguments)

	for _, fIndex := range argumentsFilterIndexes {
		v := aVal.FieldByIndex(fIndex)

		if !v.IsZero() { // flag was set iff it is non-zero
			tp := argumentsType.FieldByIndex(fIndex) // needed for the error message only!
			return errTakesNoArgument.WithMessageF(args.Command, "--"+tp.Tag.Get("long"))
		}
	}

	return nil
}

// Description returns the description of the command invoked by these arguments
func (args CommandArguments[Runtime, Parameters, Requirements]) Description() Description {
	return args.description
}
