package gg

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// NewProgram returns a new ggman program
func NewProgram() (p program.Program) {
	p.Initalizer = func(params env.EnvironmentParameters, cmdargs program.CommandArguments) (program.Runtime, error) {
		rt, err := NewRuntime(params, cmdargs)
		return rt, err
	}
	return
}

// C2E returns the environment belonging to a context.
// It should be replaced by a type parameter in the future.
func C2E(context program.Context) *env.Env {
	return context.Runtime().(*env.Env)
}
