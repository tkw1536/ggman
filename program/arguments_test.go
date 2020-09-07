package program

import (
	"reflect"
	"testing"
)

func TestArguments_Parse(t *testing.T) {
	type args struct {
		argv []string
	}
	tests := []struct {
		name       string
		args       args
		wantParsed Arguments
		wantErr    error
	}{
		{"no arguments", args{[]string{}}, Arguments{}, errParseArgsNeedOneArgument},

		{"command without arguments", args{[]string{"cmd"}}, Arguments{"cmd", "", false, false, []string{}}, nil},

		{"help with command (1)", args{[]string{"help", "cmd"}}, Arguments{"", "", true, false, []string{"cmd"}}, nil},
		{"help with command (2)", args{[]string{"--help", "cmd"}}, Arguments{"", "", true, false, []string{"cmd"}}, nil},
		{"help with command (3)", args{[]string{"-h", "cmd"}}, Arguments{"", "", true, false, []string{"cmd"}}, nil},

		{"version with command (1)", args{[]string{"version", "cmd"}}, Arguments{"", "", false, true, []string{"cmd"}}, nil},
		{"version with command (2)", args{[]string{"--version", "cmd"}}, Arguments{"", "", false, true, []string{"cmd"}}, nil},
		{"version with command (3)", args{[]string{"-v", "cmd"}}, Arguments{"", "", false, true, []string{"cmd"}}, nil},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, Arguments{"cmd", "", false, false, []string{"a1", "a2"}}, nil},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, Arguments{"cmd", "", false, false, []string{"help", "a1"}}, nil},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, Arguments{"cmd", "", false, false, []string{"--help", "a1"}}, nil},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, Arguments{"cmd", "", false, false, []string{"-h", "a1"}}, nil},

		{"command with version (1)", args{[]string{"cmd", "version", "a1"}}, Arguments{"cmd", "", false, false, []string{"version", "a1"}}, nil},
		{"command with version (2)", args{[]string{"cmd", "--version", "a1"}}, Arguments{"cmd", "", false, false, []string{"--version", "a1"}}, nil},
		{"command with version (3)", args{[]string{"cmd", "-v", "a1"}}, Arguments{"cmd", "", false, false, []string{"-v", "a1"}}, nil},

		{"only a for (1)", args{[]string{"for"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"only a for (2)", args{[]string{"--for"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"only a for (3)", args{[]string{"-f"}}, Arguments{}, errParseArgsNeedTwoAfterFor},

		{"for without command (1)", args{[]string{"for", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"for without command (2)", args{[]string{"--for", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"for without command (3)", args{[]string{"-f", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},

		{"for with command (1)", args{[]string{"for", "match", "cmd"}}, Arguments{"cmd", "match", false, false, []string{}}, nil},
		{"for with command (2)", args{[]string{"--for", "match", "cmd"}}, Arguments{"cmd", "match", false, false, []string{}}, nil},
		{"for with command (3)", args{[]string{"-f", "match", "cmd"}}, Arguments{"cmd", "match", false, false, []string{}}, nil},

		{"for with command and arguments (1)", args{[]string{"for", "match", "cmd", "a1", "a2"}}, Arguments{"cmd", "match", false, false, []string{"a1", "a2"}}, nil},
		{"for with command and arguments (2)", args{[]string{"--for", "match", "cmd", "a1", "a2"}}, Arguments{"cmd", "match", false, false, []string{"a1", "a2"}}, nil},
		{"for with command and arguments (3)", args{[]string{"-f", "match", "cmd", "a1", "a2"}}, Arguments{"cmd", "match", false, false, []string{"a1", "a2"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &Arguments{}
			if err := args.Parse(tt.args.argv); err != tt.wantErr {
				t.Errorf("Arguments.Parse() error = %#v, wantErr %#v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(args, &tt.wantParsed) {
				t.Errorf("Arguments.Parse() args = %#v, wantArgs %#v", args, &tt.wantParsed)
			}
		})
	}
}
