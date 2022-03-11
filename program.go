package ggman

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/meta"
)

// these define the ggman-specific program types
// none of these are strictly needed, they're just around for convenience
type ggmanEnv = env.Env
type ggmanParameters = env.Parameters
type ggmanRequirements = env.Requirement
type ggmanFlags = env.Flags

type Program = goprogram.Program[ggmanEnv, ggmanParameters, ggmanFlags, ggmanRequirements]
type Command = goprogram.Command[ggmanEnv, ggmanParameters, ggmanFlags, ggmanRequirements]
type Context = goprogram.Context[ggmanEnv, ggmanParameters, ggmanFlags, ggmanRequirements]
type Arguments = goprogram.Arguments[ggmanFlags]
type Description = goprogram.Description[ggmanFlags, ggmanRequirements]

// info contains information about the ggman program
var info = meta.Info{
	BuildVersion: constants.BuildVersion,
	BuildTime:    constants.BuildTime,

	Executable:  "ggman",
	Description: fmt.Sprintf("ggman manages local git repositories\n\nggman version %s\nggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information.", constants.BuildVersion),
}

var errParseArgsNeedTwoAfterFor = exit.Error{
	ExitCode: exit.ExitGeneralArguments,
	Message:  "Unable to parse arguments: At least two arguments needed after 'for' keyword. ",
}

// NewProgram returns a new ggman program
func NewProgram() (p Program) {
	p.NewEnvironment = NewRuntime
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
			return errParseArgsNeedTwoAfterFor
		}
		args.Flags.Filters = append(args.Flags.Filters, args.Pos[0])
		args.Command = args.Pos[1]
		args.Pos = args.Pos[2:]

		return nil
	})

	return
}
