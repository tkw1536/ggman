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

	// Iterate over the bytes in this string.
	// We can do this safely, because any valid scheme must be ascii.
	bytes := []byte(s)
	if len(bytes) == 0 {
		return false
	}

	var (
		idx = 0        // current index
		c   = bytes[0] // current character
	)

	// scheme must start with a letter
	// so omit the loop preamble
	goto start

nextLetter:
	// go to the next letter
	// or be done
	idx++
	if idx >= len(bytes) {
		return true
	}

	// get the current letter
	c = bytes[idx]

	if '0' <= c && c <= '9' {
		goto nextLetter
	}

	if c == '+' || c == '-' || c == '.' {
		goto nextLetter
	}

start:
	if 'a' <= c && c <= 'z' {
		goto nextLetter
	}

	if 'A' <= c && c <= 'Z' {
		goto nextLetter
	}

	return false
}
