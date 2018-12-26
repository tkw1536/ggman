package utils

import "strings"

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
