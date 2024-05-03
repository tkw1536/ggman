//spellchecker:words testutil
package testutil

//spellchecker:words errors path filepath testing github pkglib testlib
import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/pkglib/testlib"
)

// CaseSensitive checks if temporary directories exist on a case-sensitive file system.
//
// This function is untested due to unpredictability of runtime environment.
func CaseSensitive(t *testing.T) bool {
	t.Helper()
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
