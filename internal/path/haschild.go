// Package path provides path-based functions.
//
//spellchecker:words path
package path

//spellchecker:words path filepath
import (
	"os"
	"path/filepath"
)

// HasChild checks, using lexical analysis only, if the path 'parent' could contain the path 'child'.
// NOTE: HasChild does not take filesystem case sensitivity or symlink-equivalence into account.
func HasChild(parent string, child string) bool {
	parent, child = normSep(filepath.Clean(parent)), normSep(filepath.Clean(child)) // clean + normalize paths
	return child == parent || (len(child) > len(parent) && os.IsPathSeparator(child[len(parent)]) && child[:len(parent)] == parent)
}

// normSep uses a normalized separator for all paths.
func normSep(path string) string {
	// check if there are any characters that need to be replaced.
	// If there are none, return the path as is.
	{
		var changed bool

		for _, r := range path {
			if IsPathSeparator(r) && r != os.PathSeparator {
				changed = true
				break
			}
		}

		if !changed {
			return path
		}
	}

	// do the actual normalization!
	ret := []rune(path)
	for i, r := range ret {
		if IsPathSeparator(r) {
			ret[i] = os.PathSeparator
		}
	}

	// and make it a string again
	return string(ret)
}

// IsPathSeparator reports whether r is a directory separator character.
// See also "os".IsPathSeparator.
func IsPathSeparator(r rune) bool {
	// NOTE: This function is untested

	// check that it is in bounds
	return 0 <= r && r <= 255 && os.IsPathSeparator(uint8(r))
}
