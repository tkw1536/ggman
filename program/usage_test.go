package program

import (
	"testing"

	"github.com/spf13/pflag"
)

func Test_flagUsage_Flag(t *testing.T) {
	tests := []struct {
		name string
		flag pflag.Flag
		want string
	}{
		{"long flag with description", pflag.Flag{Name: "long", Usage: "a `random` argument", Value: usageFakeValue}, "--long random"},
		{"long flag without description", pflag.Flag{Name: "long", Usage: "a very long thing", Value: usageFakeValue}, "--long mock/type"},

		{"short flag without description", pflag.Flag{Shorthand: "s", Name: "long", Usage: "a `random` argument", Value: usageFakeValue}, "--long|-s random"},
		{"short flag without description", pflag.Flag{Shorthand: "s", Name: "long", Usage: "a very long thing", Value: usageFakeValue}, "--long|-s mock/type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := flagUsage{flag: &tt.flag}
			if got := f.Flag(); got != tt.want {
				t.Errorf("flagUsage.Flag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flagUsage_Description(t *testing.T) {
	tests := []struct {
		name string
		flag pflag.Flag
		want string
	}{
		{"long flag with description and no default", pflag.Flag{Name: "long", Usage: "a `random` argument", Value: usageFakeValue}, "\n\n   --long random\n       a random argument"},
		{"long flag without description and no default", pflag.Flag{Name: "long", Usage: "a very long thing", Value: usageFakeValue}, "\n\n   --long mock/type\n       a very long thing"},

		{"long flag with description and a default", pflag.Flag{Name: "long", Usage: "a `random` argument", Value: usageFakeValueDefault}, "\n\n   --long random\n       a random argument (default default)"},
		{"long flag without description and a default", pflag.Flag{Name: "long", Usage: "a very long thing", Value: usageFakeValueDefault}, "\n\n   --long mock/type\n       a very long thing (default default)"},

		{"short flag without description and no default", pflag.Flag{Shorthand: "s", Name: "long", Usage: "a `random` argument", Value: usageFakeValue}, "\n\n   -s, --long random\n       a random argument"},
		{"short flag without description and no default", pflag.Flag{Shorthand: "s", Name: "long", Usage: "a very long thing", Value: usageFakeValue}, "\n\n   -s, --long mock/type\n       a very long thing"},

		{"short flag without description and a default", pflag.Flag{Shorthand: "s", Name: "long", Usage: "a `random` argument", Value: usageFakeValueDefault}, "\n\n   -s, --long random\n       a random argument (default default)"},
		{"short flag without description and a default", pflag.Flag{Shorthand: "s", Name: "long", Usage: "a very long thing", Value: usageFakeValueDefault}, "\n\n   -s, --long mock/type\n       a very long thing (default default)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := flagUsage{flag: &tt.flag}
			f.flag.DefValue = f.flag.Value.String() // is usually done by .VarP()
			if got := f.Description(); got != tt.want {
				t.Errorf("flagUsage.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}

var usageFakeValue pflag.Value = usageFakeValueT("")
var usageFakeValueDefault pflag.Value = usageFakeValueT("default")

// usageFakeValueT is a type used to mock a value for tests
type usageFakeValueT string

func (usageFakeValueT) Set(string) error { return nil }
func (v usageFakeValueT) String() string { return string(v) }
func (usageFakeValueT) Type() string     { return "mock/type" }
