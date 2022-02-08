package env

import (
	"bufio"
	"io"
	"strings"

	"github.com/pkg/errors"
)

// CanLine represents a line within in the canonical configuration file
type CanLine struct {
	Pattern   string
	Canonical string
}

// ErrEmpty is an error representing an empty CanLine
var ErrEmpty = errors.New("CanLine is empty")

// Unmarshal reads a CanLine from a string
func (cl *CanLine) Unmarshal(s string) error {
	s = strings.TrimSpace(s)

	// if the line is empty or starts with a comment character return nothing
	if s == "" || strings.HasPrefix(s, "#") || strings.HasPrefix(s, "//") || strings.HasPrefix(s, ";") {
		return ErrEmpty
	}

	// get the fields of the string
	fields := strings.Fields(s)

	// switch based on the length
	switch len(fields) {
	case 0:
		return errors.Errorf("strings.Fields() unexpectedly returned 0-length slice")
	case 1:
		fields = []string{"", fields[0]}
	}

	cl.Pattern = fields[0]
	cl.Canonical = fields[1]

	return nil
}

// CanFile represents a list of CanLines
type CanFile []CanLine

// ReadFrom populates this CanFile with CanLines read from the given reader.
// It returns an error (if any occured) and the total bytes read from reader.
//
// Individual CanLines are parsed using CanLine.Unmarshal().
// If reader returns a non-EOF error or parsing fails, ReadFrom returns an appropriate error.
func (cf *CanFile) ReadFrom(reader io.Reader) (int64, error) {
	var bytes int64

	*cf = nil

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		bytes += int64(len(text))

		var line CanLine
		if err := line.Unmarshal(text); err != nil {
			return bytes, errors.Wrap(err, "Unable to parse CANFILE line")
		}
		*cf = append(*cf, line)
	}

	return bytes, errors.Wrap(scanner.Err(), "Unable to read CANFILE")
}

var defaultCanFile = []string{
	"git@^:$.git",
}

// ReadDefault loads the default CanLines into this file.
//
// If the default CanLines can not be read, calls panic().
// A call to panic() is considered a bug.
func (cf *CanFile) ReadDefault() {
	*cf = make([]CanLine, len(defaultCanFile))
	for i, cl := range defaultCanFile {
		if err := (*cf)[i].Unmarshal(cl); err != nil {
			panic("CanFile.ReadDefault(): Unable to parse default CanFile line")
		}
	}
}
