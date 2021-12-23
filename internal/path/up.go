// Package path provides path-based utility functions
package path

import (
	"strings"
)

// pathUpSegment is a segment that indiciates that a path should go up
var pathUpSegment = ".." + pathSeperator

// GoesUp checks if path is a path that goes up at least one directory.
func GoesUp(path string) bool {
	return path == ".." || strings.HasPrefix(ToOSPath(path), pathUpSegment)
}
