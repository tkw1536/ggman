package utils

import (
	"regexp"
	"strings"
)

// WrapLine wraps a single line at length, not preserving any whitespace
func WrapLine(line string, length int) (lines []string) {
	// split it into strings
	fields := strings.Fields(line)

	// wrapping at length 0 makes no sense
	if length <= 0 {
		return []string{strings.Join(fields, " ")}
	}

	// while we have a field left
	for len(fields) > 0 {
		newLine := ""

		for len(fields) > 0 {
			// the next element
			next := fields[0]

			// the space seperator we have to add
			space := " "
			if newLine == "" {
				space = ""
			}

			// if we have enough space
			// or we are at the first word to add to the line
			if (len(newLine)+1+len(next) < length) || newLine == "" {
				newLine += space + next
				fields = fields[1:]

				// else we break and skip to the next line
			} else {
				break
			}
		}

		// add the line to the buffer
		lines = append(lines, newLine)
	}

	return
}

var wsRegex = regexp.MustCompile("^(\\s*)(.*)$")

// WrapLinePreserve wraps a single line at length
// preserving leading whitespace in all resulting lines
func WrapLinePreserve(line string, length int) []string {
	// if the line is short enough return it as is
	if len(line) < length {
		return []string{line}
	}

	// split off the leading whitespace
	split := wsRegex.FindStringSubmatch(line)
	prefix := split[1]

	// wrap the lines and append the prefix everywhere
	lines := WrapLine(split[2], length-len(prefix))
	for i := range lines {
		lines[i] = prefix + lines[i]
	}
	return lines
}

// WrapString is like WrapLine except that it treats each line in the input seperatly
func WrapString(s string, length int) (lines []string) {
	// split the string into lines and treat each line seperatly
	for _, line := range strings.Split(strings.Replace(s, "\r\n", "\n", -1), "\n") {
		lines = append(lines, WrapLine(line, length)...)
	}

	return
}

// WrapStringPreserve is like WrapLinePreserve except that it treats
// each line in the input seperatly
func WrapStringPreserve(s string, length int) (lines []string) {
	// split the string into lines and treat each line seperatly
	for _, line := range strings.Split(strings.Replace(s, "\r\n", "\n", -1), "\n") {
		lines = append(lines, WrapLinePreserve(line, length)...)
	}

	return
}

// WrapStringPreserveJ is like WrapStringPreserve except that it
// joins the resulting lines instead of returning an array
func WrapStringPreserveJ(s string, length int) string {
	return strings.Join(WrapStringPreserve(s, length), "\n")
}
