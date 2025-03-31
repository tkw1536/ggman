//spellchecker:words pattern
package pattern

//spellchecker:words reflect strings testing
import (
	"reflect"
	"strings"
	"testing"
)

//spellchecker:words aaaab

func TestNewGlobPattern(t *testing.T) {
	type args struct {
		s     string
		fuzzy bool
	}
	tests := []struct {
		name string
		args args
		want Pattern
	}{
		{"empty non-fuzzy pattern", args{"", false}, AnyStringPattern{}},
		{"constant non-fuzzy pattern", args{"hello world", false}, EqualityFoldPattern("hello world")},
		{"glob non-fuzzy pattern", args{"a*b", false}, GlobPattern("a*b")},

		{"empty fuzzy pattern", args{"", true}, AnyStringPattern{}},
		{"constant fuzzy pattern", args{"hello world", true}, FuzzyFoldPattern("hello world")},
		{"glob fuzzy pattern", args{"a*b", true}, GlobPattern("a*b")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGlobPattern(tt.args.s, tt.args.fuzzy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGlobPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyStringPattern_Score(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			"empty string match",
			args{""},
			1,
		},

		{
			"hello world string match",
			args{"hello world"},
			1,
		},
		{
			"$*? string match",
			args{"$*?"},
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyStringPattern{}
			if got := a.Score(tt.args.s); got != tt.want {
				t.Errorf("AnyStringPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualityFoldPattern_Score(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		p    EqualityFoldPattern
		args args
		want float64
	}{
		{
			"pattern matches exactly",
			EqualityFoldPattern("test"),
			args{"test"},
			1,
		},

		{
			"pattern matches case",
			EqualityFoldPattern("test"),
			args{"tEsT"},
			1,
		},

		{
			"pattern does not match",
			EqualityFoldPattern("test"),
			args{"not-match"},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Score(tt.args.s); got != tt.want {
				t.Errorf("EqualityFoldPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFuzzyFoldPattern_Score(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		p    FuzzyFoldPattern
		args args
		want float64
	}{
		{
			"pattern matches exactly",
			FuzzyFoldPattern("test"),
			args{"test"},
			1,
		},

		{
			"pattern matches case",
			FuzzyFoldPattern("test"),
			args{"tEsT"},
			1,
		},

		{
			"pattern matches fuzzy",
			FuzzyFoldPattern("tst"),
			args{"test"},
			0.75,
		},

		{
			"pattern matches fuzzy case",
			FuzzyFoldPattern("TsT"),
			args{"TeSt"},
			0.75,
		},

		{
			"pattern does not match",
			FuzzyFoldPattern("test"),
			args{"not-match"},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Score(tt.args.s); got != tt.want {
				t.Errorf("FuzzyFoldPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobPattern_Score(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		p    GlobPattern
		args args
		want float64
	}{
		{
			"pattern matches exactly",
			GlobPattern("a*b"),
			args{"aaaab"},
			1,
		},

		{
			"pattern matches case",
			GlobPattern("a*b"),
			args{"AaAaB"},
			1,
		},

		{
			"pattern does not match",
			GlobPattern("a*b"),
			args{"1234"},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Score(tt.args.s); got != tt.want {
				t.Errorf("GlobPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSplitGlobPattern(t *testing.T) {
	type args struct {
		pattern  string
		splitter func(string) []string
		fuzzy    bool
	}

	simpleSplitter := func(s string) []string {
		return strings.Split(s, ";")
	}

	tests := []struct {
		name string
		args args
		want SplitPattern
	}{
		{
			"simple non-fuzzy splitter",
			args{"a;a*b;;", simpleSplitter, false},
			SplitPattern{
				Split: simpleSplitter,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					GlobPattern("a*b"),
					AnyStringPattern{},
					AnyStringPattern{},
				},
			},
		},
		{
			"match at start without fuzzy set",
			args{"^a;a*b", simpleSplitter, false},
			SplitPattern{
				Split: simpleSplitter,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					GlobPattern("a*b"),
				},
				MatchAtStart: true,
			},
		},
		{
			"match at start with fuzzy set",
			args{"^a;a*b", simpleSplitter, true},
			SplitPattern{
				Split: simpleSplitter,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					GlobPattern("a*b"),
				},
				MatchAtStart: true,
			},
		},
		{
			"match at end without fuzzy set",
			args{"a;a*b$", simpleSplitter, false},
			SplitPattern{
				Split: simpleSplitter,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					GlobPattern("a*b"),
				},
				MatchAtEnd: true,
			},
		},
		{
			"match at end with fuzzy set",
			args{"a;a*b$", simpleSplitter, true},
			SplitPattern{
				Split: simpleSplitter,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					GlobPattern("a*b"),
				},
				MatchAtEnd: true,
			},
		},

		{
			"match at start and end without fuzzy set",
			args{"^a;a*b$", simpleSplitter, false},
			SplitPattern{
				Split: simpleSplitter,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					GlobPattern("a*b"),
				},
				MatchAtStart: true,
				MatchAtEnd:   true,
			},
		},
		{
			"match at start and end without fuzzy set",
			args{"^a;a*b$", simpleSplitter, true},
			SplitPattern{
				Split: simpleSplitter,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					GlobPattern("a*b"),
				},
				MatchAtStart: true,
				MatchAtEnd:   true,
			},
		},

		{
			"simple fuzzy splitter",
			args{"a;a*b;;", simpleSplitter, true},
			SplitPattern{
				Split: simpleSplitter,
				Patterns: []Pattern{
					FuzzyFoldPattern("a"),
					GlobPattern("a*b"),
					AnyStringPattern{},
					AnyStringPattern{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSplitGlobPattern(tt.args.pattern, tt.args.splitter, tt.args.fuzzy)

			gotPointer := reflect.ValueOf(got.Split).Pointer()
			wantPointer := reflect.ValueOf(tt.want.Split).Pointer()

			if !reflect.DeepEqual(got.Patterns, tt.want.Patterns) || gotPointer != wantPointer {
				t.Errorf("NewSplitGlobPattern() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestSplitPattern_Score(t *testing.T) {
	type fields struct {
		Split        func(s string) []string
		Patterns     []Pattern
		MatchAtStart bool
		MatchAtEnd   bool
	}
	type args struct {
		s string
	}

	neverCalled := func(s string) []string {
		panic("never called")
	}

	splitSemicolon := func(s string) []string {
		return strings.Split(s, ";")
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			"empty Split pattern matches anything (1)",
			fields{
				Split:    neverCalled,
				Patterns: nil,
			},
			args{"a"},
			1,
		},

		{
			"empty Split pattern matches anything (2)",
			fields{
				Split:    neverCalled,
				Patterns: nil,
			},
			args{"a*b"},
			1,
		},

		{
			"empty Split pattern matches anything (3)",
			fields{
				Split:    neverCalled,
				Patterns: nil,
			},
			args{""},
			1,
		},

		{
			"exact match",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
			},
			args{"a;b;c"},
			1,
		},

		{
			"start match",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
			},
			args{"a;b;c;c;c"},
			1,
		},

		{
			"middle match",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
			},
			args{"a;a;b;c;c"},
			0.5,
		},
		{
			"end match",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
			},
			args{"a;a;a;b;c"},
			1,
		},

		{
			"no match (too short 1)",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
			},
			args{"a;b"},
			-1,
		},

		{
			"no match (too short 2)",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
			},
			args{"b;c"},
			-1,
		},

		{
			"no match (exact)",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
			},
			args{"b;b;c"},
			-1,
		},
		{
			"no match (long)",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
			},
			args{"a;a;b;b;c;c"},
			-1,
		},

		{
			"exact match at start",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtStart: true,
			},
			args{"a;b;c"},
			1,
		},

		{
			"exact match as prefix",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtStart: true,
			},
			args{"a;b;c;b"},
			1,
		},

		{
			"no match when forced at start",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtStart: true,
			},
			args{"d;a;b;c"},
			-1,
		},

		{
			"exact match at end",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtEnd: true,
			},
			args{"a;b;c"},
			1,
		},

		{
			"exact match as suffix",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtEnd: true,
			},
			args{"d;a;b;c"},
			1,
		},

		{
			"no match when forced as suffix",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtEnd: true,
			},
			args{"a;b;c;d"},
			-1,
		},

		{
			"exact match at with start and end",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtStart: true,
				MatchAtEnd:   true,
			},
			args{"a;b;c"},
			1,
		},

		{
			name: "no match when exact match required, but only suffix",
			fields: fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtStart: true,
				MatchAtEnd:   true,
			},
			args: args{"d;a;b;c"},
			want: -1,
		},

		{
			"no match when exact match required, but only prefix",
			fields{
				Split: splitSemicolon,
				Patterns: []Pattern{
					EqualityFoldPattern("a"),
					EqualityFoldPattern("b"),
					EqualityFoldPattern("c"),
				},
				MatchAtStart: true,
				MatchAtEnd:   true,
			},
			args{"a;b;c;d"},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := SplitPattern{
				Split:        tt.fields.Split,
				Patterns:     tt.fields.Patterns,
				MatchAtStart: tt.fields.MatchAtStart,
				MatchAtEnd:   tt.fields.MatchAtEnd,
			}
			if got := sp.Score(tt.args.s); got != tt.want {
				t.Errorf("SplitPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}
