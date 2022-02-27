package program

import (
	"reflect"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/program/exit"
	"github.com/tkw1536/ggman/program/lib/slice"
	"github.com/tkw1536/ggman/program/meta"
)

// Command represents a command associated with a program.
// It takes the same type parameters as a program.
//
// Each command is first initialized using any means by the user.
// Next, it is registered with a program using the Program.Register Method.
// Once a program is called with this command, the arguments for it are parsed by making use of the Description.
// Eventually the Run method of this command is invoked, using an appropriate new Context.
// See Context for details on the latter.
//
// A command must be implemented as a struct or pointer to a struct.
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
type Command[E any, P any, F any, R Requirement[F]] interface {
	// Run runs this command in the given context.
	//
	// It is called only once and must return either nil or an error of type Error.
	Run(context Context[E, P, F, R]) error

	// BeforeRegister is called right before this command is registered with a program.
	// In particular it is called before any other function on this command is called.
	//
	// It is never called more than once for a single instance of a command.
	BeforeRegister(program *Program[E, P, F, R])

	// Description returns a description of this command.
	// It may be called multiple times.
	Description() Description[F, R]

	// AfterParse is called after arguments have been parsed, but before the command is being run.
	// It may perform additional argument checking and should return an error if needed.
	//
	// It is called only once and must return either nil or an error of type Error.
	AfterParse() error
}

// Description describes a command, and specifies any potential requirements.
type Description[F any, R Requirement[F]] struct {
	// Command and Description the name and human-readable description of this command.
	// Command must not be taken by any other command registered with the corresponding program.
	Command     string
	Description string

	// Positional holds information about positional arguments for this command.
	Positional meta.Positional

	// Requirements on the environment to be able to run the command
	Requirements R
}

// Requirement describes a requirement on a type of Flags F.
type Requirement[F any] interface {
	// AllowsFlag checks if the provided flag may be passed to fullfill this requirement
	// By default it is used only for help page generation, and may be inaccurate.
	AllowsFlag(flag meta.Flag) bool

	// Validate validates if this requirement is fullfilled for the provided global flags.
	// It should return either nil, or an error of type exit.Error.
	//
	// Validate does not take into account AllowsOption, see ValidateAllowedOptions.
	// TODO: Make this take context
	Validate(arguments Arguments[F]) error
}

// Register registers a command c with this program.
// It calls the BeforeRegister method on c, and then register.
//
// It expects that the command does not have a name that is already taken.
func (p *Program[R, P, Flags, Requirements]) Register(c Command[R, P, Flags, Requirements]) {
	if p.commands == nil {
		p.commands = make(map[string]Command[R, P, Flags, Requirements])
	}

	c.BeforeRegister(p)

	Name := c.Description().Command

	if _, ok := p.commands[Name]; ok {
		panic("Register(): Command already registered")
	}

	p.commands[Name] = c
}

// Commands returns a list of known commands
func (p Program[E, P, F, R]) Commands() []string {
	commands := make([]string, 0, len(p.commands))
	for cmd := range p.commands {
		commands = append(commands, cmd)
	}
	slice.Sort(commands)
	return commands
}

// FmtCommands returns a human readable string describing the commands.
// See also Commands.
func (p Program[E, P, F, R]) FmtCommands() string {
	return meta.JoinCommands(p.Commands())
}

// CloneCommand returns a new Command that behaves exactly like Command.
//
// This function is mostly intended to be used when a command should be called multiple times for testing.
// during a single run of ggman.
func CloneCommand[E any, P any, F any, R Requirement[F]](command Command[E, P, F, R]) (cmd Command[E, P, F, R]) {
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

var errTakesNoArgument = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "Wrong number of arguments: '%s' takes no '%s' argument. ",
}

// Validate validates that every flag f in args.flags either passes the AllowsOption method of the given requirement, or has the zero value.
// If this is not the case returns an error of type ValidateAllowedFlags.
//
// This function is intended to be used to implement the validate method of a Requirement.
func ValidateAllowedFlags[F any](r Requirement[F], args Arguments[F]) error {
	fVal := reflect.ValueOf(args.Flags)

	for _, flag := range globalFlags[F]() {
		if r.AllowsFlag(flag) {
			continue
		}

		v := fVal.FieldByName(flag.FieldName)
		if !v.IsZero() { // flag was set!
			name := flag.Long
			if len(name) == 0 {
				name = []string{""}
			}
			return errTakesNoArgument.WithMessageF(args.Command, "--"+name[0])
		}
	}

	return nil

}

var universalOpts = meta.AllFlags(flags.NewParser(&Universals{}, flags.None))

// globalOptions returns a list of global options for a command with the provided flag type
func globalOptions[F any]() (flags []meta.Flag) {
	flags = append(flags, universalOpts...)
	flags = append(flags, globalFlags[F]()...)
	return
}

// globalFlagsFor returns a list of global options for a command with the provided flag type
func globalFlagsFor[F any](r Requirement[F]) (flags []meta.Flag) {
	// filter options to be those that are allowed
	gFlags := globalFlags[F]()
	n := 0
	for _, flag := range gFlags {
		if !r.AllowsFlag(flag) {
			continue
		}
		gFlags[n] = flag
		n++
	}
	gFlags = gFlags[:n]

	// concat universal flags and normal flags
	flags = append(flags, universalOpts...)
	flags = append(flags, gFlags...)
	return
}

// globalFlags returns a list of flags for the provided flag type.
func globalFlags[F any]() []meta.Flag {
	return meta.AllFlags(flags.NewParser(new(F), flags.None))
}
