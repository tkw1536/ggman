//spellchecker:words ggman
package ggman

//spellchecker:words ggman
import (
	"fmt"

	"go.tkw01536.de/ggman/env"
)

// TODO: Rework this once we have fully ported

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
