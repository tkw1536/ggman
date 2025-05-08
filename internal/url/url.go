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
	if s == "" || len(s) > maxPortLen {
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

// SplitURLScheme splits off the URL scheme from s, returning the rest of the string in rest.
// If it does not contain a valid scheme, returns "", s.
//
// A scheme is of the form 'scheme://rest'.
// Scheme must match the regular expression `^[a-zA-Z][a-zA-Z0-9+\-\.]*$`.
func SplitURLScheme(s string) (scheme string, rest string) {
	// An obvious implementation of this function would simply match against the
	// regular expression.
	//
	// However such an implementation would be a lot slower when called repeatedly.
	// Instead we directly build code that directly implements the trimming.

	// Iterate over the bytes in this string.
	// We can do this safely, because any valid scheme must be ascii.
	bytes := []byte(s)
	if len(bytes) == 0 {
		return "", s
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
		goto noScheme
	}

	// get the current letter
	c = bytes[idx]

	// reached end of the scheme
	// we can make use of this because no valid scheme has a ':'.
	if c == ':' {
		goto scheme
	}

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

noScheme:
	// not a valid letter
	return "", s

scheme:
	// split into scheme and rest
	scheme = string(bytes[:idx])
	rest = string(bytes[idx+1:])

	// check that the rest starts with "//"
	if len(rest) < 2 || rest[0] != '/' || rest[1] != '/' {
		goto noScheme
	}

	// and trim off the valid prefix
	rest = rest[2:]
	return
}
