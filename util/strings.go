package util

import "strings"

// SplitBefore splits the string s into a part before a seperator, called the prefix, and a part after the seperator, called the suffix.
// If seperator is not contained in the source string, prefix is empty and suffix is equal to the input string.
//
// See also SplitAfter.
func SplitBefore(s, sep string) (prefix, suffix string) {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 {
		suffix = s
		return
	}
	return parts[0], parts[1]
}

// SplitAfter splits the string s into a part before a seperator, called the prefix, and a part after the seperator, called the suffix.
// If seperator is not contained in the source string, suffix is empty and prefix is equal to the input string.
//
// See also SplitBefore.
func SplitAfter(s, sep string) (prefix, suffix string) {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 {
		prefix = s
		return
	}
	return parts[0], parts[1]
}

// TrimSuffixWhile repeatedly removes a suffix from a string, until it is no longer a suffix of a string
// If suffix is the empty string, returns s unchanged.
//
// See also TrimPrefixWhile.
func TrimSuffixWhile(s, suffix string) string {
	// A straightforward implementation of TrimSuffixWhile would be:
	//
	// for strings.HasSuffix(s, suffix) {
	// 	s = strings.TrimSuffix(s, suffix)
	// }
	//
	// However as TrimSuffix internally calls HasSuffix
	// (to check if it actually needs to do something),
	// this implementation is somewhat inefficient.
	// It furthermore requires a special case for the empty suffix.
	//
	// Instead we just keep trying to trim the suffix of the string.
	// If s does not have suffix, TrimSuffix() will leave it unchanged.
	// Furthermore TrimSuffix changes s iff it changes s's length.
	//
	// To furthermore improve performance, we only compute the length of s
	// at most once per iteration and store the previous length in the preLen variable.
	//
	// Note that in the first iteration, the for loop just checks that s is not the empty string.
	// This is ok, as the empty string has only a single prefix (the empty string) and in that case we do not need to change the string.
	var prevLen int

	for curLen := len(s); prevLen != curLen; curLen = len(s) {
		prevLen = curLen
		s = strings.TrimSuffix(s, suffix)
	}

	return s
}

// TrimPrefixWhile repeatedly removes a prefix from a string, until it is no longer a prefix of a string.
// If prefix is the empty string, returns s unchanged.
//
// See also TrimSuffixWhile.
func TrimPrefixWhile(s, prefix string) string {

	// This function exactly mirrors TrimSuffixWhile.
	// See comments above for why we use this trick.

	var prevLen int

	for curLen := len(s); prevLen != curLen; curLen = len(s) {
		prevLen = curLen
		s = strings.TrimPrefix(s, prefix)
	}

	return s
}
