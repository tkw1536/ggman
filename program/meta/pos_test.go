package meta

import (
	"strings"
	"testing"
)

func TestPostional_WriteSpecTo(t *testing.T) {
	tests := []struct {
		name string
		pos  Positional
		want string
	}{
		{"arg 0, 0", Positional{Value: "arg", Min: 0, Max: 0}, ""},
		{"arg 0, -1", Positional{Value: "arg", Min: 0, Max: -1}, "[arg ...]"},
		{"arg 0, 3", Positional{Value: "arg", Min: 0, Max: 3}, "[arg [arg [arg]]]"},

		{"no name 0, 0", Positional{Value: "", Min: 0, Max: 0}, ""},
		{"no name 0, -1", Positional{Value: "", Min: 0, Max: -1}, "[ARGUMENT ...]"},
		{"no name 0, 3", Positional{Value: "", Min: 0, Max: 3}, "[ARGUMENT [ARGUMENT [ARGUMENT]]]"},

		{"arg 2, 2", Positional{Value: "arg", Min: 2, Max: 2}, "arg arg"},
		{"arg 2, 4", Positional{Value: "arg", Min: 2, Max: 4}, "arg arg [arg [arg]]"},
		{"arg 2, -1", Positional{Value: "arg", Min: 2, Max: -1}, "arg arg [arg ...]"},

		{"no name 2, 2", Positional{Value: "", Min: 2, Max: 2}, "ARGUMENT ARGUMENT"},
		{"no name 2, 4", Positional{Value: "", Min: 2, Max: 4}, "ARGUMENT ARGUMENT [ARGUMENT [ARGUMENT]]"},
		{"no name 2, -1", Positional{Value: "", Min: 2, Max: -1}, "ARGUMENT ARGUMENT [ARGUMENT ...]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builder strings.Builder
			tt.pos.WriteSpecTo(&builder)

			if got := builder.String(); got != tt.want {
				t.Errorf("Positional.WriteSpecTo() = %q, want %q", got, tt.want)
			}
		})
	}
}
