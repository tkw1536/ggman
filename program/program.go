package program

import (
	"flag"
	"fmt"
	"sort"
	"strings"

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
	Options(flagset *flag.FlagSet) Options

	// AfterParse is called after arguments have been parsed, but before the command is being run.
	// It is intended to perform any additional error checking on arguments, and return an error if needed.
	// It is expected to return either nil or type Error.
	AfterParse() error

	// Run runs this command in the given context.
	// This function should assume that flagset.Parse() has been called.
	// The error returned should be either nil or of type ggman.Error
	Run(context Context) error
}

var errProgramUnknownCommand = ggman.Error{
	ExitCode: ggman.ExitUnknownCommand,
	Message:  "Unknown command. Must be one of %s. ",
}

// Main is the entry point to this program.
// When an error occurs, returns an error of type Error and writes the error to context.Stderr.
func (p Program) Main(argv []string) (err error) {

	// whenever an error occurs, we want it printed
	defer func() {
		err = p.Die(err)
	}()

	// parse the general arguments
	var args Arguments
	if err := (&args).Parse(argv); err != nil {
		return err
	}

	// handle special cases
	switch {
	case args.Help:
		p.printHelp()
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

	// get it's options
	flagset := flag.NewFlagSet("ggman "+args.Command, flag.ContinueOnError)
	options := command.Options(flagset)

	// parse the command arguments
	cmdargs := CommandArguments{
		Flagset:   flagset,
		Options:   options,
		Arguments: args,
	}

	if err := (&cmdargs).Parse(); err != nil {
		return err
	}

	// special cases of arguments
	switch {
	case cmdargs.Help:
		p.StdoutWriteWrap(options.Usage(args.Command, flagset))
		return nil
	}

	// create a new context and make an environment for it
	context := Context{
		IOStream:         p.IOStream,
		CommandArguments: cmdargs,
	}
	if context.Env, err = env.NewEnv(options.Environment, cmdargs.For); err != nil {
		return err
	}

	return command.Run(context)
}

// knownCommandsString returns a string containing all the known commands
func (p Program) knownCommandsString() string {
	keys := make([]string, 0, len(p.commands))
	for k := range p.commands {
		keys = append(keys, "'"+k+"'")
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

const stringUsage = `ggman version %s
(built %s)

Usage:
    ggman [help|--help|-h] [version|--version|-v] [for|--for|-f FILTER] COMMAND [ARGS...]

    help, --help, -h
        Print this usage dialog and exit
    
    version|--version|-v
		Print version message and exit. 
	
    for FILTER, --for FILTER, -f FILTER
        Filter the list of repositories to apply command to by FILTER. 
	
    COMMAND [ARGS...]
	    Command to call. One of %s. See individual commands for more help. 

ggman is licensed under the terms of the MIT License. Use 'ggman license'
to view licensing information. `

// printHelp prints help for this program
func (p Program) printHelp() {

	// and write it to the context
	p.StdoutWriteWrap(fmt.Sprintf(stringUsage, constants.BuildVersion, constants.BuildTime, p.knownCommandsString()))
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
