package ggman

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
	program "github.com/tkw1536/ggman/goprogram"
	"github.com/tkw1536/ggman/goprogram/exit"
	"github.com/tkw1536/ggman/goprogram/meta"
)

// these define the ggman-specific program types
// none of these are strictly needed, they're just around for convenience
type ggmanRuntime = *env.Env
type ggmanParameters = env.Parameters
type ggmanRequirements = env.Requirement
type ggmanFlags = env.Flags

type Program = program.Program[ggmanRuntime, ggmanParameters, ggmanFlags, ggmanRequirements]
type Command = program.Command[ggmanRuntime, ggmanParameters, ggmanFlags, ggmanRequirements]
type Context = program.Context[ggmanRuntime, ggmanParameters, ggmanFlags, ggmanRequirements]
type Arguments = program.Arguments[ggmanFlags]
type Description = program.Description[ggmanFlags, ggmanRequirements]

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
