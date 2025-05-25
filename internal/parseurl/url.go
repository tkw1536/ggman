// Package parseurl provides utilities for URL parsing
//
//spellchecker:words parseurl
package parseurl

//spellchecker:words errors math
import (
	"errors"
	"math"
	"strings"
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

const schemeSlashes = "//"

// SplitScheme splits off the scheme from a URL s, returning the rest of the string in rest.
// If it does not contain a valid scheme, returns "", s.
//
// A scheme is of the form 'scheme://rest'.
// Scheme must match the regular expression `^[a-zA-Z][a-zA-Z0-9+\-\.]*$`.
func SplitScheme(s string) (scheme string, rest string) {
	bytes := []byte(s) // safe because any valid scheme is only single-byte runes
	if len(bytes) == 0 {
		return "", s
	}

	var firstInvalidColonIndex = len(bytes) - len(schemeSlashes)
	if firstInvalidColonIndex < 0 {
		return "", s
	}

	var (
		index       = 0
		currentByte = bytes[0]
	)

	goto checkFirstByte

advanceByteAndCheck:
	index++
	if index >= firstInvalidColonIndex {
		return "", s
	}
	currentByte = bytes[index]

	if currentByte == ':' {
		goto sawColon
	}

	if currentByte == '+' || currentByte == '-' || currentByte == '.' {
		goto advanceByteAndCheck
	}

	if '0' <= currentByte && currentByte <= '9' {
		goto advanceByteAndCheck
	}

checkFirstByte:
	if 'a' <= currentByte && currentByte <= 'z' {
		goto advanceByteAndCheck
	}

	if 'A' <= currentByte && currentByte <= 'Z' {
		goto advanceByteAndCheck
	}

	return "", s

sawColon:
	rest = string(bytes[index+1:])
	if !strings.HasPrefix(rest, schemeSlashes) {
		return "", s
	}

	scheme = string(bytes[:index])
	rest = rest[len(schemeSlashes):]
	return
}
