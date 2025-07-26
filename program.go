//spellchecker:words ggman
package ggman

//spellchecker:words ggman constants goprogram exit meta
import (
	"fmt"

	"go.tkw01536.de/ggman/constants"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/goprogram/meta"
)

// TODO: Rework this once we have fully ported

// info contains information about the ggman program.
var info = meta.Info{
	BuildVersion: constants.BuildVersion,
	BuildTime:    constants.BuildTime,

	Executable:  "ggman",
	Description: "ggman manages local git repositories\n\nggman version " + constants.BuildVersion + "\nggman is licensed under the terms of the MIT License.\nuse 'ggman license' to view licensing information.",
}

// newEnvironment makes a new runtime for ggman.
func newEnvironmentLegacy(params env.Parameters, context Context) (env.Env, error) {
	return newEnvironment(context.Description.Requirements, params, context.Args.Flags)
}

func newEnvironment(requirements env.Requirement, params env.Parameters, flags env.Flags) (env.Env, error) {
	// create a new environment
	e, err := env.NewEnv(requirements, params)
	if err != nil {
		return env.Env{}, fmt.Errorf("error creating env: %w", err)
	}

	// setup a filter for it!
	f, err := env.NewFilter(flags, &e)
	if err != nil {
		return e, fmt.Errorf("error creating filter: %w", err)
	}
	e.Filter = f

	return e, nil
}

var errParseArgsNeedTwoAfterFor = exit.NewErrorWithCode("unable to parse arguments: at least two arguments needed after `for' keyword", exit.ExitGeneralArguments)

// NewProgram returns a new ggman program.
func NewProgram() (p Program) {
	p.NewEnvironment = newEnvironmentLegacy
	p.Info = info

	p.RegisterKeyword("help", func(args *Arguments, pos *[]string) error {
		args.Command = ""
		args.Universals.Help = true
		return nil
	})

	p.RegisterKeyword("version", func(args *Arguments, pos *[]string) error {
		args.Command = ""
		args.Universals.Version = true
		return nil
	})

	p.RegisterKeyword("for", func(args *Arguments, pos *[]string) error {
		if len(*pos) < 2 {
			return errParseArgsNeedTwoAfterFor
		}
		args.Flags.For = append(args.Flags.For, (*pos)[0])
		args.Command = (*pos)[1]
		*pos = (*pos)[2:]
		return nil
	})

	return
}
