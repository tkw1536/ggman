// Package split provides Before and After.
package split

import (
	"strings"
)

// Before splits the string s into a part before a separator, called the prefix, and a part after the separator, called the suffix.
// If separator is not contained in the source string, prefix is empty and suffix is equal to the input string.
//
// See also After.
func Before(s, sep string) (prefix, suffix string) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):]
	}
	return "", s
}

// After splits the string s into a part before a separator, called the prefix, and a part after the separator, called the suffix.
// If separator is not contained in the source string, suffix is empty and prefix is equal to the input string.
//
// See also Before.
func After(s, sep string) (prefix, suffix string) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):]
	}
	return s, ""
}
