package program

import (
	"reflect"
	"testing"

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

		{"command without arguments", args{[]string{"cmd"}}, Arguments{"cmd", env.NoFilter, false, false, []string{}, nil}, nil},

		{"help with command (1)", args{[]string{"help", "cmd"}}, Arguments{"", env.NoFilter, true, false, []string{"cmd"}, nil}, nil},
		{"help with command (2)", args{[]string{"--help", "cmd"}}, Arguments{"", env.NoFilter, true, false, []string{"cmd"}, nil}, nil},
		{"help with command (3)", args{[]string{"-h", "cmd"}}, Arguments{"", env.NoFilter, true, false, []string{"cmd"}, nil}, nil},

		{"help without command (1)", args{[]string{"help"}}, Arguments{"", env.NoFilter, true, false, []string{}, nil}, nil},
		{"help without command (2)", args{[]string{"--help"}}, Arguments{"", env.NoFilter, true, false, []string{}, nil}, nil},
		{"help without command (3)", args{[]string{"-h"}}, Arguments{"", env.NoFilter, true, false, []string{}, nil}, nil},

		{"version with command (1)", args{[]string{"version", "cmd"}}, Arguments{"", env.NoFilter, false, true, []string{"cmd"}, nil}, nil},
		{"version with command (2)", args{[]string{"--version", "cmd"}}, Arguments{"", env.NoFilter, false, true, []string{"cmd"}, nil}, nil},
		{"version with command (3)", args{[]string{"-v", "cmd"}}, Arguments{"", env.NoFilter, false, true, []string{"cmd"}, nil}, nil},

		{"version without command (1)", args{[]string{"version"}}, Arguments{"", env.NoFilter, false, true, []string{}, nil}, nil},
		{"version without command (2)", args{[]string{"--version"}}, Arguments{"", env.NoFilter, false, true, []string{}, nil}, nil},
		{"version without command (3)", args{[]string{"-v"}}, Arguments{"", env.NoFilter, false, true, []string{}, nil}, nil},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, Arguments{"cmd", env.NoFilter, false, false, []string{"a1", "a2"}, nil}, nil},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, Arguments{"cmd", env.NoFilter, false, false, []string{"help", "a1"}, nil}, nil},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, Arguments{"cmd", env.NoFilter, false, false, []string{"--help", "a1"}, nil}, nil},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, Arguments{"cmd", env.NoFilter, false, false, []string{"-h", "a1"}, nil}, nil},

		{"command with version (1)", args{[]string{"cmd", "version", "a1"}}, Arguments{"cmd", env.NoFilter, false, false, []string{"version", "a1"}, nil}, nil},
		{"command with version (2)", args{[]string{"cmd", "--version", "a1"}}, Arguments{"cmd", env.NoFilter, false, false, []string{"--version", "a1"}, nil}, nil},
		{"command with version (3)", args{[]string{"cmd", "-v", "a1"}}, Arguments{"cmd", env.NoFilter, false, false, []string{"-v", "a1"}, nil}, nil},

		{"only a for (1)", args{[]string{"for"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"only a for (2)", args{[]string{"--for"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"only a for (3)", args{[]string{"-f"}}, Arguments{}, errParseArgsNeedTwoAfterFor},

		{"for without command (1)", args{[]string{"for", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"for without command (2)", args{[]string{"--for", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},
		{"for without command (3)", args{[]string{"-f", "match"}}, Arguments{}, errParseArgsNeedTwoAfterFor},

		{"for with command (1)", args{[]string{"for", "match", "cmd"}}, Arguments{"cmd", env.NewFilter("match"), false, false, []string{}, nil}, nil},
		{"for with command (2)", args{[]string{"--for", "match", "cmd"}}, Arguments{"cmd", env.NewFilter("match"), false, false, []string{}, nil}, nil},
		{"for with command (3)", args{[]string{"-f", "match", "cmd"}}, Arguments{"cmd", env.NewFilter("match"), false, false, []string{}, nil}, nil},

		{"for with command and arguments (1)", args{[]string{"for", "match", "cmd", "a1", "a2"}}, Arguments{"cmd", env.NewFilter("match"), false, false, []string{"a1", "a2"}, nil}, nil},
		{"for with command and arguments (2)", args{[]string{"--for", "match", "cmd", "a1", "a2"}}, Arguments{"cmd", env.NewFilter("match"), false, false, []string{"a1", "a2"}, nil}, nil},
		{"for with command and arguments (3)", args{[]string{"-f", "match", "cmd", "a1", "a2"}}, Arguments{"cmd", env.NewFilter("match"), false, false, []string{"a1", "a2"}, nil}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &Arguments{}
			if err := args.Parse(tt.args.argv); err != tt.wantErr {
				t.Errorf("Arguments.Parse() error = %#v, wantErr %#v", err, tt.wantErr)
			}

			// when an error occured, we don't care about the returned value
			// and the behaviour is unspecified
			if tt.wantErr != nil {
				return
			}
			// ignore flagset during comparison
			args.flagsetGlobal = nil
			tt.wantParsed.flagsetGlobal = nil

			if !reflect.DeepEqual(args, &tt.wantParsed) {
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
			Arguments{Command: "example", For: env.NoFilter},
			"",
		},
		{
			"for not allowed, for given",
			Options{Environment: env.Requirement{AllowsFilter: false}},
			Arguments{Command: "example", For: env.NewFilter("pattern")},
			"Wrong number of arguments: 'example' takes no 'for' argument. ",
		},

		{
			"for allowed, for not given",
			Options{Environment: env.Requirement{AllowsFilter: true}},
			Arguments{Command: "example", For: env.NoFilter},
			"",
		},
		{
			"for allowed, for given",
			Options{Environment: env.Requirement{AllowsFilter: true}},
			Arguments{Command: "example", For: env.NewFilter("pattern")},
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
