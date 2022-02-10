package ggman

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// these define the ggman-specific program types
// none of these are strictly needed, they're just around for convenience
type ggmanRuntime = *env.Env
type ggmanParameters = env.EnvironmentParameters
type ggmanRequirements = env.Requirement

type Program = program.Program[ggmanRuntime, ggmanParameters, ggmanRequirements]
type Command = program.Command[ggmanRuntime, ggmanParameters, ggmanRequirements]
type Context = program.Context[ggmanRuntime, ggmanParameters, ggmanRequirements]
type CommandArguments = program.CommandArguments[ggmanRuntime, ggmanParameters, ggmanRequirements]
type Description = program.Description[ggmanRequirements]

// info contains information about the ggman program
var info = program.Info{
	BuildVersion: constants.BuildVersion,
	BuildTime:    constants.BuildTime,

	MainName:    "ggman",
	Description: fmt.Sprintf("ggman manages local git repositories\n\nggman version %s\nggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information.", constants.BuildVersion),
}

// NewProgram returns a new ggman program
func NewProgram() (p Program) {
	p.NewRuntime = func(params env.EnvironmentParameters, cmdargs CommandArguments) (*env.Env, error) {
		rt, err := NewRuntime(params, cmdargs)
		return rt, err
	}
	p.Info = info

	return
}
