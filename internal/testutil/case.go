package testutil

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/goprogram/lib/testlib"
)

// CaseSensitive checks if temporary directories exist on a case-sensitive file system.
//
// This function is untested due to unpredictability of runtime environment.
func CaseSensitive(t *testing.T) bool {
	temp := testlib.TempDirAbs(t)

	// create lower case
	lower := filepath.Join(temp, "test")
	if err := os.Mkdir(lower, os.ModeDir|os.ModePerm); err != nil {
		panic(err)
	}

	upper := filepath.Join(temp, "TEST")
	err := os.Mkdir(upper, os.ModeDir|os.ModePerm)

	switch {
	case errors.Is(err, fs.ErrExist):
		return false // directory already exists!
	case err == nil:
		return true // both were created ok
	}
	panic(err)
}
