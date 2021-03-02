package usagefmt

import (
	"testing"
)

func TestFmtSpecShort(t *testing.T) {
	tests := []struct {
		name string
		opt  Opt
		want string
	}{
		{
			"long only optional option",
			FakeOpt{long: []string{"long"}},
			"[--long]",
		},
		{
			"short and long optional option",
			FakeOpt{short: []string{"s"}, long: []string{"long"}},
			"[--long|-s]",
		},
		{
			"short and long named optional option",
			FakeOpt{value: "name", short: []string{"s"}, long: []string{"long"}},
			"[--long|-s name]",
		},

		{
			"long only required option",
			FakeOpt{long: []string{"long"}, required: true},
			"--long",
		},
		{
			"short and long required option",
			FakeOpt{short: []string{"s"}, long: []string{"long"}, required: true},
			"--long|-s",
		},
		{
			"short and long named required option",
			FakeOpt{value: "name", short: []string{"s"}, long: []string{"long"}, required: true},
			"--long|-s name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FmtSpecShort(tt.opt); got != tt.want {
				t.Errorf("FmtSpecShort() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFmtSpecLong(t *testing.T) {
	tests := []struct {
		name string
		opt  Opt
		want string
	}{
		{
			"long only option",
			FakeOpt{long: []string{"long"}},
			"--long",
		},
		{
			"short and long option",
			FakeOpt{short: []string{"s"}, long: []string{"long"}},
			"-s, --long",
		},
		{
			"short and long named option",
			FakeOpt{value: "name", short: []string{"s"}, long: []string{"long"}},
			"-s, --long name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FmtSpecLong(tt.opt); got != tt.want {
				t.Errorf("FmtSpecLong() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFmtMessage(t *testing.T) {
	tests := []struct {
		name string
		opt  Opt
		want string
	}{
		{
			"long only option",
			FakeOpt{usage: "a long option", long: []string{"long"}},
			"\n\n   --long\n      a long option",
		},
		{
			"short and long option",
			FakeOpt{usage: "a long or short option", short: []string{"s"}, long: []string{"long"}},
			"\n\n   -s, --long\n      a long or short option",
		},
		{
			"short and long named option",
			FakeOpt{usage: "this one is named", value: "name", short: []string{"s"}, long: []string{"long"}},
			"\n\n   -s, --long name\n      this one is named",
		},
		{
			"short and long named option with default",
			FakeOpt{usage: "this one is named", value: "name", short: []string{"s"}, long: []string{"long"}, dflt: "default"},
			"\n\n   -s, --long name\n      this one is named (default default)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FmtMessage(tt.opt); got != tt.want {
				t.Errorf("FmtMessage() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFmtSpecPositional(t *testing.T) {
	type args struct {
		name string
		min  int
		max  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"arg 0, 0", args{"arg", 0, 0}, ""},
		{"arg 0, -1", args{"arg", 0, -1}, "[arg ...]"},
		{"arg 0, 3", args{"arg", 0, 3}, "[arg [arg [arg]]]"},

		{"no name 0, 0", args{"", 0, 0}, ""},
		{"no name 0, -1", args{"", 0, -1}, "[ARGUMENT ...]"},
		{"no name 0, 3", args{"", 0, 3}, "[ARGUMENT [ARGUMENT [ARGUMENT]]]"},

		{"arg 2, 2", args{"arg", 2, 2}, "arg arg"},
		{"arg 2, 4", args{"arg", 2, 4}, "arg arg [arg [arg]]"},
		{"arg 2, -1", args{"arg", 2, -1}, "arg arg [arg ...]"},

		{"no name 2, 2", args{"", 2, 2}, "ARGUMENT ARGUMENT"},
		{"no name 2, 4", args{"", 2, 4}, "ARGUMENT ARGUMENT [ARGUMENT [ARGUMENT]]"},
		{"no name 2, -1", args{"", 2, -1}, "ARGUMENT ARGUMENT [ARGUMENT ...]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FmtSpecPositional(tt.args.name, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("FmtSpecPositional() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFmtCommands(t *testing.T) {
	type args struct {
		commands []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no commands", args{nil}, ""},
		{"single command", args{[]string{"a"}}, `"a"`},
		{"multiple commands", args{[]string{"a", "b", "c"}}, `"a", "b", "c"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FmtCommands(tt.args.commands); got != tt.want {
				t.Errorf("FmtCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}
