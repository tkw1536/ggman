package util

import (
	"reflect"
	"strings"
	"testing"
)

func TestNewGlobPattern(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want Pattern
	}{
		{"empty pattern", args{""}, AnyStringPattern{}},
		{"constant pattern", args{"hello world"}, EqualityFoldPattern("hello world")},
		{"glob pattern", args{"a*b"}, GlobPattern("a*b")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGlobPattern(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGlobPattern() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyStringPattern_Match(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"empty string match",
			args{""},
			true,
		},

		{
			"hello world string match",
			args{"hello world"},
			true,
		},
		{
			"$*? string match",
			args{"$*?"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := AnyStringPattern{}
			if got := a.Match(tt.args.s); got != tt.want {
				t.Errorf("AnyStringPattern.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEqualityFoldPattern_Match(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		p    EqualityFoldPattern
		args args
		want bool
	}{
		{
			"pattern matches exactly",
			EqualityFoldPattern("test"),
			args{"test"},
			true,
		},

		{
			"pattern matches case",
			EqualityFoldPattern("test"),
			args{"tEsT"},
			true,
		},

		{
			"pattern does not match",
			EqualityFoldPattern("test"),
			args{"not-match"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Match(tt.args.s); got != tt.want {
				t.Errorf("EqualityFoldPattern.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGlobPattern_Match(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		p    GlobPattern
		args args
		want bool
	}{
		{
			"pattern matches exactly",
			GlobPattern("a*b"),
			args{"aaaab"},
			true,
		},

		{
			"pattern matches case",
			GlobPattern("a*b"),
			args{"AaAaB"},
			true,
		},

		{
			"pattern does not match",
			GlobPattern("a*b"),
			args{"1234"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Match(tt.args.s); got != tt.want {
				t.Errorf("GlobPattern.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSplitGlobPattern(t *testing.T) {
	type args struct {
		pattern  string
		splitter func(string) []string
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
			"simple spliter",
			args{"a;a*b;;", simpleSplitter},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSplitGlobPattern(tt.args.pattern, tt.args.splitter)

			gotPointer := reflect.ValueOf(got.Split).Pointer()
			wantPointer := reflect.ValueOf(tt.want.Split).Pointer()

			if !reflect.DeepEqual(got.Patterns, tt.want.Patterns) || gotPointer != wantPointer {
				t.Errorf("NewSplitGlobPattern() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestSplitPattern_Match(t *testing.T) {
	type fields struct {
		Split    func(s string) []string
		Patterns []Pattern
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
		want   bool
	}{
		{
			"empty Split pattern matches anything (1)",
			fields{
				Split:    neverCalled,
				Patterns: nil,
			},
			args{"a"},
			true,
		},

		{
			"empty Split pattern matches anything (2)",
			fields{
				Split:    neverCalled,
				Patterns: nil,
			},
			args{"a*b"},
			true,
		},

		{
			"empty Split pattern matches anything (3)",
			fields{
				Split:    neverCalled,
				Patterns: nil,
			},
			args{""},
			true,
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
			true,
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
			true,
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
			true,
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
			true,
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
			false,
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
			false,
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
			false,
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
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := SplitPattern{
				Split:    tt.fields.Split,
				Patterns: tt.fields.Patterns,
			}
			if got := sp.Match(tt.args.s); got != tt.want {
				t.Errorf("SplitPattern.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}
