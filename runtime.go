package ggman

import (
	"github.com/tkw1536/ggman/env"
)

// NewRuntime makes a new runtime for ggman
func NewRuntime(params env.Parameters, context Context) (*env.Env, error) {
	// create a new environment
	e, err := env.NewEnv(context.Description.Requirements, params)
	if err != nil {
		return nil, err
	}

	// setup a filter for it!
	f, err := env.NewFilter(context.Args.Flags, e)
	if err != nil {
		return nil, err
	}
	e.Filter = f

	return e, nil

}
