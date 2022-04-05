package cmd

import (
	"fmt"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/meta"
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
//  --list
// Instead of printing "name=value" pairs, print only the name.
//
//  --describe
// Instead of printing "name=value" pairs, print "name: description" pairs.
// The description explains what the value does.
//
//  --raw
// Instead of printing "name=value" pairs, print only the raw, unescaped value.
var Env ggman.Command = &_env{}

type _env struct {
	info meta.Info

	List     bool `short:"l" long:"list" description:"Instead of 'name=value' pairs print only the variable"`
	Describe bool `short:"d" long:"describe" description:"Instead of 'name=value' pairs print 'name: description' pairs describing the use of variables"`
	Raw      bool `short:"r" long:"raw" description:"Instead of 'name=value' pairs print only the unescaped value"`
}

func (e *_env) BeforeRegister(program *ggman.Program) {
	e.info = program.Info
}

func (e _env) Description() ggman.Description {
	uvs := env.GetUserVariables()
	keys := make([]string, len(uvs))
	for i, uv := range uvs {
		keys[i] = "'" + uv.Key + "'"
	}
	key_choices := strings.Join(keys, ", ")

	return ggman.Description{
		Command:     "env",
		Description: "Print information about the ggman environment",

		Positional: meta.Positional{
			Min: 0,
			Max: -1,

			Value:       "VAR",
			Description: fmt.Sprintf("Print only information about specified variables. One of %s, matched without case-sensitivity. ", key_choices),
		},

		Requirements: env.Requirement{
			NeedsRoot: true,
		},
	}
}

func (_env) AfterParse() error {
	return nil
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
			context.Println(v.Get(context.Environment, e.info))
		case e.Describe:
			context.Printf("%s: %s\n", v.Key, v.Description)
		default:
			value := shellescape.Quote(v.Get(context.Environment, e.info))
			context.Printf("%s=%s\n", v.Key, value)
		}
	}

	return nil
}

func (e _env) variables(context ggman.Context) ([]env.UserVariable, error) {
	// no variables provided => use all of them
	if len(context.Args.Pos) == 0 {
		return env.GetUserVariables(), nil
	}

	// names provided => make sure each exists
	variables := make([]env.UserVariable, len(context.Args.Pos))
	var ok bool
	for i, name := range context.Args.Pos {
		variables[i], ok = env.GetUserVariable(name)
		if !ok {
			return nil, errEnvInvalidVar.WithMessageF(name)
		}
	}
	return variables, nil
}
