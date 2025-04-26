package env

//spellchecker:words bufio errors strings
import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

//spellchecker:words canfile unmarshals

// CanLine represents a line within in the canonical configuration file.
type CanLine struct {
	Pattern   string
	Canonical string
}

var errCanLineEmpty = errors.New("CanLine.UnmarshalText: CanLine is empty")

// UnmarshalText unmarshals a text representation of itself.
// An empty line returns [errCanLineEmpty].
func (cl *CanLine) UnmarshalText(text []byte) error {
	s := strings.TrimSpace(string(text))

	// if the line is empty or starts with a comment character return nothing
	if s == "" || strings.HasPrefix(s, "#") || strings.HasPrefix(s, "//") || strings.HasPrefix(s, ";") {
		return errCanLineEmpty
	}

	// get the fields of the string
	fields := strings.Fields(s)

	// switch based on the length
	switch len(fields) {
	case 0:
		// strings.TrimSpace(s) != "" but strings.Fields(s) == []
		panic("never reached")
	case 1:
		cl.Pattern = ""
		cl.Canonical = fields[0]
	default:
		cl.Pattern = fields[0]
		cl.Canonical = fields[1]
	}

	return nil
}

// CanFile represents a list of CanLines.
type CanFile []CanLine

// ReadFrom populates this CanFile with CanLines read from the given reader.
// It returns an error (if any occurred) and the total bytes read from reader.
//
// Individual CanLines are parsed using CanLine.Unmarshal().
// If reader returns a non-EOF error or parsing fails, ReadFrom returns an appropriate error.
func (cf *CanFile) ReadFrom(reader io.Reader) (int64, error) {
	var bytes int64

	*cf = nil

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Bytes()
		bytes += int64(len(text))

		var line CanLine

		err := line.UnmarshalText(text)
		if errors.Is(err, errCanLineEmpty) {
			continue
		}
		if err != nil {
			return bytes, fmt.Errorf("unable to parse CANFILE line: %w", err)
		}
		*cf = append(*cf, line)
	}

	if err := scanner.Err(); err != nil {
		return bytes, fmt.Errorf("unable to read CANFILE: %w", err)
	}
	return bytes, nil
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
		if err := (*cf)[i].UnmarshalText([]byte(cl)); err != nil {
			panic("CanFile.ReadDefault: Unable to parse default CanFile line")
		}
	}
}
