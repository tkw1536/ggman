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

func TestPositional_Validate(t *testing.T) {
	tests := []struct {
		name string

		Positional Positional
		Count      int

		wantErr string
	}{
		// taking 0 args
		{
			"no arguments",
			Positional{Min: 0, Max: 0},
			0,
			"",
		},

		// taking 1 arg
		{
			"one argument, too few",
			Positional{Min: 1, Max: 1},
			0,
			"Exactly 1 argument(s) required",
		},
		{
			"one argument, exactly enough",
			Positional{Min: 1, Max: 1},
			1,
			"",
		},
		{
			"one argument, too many",
			Positional{Min: 1, Max: 1},
			2,
			"Exactly 1 argument(s) required",
		},

		// taking 1 or 2 args
		{
			"1-2 arguments, too few",
			Positional{Min: 1, Max: 2},
			0,
			"Between 1 and 2 argument(s) required",
		},
		{
			"1-2 arguments, enough",
			Positional{Min: 1, Max: 2},
			1,
			"",
		},
		{
			"1-2 arguments, enough (2)",
			Positional{Min: 1, Max: 2},
			2,
			"",
		},
		{
			"1-2 arguments, too many",
			Positional{Min: 1, Max: 2},
			3,
			"Between 1 and 2 argument(s) required",
		},

		// taking 2 args
		{
			"two arguments, too few",
			Positional{Min: 2, Max: 2},
			0,
			"Exactly 2 argument(s) required",
		},
		{
			"two arguments, too few (2)",
			Positional{Min: 2, Max: 2},
			1,
			"Exactly 2 argument(s) required",
		},
		{
			"two arguments, enough",
			Positional{Min: 2, Max: 2},
			2,
			"",
		},
		{
			"two arguments, too many",
			Positional{Min: 2, Max: 2},
			3,
			"Exactly 2 argument(s) required",
		},

		// at least one argument
		{
			"at least 1 arguments, not enough",
			Positional{Min: 1, Max: -1},
			0,
			"At least 1 argument(s) required",
		},
		{
			"at least 1 arguments, enough",
			Positional{Min: 1, Max: -1},
			1,
			"",
		},
		{
			"at least 1 arguments, enough (2)",
			Positional{Min: 1, Max: -1},
			2,
			"",
		},
		{
			"at least 1 arguments, enough (3)",
			Positional{Min: 1, Max: -1},
			3,
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := tt.Positional.Validate(tt.Count)
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("Positional.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
