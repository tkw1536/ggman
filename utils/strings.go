package utils

import "strings"

// SplitBefore splits into two parts given a seperator
// if the seperator does not exist in the source string, only suffix is returned
func SplitBefore(s string, sep string) (prefix string, suffix string) {
	return splitTwo(s, sep, false)
}

// SplitAfter splits into two parts given a seperator
// if the seperator does not exist in the source string, only prefix is returned
func SplitAfter(s string, sep string) (prefix string, suffix string) {
	return splitTwo(s, sep, true)
}

func splitTwo(s string, sep string, returnPrefix bool) (prefix string, suffix string) {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 1 {
		if returnPrefix {
			prefix = parts[0]
		} else {
			suffix = parts[0]
		}
	} else {
		prefix = parts[0]
		suffix = parts[1]
	}
	return
}

// TrimSuffixWhile repeatedly trims a suffix from a string
func TrimSuffixWhile(s string, suffix string) (trimmed string) {
	trimmed = s
	if suffix != "" {
		for strings.HasSuffix(trimmed, suffix) {
			trimmed = strings.TrimSuffix(trimmed, suffix)
		}
	}
	return
}

// TrimPrefixWhile repeatedly trims a prefix of a string
func TrimPrefixWhile(s string, prefix string) (trimmed string) {
	trimmed = s
	if prefix != "" {
		for strings.HasPrefix(trimmed, prefix) {
			trimmed = strings.TrimPrefix(trimmed, prefix)
		}
	}
	return
}
