//spellchecker:words split
package split

//spellchecker:words testing
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

func TestAfterRune(t *testing.T) {
	type args struct {
		s   string
		sep rune
	}
	tests := []struct {
		name       string
		args       args
		wantPrefix string
		wantSuffix string
	}{
		{"splitFoundOnce", args{"a;b", ';'}, "a", "b"},
		{"splitFoundMultiple", args{"a;b;c", ';'}, "a", "b;c"},
		{"splitNotFound", args{"aaa", ';'}, "aaa", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPrefix, gotSuffix := AfterRune(tt.args.s, tt.args.sep)
			if gotPrefix != tt.wantPrefix {
				t.Errorf("AfterRune() gotPrefix = %v, want %v", gotPrefix, tt.wantPrefix)
			}
			if gotSuffix != tt.wantSuffix {
				t.Errorf("AfterRune() gotSuffix = %v, want %v", gotSuffix, tt.wantSuffix)
			}
		})
	}
}

func Benchmark_AfterRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		AfterRune("a;b", ';')
		AfterRune("a;b;c", ';')
		AfterRune("aaa", ';')
	}
}
