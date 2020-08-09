package repos

import (
	"regexp"
	"strings"

	"github.com/tkw1536/ggman/util"
)

var specReplace = regexp.MustCompile("[\\^\\%]")

// Canonical returns the canonical version of this URI given a canonical specification
// the canonical specification can contain any character, except for three special ones
// ^ -- replaced by the first un-used component of the URI
// % -- replaced by the second un-used component of the URI (commonly username)
// $ -- replaced by all remaining components in the URI joined with a '/'. Also stops all processing afterwards
// If $ does not exist in the cspec, it is assumed to be at the end of the cspec.
func (rURI *RepoURL) Canonical(cspec string) (canonical string) {
	// get the components of the URI
	components := rURI.Components()

	// split into prefix and suffix
	prefix, suffix := util.SplitAfter(cspec, "$")

	prefix = specReplace.ReplaceAllStringFunc(prefix, func(s string) string {
		// if everything is empty, return the string as is
		if len(components) == 0 {
			return s
		}

		// replace by the first component
		if s == "^" {
			defer func() { components = components[1:] }()
			return components[0]
		}

		// else we want to replace by the second component
		// so we need to make sure we have that many
		if len(components) < 2 {
			return s
		}

		// do the replacement
		defer func() { components = append(components[:1], components[2:]...) }()
		return components[1]
	})

	// add the remaining components
	return prefix + strings.Join(components, "/") + suffix
}

// CanonicalWith returns the canonical url given a set of lines
func (rURI *RepoURL) CanonicalWith(lines CanFile) (canonical string) {
	for _, line := range lines {
		if rURI.Matches(line.Pattern) {
			return rURI.Canonical(line.Canonical)
		}
	}

	return
}

// IsCanonicalURI checks if a given URI is in canonical form
// using a specific canonical specification
func IsCanonicalURI(s string, cspec string) bool {
	return ParseRepoURL(s).Canonical(cspec) == s
}
