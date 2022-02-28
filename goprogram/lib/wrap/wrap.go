// Package wrap provides facilities to wrap text
package wrap

import (
	"bufio"
	"io"
	"strings"
	"sync"
	"unicode"

	"github.com/tkw1536/ggman/goprogram/lib/text"
)

var newLine = []byte("\n")

// Wrapper provides methods to write hard-wrapped lines to Writer.
type Wrapper struct {
	Writer io.Writer
	Length int
}

// WriteLine writes a wrapped line to the output
func (w Wrapper) WriteLine(line string) (int, error) {
	return w.write("", line)
}

// WritePrefix writes a line with a prefix
func (w Wrapper) WritePrefix(prefix, line string) (int, error) {
	w.Length -= len(prefix)
	if w.Length < 0 {
		w.Length = 1
	}
	return w.write(prefix, line)
}

// WriteIndent determines the largest space-only prefix of line, and uses it to call WriteLinePrefix().
func (w Wrapper) WriteIndent(line string) (int, error) {
	length := strings.IndexFunc(line, func(r rune) bool { return !unicode.IsSpace(r) })
	if length == -1 {
		length = len(line)
	}
	return w.WritePrefix(line[:length], line[length:])
}

// WriteString splits s into lines, and then passes each line into WriteIndent.
// It also inserts newlines in between each line passed to WriteIndent.
func (w Wrapper) WriteString(s string) (n int, err error) {
	scanner := bufio.NewScanner(strings.NewReader(s))

	// write the first line
	if scanner.Scan() {
		n, err = w.WriteIndent(scanner.Text())
		if err != nil {
			return
		}
	}

	// write subsequent lines followed by newlines
	for scanner.Scan() {
		w.Writer.Write(newLine)
		m, err := w.WriteIndent(scanner.Text())
		n += m
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

// write implements WriteLine and WriteLinePrefix.
// It writes line, wrapped at w.Length, to the output.
// Each line is prefixed by prefix.
func (w Wrapper) write(prefix, line string) (n int, err error) {
	// split the line into words!
	words := strings.Fields(line)
	if w.Length <= 0 {
		n, err = io.WriteString(w.Writer, prefix)
		if err != nil {
			return
		}

		m, err := text.Join(w.Writer, words, " ")
		n += m
		return n, err
	}

	// when there are no words, only write the prefix!
	if len(words) == 0 {
		return io.WriteString(w.Writer, prefix)
	}

	for len(words) > 0 {
		// find the word count and length of the current line!

		// always pick the first word!
		ll := len(words[0]) // current length
		wc := 1             // word count

		// keep picking words while there is space left in the line
		for ; len(words) > wc; wc++ {
			ll += len(words[wc]) + 1
			if ll >= w.Length {
				break
			}
		}

		// if there are words left, then we need to write a newline
		// so we want to allocate space for that too
		if len(words) > wc {
			ll++
		}
		text.Grow(w.Writer, ll+len(prefix))

		m, err := io.WriteString(w.Writer, prefix)
		n += m
		if err != nil {
			return n, err
		}

		// io.WriteString(w.Writer, strings.Join(" ", words[:wc]))
		m, err = io.WriteString(w.Writer, words[0])
		n += m
		if err != nil {
			return n, err
		}
		for _, word := range words[1:wc] {
			m, err = io.WriteString(w.Writer, " ")
			n += m
			if err != nil {
				return n, err
			}

			m, err = io.WriteString(w.Writer, word)
			n += m
			if err != nil {
				return n, err
			}
		}

		// write a newline if there are still words left!
		words = words[wc:]
		if len(words) > 0 {
			m, err = w.Writer.Write(newLine)
			n += m
			if err != nil {
				return n, err
			}
		}
	}

	return n, nil
}

var wrapperPool = &sync.Pool{
	New: func() interface{} {
		return new(Wrapper)
	},
}

// WriteString is a convenience method that creates a new wrapper and calls WriteString(s) on it.
//
// This method is untested because Wrapper.WriteString is tested.
func WriteString(writer io.Writer, length int, s string) (int, error) {
	wrapper := wrapperPool.Get().(*Wrapper)
	wrapper.Writer = writer
	wrapper.Length = length

	// avoid leaking writer
	defer func() {
		wrapper.Writer = nil
		wrapperPool.Put(wrapper)
	}()

	return wrapper.WriteString(s)
}

var builderPool = &sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

// WrapString is a convenience method that creates a new Wrapper, writes s to it, and then returns the written data.
// When s has a trailing newline, also adds a trailing newline to the return value.
//
// This method is untested because Wrapper.WriteString is tested.
func WrapString(length int, s string) string {
	builder := builderPool.Get().(*strings.Builder)
	builder.Reset()
	defer builderPool.Put(builder)

	WriteString(builder, length, s)
	if strings.HasSuffix(s, "\n") {
		builder.Write(newLine)
	}

	return builder.String()
}
