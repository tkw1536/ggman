package util

import (
	"errors"
	"strconv"
)

var errNoPlus = errors.New("ParsePort: s may not start with '+'")
var errInvalidRange = errors.New("ParsePort: s out of range")

// ParsePort parses a string into a valid port.
// A port is between 0 and 65535 (inclusive).
// It may not start with "+" and must only consist of digits.
//
// When a port can not be parsed, returns 0 and an error.
func ParsePort(s string) (uint16, error) {
	// This function could use strconv.ParseUint(s, 10, 16) and then do a conversion on the result.
	// Instead we call strconv.AtoI and check that the result is in the right range.
	//
	// In Benchmarking it turned out that this implementation is about 33% faster for common port usecases.
	// This is likely because strconv.AtoI has a 'fast' path for small integers, and most port candidates are small enough to fit.

	if len(s) > 0 && s[0] == '+' {
		return 0, errNoPlus
	}
	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, err
	}
	if port < 0 || port > 65535 {
		return 0, errInvalidRange
	}
	return uint16(port), nil
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
