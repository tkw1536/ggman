package ggman

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// these define the ggman-specific program types
// none of these are strictly needed, they're just around for convenience
type runtimeT = *env.Env

type Program = program.Program[runtimeT]
type Command = program.Command[runtimeT]
type Context = program.Context[runtimeT]
type CommandArguments = program.CommandArguments[runtimeT]

// info contains information about the ggman program
var info = program.Info{
	BuildVersion: constants.BuildVersion,
	BuildTime:    constants.BuildTime,

	MainName:    "ggman",
	Description: fmt.Sprintf("ggman manages local git repositories\n\nggman version %s\nggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information.", constants.BuildVersion),
}

// NewProgram returns a new ggman program
func NewProgram() (p program.Program[*env.Env]) {
	p.Initalizer = func(params env.EnvironmentParameters, cmdargs CommandArguments) (*env.Env, error) {
		rt, err := NewRuntime(params, cmdargs)
		return rt, err
	}
	p.Info = info

	return
}
