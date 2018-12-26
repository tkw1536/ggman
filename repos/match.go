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
	components, es := Components(s)
	if es != nil {
		return false
	}
	componentsLength := len(components)

	// parse components of strings and ignore any casing
	var patternComponents []string
	if strings.Contains(pattern, ":") {
		patternComponents, es = Components(pattern)
		if es != nil {
			return false
		}
	} else {
		patternComponents, es = Components(":" + pattern)
		if es != nil {
			return false
		}
		patternComponents = patternComponents[1:]
	}

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
