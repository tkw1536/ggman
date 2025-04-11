//spellchecker:words pattern
package pattern_test

//spellchecker:words reflect strings testing github ggman internal pattern
import (
	"reflect"
	"strings"
	"testing"

	"github.com/tkw1536/ggman/internal/pattern"
)

//spellchecker:words aaaab

func TestNewGlobPattern(t *testing.T) {
	t.Parallel()

	type args struct {
		s     string
		fuzzy bool
	}
	tests := []struct {
		name string
		args args
		want pattern.Pattern
	}{
		{"empty non-fuzzy pattern", args{"", false}, pattern.AnyStringPattern{}},
		{"constant non-fuzzy pattern", args{"hello world", false}, pattern.EqualityFoldPattern("hello world")},
		{"glob non-fuzzy pattern", args{"a*b", false}, pattern.GlobPattern("a*b")},

		{"empty fuzzy pattern", args{"", true}, pattern.AnyStringPattern{}},
		{"constant fuzzy pattern", args{"hello world", true}, pattern.FuzzyFoldPattern("hello world")},
		{"glob fuzzy pattern", args{"a*b", true}, pattern.GlobPattern("a*b")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := pattern.NewGlobPattern(tt.args.s, tt.args.fuzzy); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGlobPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyStringPattern_Score(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			a := pattern.AnyStringPattern{}
			if got := a.Score(tt.args.s); got != tt.want {
				t.Errorf("AnyStringPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualityFoldPattern_Score(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}
	tests := []struct {
		name string
		p    pattern.EqualityFoldPattern
		args args
		want float64
	}{
		{
			"pattern matches exactly",
			pattern.EqualityFoldPattern("test"),
			args{"test"},
			1,
		},

		{
			"pattern matches case",
			pattern.EqualityFoldPattern("test"),
			args{"tEsT"},
			1,
		},

		{
			"pattern does not match",
			pattern.EqualityFoldPattern("test"),
			args{"not-match"},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.p.Score(tt.args.s); got != tt.want {
				t.Errorf("EqualityFoldPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFuzzyFoldPattern_Score(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}
	tests := []struct {
		name string
		p    pattern.FuzzyFoldPattern
		args args
		want float64
	}{
		{
			"pattern matches exactly",
			pattern.FuzzyFoldPattern("test"),
			args{"test"},
			1,
		},

		{
			"pattern matches case",
			pattern.FuzzyFoldPattern("test"),
			args{"tEsT"},
			1,
		},

		{
			"pattern matches fuzzy",
			pattern.FuzzyFoldPattern("tst"),
			args{"test"},
			0.75,
		},

		{
			"pattern matches fuzzy case",
			pattern.FuzzyFoldPattern("TsT"),
			args{"TeSt"},
			0.75,
		},

		{
			"pattern does not match",
			pattern.FuzzyFoldPattern("test"),
			args{"not-match"},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.p.Score(tt.args.s); got != tt.want {
				t.Errorf("FuzzyFoldPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobPattern_Score(t *testing.T) {
	t.Parallel()

	type args struct {
		s string
	}
	tests := []struct {
		name string
		p    pattern.GlobPattern
		args args
		want float64
	}{
		{
			"pattern matches exactly",
			pattern.GlobPattern("a*b"),
			args{"aaaab"},
			1,
		},

		{
			"pattern matches case",
			pattern.GlobPattern("a*b"),
			args{"AaAaB"},
			1,
		},

		{
			"pattern does not match",
			pattern.GlobPattern("a*b"),
			args{"1234"},
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.p.Score(tt.args.s); got != tt.want {
				t.Errorf("GlobPattern.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSplitGlobPattern(t *testing.T) {
	t.Parallel()

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
		want pattern.SplitPattern
	}{
		{
			"simple non-fuzzy splitter",
			args{"a;a*b;;", simpleSplitter, false},
			pattern.SplitPattern{
				Split: simpleSplitter,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.GlobPattern("a*b"),
					pattern.AnyStringPattern{},
					pattern.AnyStringPattern{},
				},
			},
		},
		{
			"match at start without fuzzy set",
			args{"^a;a*b", simpleSplitter, false},
			pattern.SplitPattern{
				Split: simpleSplitter,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.GlobPattern("a*b"),
				},
				MatchAtStart: true,
			},
		},
		{
			"match at start with fuzzy set",
			args{"^a;a*b", simpleSplitter, true},
			pattern.SplitPattern{
				Split: simpleSplitter,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.GlobPattern("a*b"),
				},
				MatchAtStart: true,
			},
		},
		{
			"match at end without fuzzy set",
			args{"a;a*b$", simpleSplitter, false},
			pattern.SplitPattern{
				Split: simpleSplitter,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.GlobPattern("a*b"),
				},
				MatchAtEnd: true,
			},
		},
		{
			"match at end with fuzzy set",
			args{"a;a*b$", simpleSplitter, true},
			pattern.SplitPattern{
				Split: simpleSplitter,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.GlobPattern("a*b"),
				},
				MatchAtEnd: true,
			},
		},

		{
			"match at start and end without fuzzy set",
			args{"^a;a*b$", simpleSplitter, false},
			pattern.SplitPattern{
				Split: simpleSplitter,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.GlobPattern("a*b"),
				},
				MatchAtStart: true,
				MatchAtEnd:   true,
			},
		},
		{
			"match at start and end without fuzzy set",
			args{"^a;a*b$", simpleSplitter, true},
			pattern.SplitPattern{
				Split: simpleSplitter,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.GlobPattern("a*b"),
				},
				MatchAtStart: true,
				MatchAtEnd:   true,
			},
		},

		{
			"simple fuzzy splitter",
			args{"a;a*b;;", simpleSplitter, true},
			pattern.SplitPattern{
				Split: simpleSplitter,
				Patterns: []pattern.Pattern{
					pattern.FuzzyFoldPattern("a"),
					pattern.GlobPattern("a*b"),
					pattern.AnyStringPattern{},
					pattern.AnyStringPattern{},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := pattern.NewSplitGlobPattern(tt.args.pattern, tt.args.splitter, tt.args.fuzzy)

			gotPointer := reflect.ValueOf(got.Split).Pointer()
			wantPointer := reflect.ValueOf(tt.want.Split).Pointer()

			if !reflect.DeepEqual(got.Patterns, tt.want.Patterns) || gotPointer != wantPointer {
				t.Errorf("NewSplitGlobPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSplitPattern_Score(t *testing.T) {
	t.Parallel()

	type fields struct {
		Split        func(s string) []string
		Patterns     []pattern.Pattern
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
				},
			},
			args{"a;b;c"},
			1,
		},

		{
			"start match",
			fields{
				Split: splitSemicolon,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
				},
			},
			args{"a;b;c;c;c"},
			1,
		},

		{
			"middle match",
			fields{
				Split: splitSemicolon,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
				},
			},
			args{"a;a;b;c;c"},
			0.5,
		},
		{
			"end match",
			fields{
				Split: splitSemicolon,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
				},
			},
			args{"a;a;a;b;c"},
			1,
		},

		{
			"no match (too short 1)",
			fields{
				Split: splitSemicolon,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
				},
			},
			args{"a;b"},
			-1,
		},

		{
			"no match (too short 2)",
			fields{
				Split: splitSemicolon,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
				},
			},
			args{"b;c"},
			-1,
		},

		{
			"no match (exact)",
			fields{
				Split: splitSemicolon,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
				},
			},
			args{"b;b;c"},
			-1,
		},
		{
			"no match (long)",
			fields{
				Split: splitSemicolon,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
				},
			},
			args{"a;a;b;b;c;c"},
			-1,
		},

		{
			"exact match at start",
			fields{
				Split: splitSemicolon,
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
				Patterns: []pattern.Pattern{
					pattern.EqualityFoldPattern("a"),
					pattern.EqualityFoldPattern("b"),
					pattern.EqualityFoldPattern("c"),
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
			t.Parallel()

			sp := pattern.SplitPattern{
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
