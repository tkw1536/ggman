package split

import (
	"testing"
)

func TestBefore(t *testing.T) {
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
			gotPrefix, gotSuffix := Before(tt.args.s, tt.args.sep)
			if gotPrefix != tt.wantPrefix {
				t.Errorf("Before() gotPrefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
			if gotSuffix != tt.wantSuffix {
				t.Errorf("Before() gotSuffix = %v, want %v", gotSuffix, tt.wantSuffix)
			}
		})
	}
}

func Benchmark_Before(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Before("a;b", ";")
		Before("a;b;c", ";")
		Before("aaa", ";")
	}
}

func TestAfter(t *testing.T) {
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
			gotPrefix, gotSuffix := After(tt.args.s, tt.args.sep)
			if gotPrefix != tt.wantPrefix {
				t.Errorf("After() gotPrefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
			if gotSuffix != tt.wantSuffix {
				t.Errorf("After() gotSuffix = %v, want %v", gotSuffix, tt.wantSuffix)
			}
		})
	}
}

func Benchmark_After(b *testing.B) {
	for i := 0; i < b.N; i++ {
		After("a;b", ";")
		After("a;b;c", ";")
		After("aaa", ";")
	}
}
