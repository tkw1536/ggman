// Package url provides utilities for URL parsing
package url

//spellchecker:words errors math
import (
	"errors"
	"math"
)

var (
	errNotANumber = errors.New("ParsePort: not a valid port number")
	errOutOfRange = errors.New("ParsePort: port not in range")
)

const (
	maxValid      = math.MaxUint16      // maximal port (as a number)
	maxMultiply10 = maxValid / 10       // maximum value before multiplication by 10
	maxPortStr    = "65535"             // maximal port (as a string)
	maxPortLen    = len(maxPortStr)     // maximal port length
	maxPortIndex  = len(maxPortStr) - 1 // maximal valid index
)

// ParsePort parses a string into a valid port.
// A port is between 0 and 65535 (inclusive).
// It may not start with "+", may not be empty, and must only consist of digits.
//
// When a port can not be parsed, returns 0 and an error.
func ParsePort(s string) (port uint16, err error) {
	if s == "" {
		return 0, errNotANumber
	}

	for i, ch := range []byte(s) {
		// determine the current digit
		digit := uint16(ch - '0')
		if digit > 9 { // not a digit!
			return 0, errNotANumber
		}

		// prior to the last digit, we can just add the digit.
		// because it's guaranteed to be in bounds!
		if i < maxPortIndex {
			port = 10 * port
			port += digit
			continue
		}

		// we're at (or beyond) the last index
		// again add the digit, but do a bounds check!
		if port > maxMultiply10 {
			return 0, errOutOfRange
		}
		port = 10 * port

		if port > maxValid-digit {
			return 0, errOutOfRange
		}
		port += digit
	}

	// and we're done!
	return port, nil
}

// IsValidURLScheme checks if a string represents a valid URL scheme.
//
// A valid url scheme is one that matches the regex [a-zA-Z][a-zA-Z0-9+\-\.]*.
func IsValidURLScheme(s string) bool {
	// An obvious implementation of this function would be:
	//   regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+\-\.]*$`).MatchString(s)
	// However such an implementation would be a lot slower when called repeatedly.
	// Instead we directly build code that implements this regex.
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
		if !(('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')) && // more efficient version of !@unicode.isLetter(r) for ascii
			r != '+' &&
			r != '-' &&
			r != '.' {
			return false
		}
	}

	return true
}
