package ggman

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/program/exit"
)

// these define the ggman-specific program types
// none of these are strictly needed, they're just around for convenience
type ggmanRuntime = *env.Env
type ggmanParameters = env.EnvironmentParameters
type ggmanRequirements = env.Requirement
type ggmanFlags = env.Flags

type Program = program.Program[ggmanRuntime, ggmanParameters, ggmanFlags, ggmanRequirements]
type Command = program.Command[ggmanRuntime, ggmanParameters, ggmanFlags, ggmanRequirements]
type Context = program.Context[ggmanRuntime, ggmanParameters, ggmanFlags, ggmanRequirements]
type CommandArguments = program.CommandArguments[ggmanRuntime, ggmanParameters, ggmanFlags, ggmanRequirements]
type Arguments = program.Arguments[ggmanFlags]
type Description = program.Description[ggmanFlags, ggmanRequirements]

// info contains information about the ggman program
var info = program.Info{
	BuildVersion: constants.BuildVersion,
	BuildTime:    constants.BuildTime,

	MainName:    "ggman",
	Description: fmt.Sprintf("ggman manages local git repositories\n\nggman version %s\nggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information.", constants.BuildVersion),
}

var ErrParseArgsNeedTwoAfterFor = exit.Error{ // TODO: Public because test
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: At least two arguments needed after 'for' keyword. ",
}

// NewProgram returns a new ggman program
func NewProgram() (p Program) {
	p.NewRuntime = func(params env.EnvironmentParameters, cmdargs CommandArguments) (*env.Env, error) {
		rt, err := NewRuntime(params, cmdargs)
		return rt, err
	}
	p.Info = info

	p.RegisterKeyword("help", func(args *Arguments) error {
		args.Command = ""
		args.Universals.Help = true
		return nil
	})

	p.RegisterKeyword("version", func(args *Arguments) error {
		args.Command = ""
		args.Universals.Version = true
		return nil
	})

	p.RegisterKeyword("for", func(args *Arguments) error {
		if len(args.Pos) < 2 {
			return ErrParseArgsNeedTwoAfterFor
		}
		args.Flags.Filters = append(args.Flags.Filters, args.Pos[0])
		args.Command = args.Pos[1]
		args.Pos = args.Pos[2:]

		return nil
	})

	return
}
