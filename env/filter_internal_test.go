package env

//spellchecker:words reflect testing ggman internal pattern
import (
	"reflect"
	"testing"

	"go.tkw01536.de/ggman/internal/pattern"
)

func TestNewPatternFilter(t *testing.T) {
	t.Parallel()

	type args struct {
		value string
		fuzzy bool
	}
	tests := []struct {
		name    string
		args    args
		wantPat PatternFilter
	}{
		{
			"a/b (non-fuzzy)",
			args{"a/b", false},
			PatternFilter{
				value:   "a/b",
				pattern: pattern.NewSplitGlobPattern("a/b", ComponentsOf, false),
			},
		},
		{
			"'' (non-fuzzy)",
			args{"", false},
			PatternFilter{
				value:   "",
				pattern: pattern.NewSplitGlobPattern("", ComponentsOf, false),
			},
		},

		{
			"a/b (fuzzy)",
			args{"a/b", true},
			PatternFilter{
				value:   "a/b",
				fuzzy:   true,
				pattern: pattern.NewSplitGlobPattern("a/b", ComponentsOf, true),
			},
		},
		{
			"'' (fuzzy)",
			args{"", true},
			PatternFilter{
				value:   "",
				fuzzy:   true,
				pattern: pattern.NewSplitGlobPattern("", ComponentsOf, true),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotPat := NewPatternFilter(tt.args.value, tt.args.fuzzy)

			// .Split cannot be compared with reflect
			gotPat.pattern.Split = nil
			tt.wantPat.pattern.Split = nil

			if !reflect.DeepEqual(gotPat, tt.wantPat) {
				t.Errorf("NewPatternFilter() = %v, want %v", gotPat, tt.wantPat)
			}
		})
	}
}

func TestPatternFilter_String(t *testing.T) {
	t.Parallel()

	type fields struct {
		value   string
		pattern pattern.SplitPattern
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"empty pattern",
			fields{
				value:   "",
				pattern: pattern.SplitPattern{},
			},
			"",
		},
		{
			"a/b pattern",
			fields{
				value:   "a/b",
				pattern: pattern.SplitPattern{},
			},
			"a/b",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pat := PatternFilter{
				value:   tt.fields.value,
				pattern: tt.fields.pattern,
			}
			if got := pat.String(); got != tt.want {
				t.Errorf("PatternFilter.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
