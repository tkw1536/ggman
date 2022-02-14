// Package program provides a program abstraction that can be used to create programs
package program

import (
	"github.com/tkw1536/ggman/program/exit"
	"github.com/tkw1536/ggman/program/meta"
	"github.com/tkw1536/ggman/program/stream"
)

// Program represents an executable program.
// A program is intended to be invoked on the command line.
// Each invocation of a program executes a command.
//
// Programs have 4 type parameters:
// An environment of type E, a type of parameters P, a type of flags F and a type requirements R.
//
// The Environment type E defines a runtime environment for commands to execute in.
// An Environment is created using the NewEnvironment function, taking parameters P.
//
// The type of (global) command line flags F is backed by a struct type.
// It is jointed by a type of Requirements R which impose restrictions on flags for commands.
//
// Internally a program also contains a list of commands, keywords and aliases.
//
// See the Main method for a description of how program execution takes place.
type Program[E any, P any, F any, R Requirement[F]] struct {
	// Meta-information about the current program
	// Used to generate help and version pages
	Info meta.Info

	// The NewEnvironment function associated is used to create a new environment.
	// The returned error may be nil or of type exit.Error.
	NewEnvironment func(params P, context Context[E, P, F, R]) (E, error)

	// Commands, Keywords, and Aliases associated with this program.
	// They are expanded in order; see Main for details.
	keywords map[string]Keyword[F]
	aliases  map[string]Alias
	commands map[string]Command[E, P, F, R]
}

// Main invokes this program and returns an error of type exit.Error or nil.
//
// Main takes input / output streams, parameters for the environment and a set of command-line arguments.
//
// It first parses these into arguments for a specific command to be executed.
// Next, it executes any keywords and expands any aliases.
// Finally, it executes the requested command or displays a help or version page.
//
// For keyword actions, see Keyword.
// For alias expansion, see Alias.
// For command execution, see Command.
//
// For help pages, see MainUsage, CommandUsage, AliasUsage.
// For version pages, see FmtVersion.
func (p Program[E, P, F, R]) Main(stream stream.IOStream, params P, argv []string) (err error) {
	// whenever an error occurs, we want it printed
	defer func() {
		err = stream.Die(err)
	}()

	// parse the general arguments
	var args Arguments[F]
	if err := args.Parse(argv); err != nil {
		return err
	}

	// expand keywords
	keyword, hasKeyword := p.keywords[args.Command]
	if hasKeyword {
		if err := keyword(&args); err != nil {
			return err
		}
	}

	// handle special global flags!
	switch {
	case args.Universals.Help:
		stream.StdoutWriteWrap(p.MainUsage().String())
		return nil
	case args.Universals.Version:
		stream.StdoutWriteWrap(p.Info.FmtVersion())
		return nil
	}

	// expand the command (if any)
	alias, hasAlias := p.aliases[args.Command]
	if hasAlias {
		args.Command, args.Pos = alias.Invoke(args.Pos)
	}

	// load the command if we have it
	command, hasCommand := p.commands[args.Command]
	if !hasCommand {
		return errProgramUnknownCommand.WithMessageF(p.FmtCommands())
	}

	// create a new context and make an environment for it
	context := Context[E, P, F, R]{
		IOStream: stream,
	}

	// parse the command arguments
	if err := context.Parse(command, args); err != nil {
		return err
	}

	// special cases of arguments
	if context.Args.Universals.Help {
		if hasAlias {
			stream.StdoutWriteWrap(p.AliasUsage(context, alias).String())
			return nil
		}
		stream.StdoutWriteWrap(p.CommandUsage(context).String())
		return nil
	}

	// create the environment
	if context.Environment, err = p.NewEnvironment(params, context); err != nil {
		return err
	}

	// do the command!
	return command.Run(context)
}

var errProgramUnknownCommand = exit.Error{
	ExitCode: exit.ExitUnknownCommand,
	Message:  "Unknown command. Must be one of %s. ",
}

var errInitContext = exit.Error{
	ExitCode: exit.ExitInvalidEnvironment,
	Message:  "Unable to initialize context: %s",
}
