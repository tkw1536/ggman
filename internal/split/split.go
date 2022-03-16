// Package split provides Before and After.
package split

import (
	"strings"
)

// Before splits the string s into a part before a seperator, called the prefix, and a part after the seperator, called the suffix.
// If seperator is not contained in the source string, prefix is empty and suffix is equal to the input string.
//
// See also After.
func Before(s, sep string) (prefix, suffix string) {
	prefix, suffix, found := strings.Cut(s, sep)
	if !found {
		return "", s
	}
	return prefix, suffix
}

// After splits the string s into a part before a seperator, called the prefix, and a part after the seperator, called the suffix.
// If seperator is not contained in the source string, suffix is empty and prefix is equal to the input string.
//
// See also Before.
func After(s, sep string) (prefix, suffix string) {
	prefix, suffix, _ = strings.Cut(s, sep)
	return
}
