package cmd

//spellchecker:words essio shellescape ggman goprogram exit pkglib collection
import (
	"fmt"

	"al.essio.dev/pkg/shellescape"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/collection"
)

//spellchecker:words positionals nolint wrapcheck

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
		Vars []string `description:"print only information about specified variables" positional-arg-name:"VAR"`
	} `positional-args:"true"`

	List     bool `description:"instead of \"name=value\" pairs print only the variable"                                           long:"list"     short:"l"`
	Describe bool `description:"instead of \"name=value\" pairs print \"name: description\" pairs describing the use of variables" long:"describe" short:"d"`
	Raw      bool `description:"instead of \"name=value\" pairs print only the unescaped value"                                    long:"raw"      short:"r"`
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

var (
	errEnvInvalidVar        = exit.NewErrorWithCode("unknown environment variable", exit.ExitCommandArguments)
	errEnvModesIncompatible = exit.NewErrorWithCode("at most one of `--raw`, `--list` and `--describe` may be given", exit.ExitCommandArguments)
)

func (e _env) AfterParse() error {
	// check that at most one mode was given
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
		return errEnvModesIncompatible
	}
	return nil
}

func (e _env) Run(context ggman.Context) error {
	variables, err := e.variables()
	if err != nil {
		return err
	}

	for _, v := range variables {
		switch {
		case e.List:
			if _, err := context.Println(v.Key); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		case e.Raw:
			if _, err := context.Println(v.Get(context.Environment, context.Program.Info)); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		case e.Describe:
			if _, err := context.Printf("%s: %s\n", v.Key, v.Description); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		default:
			value := shellescape.Quote(v.Get(context.Environment, context.Program.Info))
			if _, err := context.Printf("%s=%s\n", v.Key, value); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		}
	}

	return nil
}

func (e _env) variables() ([]env.UserVariable, error) {
	// no variables provided => use all of them
	if len(e.Positionals.Vars) == 0 {
		return env.GetUserVariables(), nil
	}

	var invalid string
	variables := collection.MapSlice(e.Positionals.Vars, func(name string) env.UserVariable {
		value, ok := env.GetUserVariable(name)
		if !ok && invalid == "" { // store an invalid name!
			invalid = name
		}
		return value
	})

	if invalid != "" {
		return nil, fmt.Errorf("%q: %w", invalid, errEnvInvalidVar)
	}

	return variables, nil
}
