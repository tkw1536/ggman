// Package text provides functions similar to strings.Join, but based on writers as opposed to strings
package text

import (
	"io"
)

// Join writes the elements of elem into writer, seperated by sep.
// Returns the number of runes written and a nil error.
//
// It is like strings.Join, but writes into a writer instead of allocating a strings.Builder.
func Join(writer io.Writer, elems []string, sep string) (n int, err error) {
	// this function has been adapted from strings.Join

	switch len(elems) {
	case 0:
		return
	case 1:
		return io.WriteString(writer, elems[0])
	}

	n = len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}
	Grow(writer, n)

	io.WriteString(writer, elems[0])
	for _, s := range elems[1:] {
		io.WriteString(writer, sep)
		io.WriteString(writer, s)
	}

	return
}

// RepeatJoin writes s, followed by (count -1) instances of sep + s into w.
// It returns the number of runes written and a nil error.
//
// When count <= 0, no instances of s or sep are written into count.
func RepeatJoin(w io.Writer, s, sep string, count int) (n int, err error) {
	if count <= 0 {
		return
	}

	n = len(s)*count + len(sep)*(count-1)
	Grow(w, n)

	io.WriteString(w, s)
	Repeat(w, sep+s, count-1)

	return
}

// Repeat writes count instances of s into w.
// It returns the number of runes written and a nil error.
// When count would cause an overflow, calls panic().
//
// It is similar to strings.Repeat, but writes into an existing builder without allocating a new one.
//
// When s is empty or count <= 0, no instances of s are written.
func Repeat(w io.Writer, s string, count int) (n int, err error) {
	// this function has been adapted from strings.Repeat
	// with the only significant change being that we track an additional offset in builder!

	if count <= 0 || s == "" {
		return
	}

	if len(s)*count/count != len(s) {
		panic("Repeat: Repeat count causes overflow")
	}

	// grow the buffer by the overall number of bytes needed
	n = len(s) * count
	Grow(w, n)

	// write the string into w repeatedly
	// only compute the number of bytes written if something goes wrong
	for i := 0; i < count; i++ {
		if m, err := io.WriteString(w, s); err != nil {
			return len(s)*i + m, err
		}
	}

	return
}
