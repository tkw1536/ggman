package util

import (
	"os"
	"path/filepath"
	"strings"
)

var pathSeperator = string(filepath.Separator)

// defaultVolumePrefix is a prefix to use for the default volume
var defaultVolumePrefix = filepath.VolumeName(os.TempDir())

// ToOSPath turns a path that is seperated via "/"s into a path seperated by the current os-seperator.
//
// When path starts with "/", the default volume is prefixed on windows.
// The default volume is defined as the volume the temporary directory resides on.
// On Windows, this is usually the 'C:' volume, but not guaranteed to be so.
func ToOSPath(path string) (result string) {
	if len(path) > 0 && path[0] == '/' {
		result = defaultVolumePrefix
	}
	result += strings.ReplaceAll(path, "/", pathSeperator)
	return
}

// ToOSPaths is like ToOSPath, but applies to each value in a slice or array.
// ToOSPaths modifies the slice in-place.
//
// This function is untested.
func ToOSPaths(paths []string) {
	for i := range paths {
		paths[i] = ToOSPath(paths[i])
	}
}
