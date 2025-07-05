//spellchecker:words split
package split_test

//spellchecker:words testing github ggman internal split
import (
	"testing"

	"go.tkw01536.de/ggman/internal/split"
)

func TestAfterRune(t *testing.T) {
	t.Parallel()

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
			t.Parallel()

			gotPrefix, gotSuffix := split.AfterRune(tt.args.s, tt.args.sep)
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
	for b.Loop() {
		split.AfterRune("a;b", ';')
		split.AfterRune("a;b;c", ';')
		split.AfterRune("aaa", ';')
	}
}
