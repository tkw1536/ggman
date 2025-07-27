package env

//spellchecker:words context github cobra
import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

type cobraKey int

const (
	flagsKey cobraKey = iota
	parametersKey
)

// GetEnv gets the environment for a specific command.
// [SetFlags] and [SetParameters] must have been called.
func GetEnv(cmd *cobra.Command, requirements Requirement) (Env, error) {
	// create a new environment
	ne, err := NewEnv(requirements, get[Parameters](cmd, parametersKey))
	if err != nil {
		return Env{}, fmt.Errorf("error creating environment: %w", err)
	}

	// setup a filter for it!
	f, err := NewFilter(get[Flags](cmd, flagsKey), &ne)
	if err != nil {
		return ne, fmt.Errorf("error creating filter: %w", err)
	}
	ne.Filter = f

	// and return
	return ne, nil
}

// SetFlags sets the value for a cobra command from a set of flags.
func SetFlags(cmd *cobra.Command, flags *Flags) {
	set(cmd, flagsKey, flags)
}

// SetParameters sets parameters for a cobra command.
func SetParameters(cmd *cobra.Command, params *Parameters) {
	set(cmd, parametersKey, params)
}

func set[T any](cmd *cobra.Command, key cobraKey, data *T) {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	cmd.SetContext(context.WithValue(ctx, key, data))
}

func get[T any](cmd *cobra.Command, key cobraKey) T {
	flags := cmd.Context().Value(key)
	data, ok := flags.(*T)
	if !ok {
		var zero T
		return zero
	}
	return *data
}
