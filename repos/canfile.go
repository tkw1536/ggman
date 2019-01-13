package repos

import (
	"bufio"
	"os"
	"strings"
)

// CanLine represents a line within in the canonical configuration file
type CanLine struct {
	Pattern   string
	Canonical string
}

// ReadCanFile reads lines from the canonical file
func ReadCanFile(filename string) (lines []CanLine, err error) {
	// open the file
	file, err := os.Open(filename)
	defer file.Close()

	// scan through it
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// read every line
		newLine := ReadCanLine(scanner.Text())
		// add every valid line
		if newLine != nil {
			lines = append(lines, *newLine)
		}
	}

	// if there was an error, return the error
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// else return the lines
	return
}

// ReadCanLine reads a single canon
func ReadCanLine(line string) (cl *CanLine) {
	// remove all the spaces
	trimmed := strings.TrimSpace(line)

	// if the line is empty or starts with a comment character return nothing
	if trimmed == "" || strings.HasPrefix(trimmed, "#") || strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, ";") {
		return
	}

	// get the fields of the string
	fields := strings.Fields(trimmed)

	// if we have only one field, assume it is the default
	if len(fields) == 1 {
		cl = &CanLine{"", fields[0]}

		// else take the first two fields
	} else {
		cl = &CanLine{fields[0], fields[1]}
	}

	return
}
