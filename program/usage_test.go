package program

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/tkw1536/ggman/env"
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

func (usageFakeValueT) Set(string) error { panic("usageFakeValueT: not implemented") }
func (v usageFakeValueT) String() string { return string(v) }
func (usageFakeValueT) Type() string     { return "mock/type" }

func TestProgram_knownCommandsString(t *testing.T) {
	got := usageFakeProgram.knownCommandsString()
	want := "'a', 'b', 'c'"
	if got != want {
		t.Errorf("Program.knownCommandsString() = %v, want %v", got, want)
	}
}

func TestProgram_Usage(t *testing.T) {
	flagset := pflag.NewFlagSet("ggman", pflag.ContinueOnError)
	flagset.BoolP("bool", "b", false, "a `random` boolean argument with short")
	flagset.Int("int", 12, "a `dummy` integer flag")

	got := usageFakeProgram.Usage(flagset)
	want := "ggman version v0.0.0-unknown\n\nUsage: ggman [--bool|-b random] [--int dummy] [--] COMMAND [ARGS...]\n\n   -b, --bool random\n       a random boolean argument with short\n\n   --int dummy\n       a dummy integer flag (default 12)\n\n   COMMAND [ARGS...]\n       Command to call. One of 'a', 'b', 'c'. See individual commands for more help.\n\nggman is licensed under the terms of the MIT License. Use 'ggman license' to view licensing information."
	if got != want {
		t.Errorf("Program.Usage() = %v, want %v", got, want)
	}
}

func TestOptions_Usage(t *testing.T) {

	flagset := pflag.NewFlagSet("ggman", pflag.ContinueOnError)
	flagset.BoolP("bool", "b", false, "a `random` boolean argument with short")
	flagset.Int("int", 12, "a `dummy` integer flag")

	type fields struct {
		Environment      env.Requirement
		MinArgs          int
		MaxArgs          int
		Metavar          string
		UsageDescription string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantUsage string
	}{
		{
			"command without args and allowing filter",
			fields{Environment: env.Requirement{AllowsFilter: true}, UsageDescription: "usage"},
			args{"a"},
			"Usage: ggman [--for|-f FILTER] [--here|-H] [global arguments] [--] a [--help|-h] [--bool|-b random] [--int dummy]\n\n   -h, --help\n       Print this usage message and exit.\n\n   -f, --for filter\n       Filter the list of repositories to apply command to by FILTER.\n\n   -h, --here\n       Filter the list of repositories to apply command to only contain the current repository. \n\n   global arguments\n       Global arguments for ggman. See ggman --help for more information.\n\n   -b, --bool random\n       a random boolean argument with short\n\n   --int dummy\n       a dummy integer flag (default 12)\n\n   \n       usage",
		},

		{
			"command without args and not allowing filter",
			fields{Environment: env.Requirement{}, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [global arguments] [--] a [--help|-h] [--bool|-b random] [--int dummy]\n\n   -h, --help\n       Print this usage message and exit.\n\n   global arguments\n       Global arguments for ggman. See ggman --help for more information.\n\n   -b, --bool random\n       a random boolean argument with short\n\n   --int dummy\n       a dummy integer flag (default 12)\n\n   \n       usage",
		},

		{
			"command with max finite args",
			fields{Environment: env.Requirement{}, MaxArgs: 4, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [global arguments] [--] a [--help|-h] [--bool|-b random] [--int dummy] [META [META [META [META]]]]\n\n   -h, --help\n       Print this usage message and exit.\n\n   global arguments\n       Global arguments for ggman. See ggman --help for more information.\n\n   -b, --bool random\n       a random boolean argument with short\n\n   --int dummy\n       a dummy integer flag (default 12)\n\n   [META [META [META [META]]]]\n       usage",
		},

		{
			"command with finite args",
			fields{Environment: env.Requirement{}, MinArgs: 1, MaxArgs: 2, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [global arguments] [--] a [--help|-h] [--bool|-b random] [--int dummy] [--] META [META]\n\n   -h, --help\n       Print this usage message and exit.\n\n   global arguments\n       Global arguments for ggman. See ggman --help for more information.\n\n   -b, --bool random\n       a random boolean argument with short\n\n   --int dummy\n       a dummy integer flag (default 12)\n\n   META [META]\n       usage",
		},

		{
			"command with infinite args",
			fields{Environment: env.Requirement{}, MinArgs: 1, MaxArgs: -1, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [global arguments] [--] a [--help|-h] [--bool|-b random] [--int dummy] [--] META [META ... ]\n\n   -h, --help\n       Print this usage message and exit.\n\n   global arguments\n       Global arguments for ggman. See ggman --help for more information.\n\n   -b, --bool random\n       a random boolean argument with short\n\n   --int dummy\n       a dummy integer flag (default 12)\n\n   META [META ... ]\n       usage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt := Options{
				Environment:      tt.fields.Environment,
				MinArgs:          tt.fields.MinArgs,
				MaxArgs:          tt.fields.MaxArgs,
				Metavar:          tt.fields.Metavar,
				UsageDescription: tt.fields.UsageDescription,
			}
			if gotUsage := opt.Usage(tt.args.name, flagset); gotUsage != tt.wantUsage {
				t.Errorf("Options.Usage() = %q, want %q", gotUsage, tt.wantUsage)
			}
		})
	}
}

// usageFakeProgram represents a fake 'program' instance that has three dummy commands
// 'a', 'c' and 'b'
var usageFakeProgram Program

func init() {
	usageFakeProgram.Register(usageFakeCommandT("a"))
	usageFakeProgram.Register(usageFakeCommandT("c"))
	usageFakeProgram.Register(usageFakeCommandT("b"))
}

// usageFakeCommandT represents a dummy command with the given name
type usageFakeCommandT string

func (u usageFakeCommandT) Name() string                 { return string(u) }
func (usageFakeCommandT) Options(*pflag.FlagSet) Options { panic("usageFakeCommandT: not implemented") }
func (usageFakeCommandT) AfterParse() error              { panic("usageFakeCommandT: not implemented") }
func (usageFakeCommandT) Run(Context) error              { panic("usageFakeCommandT: not implemented") }
