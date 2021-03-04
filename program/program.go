package program

import (
	"fmt"
	"reflect"
	"runtime"
	"sort"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/usagefmt"
)

// Program represents an executable program with a list of subcommands.
// the zero value is ready to use.
type Program struct {
	ggman.IOStream

	commands map[string]Command
}

// Commands returns a list of known commands
func (p Program) Commands() []string {
	commands := make([]string, 0, len(p.commands))
	for cmd := range p.commands {
		commands = append(commands, cmd)
	}
	sort.Strings(commands)
	return commands
}

// FmtCommands returns a human readable string describing the commands.
// See also Commands.
func (p Program) FmtCommands() string {
	return usagefmt.FmtCommands(p.Commands())
}

// Command represents a single command to be parsed.
//
// A command may contain state representing different flags.
// Flag parsing is implemented using the "github.com/jessevdk/go-flags" package.
// A Command implementation that is not a pointer to a struct is assumed to be flagless.
//
// Typically command contains state that represents the parsed options.
// This would prevent a single value of type command to run multiple times.
// To work around this, the CloneCommand method exists.
//
// In order for the CloneCommand method to work correctly, a Command must fullfill the following:
// If it is not implemented as a pointer receiver, the zero value is expected to be ready to use.
// Otherwise the zero value of the element struct is expected to be ready to use.
// See also CloneCommand.
type Command interface {
	// BeforeRegister is called right before this command is registered with a program.
	// In particular it is called before any other function on this command is called.
	//
	// It is never called more than once for a single instance of a command.
	BeforeRegister()

	// Description returns a description of this command.
	// It may be called multiple times.
	Description() Description

	// AfterParse is called after arguments have been parsed, but before the command is being run.
	// It may perform additional argument checking and should return an error if needed.
	//
	// It is called only once and must return either nil or an error of type Error.
	AfterParse() error

	// Run runs this command in the given context.
	//
	// It is called only once and must return either nil or an error of type Error.
	Run(context Context) error
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

// CloneCommand returns a new Command that behaves exactly like Command,
// except that it does not modify any internal state of Command.
//
// This function is mostly intended to be used when a command should be called multiple times
// during a single run of ggman.
func CloneCommand(command Command) (cmd Command) {
	cmdStruct := reflect.ValueOf(command) // cmd.CommandStruct

	// clone := cmd.CommandStruct{...zero...}
	var clone reflect.Value
	if cmdStruct.Type().Kind() == reflect.Ptr {
		clone = reflect.New(cmdStruct.Type().Elem())
	} else {
		clone = reflect.Zero(cmdStruct.Type())
	}

	// command = clone
	reflect.ValueOf(&cmd).Elem().Set(clone)
	return cmd
}

// Description represents the description of a command
type Description struct {
	Name        string // name this command can be invoked under
	Description string // human readable description of the command

	Environment env.Requirement // environment requirements for this command

	SkipUnknownOptions bool // do not complain about unkown options and add the to positionals instead

	// description of the positional arguments this command takes in addition to the regular option parsing.

	PosArgName        string // used only in help page, defaults to "ARGUMENT"
	PosArgDescription string // used only in help page, human readable
	PosArgsMin        int    // minimum number of positional arguments taken >= 0
	PosArgsMax        int    // maximal number of positional arguments taken, set to -1 for unlimited arguments
}

var errProgramUnknownCommand = ggman.Error{
	ExitCode: ggman.ExitUnknownCommand,
	Message:  "Unknown command. Must be one of %s. ",
}

var errInitContext = ggman.Error{
	ExitCode: ggman.ExitInvalidEnvironment,
	Message:  "Unable to initialize context: %s",
}

// Main is the entry point to this program.
// When an error occurs, returns an error of type Error and writes the error to context.Stderr.
func (p Program) Main(params env.EnvironmentParameters, argv []string) (err error) {
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
		p.StdoutWriteWrap(p.UsagePage().String())
		return nil
	case args.Version:
		p.printVersion()
		return nil
	}

	// load the command if we have it
	command, hasCommand := p.commands[args.Command]
	if !hasCommand {
		return errProgramUnknownCommand.WithMessageF(p.FmtCommands())
	}

	// parse the command arguments
	cmdargs := &CommandArguments{}
	if err := cmdargs.Parse(command, *args); err != nil {
		return err
	}

	// special cases of arguments
	switch {
	case cmdargs.Help:
		p.StdoutWriteWrap(cmdargs.UsagePage().String())
		return nil
	}

	// create a new context and make an environment for it
	context := &Context{
		IOStream:         p.IOStream,
		CommandArguments: *cmdargs,
	}
	if context.Env, err = env.NewEnv(cmdargs.description.Environment, params); err != nil {
		return err
	}

	// initialize the context
	if err := context.init(); err != nil {
		return errInitContext.WithMessageF(err)
	}

	return command.Run(*context)
}

const stringVersion = "ggman version %s, built %s, using %s"

// printVersion prints version information for this program
func (p Program) printVersion() {
	p.StdoutWriteWrap(fmt.Sprintf(stringVersion, constants.BuildVersion, constants.BuildTime, runtime.Version()))
}

// Register registers a new command with this program.
// It expects that the command does not have a name that is already taken.
func (p *Program) Register(c Command) {
	if p.commands == nil {
		p.commands = make(map[string]Command)
	}

	c.BeforeRegister()
	Name := c.Description().Name

	if _, ok := p.commands[Name]; ok {
		panic("Register(): Command already registered")
	}

	p.commands[Name] = c
}
