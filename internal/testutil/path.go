//spellchecker:words testutil
package testutil

//spellchecker:words path filepath
import (
	"os"
	"path/filepath"
)

// defaultVolumeName is a consistent volume name to prefix to paths.
//
// We use the volume that the temporary directory resides on.
//
// NOTE: When updating this, also update path_test.go.
var defaultVolumeName = filepath.VolumeName(os.TempDir())

// ToOSPath turns a path that is separated via "/"s into a path separated by the current os-separator.
//
// When path starts with "/", the path is guaranteed to contain a volume name.
func ToOSPath(path string) (result string) {
	path = filepath.FromSlash(path)
	if len(path) > 0 && os.IsPathSeparator(path[0]) {
		return defaultVolumeName + path
	}
	return path
}

// ToOSPaths is like ToOSPath, but applies to each value in a slice or array.
// ToOSPaths modifies the slice in-place and returns it for convenience.
//
// This function is untested.
func ToOSPaths(paths []string) []string {
	for i := range paths {
		paths[i] = ToOSPath(paths[i])
	}
	return paths
}
