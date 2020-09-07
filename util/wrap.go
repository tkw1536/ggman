package util

import (
	"strings"
	"unicode"
)

// WrapLine wraps a line by splitting it into multiple lines, with each at most length length.
// This furthermore condenses any sequence of whitespace characters (as defined by unicode.IsSpace) into a single space character.
// If a single non-whitespace containing substring is longer than length, it will be on it's own line.
// When length is <= 0, a single line is returned.
func WrapLine(line string, length int) []string {

	// We first use strings.Fields to split the input into a set of space-seperated words.
	// We now work only with this, and not the input line.
	words := strings.Fields(line)

	// invalid values of length: return a single line
	if length <= 0 {
		return []string{strings.Join(words, " ")}
	}

	// allocate an array of lines, worst case: 1 word per line
	lines := make([]string, 0, len(words))
	for len(words) > 0 {
		var currentLine string

		for len(words) > 0 {

			// take the first word from the array
			// and if neccessary, add a space before it.
			thisWord := words[0]
			if currentLine != "" {
				thisWord = " " + thisWord
			}

			// if the current word does not fit the line and the currentLine is not empty
			// this line is full
			if len(currentLine)+len(thisWord) >= length && currentLine != "" {
				break
			}

			// add this word to the current line
			currentLine += thisWord
			words = words[1:]
		}

		lines = append(lines, currentLine)
	}

	return lines
}

// WrapLinePrefix is like WrapLine, except that it preserves any prefix of whitespace at the beginning of line
// and prefixes it to every returned line. The wrapping length passed to WrapLine will be length - len(leading whitespace), or 1, whichever is bigger.
// If line only contains whitespace, the return value will contain only the input line.
func WrapLinePrefix(line string, length int) []string {

	// if line is already very short, then we can return it directly
	// and don't need to do any computation.
	if len(line) <= length {
		return []string{line}
	}

	// trim off the whitespace prefix from the string
	// by repeatedly checking if a rune is a space character
	var prefix string
	trimmedPrefix := false
	for i, r := range line {
		if !unicode.IsSpace(r) {
			trimmedPrefix = true
			prefix = line[:i]
			line = line[i:]
			break
		}
	}

	// if we did not trim the prefix, then the string is all whitespace
	// so the prefix is everything and the line is empty.
	if !trimmedPrefix {
		prefix = line
		line = ""
	}

	// compute the new length by subtracting the length of prefix.
	// Make sure that it is at least 1, so that we actually do some wrapping.
	length -= len(prefix)
	if length < 1 {
		length = 1
	}

	// finally call WrapLine and prefix every line with the prefix
	lines := WrapLine(line, length)
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return lines
}

// WrapStringPrefix is like WrapStringPrefix except that it first splits the input into newline seperated strings.
// It then treats each line seperatly.
func WrapStringPrefix(s string, length int) (lines []string) {
	for _, line := range strings.Split(strings.Replace(s, "\r\n", "\n", -1), "\n") {
		lines = append(lines, WrapLinePrefix(line, length)...)
	}

	return
}

// WrapStringsPrefix is like WrapStringPrefix except that it joins all resulting lines into a single string seperated by newlines.
func WrapStringsPrefix(s string, length int) string {
	return strings.Join(WrapStringPrefix(s, length), "\n")
}
