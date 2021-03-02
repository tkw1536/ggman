package text

import "strings"

// Join writes the elements of elem into builder, seperated by sep.
// Returns the number of runes written and a nil error.
//
// It is like strings.Join, but writes into a builder instead of allocating a new one.
func Join(builder *strings.Builder, elems []string, sep string) (n int, err error) {
	// this function has been adapted from strings.Join

	switch len(elems) {
	case 0:
		return
	case 1:
		return builder.WriteString(elems[0])
	}

	n = len(sep) * (len(elems) - 1)
	for i := 0; i < len(elems); i++ {
		n += len(elems[i])
	}
	builder.Grow(n)

	builder.WriteString(elems[0])
	for _, s := range elems[1:] {
		builder.WriteString(sep)
		builder.WriteString(s)
	}

	return
}

// RepeatJoin writes s, followed by (count -1) instances of sep + s into builder.
// It returns the number of runes written and a nil error.
//
// When count <= 0, no instances of s or sep are written into count.
func RepeatJoin(builder *strings.Builder, s, sep string, count int) (n int, err error) {
	if count <= 0 {
		return
	}

	n = len(s)*count + len(sep)*(count-1)
	builder.Grow(n)

	builder.WriteString(s)
	Repeat(builder, sep+s, count-1)

	return
}

// Repeat writes count instances of s into builder.
// It returns the number of runes written and a nil error.
// When count would cause an overflow, calls panic().
//
// It is similar to strings.Repeat, but writes into an existing builder without allocating a new one.
//
// When s is empty or count <= 0, no instances of s are written.
func Repeat(builder *strings.Builder, s string, count int) (n int, err error) {
	// this function has been adapted from strings.Repeat
	// with the only significant change being that we track an additional offset in builder!

	if count <= 0 || s == "" {
		return
	}

	if len(s)*count/count != len(s) {
		panic("Repeat: Repeat count causes overflow")
	}

	n = len(s) * count
	builder.Grow(n)

	off := builder.Len()
	builder.WriteString(s)

	// as opposed to strings.Repeat, we need to take care of an offset

	for l := len(s); l < n; l = builder.Len() - off {
		if l <= n/2 {
			builder.WriteString(builder.String()[off:])
		} else {
			builder.WriteString(builder.String()[off : n-l+off])
			break
		}
	}

	return
}
