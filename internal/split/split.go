// Package split provides Before and After.
package split

import (
	"strings"
	"unicode/utf8"
)

// Before splits the string s into a part before a separator, called the prefix, and a part after the separator, called the suffix.
// If separator is not contained in the source string, prefix is empty and suffix is equal to the input string.
//
// See also AfterRune.
func Before(s, sep string) (prefix, suffix string) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):]
	}
	return "", s
}

// AfterRune splits the string s into a part before a separator, called the prefix, and a part after the separator, called the suffix.
// If separator is not contained in the source string, suffix is empty and prefix is equal to the input string.
//
// See also Before.
func AfterRune(s string, sep rune) (prefix, suffix string) {
	// NOTE(twiesing): This uses sep as a rune, because nothing else is required.
	// And that turns out to be the most efficient variant.
	if i := strings.IndexRune(s, sep); i >= 0 {
		return s[:i], s[i+utf8.RuneLen(sep):]
	}
	return s, ""
}
