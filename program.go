//spellchecker:words ggman
package ggman

//spellchecker:words github ggman constants goprogram exit meta
import (
	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/meta"
)

// info contains information about the ggman program
var info = meta.Info{
	BuildVersion: constants.BuildVersion,
	BuildTime:    constants.BuildTime,

	Executable:  "ggman",
	Description: "ggman manages local git repositories\n\nggman version " + constants.BuildVersion + "\nggman is licensed under the terms of the MIT License.\nuse 'ggman license' to view licensing information.",
}

// newEnvironment makes a new runtime for ggman
func newEnvironment(params env.Parameters, context Context) (env.Env, error) {
	// create a new environment
	e, err := env.NewEnv(context.Description.Requirements, params)
	if err != nil {
		return env.Env{}, err
	}

	// setup a filter for it!
	f, err := env.NewFilter(context.Args.Flags, &e)
	if err != nil {
		return e, err
	}
	e.Filter = f

	return e, nil

}

var errParseArgsNeedTwoAfterFor = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "unable to parse arguments: at least two arguments needed after `for' keyword",
}

// NewProgram returns a new ggman program
func NewProgram() (p Program) {
	p.NewEnvironment = newEnvironment
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
