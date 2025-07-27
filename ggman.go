// Package ggman serves as the implementation of the ggman program.
// See documentation of the ggman command as an entry point into the documentation.
//
// Note that this package and it's sub-packages are not intended to be consumed by other go packages.
// The public interface of the ggman is defined only by the ggman executable.
// This package is not considered part of the public interface as such and not subject to Semantic Versioning.
//
// The top-level ggman package is considered to be stand-alone, and (with the exception of 'env') does not directly depend on any of its' sub-packages.
// As such it can be safely used by any sub-package without cyclic imports.
//
//spellchecker:words ggman
package ggman

//spellchecker:words context github cobra ggman goprogram exit
import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
)

type cobraKey int

const (
	flagsKey cobraKey = iota
	envKey
	parametersKey
)

// SetFlags sets the flags for a cobra.Command.
func SetFlags(cmd *cobra.Command, flags *env.Flags) {
	setType(cmd, flagsKey, flags)
}

// GetFlags gets the flags from a cobra command.
func GetFlags(cmd *cobra.Command) env.Flags {
	return getType[env.Flags](cmd, flagsKey)
}

func SetParameters(cmd *cobra.Command, params *env.Parameters) {
	setType(cmd, parametersKey, params)
}
func GetParameters(cmd *cobra.Command) env.Parameters {
	return getType[env.Parameters](cmd, parametersKey)
}

func setType[T any](cmd *cobra.Command, key cobraKey, data *T) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cmd.SetContext(context.WithValue(ctx, key, data))
}

func getType[T any](cmd *cobra.Command, key cobraKey) T {
	flags := cmd.Context().Value(key)
	data, ok := flags.(*T)
	if !ok {
		var zero T
		return zero
	}
	return *data
}

var ErrGenericEnvironment = exit.NewErrorWithCode("failed to initialize environment", env.ExitInvalidEnvironment)

// GetEnv gets the environment of a ggman command.
func GetEnv(cmd *cobra.Command, requirements env.Requirement) (env.Env, error) {
	flags := cmd.Context().Value(envKey)
	data, ok := flags.(env.Env)
	if ok {
		return data, nil
	}

	// make a new environment
	ne, err := newEnvironment(requirements, GetParameters(cmd), GetFlags(cmd))
	if err != nil {
		return env.Env{}, fmt.Errorf("%w: %w", ErrGenericEnvironment, err)
	}

	// store the environment for future usage
	cmd.SetContext(context.WithValue(cmd.Context(), envKey, ne))
	return ne, nil
}

// ErrGenericOutput indicates that a generic output error occurred.
var ErrGenericOutput = exit.NewErrorWithCode("unknown output error", exit.ExitGeneric)
