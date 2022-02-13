package text

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
