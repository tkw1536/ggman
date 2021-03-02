package program

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/env"
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
		{"command without arguments", args{[]string{"cmd"}}, Arguments{Command: "cmd", Args: []string{}}, nil},

		{"help with command (1)", args{[]string{"help", "cmd"}}, Arguments{Help: true, Args: []string{"cmd"}}, nil},
		{"help with command (2)", args{[]string{"--help", "cmd"}}, Arguments{Help: true, Args: []string{"cmd"}}, nil},
		{"help with command (3)", args{[]string{"-h", "cmd"}}, Arguments{Help: true, Args: []string{"cmd"}}, nil},

		{"help without command (1)", args{[]string{"help"}}, Arguments{Help: true, Args: []string{}}, nil},
		{"help without command (2)", args{[]string{"--help"}}, Arguments{Help: true, Args: []string{}}, nil},
		{"help without command (3)", args{[]string{"-h"}}, Arguments{Help: true, Args: []string{}}, nil},

		{"version with command (1)", args{[]string{"version", "cmd"}}, Arguments{Version: true, Args: []string{"cmd"}}, nil},
		{"version with command (2)", args{[]string{"--version", "cmd"}}, Arguments{Version: true, Args: []string{"cmd"}}, nil},
		{"version with command (3)", args{[]string{"-v", "cmd"}}, Arguments{Version: true, Args: []string{"cmd"}}, nil},

		{"version without command (1)", args{[]string{"version"}}, Arguments{Version: true, Args: []string{}}, nil},
		{"version without command (2)", args{[]string{"--version"}}, Arguments{Version: true, Args: []string{}}, nil},
		{"version without command (3)", args{[]string{"-v"}}, Arguments{Version: true, Args: []string{}}, nil},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, Arguments{Command: "cmd", Args: []string{"a1", "a2"}}, nil},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, Arguments{Command: "cmd", Args: []string{"help", "a1"}}, nil},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, Arguments{Command: "cmd", Args: []string{"--help", "a1"}}, nil},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, Arguments{Command: "cmd", Args: []string{"-h", "a1"}}, nil},

		{"command with version (1)", args{[]string{"cmd", "version", "a1"}}, Arguments{Command: "cmd", Args: []string{"version", "a1"}}, nil},
		{"command with version (2)", args{[]string{"cmd", "--version", "a1"}}, Arguments{Command: "cmd", Args: []string{"--version", "a1"}}, nil},
		{"command with version (3)", args{[]string{"cmd", "-v", "a1"}}, Arguments{Command: "cmd", Args: []string{"-v", "a1"}}, nil},
		{"only a for (1)", args{[]string{"for"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"only a for (2)", args{[]string{"--for"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"only a for (3)", args{[]string{"-f"}}, Arguments{}, errParseArgsNeedTwoAfterFor},

		{"only a here (1)", args{[]string{"--here"}}, Arguments{}, errParseArgsNeedOneArgument},
		{"only a here (2)", args{[]string{"-H"}}, Arguments{}, errParseArgsNeedOneArgument},

		{"for without command (1)", args{[]string{"for", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"for without command (2)", args{[]string{"--for", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"for without command (3)", args{[]string{"-f", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},

		{"for with command (1)", args{[]string{"for", "match", "cmd"}}, Arguments{Command: "cmd", Filters: []string{"match"}, Args: []string{}}, nil},
		{"for with command (2)", args{[]string{"--for", "match", "cmd"}}, Arguments{Command: "cmd", Filters: []string{"match"}, Args: []string{}}, nil},
		{"for with command (3)", args{[]string{"-f", "match", "cmd"}}, Arguments{Command: "cmd", Filters: []string{"match"}, Args: []string{}}, nil},

		{"here with command (1)", args{[]string{"--here", "cmd"}}, Arguments{Command: "cmd", Here: true, Args: []string{}}, nil},
		{"here with command (2)", args{[]string{"-H", "cmd"}}, Arguments{Command: "cmd", Here: true, Args: []string{}}, nil},

		{"for with command and arguments (1)", args{[]string{"for", "match", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Filters: []string{"match"}, Args: []string{"a1", "a2"}}, nil},
		{"for with command and arguments (2)", args{[]string{"--for", "match", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Filters: []string{"match"}, Args: []string{"a1", "a2"}}, nil},
		{"for with command and arguments (3)", args{[]string{"-f", "match", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Filters: []string{"match"}, Args: []string{"a1", "a2"}}, nil},

		{"here with command and arguments (1)", args{[]string{"--here", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Here: true, Args: []string{"a1", "a2"}}, nil},
		{"here with command and arguments (2)", args{[]string{"-H", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Here: true, Args: []string{"a1", "a2"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &Arguments{}
			if err := args.Parse(tt.args.argv); err != tt.wantErr {
				t.Errorf("Arguments.Parse() error = %#v, wantErr %#v", err, tt.wantErr)
			}

			if tt.wantErr != nil { // ignore checks when an error is returned; we don't care
				return
			}

			if !reflect.DeepEqual(*args, tt.wantParsed) {
				t.Errorf("Arguments.Parse() args = %#v, wantArgs %#v", args, &tt.wantParsed)
			}
		})
	}
}

func TestCommandArguments_checkArgumentCount(t *testing.T) {
	tests := []struct {
		name string

		options   Options
		arguments Arguments

		wantErr string
	}{
		// taking 0 args
		{
			"no arguments",
			Options{MinArgs: 0, MaxArgs: 0},
			Arguments{Command: "example", Args: []string{}},
			"",
		},

		// taking 1 arg
		{
			"one argument, too few",
			Options{MinArgs: 1, MaxArgs: 1},
			Arguments{Command: "example", Args: []string{}},
			"Wrong number of arguments: 'example' takes exactly 1 argument(s). ",
		},
		{
			"one argument, exactly enough",
			Options{MinArgs: 1, MaxArgs: 1},
			Arguments{Command: "example", Args: []string{"world"}},
			"",
		},
		{
			"one argument, too many",
			Options{MinArgs: 1, MaxArgs: 1},
			Arguments{Command: "example", Args: []string{"hello", "world"}},
			"Wrong number of arguments: 'example' takes exactly 1 argument(s). ",
		},

		// taking 1 or 2 args
		{
			"1-2 arguments, too few",
			Options{MinArgs: 1, MaxArgs: 2},
			Arguments{Command: "example", Args: []string{}},
			"Wrong number of arguments: 'example' takes between 1 and 2 arguments. ",
		},
		{
			"1-2 arguments, enough",
			Options{MinArgs: 1, MaxArgs: 2},
			Arguments{Command: "example", Args: []string{"world"}},
			"",
		},
		{
			"1-2 arguments, enough (2)",
			Options{MinArgs: 1, MaxArgs: 2},
			Arguments{Command: "example", Args: []string{"hello", "world"}},
			"",
		},
		{
			"1-2 arguments, too many",
			Options{MinArgs: 1, MaxArgs: 2},
			Arguments{Command: "example", Args: []string{"hello", "world", "you"}},
			"Wrong number of arguments: 'example' takes between 1 and 2 arguments. ",
		},

		// taking 2 args
		{
			"two arguments, too few",
			Options{MinArgs: 2, MaxArgs: 2},
			Arguments{Command: "example", Args: []string{}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},
		{
			"two arguments, too few (2)",
			Options{MinArgs: 2, MaxArgs: 2},
			Arguments{Command: "example", Args: []string{"world"}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},
		{
			"two arguments, enough",
			Options{MinArgs: 2, MaxArgs: 2},
			Arguments{Command: "example", Args: []string{"hello", "world"}},
			"",
		},
		{
			"two arguments, too many",
			Options{MinArgs: 2, MaxArgs: 2},
			Arguments{Command: "example", Args: []string{"hello", "world", "you"}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},

		// at least one argument
		{
			"at least 1 arguments, not enough",
			Options{MinArgs: 1, MaxArgs: -1},
			Arguments{Command: "example", Args: []string{}},
			"Wrong number of arguments: 'example' takes at least 1 argument(s). ",
		},
		{
			"at least 1 arguments, enough",
			Options{MinArgs: 1, MaxArgs: -1},
			Arguments{Command: "example", Args: []string{"hello"}},
			"",
		},
		{
			"at least 1 arguments, enough (2)",
			Options{MinArgs: 1, MaxArgs: -1},
			Arguments{Command: "example", Args: []string{"hello", "cruel", "world"}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &CommandArguments{
				options:   tt.options,
				Arguments: tt.arguments,
			}
			err := args.checkArgumentCount()
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("CommandArguments.checkArgumentCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandArguments_checkFilterArgument(t *testing.T) {
	tests := []struct {
		name      string
		options   Options
		arguments Arguments

		wantErr string
	}{
		{
			"for not allowed, for not given",
			Options{Environment: env.Requirement{AllowsFilter: false}},
			Arguments{Command: "example"},
			"",
		},
		{
			"for not allowed, for given",
			Options{Environment: env.Requirement{AllowsFilter: false}},
			Arguments{Command: "example", Filters: []string{"pattern"}},
			"Wrong number of arguments: 'example' takes no 'for' argument. ",
		},

		{
			"for allowed, for not given",
			Options{Environment: env.Requirement{AllowsFilter: true}},
			Arguments{Command: "example", Filters: nil},
			"",
		},
		{
			"for allowed, for given",
			Options{Environment: env.Requirement{AllowsFilter: true}},
			Arguments{Command: "example", Filters: []string{"pattern"}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := CommandArguments{
				options:   tt.options,
				Arguments: tt.arguments,
			}
			err := args.checkFilterArgument()
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("CommandArguments.checkFilterArgument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

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
