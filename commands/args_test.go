package commands

import (
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/constants"
)

func TestParseArgs(t *testing.T) {
	type args struct {
		args []string
	}
	tests := []struct {
		name       string
		args       args
		wantParsed *GGArgs
		wantErr    string
	}{
		{"no arguments", args{[]string{}}, nil, constants.StringNeedOneArgument},

		{"command without arguments", args{[]string{"cmd"}}, &GGArgs{"cmd", "", false, []string{}}, ""},

		{"command with help (1)", args{[]string{"help", "cmd"}}, &GGArgs{"", "", true, []string{"cmd"}}, ""},
		{"command with help (2)", args{[]string{"--help", "cmd"}}, &GGArgs{"", "", true, []string{"cmd"}}, ""},
		{"command with help (3)", args{[]string{"-h", "cmd"}}, &GGArgs{"", "", true, []string{"cmd"}}, ""},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, &GGArgs{"cmd", "", false, []string{"a1", "a2"}}, ""},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, &GGArgs{"cmd", "", false, []string{"help", "a1"}}, ""},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, &GGArgs{"cmd", "", false, []string{"--help", "a1"}}, ""},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, &GGArgs{"cmd", "", false, []string{"-h", "a1"}}, ""},

		{"only a for (1)", args{[]string{"for"}}, nil, constants.StringNeedTwoAfterFor},
		{"only a for (2)", args{[]string{"--for"}}, nil, constants.StringNeedTwoAfterFor},
		{"only a for (3)", args{[]string{"-f"}}, nil, constants.StringNeedTwoAfterFor},

		{"for without command (1)", args{[]string{"for", "match"}}, nil, constants.StringNeedTwoAfterFor},
		{"for without command (2)", args{[]string{"--for", "match"}}, nil, constants.StringNeedTwoAfterFor},
		{"for without command (3)", args{[]string{"-f", "match"}}, nil, constants.StringNeedTwoAfterFor},

		{"for with command (1)", args{[]string{"for", "match", "cmd"}}, &GGArgs{"cmd", "match", false, []string{}}, ""},
		{"for with command (2)", args{[]string{"--for", "match", "cmd"}}, &GGArgs{"cmd", "match", false, []string{}}, ""},
		{"for with command (3)", args{[]string{"-f", "match", "cmd"}}, &GGArgs{"cmd", "match", false, []string{}}, ""},

		{"for with command and arguments (1)", args{[]string{"for", "match", "cmd", "a1", "a2"}}, &GGArgs{"cmd", "match", false, []string{"a1", "a2"}}, ""},
		{"for with command and arguments (2)", args{[]string{"--for", "match", "cmd", "a1", "a2"}}, &GGArgs{"cmd", "match", false, []string{"a1", "a2"}}, ""},
		{"for with command and arguments (3)", args{[]string{"-f", "match", "cmd", "a1", "a2"}}, &GGArgs{"cmd", "match", false, []string{"a1", "a2"}}, ""},
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
		name      string
		fields    fields
		args      args
		wantValue bool
		wantErr   bool
	}{
		// giving no arguments
		{"no arguments given", fields{"cmd", "", []string{}}, args{"--test"}, false, false},
		{"right argument given", fields{"cmd", "", []string{"--test"}}, args{"--test"}, true, false},
		{"wrong argument given", fields{"cmd", "", []string{"--fake"}}, args{"--test"}, false, true},
		{"too many arguments", fields{"cmd", "", []string{"--fake", "--untrue"}}, args{"--test"}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := &GGArgs{
				Command: tt.fields.Command,
				Pattern: tt.fields.Pattern,
				Args:    tt.fields.Args,
			}
			gotValue, gotErr := parsed.ParseSingleFlag(tt.args.flag)
			if gotValue != tt.wantValue {
				t.Errorf("GGArgs.ParseSingleFlag() gotValue = %v, want %v", gotValue, tt.wantValue)
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
			parsed := &GGArgs{
				Command: tt.fields.Command,
				Pattern: tt.fields.Pattern,
				Help:    tt.fields.Help,
				Args:    tt.fields.Args,
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
