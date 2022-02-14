package meta

import (
	"io"

	"github.com/tkw1536/ggman/internal/text"
)

// Positional holds meta-information about a positional argument.
type Positional struct {
	// Name and Description of the positional in help texts
	Value       string // defaults to "ARGUMENT"
	Description string

	// Minimal and Maximum number of times the positonal argument might be passed.
	// Min must be >= 0. Max must be either Min, or -1.
	// Max == -1 inidicates an unlimited number of repeats.
	Min, Max int
}

// defaultPositionalValue is the default name used for a positional argument.
const defaultPositionalValue = "ARGUMENT"

// WriteSpecTo writes a specification of this argument into w.
// A specification looks like "arg [arg...]".
func (pos Positional) WriteSpecTo(w io.Writer) {
	extra := pos.Max - pos.Min

	if pos.Min < 0 || (pos.Max > 0 && extra < 0) { // invalid arguments
		panic("Positional: negative min or out of range max")
	}

	if pos.Value == "" {
		pos.Value = defaultPositionalValue
	}

	// nothing to generate!
	if pos.Max == 0 && extra == 0 {
		return
	}

	// arg arg arg
	text.RepeatJoin(w, pos.Value, " ", pos.Min)
	if pos.Min > 0 && extra != 0 {
		io.WriteString(w, " ")
	}

	if pos.Max < 0 {
		// [arg ...]
		io.WriteString(w, "[")
		io.WriteString(w, pos.Value)
		io.WriteString(w, " ...]")
		return
	}

	// [arg [arg]]
	text.RepeatJoin(w, "["+pos.Value, " ", extra)
	text.Repeat(w, "]", extra)
}
