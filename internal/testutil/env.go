//spellchecker:words testutil
package testutil

import (
	"fmt"
	"os"
)

// MockVariables sets environment variables as configured in map.
// It returns a function that can be used to revert the environment variables to their previous values.
//
// It should be called like:
//
//	defer MockEnv(values)()
func MockVariables(values map[string]string) (revert func()) {
	originals := make(map[string]string, len(values))
	for k, v := range values {
		originals[k] = os.Getenv(k)
		if err := os.Setenv(k, v); err != nil {
			panic(fmt.Errorf("failed to set variable %q: %w", k, err))
		}
	}
	return func() {
		for k, v := range originals {
			if err := os.Setenv(k, v); err != nil {
				panic(fmt.Errorf("failed to set variable %q: %w", k, err))
			}
		}
	}
}
