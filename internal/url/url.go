// Package url provides utilities for URL parsing
package url

import (
	"errors"
	"math"
)

var (
	errNotANumber = errors.New("ParsePort: not a valid port number")
	errOutOfRange = errors.New("ParsePort: port not in range")
)

const (
	maxValid      = math.MaxUint16  // maximal port (as a number)
	maxMultiply10 = maxValid / 10   // maximum value before multiplication by 10
	maxPortStr    = "65535"         // maximal port (as a string)
	maxPortLen    = len(maxPortStr) // maximal port length
)

// ParsePort parses a string into a valid port.
// A port is between 0 and 65535 (inclusive).
// It may not start with "+" and must only consist of digits.
//
// When a port can not be parsed, returns 0 and an error.
func ParsePort(s string) (v uint16, err error) {

	// hadOverflow determines if we had an overflow
	hadOverflow := false

	for _, ch := range []byte(s) {
		if '0' > ch || ch > '9' { // not a digit!
			return 0, errNotANumber
		}

		// if we had an overflow, we just check for digits
		if hadOverflow {
			continue
		}

		// multiply the previous digits by 10
		if v > maxMultiply10 {
			hadOverflow = true
			continue
		}
		v = 10 * v

		// add the current digit
		digit := uint16(ch - '0')
		if v > maxValid-digit {
			hadOverflow = true
			continue
		}
		v += digit
	}

	// we had an overflow
	if hadOverflow {
		return 0, errOutOfRange
	}

	return uint16(v), nil
}

// IsValidURLScheme checks if a string represents a valid URL scheme.
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
