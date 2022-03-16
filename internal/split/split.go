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
	// a perfectly valid implementation of this function could make use of strings.Split or strings.SplitN()
	// but both of those allocate an array, and we do not need that here because we have a special situation.
	// It's much more efficient to just check the index of the seperator and trim the string if found.
	index := strings.Index(s, sep)
	if index == -1 {
		return "", s
	}
	return s[:index], s[index+len(sep):]
}

// After splits the string s into a part before a seperator, called the prefix, and a part after the seperator, called the suffix.
// If seperator is not contained in the source string, suffix is empty and prefix is equal to the input string.
//
// See also Before.
func After(s, sep string) (prefix, suffix string) {
	index := strings.Index(s, sep)
	if index == -1 {
		return s, ""
	}
	return s[:index], s[index+len(sep):]
}