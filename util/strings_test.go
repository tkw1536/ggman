package util

import (
	"testing"
)

func TestSplitBefore(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name       string
		args       args
		wantPrefix string
		wantSuffix string
	}{
		{"splitFoundOnce", args{"a;b", ";"}, "a", "b"},
		{"splitFoundMultiple", args{"a;b;c", ";"}, "a", "b;c"},
		{"splitNotFound", args{"aaa", ";"}, "", "aaa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefix, gotSuffix := SplitBefore(tt.args.s, tt.args.sep)
			if gotPrefix != tt.wantPrefix {
				t.Errorf("SplitBefore() gotPrefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
			if gotSuffix != tt.wantSuffix {
				t.Errorf("SplitBefore() gotSuffix = %v, want %v", gotSuffix, tt.wantSuffix)
			}
		})
	}
}

func Benchmark_SplitBefore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SplitBefore("a;b", ";")
		SplitBefore("a;b;c", ";")
		SplitBefore("aaa", ";")
	}
}

func TestSplitAfter(t *testing.T) {
	type args struct {
		s   string
		sep string
	}
	tests := []struct {
		name       string
		args       args
		wantPrefix string
		wantSuffix string
	}{
		{"splitFoundOnce", args{"a;b", ";"}, "a", "b"},
		{"splitFoundMultiple", args{"a;b;c", ";"}, "a", "b;c"},
		{"splitNotFound", args{"aaa", ";"}, "aaa", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefix, gotSuffix := SplitAfter(tt.args.s, tt.args.sep)
			if gotPrefix != tt.wantPrefix {
				t.Errorf("SplitAfterTwo() gotPrefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
			if gotSuffix != tt.wantSuffix {
				t.Errorf("SplitAfterTwo() gotSuffix = %v, want %v", gotSuffix, tt.wantSuffix)
			}
		})
	}
}

func Benchmark_SplitAfter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SplitAfter("a;b", ";")
		SplitAfter("a;b;c", ";")
		SplitAfter("aaa", ";")
	}
}

func TestTrimSuffixWhile(t *testing.T) {
	type args struct {
		s      string
		suffix string
	}
	tests := []struct {
		name        string
		args        args
		wantTrimmed string
	}{
		{"trimSingleChar", args{"abcd", "d"}, "abc"},
		{"trimNonExistingChar", args{"abc", "d"}, "abc"},
		{"trimRepeatingChar", args{"abcddd", "d"}, "abc"},
		{"trimEmpty", args{"abc def", ""}, "abc def"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTrimmed := TrimSuffixWhile(tt.args.s, tt.args.suffix); gotTrimmed != tt.wantTrimmed {
				t.Errorf("TrimSuffixWhile() = %v, want %v", gotTrimmed, tt.wantTrimmed)
			}
		})
	}
}

func BenchmarkTrimSuffixWhile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		TrimSuffixWhile("abcabcabcabcabcabcabcabcabcabcabcabcabcabcabcdabc", "abc")
	}
}

func TestTrimPrefixWhile(t *testing.T) {
	type args struct {
		s      string
		prefix string
	}
	tests := []struct {
		name        string
		args        args
		wantTrimmed string
	}{
		{"trimSingleChar", args{"abcd", "a"}, "bcd"},
		{"trimNonExistingChar", args{"bcd", "a"}, "bcd"},
		{"trimRepeatingChar", args{"aaabcd", "a"}, "bcd"},
		{"trimEmpty", args{"abc def", ""}, "abc def"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotTrimmed := TrimPrefixWhile(tt.args.s, tt.args.prefix); gotTrimmed != tt.wantTrimmed {
				t.Errorf("TrimPrefixWhile() = %v, want %v", gotTrimmed, tt.wantTrimmed)
			}
		})
	}
}
