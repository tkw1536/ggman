// Package url provides utilities for URL parsing
package url

import (
	"errors"
)

var errNoPort = errors.New("ParsePort: input is not a number")
var errInvalidRange = errors.New("ParsePort: port number out of range")

const maxValidPort = 65535         // maximal port (as a number)
const maxPortStr = "65535"         // maximal port (as a string)
const maxPortLen = len(maxPortStr) // maximal port length

// ParsePort parses a string into a valid port.
// A port is between 0 and 65535 (inclusive).
// It may not start with "+" and must only consist of digits.
//
// When a port can not be parsed, returns 0 and an error.
func ParsePort(s string) (uint16, error) {

	// when the input string is too long, we don't even need to try
	// parsding can just fail immediatly.
	if len(s) > maxPortLen {
		return 0, errInvalidRange
	}

	// an inlined version of ParseInt(s, 10, 16)

	// we first parse into a uint32
	// so that we can afterwards check for overflow

	var v uint32
	for _, ch := range []byte(s) {
		if '0' > ch || ch > '9' { // invalid digit
			return 0, errNoPort
		}
		v = v*10 + uint32(ch-'0')
	}

	if v > maxValidPort {
		return 0, errInvalidRange
	}

	return uint16(v), nil
}

// IsValidURLScheme checks if a string reprents a valid URL scheme.
//
// A valid url scheme is one that matches the regex [a-zA-Z][a-zA-Z0-9+\-\.]*.
func IsValidURLScheme(s string) bool {
	// An obvious implementation of this function would be:
	//   regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9+\-\.]*$").MatchString(s)
	// However such an implementation would be relatively slow when called a lot of times.
	// Instead we directly build code that implements this RegEx.
	//
	// For this we first check that the string is non-empty.
	// Then we check that the first character is not a '+', '-' or '.'
	// Then we check that each character is either alphanumeric or a '+', '-' or '.'

	if len(s) == 0 {
		return false
	}

	// if the first rune is invalid, no need to check everything else
	if s[0] == '+' || s[0] == '-' || s[0] == '.' {
		return false
	}

	for _, r := range s {
		if !(('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')) && // more efficient version of !@unicode.isLetter(r)
			r != '+' &&
			r != '-' &&
			r != '.' {
			return false
		}
	}

	return true
}
