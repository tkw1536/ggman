// Package text contains utilities for processing strings
package text

import (
	"sort"
	"strings"
)

// SplitBefore splits the string s into a part before a seperator, called the prefix, and a part after the seperator, called the suffix.
// If seperator is not contained in the source string, prefix is empty and suffix is equal to the input string.
//
// See also SplitAfter.
func SplitBefore(s, sep string) (prefix, suffix string) {
	// a perfectly valid implementation of this function could make use of strings.Split or strings.SplitN()
	// but both of those allocate an array, and we do not need that here because we have a special situation.
	// It's much more efficient to just check the index of the seperator and trim the string if found.
	index := strings.Index(s, sep)
	if index == -1 {
		return "", s
	}
	return s[:index], s[index+len(sep):]
}

// SplitAfter splits the string s into a part before a seperator, called the prefix, and a part after the seperator, called the suffix.
// If seperator is not contained in the source string, suffix is empty and prefix is equal to the input string.
//
// See also SplitBefore.
func SplitAfter(s, sep string) (prefix, suffix string) {
	index := strings.Index(s, sep)
	if index == -1 {
		return s, ""
	}
	return s[:index], s[index+len(sep):]
}

// RemoveEmpty returns a slice that is like s, but with empty strings removed.
// This function will invalidate the previous value of s.
//
// It is recommended to store the return value of this function in the original variable.
// The call should look something like:
//
//  s = RemoveEmpty(s)
//
func RemoveEmpty(s []string) []string {

	// Because t is backed by the same slice as s, this function will never re-allocate.
	// Furthermore, because strings are immutable, copying data over is cheap.

	t := s[:0]
	for _, v := range s {
		if v != "" {
			t = append(t, v)
		}
	}
	return t
}

// RemoveDuplicates removes duplicates in s and sorts the result.
// This function will invalidate the previous value of s.
//
// It is recommended to store the return value of this function in the original variable.
// The call should look something like:
//
//  s = RemoveDuplicates(s)
//
func RemoveDuplicates(s []string) []string {
	if len(s) == 0 {
		return s
	}

	// adapted from https://github.com/golang/go/wiki/SliceTricks#in-place-deduplicate-comparable
	sort.Strings(s)

	j := 0
	for i := 1; i < len(s); i++ {
		if s[j] == s[i] {
			continue
		}
		j++

		s[j] = s[i]
	}

	return s[:j+1]
}
