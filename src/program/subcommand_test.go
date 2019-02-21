package program

import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/src/constants"
)

func TestParseArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name       string
		args       args
		wantParsed *SubCommandArgs
		wantErr    string
	}{
		{"no arguments", args{[]string{}}, nil, constants.StringNeedOneArgument},

		{"command without arguments", args{[]string{"cmd"}}, &SubCommandArgs{"cmd", "", false, []string{}}, ""},

		{"command with help (1)", args{[]string{"help", "cmd"}}, &SubCommandArgs{"", "", true, []string{"cmd"}}, ""},
		{"command with help (2)", args{[]string{"--help", "cmd"}}, &SubCommandArgs{"", "", true, []string{"cmd"}}, ""},
		{"command with help (3)", args{[]string{"-h", "cmd"}}, &SubCommandArgs{"", "", true, []string{"cmd"}}, ""},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, &SubCommandArgs{"cmd", "", false, []string{"a1", "a2"}}, ""},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, &SubCommandArgs{"cmd", "", false, []string{"help", "a1"}}, ""},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, &SubCommandArgs{"cmd", "", false, []string{"--help", "a1"}}, ""},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, &SubCommandArgs{"cmd", "", false, []string{"-h", "a1"}}, ""},

		{"only a for (1)", args{[]string{"for"}}, nil, constants.StringNeedTwoAfterFor},
		{"only a for (2)", args{[]string{"--for"}}, nil, constants.StringNeedTwoAfterFor},
		{"only a for (3)", args{[]string{"-f"}}, nil, constants.StringNeedTwoAfterFor},

		{"for without command (1)", args{[]string{"for", "match"}}, nil, constants.StringNeedTwoAfterFor},
		{"for without command (2)", args{[]string{"--for", "match"}}, nil, constants.StringNeedTwoAfterFor},
		{"for without command (3)", args{[]string{"-f", "match"}}, nil, constants.StringNeedTwoAfterFor},

		{"for with command (1)", args{[]string{"for", "match", "cmd"}}, &SubCommandArgs{"cmd", "match", false, []string{}}, ""},
		{"for with command (2)", args{[]string{"--for", "match", "cmd"}}, &SubCommandArgs{"cmd", "match", false, []string{}}, ""},
		{"for with command (3)", args{[]string{"-f", "match", "cmd"}}, &SubCommandArgs{"cmd", "match", false, []string{}}, ""},

		{"for with command and arguments (1)", args{[]string{"for", "match", "cmd", "a1", "a2"}}, &SubCommandArgs{"cmd", "match", false, []string{"a1", "a2"}}, ""},
		{"for with command and arguments (2)", args{[]string{"--for", "match", "cmd", "a1", "a2"}}, &SubCommandArgs{"cmd", "match", false, []string{"a1", "a2"}}, ""},
		{"for with command and arguments (3)", args{[]string{"-f", "match", "cmd", "a1", "a2"}}, &SubCommandArgs{"cmd", "match", false, []string{"a1", "a2"}}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotParsed, gotErr := ParseArgs(tt.args.args)
			if !reflect.DeepEqual(gotParsed, tt.wantParsed) {
				t.Errorf("ParseArgs() gotParsed = %v, want %v", gotParsed, tt.wantParsed)
			}
			if gotErr != tt.wantErr {
				t.Errorf("ParseArgs() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGGArgs_ParseSingleFlag(t *testing.T) {
	type fields struct {
		Command string
		Pattern string
		Args    []string
	}
	type args struct {
		flag string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantValue  bool
		wantRetval int
		wantErr    string
	}{
		// giving no arguments
		{"no arguments given", fields{"cmd", "", []string{}}, args{"--test"}, false, 0, ""},
		{"right argument given", fields{"cmd", "", []string{"--test"}}, args{"--test"}, true, 0, ""},
		{"wrong argument given", fields{"cmd", "", []string{"--fake"}}, args{"--test"}, false, constants.ErrorSpecificParseArgs, "Unknown argument: 'cmd' must be called with either '--test' or no arguments. "},
		{"too many arguments", fields{"cmd", "", []string{"--fake", "--untrue"}}, args{"--test"}, false, constants.ErrorSpecificParseArgs, "Unknown argument: 'cmd' must be called with either '--test' or no arguments. "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &SubCommandArgs{
				Command: tt.fields.Command,
				Pattern: tt.fields.Pattern,
				args:    tt.fields.Args,
			}
			gotValue, gotRetval, gotErr := parsed.ParseSingleFlag(tt.args.flag)
			if gotValue != tt.wantValue {
				t.Errorf("GGArgs.ParseSingleFlag() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotRetval != tt.wantRetval {
				t.Errorf("GGArgs.ParseSingleFlag() gotRetval = %v, want %v", gotRetval, tt.wantRetval)
			}
			if gotErr != tt.wantErr {
				t.Errorf("GGArgs.ParseSingleFlag() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestGGArgs_EnsureNoFor(t *testing.T) {
	type fields struct {
		Command string
		Pattern string
		Help    bool
		Args    []string
	}
	tests := []struct {
		name       string
		fields     fields
		wantRetval int
		wantErr    string
	}{
		{"no for", fields{"example", "", false, []string{}}, 0, ""},
		{"provided filter", fields{"example", "test", false, []string{}}, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes no 'for' argument. "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &SubCommandArgs{
				Command: tt.fields.Command,
				Pattern: tt.fields.Pattern,
				Help:    tt.fields.Help,
				args:    tt.fields.Args,
			}
			gotRetval, gotErr := parsed.EnsureNoFor()
			if gotRetval != tt.wantRetval {
				t.Errorf("GGArgs.EnsureNoFor() gotRetval = %v, want %v", gotRetval, tt.wantRetval)
			}
			if gotErr != tt.wantErr {
				t.Errorf("GGArgs.EnsureNoFor() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestSubCommandArgs_EnsureNoArguments(t *testing.T) {
	type fields struct {
		Command string
		Pattern string
		Help    bool
		Args    []string
	}
	tests := []struct {
		name       string
		fields     fields
		wantRetval int
		wantErr    string
	}{
		{"no arguments", fields{"example", "", false, []string{}}, 0, ""},
		{"only a for", fields{"example", "filter", false, []string{}}, 0, ""},

		{"some arguments", fields{"example", "", false, []string{"hello"}}, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes no arguments. "},
		{"arguments and a for", fields{"example", "filter", false, []string{"hello"}}, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes no arguments. "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &SubCommandArgs{
				Command: tt.fields.Command,
				Pattern: tt.fields.Pattern,
				Help:    tt.fields.Help,
				args:    tt.fields.Args,
			}
			gotRetval, gotErr := parsed.EnsureNoArguments()
			if gotRetval != tt.wantRetval {
				t.Errorf("SubCommandArgs.EnsureNoArguments() gotRetval = %v, want %v", gotRetval, tt.wantRetval)
			}
			if gotErr != tt.wantErr {
				t.Errorf("SubCommandArgs.EnsureNoArguments() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}

func TestSubCommandArgs_EnsureArguments(t *testing.T) {
	type fields struct {
		Command string
		Pattern string
		Help    bool
		Args    []string
	}
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantArgc   int
		wantArgv   []string
		wantRetval int
		wantErr    string
	}{
		// taking 0 args
		{"no arguments", fields{"example", "", false, []string{}}, args{0, 0}, 0, []string{}, 0, ""},

		// taking 1 arg
		{"one argument, too few", fields{"example", "", false, []string{}}, args{1, 1}, 0, nil, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes exactly 1 argument(s). "},
		{"one argument, exactly enough", fields{"example", "", false, []string{"world"}}, args{1, 1}, 1, []string{"world"}, 0, ""},
		{"one argument, too many", fields{"example", "", false, []string{"hello", "world"}}, args{1, 1}, 0, nil, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes exactly 1 argument(s). "},

		// taking 1 or 2 args
		{"1-2 arguments, too few", fields{"example", "", false, []string{}}, args{1, 2}, 0, nil, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes between 1 and 2 arguments. "},
		{"1-2 arguments, enough", fields{"example", "", false, []string{"world"}}, args{1, 2}, 1, []string{"world"}, 0, ""},
		{"1-2 arguments, enough (2)", fields{"example", "", false, []string{"hello", "world"}}, args{1, 2}, 2, []string{"hello", "world"}, 0, ""},
		{"1-2 arguments, too many", fields{"example", "", false, []string{"hello", "world", "you"}}, args{1, 2}, 0, nil, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes between 1 and 2 arguments. "},

		{"2 arguments, too few", fields{"example", "", false, []string{}}, args{2, 2}, 0, nil, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes exactly 2 argument(s). "},
		{"2 arguments, too few", fields{"example", "", false, []string{"world"}}, args{2, 2}, 0, nil, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes exactly 2 argument(s). "},
		{"2 arguments, enough (2)", fields{"example", "", false, []string{"hello", "world"}}, args{2, 2}, 2, []string{"hello", "world"}, 0, ""},
		{"2 arguments, too many", fields{"example", "", false, []string{"hello", "world", "you"}}, args{2, 2}, 0, nil, constants.ErrorSpecificParseArgs, "Wrong number of arguments: 'example' takes exactly 2 argument(s). "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &SubCommandArgs{
				Command: tt.fields.Command,
				Pattern: tt.fields.Pattern,
				Help:    tt.fields.Help,
				args:    tt.fields.Args,
			}
			gotArgc, gotArgv, gotRetval, gotErr := parsed.EnsureArguments(tt.args.min, tt.args.max)
			if gotArgc != tt.wantArgc {
				t.Errorf("SubCommandArgs.EnsureArguments() gotArgc = %v, want %v", gotArgc, tt.wantArgc)
			}
			if !reflect.DeepEqual(gotArgv, tt.wantArgv) {
				t.Errorf("SubCommandArgs.EnsureArguments() gotArgv = %v, want %v", gotArgv, tt.wantArgv)
			}
			if gotRetval != tt.wantRetval {
				t.Errorf("SubCommandArgs.EnsureArguments() gotRetval = %v, want %v", gotRetval, tt.wantRetval)
			}
			if gotErr != tt.wantErr {
				t.Errorf("SubCommandArgs.EnsureArguments() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
