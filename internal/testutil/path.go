package testutil

import (
	"os"
	"path/filepath"
	"strings"
)

// defaultVolumeName is a consistent volume name to prefix to paths.
//
// We use the volume that the temporary directory resides on.
var defaultVolumeName = filepath.VolumeName(os.TempDir())

// ToOSPath turns a path that is separated via "/"s into a path separated by the current os-separator.
//
// When path starts with "/", the path is guaranteed to contain a volume name.
func ToOSPath(path string) (result string) {
	if path == "" {
		return ""
	}

	parts := strings.Split(path, "/")
	if parts[0] == "" {
		parts[0] = defaultVolumeName
	}
	return strings.Join(parts, string(os.PathSeparator))
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
