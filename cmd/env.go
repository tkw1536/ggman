package cmd

import (
	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
)

// Env is the 'ggman env' command.
//
// Env prints "name=value" pairs about the environment the ggman command is running in to standard output.
// value is escaped for use in a shell.
//
// By default, env prints information about all known variables.
// To print information about a subset of variables, they can be provided as positional arguments.
// Variables names are matched case-insensitively.
//
//	--list
//
// Instead of printing "name=value" pairs, print only the name.
//
//	--describe
//
// Instead of printing "name=value" pairs, print "name: description" pairs.
// The description explains what the value does.
//
//	--raw
//
// Instead of printing "name=value" pairs, print only the raw, unescaped value.
var Env ggman.Command = _env{}

type _env struct {
	Positionals struct {
		Vars []string `positional-arg-name:"VAR" description:"print only information about specified variables"`
	} `positional-args:"true"`

	List     bool `short:"l" long:"list" description:"instead of \"name=value\" pairs print only the variable"`
	Describe bool `short:"d" long:"describe" description:"instead of \"name=value\" pairs print \"name: description\" pairs describing the use of variables"`
	Raw      bool `short:"r" long:"raw" description:"instead of \"name=value\" pairs print only the unescaped value"`
}

func (_env) Description() ggman.Description {
	return ggman.Description{
		Command:     "env",
		Description: "print information about the ggman environment",

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

var errEnvInvalidVar = exit.Error{
	Message:  "Unknown environment variable %q",
	ExitCode: exit.ExitCommandArguments,
}

var errModesIncompatible = exit.Error{
	Message:  "At most one of '--raw', '--list' and '--describe' may be given",
	ExitCode: exit.ExitCommandArguments,
}

func (e _env) Run(context ggman.Context) error {
	// check that at most one mode was provided
	count := 0
	if e.Describe {
		count++
	}
	if e.Raw {
		count++
	}
	if e.List {
		count++
	}
	if count > 1 {
		return errModesIncompatible
	}

	variables, err := e.variables(context)
	if err != nil {
		return err
	}

	for _, v := range variables {
		switch {
		case e.List:
			context.Println(v.Key)
		case e.Raw:
			context.Println(v.Get(context.Environment, context.Program.Info))
		case e.Describe:
			context.Printf("%s: %s\n", v.Key, v.Description)
		default:
			value := shellescape.Quote(v.Get(context.Environment, context.Program.Info))
			context.Printf("%s=%s\n", v.Key, value)
		}
	}

	return nil
}

func (e _env) variables(context ggman.Context) ([]env.UserVariable, error) {
	// no variables provided => use all of them
	if len(e.Positionals.Vars) == 0 {
		return env.GetUserVariables(), nil
	}

	// names provided => make sure each exists
	variables := make([]env.UserVariable, len(e.Positionals.Vars))
	var ok bool
	for i, name := range e.Positionals.Vars {
		variables[i], ok = env.GetUserVariable(name)
		if !ok {
			return nil, errEnvInvalidVar.WithMessageF(name)
		}
	}
	return variables, nil
}
