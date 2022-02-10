package ggman

import (
	"fmt"

	"github.com/tkw1536/ggman/constants"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// info contains information about the ggman program
var info = program.Info{
	BuildVersion: constants.BuildVersion,
	BuildTime:    constants.BuildTime,

	MainName:    "ggman",
	Description: fmt.Sprintf("ggman manages local git repositories\n\nggman version %s\nggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information.", constants.BuildVersion),
}

// NewProgram returns a new ggman program
func NewProgram() (p program.Program) {
	p.Initalizer = func(params env.EnvironmentParameters, cmdargs program.CommandArguments) (program.Runtime, error) {
		rt, err := NewRuntime(params, cmdargs)
		return rt, err
	}
	p.Info = info

	return
}

// C2E returns the environment belonging to a context.
// TODO: Type parameter
func C2E(context program.Context) *env.Env {
	return context.Runtime().(*env.Env)
}
