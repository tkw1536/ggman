package program

import (
	"fmt"
	"strings"

	"github.com/tkw1536/ggman/src/constants"
	"github.com/tkw1536/ggman/src/utils"

	"github.com/tkw1536/ggman/src/repos"
)

// SubOptions Represents the options a sub-command takes
type SubOptions struct {
	// boolean indicating if the command takes a 'for'
	ForArgument int

	// minimum and maximum number of arguments
	MinArgs int
	MaxArgs int

	// the name of the metavar to use for the usage string
	Metavar string

	// if set, the command is assumed to take a single flag of the given name
	Flag string

	// Description of the flag or the argument
	UsageDescription string

	// environment configuration the command needs
	NeedsRoot    bool
	NeedsCANFILE bool
}

const (
	// NoFor indicates that no For is allowed for the commanand
	NoFor = iota
	// OptionalFor allows an optional 'for' for the command, but does not require it
	OptionalFor
)

// SubRuntime are the runtime options passed to a sub-command
type SubRuntime struct {
	// the original sub command arguments
	Args *SubCommandArgs

	// the 'for' provided
	For string

	// the arguments and their count
	Argc int
	Argv []string

	// was the flag provided
	Flag bool

	// the root folder
	Root string

	// the CanFile
	Canfile []repos.CanLine
}

// Usage returns a usage string for this command
func (opt *SubOptions) Usage(name string) (usage string) {
	usage = "Usage: ggman"

	// the for argument
	if opt.ForArgument == OptionalFor {
		usage += " [for|--for|-f FILTER]"
	}

	// the name and help
	usage += " " + name + " [help|--help|-h]"

	flagString := ""
	if opt.Flag != "" {
		flagString += " [" + opt.Flag + "]"
	} else {
		// read the metavar
		mv := opt.Metavar
		if mv == "" {
			mv = "ARGUMENT"
		}

		// write out the argument an appropriate number of times
		flagString += strings.Repeat(" "+mv, opt.MinArgs)
		flagString += strings.Repeat(" ["+mv, opt.MaxArgs-opt.MinArgs)
		flagString += strings.Repeat("]", opt.MaxArgs-opt.MinArgs)
	}

	usage += flagString

	// start with the help argument
	usage += `

    help|--help|-h
        Print this usage message and exit.`

	// contineu with the 'for' argument
	if opt.ForArgument != NoFor {
		usage += `

    for|--for|-f FILTER
        Filter the list of repositories to apply command to by FILTER.`
	}

	// and finally add the argument description
	usage += fmt.Sprintf(`

   %s
        %s`, flagString, opt.UsageDescription)

	return
}

// Apply applies the sub-command options to sub-command arguments and runs them
func (opt *SubOptions) Apply(pgrm *Program, Args *SubCommandArgs) (runtime *SubRuntime, shouldContinue bool, retval int, err string) {
	// generate a new runtime variable
	runtime = &SubRuntime{Args: Args}

	// if we have a 'help' argument, print the usage and then exit
	if utils.SliceContainsAny(Args.args, helpLongForm, helpShortForm, helpLiteralForm) {
		pgrm.Print(opt.Usage(Args.Command))
		return
	}

	// if we do not allow a for argument, then we need to ensure that
	if opt.ForArgument == NoFor {
		retval, err = Args.EnsureNoFor()
		if retval != 0 {
			return
		}
	} else {
		runtime.For = Args.Pattern
	}

	// either read a flag
	if opt.Flag != "" {
		runtime.Flag, retval, err = Args.ParseSingleFlag(opt.Flag)
		if retval != 0 {
			return
		}

		// or read the arguments
	} else {
		runtime.Argc, runtime.Argv, retval, err = Args.EnsureArguments(opt.MinArgs, opt.MaxArgs)
		if retval != 0 {
			return
		}
	}

	// read the root folder or panic
	if opt.NeedsRoot {
		var e error
		runtime.Root, e = GetRootOrPanic()
		if e != nil {
			err = constants.StringUnableParseRootDirectory
			retval = constants.ErrorMissingConfig
			return
		}
	}

	// read the canfile or panic
	if opt.NeedsCANFILE {
		var e error
		runtime.Canfile, e = GetCanonOrPanic()
		if e != nil {
			err = constants.StringInvalidCanfile
			retval = constants.ErrorMissingConfig
			return
		}
	}

	// everything is ok, so we should continue
	shouldContinue = true
	return
}
