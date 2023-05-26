package path

import (
	"os"
	"path/filepath"
	"strings"
)

// Separator contains filepath.Separator as a string
const Separator = string(filepath.Separator)

// defaultVolumePrefix is a prefix to use for the default volume
var defaultVolumePrefix = filepath.VolumeName(os.TempDir())

// ToOSPath turns a path that is separated via "/"s into a path separated by the current os-separator.
// ToOSPath should only be used during testing.
//
// When path starts with "/", the default volume is prefixed on windows.
// The default volume is defined as the volume the temporary directory resides on.
// On Windows, this is usually the 'C:' volume, but not guaranteed to be so.
func ToOSPath(path string) (result string) {
	if len(path) > 0 && path[0] == '/' {
		result = defaultVolumePrefix
	}
	result += strings.ReplaceAll(path, "/", Separator)
	return
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
