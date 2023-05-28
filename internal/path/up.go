// Package path provides path-based utility functions
package path

import (
	"os"
	"path/filepath"
	"strings"
)

// IsPathSeparator is like [os.IsPathSeparator], but takes a rune.
func IsPathSeparator(r rune) bool {
	// NOTE(twiesing): This function is untested
	return 0 <= r && r <= 255 && os.IsPathSeparator(uint8(r))
}

// GoesUp checks if path is a path that goes up at least one directory.
func GoesUp(path string) bool {
	path = filepath.Clean(path)
	path = strings.TrimPrefix(path, filepath.VolumeName(path))
	path = strings.TrimLeftFunc(path, IsPathSeparator) // trim separators (for malformed paths like "/../path")
	return hasPrefix(path, "..")
}

// Contains checks if the path 'parent' could contain the path 'child' syntactically.
// No resolving of paths is performed.
func Contains(parent string, child string) bool {
	return hasPrefix(useStandardSep(filepath.Clean(child)), useStandardSep(filepath.Clean(parent)))
}

// hasPrefix checks if path starts with the given prefix.
// Both path and prefix are expected to be clean()ed
func hasPrefix(path string, prefix string) bool {
	return path == prefix || (len(path) > len(prefix) && strings.HasPrefix(path, prefix) && os.IsPathSeparator(path[len(prefix)]))
}

// useStandardSep replaces all non-default separators by os.PathSeparator
func useStandardSep(path string) string {
	// create a new string
	var buffer strings.Builder
	buffer.Grow(len(path))

	// by iterating over each character
	for _, r := range path {
		if IsPathSeparator(r) && r != os.PathSeparator {
			buffer.WriteRune(os.PathSeparator)
			continue
		}
		buffer.WriteRune(r)
	}

	return buffer.String()
}
