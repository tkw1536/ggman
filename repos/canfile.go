package repos

import (
	"bufio"
	"os"
	"os/user"
	"path"
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
	fieldsLength := len(fields)

	// switch based on the length
	switch fieldsLength {
	case 0:
		return errors.Errorf("strings.Fields() unexpectedly returned 0-length slice")
	case 1:
		fields = []string{"", fields[0]}
	default:
		break
	}

	cl.Pattern = fields[0]
	cl.Canonical = fields[1]

	return nil
}

// CanFile represents a list of CanLines
type CanFile []CanLine

// ReadDefault reads the default CanFile
// if it does not exist, loads the default contents
//
// This function is currently untested.
func (cf *CanFile) ReadDefault() error {

	// first determine the list of files to read
	// These are:
	// - the one pointed to by 'GGMAN_CANFILE'
	// - $HOME/.ggman

	files := make([]string, 0, 2)

	canfile := os.Getenv("GGMAN_CANFILE")
	if canfile != "" {
		files = append(files, canfile)
	}
	usr, err := user.Current()
	if err == nil {
		files = append(files, path.Join(usr.HomeDir, ".ggman"))
	}

	// Try to read each opf the files in order
	// Skipping only non-existing ones.
	// Finally fallback to the default.

	for _, file := range files {
		if _, err := os.Stat(file); !os.IsNotExist(err) {
			return errors.Wrapf(cf.unmarshalFile(file), "Error reading CanFile %q", file)
		}
	}

	return cf.loadDefault()
}

func (cf *CanFile) unmarshalFile(filename string) error {
	// start with an empty CanFile
	*cf = nil

	// open the file
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// make a text and a line
	var lineText string
	lineStruct := &CanLine{}

	// scan through it
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lineText = scanner.Text()
		err = lineStruct.Unmarshal(lineText)
		if err != nil {
			if err == ErrEmpty {
				continue
			}
			return errors.Wrapf(err, "Invalid CanLine %q", lineText)
		}

		*cf = append(*cf, *lineStruct)
	}

	// if there was an error, return the error
	if err := scanner.Err(); err != nil {
		return errors.Wrap(err, "Error scanning file")
	}

	// else return the lines
	return nil
}

var defaultCanLines = []string{
	"git@^:$.git",
}

// loadDefault loads the default lines into this CanFile
func (cf *CanFile) loadDefault() error {
	*cf = make([]CanLine, len(defaultCanLines))
	for i, cl := range defaultCanLines {
		if err := (*cf)[i].Unmarshal(cl); err != nil {
			return errors.Wrapf(err, "Unable to read default can line %q (index %d)", cl, i)
		}
	}
	return nil
}
