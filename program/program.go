package program

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
)

// Program represents an executable program with a list of subcommands.
// the zero value is ready to use.
type Program struct {
	ggman.IOStream

	commands map[string]Command
}

// Command represents a single command to be parsed
type Command interface {
	// Name returns the name of this command
	Name() string

	// Options returns the options of this command and adds appropriate flags to the flagset
	Options(flagset *pflag.FlagSet) Options

	// AfterParse is called after arguments have been parsed, but before the command is being run.
	// It is intended to perform any additional error checking on arguments, and return an error if needed.
	// It is expected to return either nil or type Error.
	AfterParse() error

	// Run runs this command in the given context.
	// This function should assume that flagset.Parse() has been called.
	// The error returned should be either nil or of type ggman.Error
	Run(context Context) error
}

// Options represent the options for a specific command
type Options struct {
	Environment env.Requirement

	// minimum and maximum number of arguments
	MinArgs int
	MaxArgs int

	// the name of the metavar to use for the usage string
	Metavar string

	// Description of the argument
	UsageDescription string
}

var errProgramUnknownCommand = ggman.Error{
	ExitCode: ggman.ExitUnknownCommand,
	Message:  "Unknown command. Must be one of %s. ",
}

// Main is the entry point to this program.
// When an error occurs, returns an error of type Error and writes the error to context.Stderr.
func (p Program) Main(vars env.Variables, argv []string) (err error) {
	// whenever an error occurs, we want it printed
	defer func() {
		err = p.Die(err)
	}()

	// parse the general arguments
	args := &Arguments{}
	if err := args.Parse(argv); err != nil {
		return err
	}

	// handle special cases
	switch {
	case args.Help:
		p.StdoutWriteWrap(p.Usage(args.flagsetGlobal))
		return nil
	case args.Version:
		p.printVersion()
		return nil
	}

	// load the command if we have it
	command, hasCommand := p.commands[args.Command]
	if !hasCommand {
		return errProgramUnknownCommand.WithMessageF(p.knownCommandsString())
	}

	// parse the command arguments
	cmdargs := &CommandArguments{}
	if err := cmdargs.Parse(command, *args); err != nil {
		return err
	}

	// special cases of arguments
	switch {
	case cmdargs.Help:
		p.StdoutWriteWrap(cmdargs.options.Usage(args.Command, cmdargs.flagsetCommand))
		return nil
	}

	// create a new context and make an environment for it
	context := Context{
		IOStream:         p.IOStream,
		CommandArguments: *cmdargs,
	}
	if context.Env, err = env.NewEnv(cmdargs.options.Environment, vars, cmdargs.For); err != nil {
		return err
	}

	return command.Run(context)
}

const stringVersion = "ggman version %s, built %s"

// printVersion prints version information for this program
func (p Program) printVersion() {
	p.StdoutWriteWrap(fmt.Sprintf(stringVersion, constants.BuildVersion, constants.BuildTime))
}

// Register registers a new command with this program.
// It expects that the command does not have a name that is already taken.
func (p *Program) Register(c Command) {
	if p.commands == nil {
		p.commands = make(map[string]Command)
	}

	p.commands[c.Name()] = c
}
