package repos

import (
	"strings"

	"github.com/danwakefield/fnmatch"
)

// Matches checks if a string matches a given repository pattern
func Matches(pattern string, s string) bool {
	// if we have an 'everything' pattern, return true immediatly
	if pattern == "" || pattern == "*" {
		return true
	}

	// get components of the input string that might match
	curi, es := NewRepoURI(s)
	if es != nil {
		return false
	}
	components := curi.Components()
	componentsLength := len(components)

	// parse components of strings and ignore any casing

	puri, es := NewRepoURI(pattern)
	if es != nil {
		return false
	}

	patternComponents := puri.Components()
	patternLength := len(patternComponents)
	thePattern := strings.Join(patternComponents, "/")

	// try and match all te sub patterns
	for i := 0; i <= componentsLength-patternLength; i++ {
		subString := strings.Join(components[i:i+patternLength], "/")
		if fnmatch.Match(thePattern, subString, fnmatch.FNM_IGNORECASE) {
			return true
		}
	}

	return false
}
