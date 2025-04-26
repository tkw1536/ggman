//spellchecker:words path
package path

//spellchecker:words errors path filepath strings
import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Normalization describes the normalization of a path.
type Normalization int

const (
	// NoNorm uses the exact spelling of the original path.
	NoNorm Normalization = iota

	// It always uses the first matching subpath, regardless of how it matches.
	FoldNorm
	// FoldPreferExactNorm is like Fold, except that it uses the exact path if it already exists.
	FoldPreferExactNorm
)

// JoinNormalized is similar to filepath.Join, except that it takes an additional normalization argument.
//
// If any sub-path of the result path already has a sibling path that is case-fold-equal, uses this instead of creating a new path.
// When an exact match for the path exists, always uses the exact match.
func JoinNormalized(n Normalization, base string, elem ...string) (string, error) {
	if n == NoNorm {
		return filepath.Join(append([]string{base}, elem...)...), nil
	}
	return joinFold(n == FoldPreferExactNorm, base, elem...)
}

func joinFold(preferExact bool, base string, elem ...string) (string, error) {
	// figure out the tail
	base = filepath.Clean(base)
	elems := filepath.Join(elem...)

	// create a new builder that is big enough
	var builder strings.Builder
	builder.Grow(len(elems) + 1 + len(base))

	// write the base into it
	builder.WriteString(base)

	exists := true // does the current path exist?
	for elem := range strings.SplitSeq(elems, string(os.PathSeparator)) {
		// if elements up to this point exist, we need to list the directory
		// and then find a matching folded name!
		if exists {
			comp, err := FindFoldedDir(builder.String(), elem, preferExact)
			switch {
			case err == nil:
				elem = comp
			case errors.Is(err, fs.ErrNotExist):
				exists = false
			default: /* err != nil */
				return "", err
			}
		}

		// add the new path to the builder
		if builder.Len() > 1 { // if root path!
			builder.WriteRune(os.PathSeparator)
		}
		builder.WriteString(elem)
	}

	return builder.String(), nil
}

// FindFoldedDir lists directory dir and finds a path that is case-folded equal to query.
//
// By default, entries are sorted alphanumerically and the first matching directory entry is returned.
// When preferExact is true, an additional check is performed to return an exact (case-sensitive) match is returned.
//
// When no matching entry exists, returns os.ErrNoExist.
func FindFoldedDir(dir string, query string, preferExact bool) (name string, err error) {
	if preferExact {
		return findFoldedDirPreferExact(dir, query)
	} else {
		return findFoldedDirNoExact(dir, query)
	}
}

// findFoldedDirNoExact implements FindFoldedDir(dir, query, false).
func findFoldedDirNoExact(dir string, query string) (name string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("%q: failed to read directory: %w", dir, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		eName := entry.Name()
		if strings.EqualFold(eName, query) {
			return eName, nil
		}
	}

	return "", os.ErrNotExist
}

// findFoldedDirPreferExact implements FindFoldedDir(dir, query, true).
func findFoldedDirPreferExact(dir string, query string) (name string, err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("%q: failed to read directory: %w", dir, err)
	}

	// iterate over the entries and check both for exact matches and inexact matches
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		eName := entry.Name()
		if eName == query {
			return query, nil
		}

		if name == "" && strings.EqualFold(eName, query) {
			// we can't return here, cause there might be an exact match later!
			name = eName
		}
	}

	if name == "" {
		return "", os.ErrNotExist
	}
	return name, nil
}
