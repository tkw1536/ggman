package util

import "strings"

// SplitBefore splits the string s into a part before a seperator, called the prefix, and a part after the seperator, caled the suffix.
// If seperator is not contained in the source string, prefix is empty and suffix is equal to the input string.
// See also SplitAfter.
func SplitBefore(s, sep string) (prefix string, suffix string) {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 {
		return "", s
	}
	return parts[0], parts[1]
}

// SplitAfter splits the string s into a part before a seperator, called the prefix, and a part after the seperator, caled the suffix.
// If seperator is not contained in the source string, suffix is empty and prefix is equal to the input string.
// See also SplitBefore.
func SplitAfter(s, sep string) (prefix string, suffix string) {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 {
		return s, ""
	}
	return parts[0], parts[1]
}

// TrimSuffixWhile repeatedly removes a suffix from a string, until it is no longer a suffix of a string.
// If suffix is the empty string, returns s unchanged.
// See also TrimSuffixWhile.
func TrimSuffixWhile(s string, suffix string) string {

	// Because the empty string is a suffix of every string, we need to treat it special.
	// It is not possible to remote the empty suffix from a string.
	if suffix == "" {
		return s
	}

	for strings.HasSuffix(s, suffix) {
		s = strings.TrimSuffix(s, suffix)
	}
	return s
}

// TrimPrefixWhile repeatedly removes a prefix from a string, until it is no longer a prefix of a string.
// If prefix is the empty string, returns s unchanged.
// See also TrimSuffixWhile.
func TrimPrefixWhile(s, prefix string) string {

	// Because the empty string is a prefix of every string, we need to treat it special.
	// It is not possible to remote the empty prefix from a string.
	if prefix == "" {
		return s
	}

	for strings.HasPrefix(s, prefix) {
		s = strings.TrimPrefix(s, prefix)
	}
	return s
}
