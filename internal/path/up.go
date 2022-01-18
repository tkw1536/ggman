// Package path provides path-based utility functions
package path

import (
	"path/filepath"
	"strings"
)

// pathUpSegment is a segment that indiciates that a path should go up
var pathUpSegment = ".." + pathSeperator

// GoesUp checks if path is a path that goes up at least one directory.
func GoesUp(path string) bool {
	return path == ".." || strings.HasPrefix(ToOSPath(path), pathUpSegment)
}

// Contains checks if the path 'parent' could contain the path 'child' synatically.
// No resolving of paths is performed.
func Contains(parent string, child string) bool {
	p := filepath.Clean(ToOSPath(parent))
	c := filepath.Clean(ToOSPath(child))
	return strings.HasPrefix(c, p)
}
