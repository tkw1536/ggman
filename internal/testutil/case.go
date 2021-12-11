package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// CaseSensitive checks if temporary directories exist on a case-sensitive file system.
//
// This function is untested due to unpredictability of runtime environment.
func CaseSensitive(t *testing.T) bool {
	temp := t.TempDir()

	// create lower case
	lower := filepath.Join(temp, "test")
	if err := os.Mkdir(lower, os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}

	upper := filepath.Join(temp, "TEST")
	err := os.Mkdir(upper, os.ModeDir|os.ModePerm)

	switch err {
	case os.ErrExist:
		return false // directory already exists!
	case nil:
		return true // both were created ok
	}
	panic(err)
}
