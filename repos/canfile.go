package repos

import (
	"bufio"
	"errors"
	"os"
	"os/user"
	"path"
	"strings"
)

// CanLine represents a line within in the canonical configuration file
type CanLine struct {
	Pattern   string
	Canonical string
}

// ReadDefaultCanFile reads the default canonical file
// or returns the default contents of the file if it does not exist
func ReadDefaultCanFile() ([]CanLine, error) {
	// the list of files to try and read
	var files []string

	// first try the GGMAN_CANFILE
	canfile := os.Getenv("GGMAN_CANFILE")
	if canfile != "" {
		files = append(files, canfile)
	}

	// and then $HOME/.ggman
	usr, err := user.Current()
	if err == nil {
		files = append(files, path.Join(usr.HomeDir, ".ggman"))
	}

	for _, file := range files {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return ReadCanFile(file)
		}
	}

	// finally fall back onto the default can file
	return defaultCanFile()
}

const canLineDefault = "git@^:$.git"

// defaultCanFile generates the default canFile
func defaultCanFile() (lines []CanLine, err error) {
	line := ReadCanLine(canLineDefault)

	// bail out if it is invalid (or a comment)
	if line == nil {
		return nil, errors.New("Invalid default CanFile")
	}

	// append the can line to the return array
	lines = append(lines, *line)

	// and return
	return
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

// ReadCanLine reads a single canonical line
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
