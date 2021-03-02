package text

import (
	"strings"
	"testing"
)

func TestJoin(t *testing.T) {
	tests := []struct {
		prefix string
		elems  []string
		sep    string
	}{

		{"", nil, ""},
		{"", nil, ","},

		{"", []string{"red"}, ""},
		{"", []string{"red"}, ","},

		{"", []string{"red", "yellow", "pink", "green", "purple", "orange", "blue"}, ""},
		{"", []string{"red", "yellow", "pink", "green", "purple", "orange", "blue"}, ","},

		{"x", nil, ""},
		{"x", nil, ","},

		{"x", []string{"red"}, ""},
		{"x", []string{"red"}, ","},

		{"x", []string{"red", "yellow", "pink", "green", "purple", "orange", "blue"}, ""},
		{"x", []string{"red", "yellow", "pink", "green", "purple", "orange", "blue"}, ","},

		{"xx", nil, ""},
		{"xx", nil, ","},

		{"xx", []string{"red"}, ""},
		{"xx", []string{"red"}, ","},

		{"xx", []string{"red", "yellow", "pink", "green", "purple", "orange", "blue"}, ""},
		{"xx", []string{"red", "yellow", "pink", "green", "purple", "orange", "blue"}, ","},

		{"xxx", nil, ""},
		{"xxx", nil, ","},

		{"xxx", []string{"red"}, ""},
		{"xxx", []string{"red"}, ","},

		{"xxx", []string{"red", "yellow", "pink", "green", "purple", "orange", "blue"}, ""},
		{"xxx", []string{"red", "yellow", "pink", "green", "purple", "orange", "blue"}, ","},
	}
	for _, tt := range tests {
		builder := &strings.Builder{}

		builder.WriteString(tt.prefix)
		gotN, gotErr := Join(builder, tt.elems, tt.sep)
		got := builder.String()

		want := tt.prefix + strings.Join(tt.elems, tt.sep)
		wantN := len(want) - len(tt.prefix)

		if got != want {
			t.Errorf("Join() = %v, want %v", got, want)
		}

		if gotN != wantN {
			t.Errorf("Join() n = %v, want %v", gotN, wantN)
		}

		if gotErr != nil {
			t.Errorf("Join() err = %s, want = nil", gotErr)
		}
	}
}

func TestRepeatJoin(t *testing.T) {
	tests := []struct {
		prefix string
		s, sep string
		count  int
	}{
		{"", "abc", ", ", 0},
		{"", "abc", ", ", 1},
		{"", "abc", ", ", 2},
		{"", "abc", ", ", 3},

		{"x", "abc", ", ", 0},
		{"x", "abc", ", ", 1},
		{"x", "abc", ", ", 2},
		{"x", "abc", ", ", 3},

		{"xx", "abc", ", ", 0},
		{"xx", "abc", ", ", 1},
		{"xx", "abc", ", ", 2},
		{"xx", "abc", ", ", 3},
	}
	for _, tt := range tests {
		builder := &strings.Builder{}

		builder.WriteString(tt.prefix)
		gotN, gotErr := RepeatJoin(builder, tt.s, tt.sep, tt.count)
		got := builder.String()

		var want string
		if tt.count > 0 {
			want = tt.prefix + tt.s + strings.Repeat(tt.sep+tt.s, tt.count-1)
		} else {
			want = tt.prefix
		}
		wantN := len(want) - len(tt.prefix)

		if got != want {
			t.Errorf("RepeatJoin() = %v, want %v", got, want)
		}

		if gotN != wantN {
			t.Errorf("RepeatJoin() n = %v, want %v", gotN, wantN)
		}

		if gotErr != nil {
			t.Errorf("RepeatJoin() err = %s, want = nil", gotErr)
		}
	}
}

func TestRepeat(t *testing.T) {
	tests := []struct {
		prefix string
		s      string
		count  int
	}{
		{"", "", 0},
		{"", "", 1},
		{"", "", 2},
		{"", "-", 0},
		{"", "-", 1},
		{"", "-", 10},
		{"", "abc ", 3},

		{"x", "", 0},
		{"x", "", 1},
		{"x", "", 2},
		{"x", "-", 0},
		{"x", "-", 1},
		{"x", "-", 10},
		{"x", "abc ", 3},

		{"xxx", "", 0},
		{"xxx", "", 1},
		{"xxx", "", 2},
		{"xxx", "-", 0},
		{"xxx", "-", 1},
		{"xxx", "-", 10},
		{"xxx", "abc ", 3},
	}
	for _, tt := range tests {
		builder := &strings.Builder{}

		builder.WriteString(tt.prefix)
		gotN, gotErr := Repeat(builder, tt.s, tt.count)
		got := builder.String()

		want := tt.prefix + strings.Repeat(tt.s, tt.count)
		wantN := len(want) - len(tt.prefix)

		if got != want {
			t.Errorf("Repeat() = %v, want %v", got, want)
		}

		if gotN != wantN {
			t.Errorf("Repeat() n = %v, want %v", gotN, wantN)
		}

		if gotErr != nil {
			t.Errorf("Repeat() err = %s, want = nil", gotErr)
		}
	}
}
