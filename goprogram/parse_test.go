package goprogram

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/goprogram/meta"
)

func Test_parseFlagNames(t *testing.T) {
	type args struct {
		err *flags.Error
	}
	tests := []struct {
		name      string
		args      args
		wantNames []string
		wantOk    bool
	}{
		{
			"unix flag message",
			args{
				&flags.Error{
					Message: "expected argument for flag `-f, --for'",
				},
			},
			[]string{"f", "for"},
			true,
		},
		{
			"windows flag message",
			args{
				&flags.Error{
					Message: "expected argument for flag `/f, /for'",
				},
			},
			[]string{"f", "for"},
			true,
		},

		{
			"no message",
			args{
				&flags.Error{
					Message: "",
				},
			},
			nil,
			false,
		},

		{
			"only beginning",
			args{
				&flags.Error{
					Message: "expected argument for flag `-f, --for",
				},
			},
			nil,
			false,
		},

		{
			"only end",
			args{
				&flags.Error{
					Message: "expected argument for flag -f, --for'",
				},
			},
			nil,
			false,
		},

		{
			"wrong order",
			args{
				&flags.Error{
					Message: "expected argument for flag '-f, --for`",
				},
			},
			nil,
			false,
		},

		{
			"empty flags",
			args{
				&flags.Error{
					Message: "expected argument for flag `'",
				},
			},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNames, gotOk := parseFlagNames(tt.args.err)
			if !reflect.DeepEqual(gotNames, tt.wantNames) {
				t.Errorf("parseFlagNames() gotNames = %v, want %v", gotNames, tt.wantNames)
			}
			if gotOk != tt.wantOk {
				t.Errorf("parseFlagNames() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestArguments_checkRequirements(t *testing.T) {

	// define requirements to allow only the Global1 (or any) arguments
	reqOne := tRequirements(func(flag meta.Flag) bool {
		return flag.FieldName == "Global1"
	})

	// define requirements to allow anything
	reqAny := tRequirements(func(flag meta.Flag) bool { return true })

	tests := []struct {
		name  string
		reqs  tRequirements
		flags tFlags

		wantErr string
	}{
		{
			"global not allowed, global not given",
			reqOne,
			tFlags{},
			"",
		},
		{
			"global not allowed, global given",
			reqOne,
			tFlags{GlobalOne: "global1"},
			"Wrong number of arguments: 'echo' takes no '--global-one' argument. ",
		},

		{
			"global allowed, global not given",
			reqAny, tFlags{},
			"",
		},
		{
			"global allowed, global given",
			reqAny, tFlags{GlobalOne: "global1"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &iContext{
				Description: iDescription{
					Requirements: tt.reqs,
				},
				Args: iArguments{
					Command: "echo",
					Flags:   tt.flags,
				},
			}
			err := context.checkRequirements()
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("Context.checkRequirements() error = %q, wantErr %q", err, tt.wantErr)
			}
		})
	}
}

func TestArguments_parseP(t *testing.T) {
	type args struct {
		argv []string
	}
	tests := []struct {
		name       string
		args       args
		wantParsed iArguments
		wantErr    error
	}{
		{"no arguments", args{[]string{}}, iArguments{}, errParseArgsNeedOneArgument},
		{"command without arguments", args{[]string{"cmd"}}, iArguments{Command: "cmd", Pos: []string{}}, nil},

		{"help with command (1)", args{[]string{"--help", "cmd"}}, iArguments{Universals: Universals{Help: true}, Pos: []string{"cmd"}}, nil},
		{"help with command (2)", args{[]string{"-h", "cmd"}}, iArguments{Universals: Universals{Help: true}, Pos: []string{"cmd"}}, nil},

		{"help without command (1)", args{[]string{"--help"}}, iArguments{Universals: Universals{Help: true}, Pos: []string{}}, nil},
		{"help without command (2)", args{[]string{"-h"}}, iArguments{Universals: Universals{Help: true}, Pos: []string{}}, nil},

		{"version with command (1)", args{[]string{"--version", "cmd"}}, iArguments{Universals: Universals{Version: true}, Pos: []string{"cmd"}}, nil},
		{"version with command (2)", args{[]string{"-v", "cmd"}}, iArguments{Universals: Universals{Version: true}, Pos: []string{"cmd"}}, nil},

		{"version without command (2)", args{[]string{"--version"}}, iArguments{Universals: Universals{Version: true}, Pos: []string{}}, nil},
		{"version without command (3)", args{[]string{"-v"}}, iArguments{Universals: Universals{Version: true}, Pos: []string{}}, nil},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, iArguments{Command: "cmd", Pos: []string{"a1", "a2"}}, nil},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, iArguments{Command: "cmd", Pos: []string{"help", "a1"}}, nil},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, iArguments{Command: "cmd", Pos: []string{"--help", "a1"}}, nil},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, iArguments{Command: "cmd", Pos: []string{"-h", "a1"}}, nil},

		{"command with version (1)", args{[]string{"cmd", "version", "a1"}}, iArguments{Command: "cmd", Pos: []string{"version", "a1"}}, nil},
		{"command with version (2)", args{[]string{"cmd", "--version", "a1"}}, iArguments{Command: "cmd", Pos: []string{"--version", "a1"}}, nil},
		{"command with version (3)", args{[]string{"cmd", "-v", "a1"}}, iArguments{Command: "cmd", Pos: []string{"-v", "a1"}}, nil},

		{"global flag without command (1)", args{[]string{"-a", "stuff"}}, iArguments{}, errParseArgsNeedOneArgument},
		{"global flag without command (2)", args{[]string{"--global-one", "stuff"}}, iArguments{}, errParseArgsNeedOneArgument},

		{"global flag with command (1)", args{[]string{"-a", "stuff", "cmd"}}, iArguments{Command: "cmd", Flags: tFlags{GlobalOne: "stuff"}, Pos: []string{}}, nil},
		{"global flag with command (2)", args{[]string{"--global-one", "stuff", "cmd"}}, iArguments{Command: "cmd", Flags: tFlags{GlobalOne: "stuff"}, Pos: []string{}}, nil},

		{"global flag with command and arguments (1)", args{[]string{"--global-two", "stuff", "cmd", "a1", "a2"}}, iArguments{Command: "cmd", Flags: tFlags{GlobalTwo: "stuff"}, Pos: []string{"a1", "a2"}}, nil},
		{"global flag with command and arguments (2)", args{[]string{"-b", "stuff", "cmd", "a1", "a2"}}, iArguments{Command: "cmd", Flags: tFlags{GlobalTwo: "stuff"}, Pos: []string{"a1", "a2"}}, nil},

		{"global looking flag", args{[]string{"--not-a-global-flag", "stuff", "command"}}, iArguments{}, errParseArgsUnknownError.WithMessageF("unknown flag `not-a-global-flag'")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var args iArguments
			err := args.parseP(tt.args.argv)

			// turn wantErr into a string
			var wantErr string
			if tt.wantErr != nil {
				wantErr = tt.wantErr.Error()
			}

			// turn gotErr into a string
			var gotErr string
			if err != nil {
				gotErr = err.Error()
			}

			// compare error messages
			if wantErr != gotErr {
				t.Errorf("Arguments.parseP() error = %#v, wantErr %#v", err, tt.wantErr)
			}

			if tt.wantErr != nil { // ignore checks when an error is returned; we don't care
				return
			}

			if !reflect.DeepEqual(args, tt.wantParsed) {
				t.Errorf("Arguments.parseP() args = %#v, wantArgs %#v", args, &tt.wantParsed)
			}
		})
	}
}
