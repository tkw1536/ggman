// Package split provides Before and After.
//
//spellchecker:words split
package split

//spellchecker:words strings unicode
import (
	"strings"
	"unicode/utf8"
)

// AfterRune splits the string s into a part before a separator, called the prefix, and a part after the separator, called the suffix.
// If separator is not contained in the source string, suffix is empty and prefix is equal to the input string.
func AfterRune(s string, sep rune) (prefix, suffix string) {
	if i := strings.IndexRune(s, sep); i >= 0 {
		return s[:i], s[i+utf8.RuneLen(sep):]
	}
	return s, ""
}
