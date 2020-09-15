package program

import (
	"github.com/spf13/pflag"
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

	flagset *pflag.FlagSet // flagset used internally
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
	fs := args.setflagset()
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
	args.Argv = fs.Args()
	if len(args.Argv) == 0 {
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
	args.Command = args.Argv[0]
	args.Argv = args.Argv[1:]

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
		if len(args.Argv) < 2 {
			return errParseArgsNeedTwoAfterFor
		}
		args.For.Set(args.Argv[0])
		args.Command = args.Argv[1]
		args.Argv = args.Argv[2:]
	}

	return nil
}

// setflagset sets flagset to a new flagset for argument parsing
func (args *Arguments) setflagset() (fs *pflag.FlagSet) {
	if args.flagset != nil {
		return args.flagset
	}
	defer func() { args.flagset = fs }()

	fs = pflag.NewFlagSet("ggman", pflag.ContinueOnError)
	fs.Usage = func() {}      // don't print a usage message on error
	fs.SetInterspersed(false) // stop at the first regular argument
	fs.SortFlags = false      // flag's shouldn't be sorted

	fs.BoolVarP(&args.Help, "help", "h", false, "Print this usage dialog and exit.")
	fs.BoolVarP(&args.Version, "version", "v", false, "Print version message and exit.")

	fs.VarP(&args.For, "for", "f", "Filter the list of repositories to apply command to by `filter`.")

	return fs
}
