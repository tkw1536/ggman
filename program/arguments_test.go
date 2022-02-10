package program_test

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

type Arguments = program.Arguments // FIXME

var errParseArgsNeedOneArgument = program.ErrParseArgsNeedOneArgument
var errParseMissingForArgument = &flags.Error{Type: flags.ErrExpectedArgument, Message: "expected argument for flag `-f, --for'"}

func TestArguments_Parse(t *testing.T) {
	type args struct {
		argv []string
	}
	tests := []struct {
		name       string
		args       args
		wantParsed program.Arguments
		wantErr    error
	}{
		{"no arguments", args{[]string{}}, Arguments{}, errParseArgsNeedOneArgument},
		{"command without arguments", args{[]string{"cmd"}}, Arguments{Command: "cmd", Args: []string{}}, nil},

		{"help with command (2)", args{[]string{"--help", "cmd"}}, Arguments{Universals: program.Universals{Help: true}, Args: []string{"cmd"}}, nil},
		{"help with command (3)", args{[]string{"-h", "cmd"}}, Arguments{Universals: program.Universals{Help: true}, Args: []string{"cmd"}}, nil},

		{"help without command (2)", args{[]string{"--help"}}, Arguments{Universals: program.Universals{Help: true}, Args: []string{}}, nil},
		{"help without command (3)", args{[]string{"-h"}}, Arguments{Universals: program.Universals{Help: true}, Args: []string{}}, nil},

		{"version with command (2)", args{[]string{"--version", "cmd"}}, Arguments{Universals: program.Universals{Version: true}, Args: []string{"cmd"}}, nil},
		{"version with command (3)", args{[]string{"-v", "cmd"}}, Arguments{Universals: program.Universals{Version: true}, Args: []string{"cmd"}}, nil},

		{"version without command (2)", args{[]string{"--version"}}, Arguments{Universals: program.Universals{Version: true}, Args: []string{}}, nil},
		{"version without command (3)", args{[]string{"-v"}}, Arguments{Universals: program.Universals{Version: true}, Args: []string{}}, nil},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, Arguments{Command: "cmd", Args: []string{"a1", "a2"}}, nil},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, Arguments{Command: "cmd", Args: []string{"help", "a1"}}, nil},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, Arguments{Command: "cmd", Args: []string{"--help", "a1"}}, nil},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, Arguments{Command: "cmd", Args: []string{"-h", "a1"}}, nil},

		{"command with version (1)", args{[]string{"cmd", "version", "a1"}}, Arguments{Command: "cmd", Args: []string{"version", "a1"}}, nil},
		{"command with version (2)", args{[]string{"cmd", "--version", "a1"}}, Arguments{Command: "cmd", Args: []string{"--version", "a1"}}, nil},
		{"command with version (3)", args{[]string{"cmd", "-v", "a1"}}, Arguments{Command: "cmd", Args: []string{"-v", "a1"}}, nil},

		{"only a for (2)", args{[]string{"--for"}}, Arguments{}, errParseMissingForArgument},
		{"only a for (3)", args{[]string{"-f"}}, Arguments{}, errParseMissingForArgument},

		{"only a here (1)", args{[]string{"--here"}}, Arguments{}, errParseArgsNeedOneArgument},
		{"only a here (2)", args{[]string{"-H"}}, Arguments{}, errParseArgsNeedOneArgument},

		{"only a path (1)", args{[]string{"--path", "p"}}, Arguments{}, errParseArgsNeedOneArgument},
		{"only a path (2)", args{[]string{"-P", "p"}}, Arguments{}, errParseArgsNeedOneArgument},

		{"for without command (2)", args{[]string{"--for", "match"}}, Arguments{}, errParseArgsNeedOneArgument},
		{"for without command (3)", args{[]string{"-f", "match"}}, Arguments{}, errParseArgsNeedOneArgument},

		{"for with command (2)", args{[]string{"--for", "match", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Filters: []string{"match"}}, Args: []string{}}, nil},
		{"for with command (3)", args{[]string{"-f", "match", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Filters: []string{"match"}}, Args: []string{}}, nil},

		{"here with command (1)", args{[]string{"--here", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Here: true}, Args: []string{}}, nil},
		{"here with command (2)", args{[]string{"-H", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Here: true}, Args: []string{}}, nil},

		{"path with command (1)", args{[]string{"--path", "P", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P"}}, Args: []string{}}, nil},
		{"path with command (2)", args{[]string{"-P", "P", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P"}}, Args: []string{}}, nil},

		{"multiple paths with command (1)", args{[]string{"--path", "P1", "--path", "P2", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P1", "P2"}}, Args: []string{}}, nil},
		{"multiple paths with command (2)", args{[]string{"-P", "P1", "--path", "P2", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P1", "P2"}}, Args: []string{}}, nil},

		{"path + here with command (1)", args{[]string{"--path", "P", "--here", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P"}, Here: true}, Args: []string{}}, nil},
		{"path + here with command (2)", args{[]string{"--path", "P", "-H", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P"}, Here: true}, Args: []string{}}, nil},
		{"path + here with command (3)", args{[]string{"-P", "P", "--here", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P"}, Here: true}, Args: []string{}}, nil},
		{"path + here with command (4)", args{[]string{"-P", "P", "-H", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P"}, Here: true}, Args: []string{}}, nil},

		{"dirty with command (1)", args{[]string{"--dirty", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Dirty: true}, Args: []string{}}, nil},
		{"dirty with command (2)", args{[]string{"-d", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Dirty: true}, Args: []string{}}, nil},
		{"clean with command (1)", args{[]string{"--clean", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Clean: true}, Args: []string{}}, nil},
		{"clean with command (2)", args{[]string{"-c", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Clean: true}, Args: []string{}}, nil},

		{"synced with command (1)", args{[]string{"--synced", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Synced: true}, Args: []string{}}, nil},
		{"synced with command (2)", args{[]string{"-s", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Synced: true}, Args: []string{}}, nil},
		{"unsynced with command (1)", args{[]string{"--unsynced", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{UnSynced: true}, Args: []string{}}, nil},
		{"unsynced with command (2)", args{[]string{"-u", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{UnSynced: true}, Args: []string{}}, nil},

		{"pristine with command (1)", args{[]string{"--pristine", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Pristine: true}, Args: []string{}}, nil},
		{"pristine with command (2)", args{[]string{"-p", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Pristine: true}, Args: []string{}}, nil},
		{"pristine with command (1)", args{[]string{"--tarnished", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Tarnished: true}, Args: []string{}}, nil},
		{"pristine with command (2)", args{[]string{"-t", "cmd"}}, Arguments{Command: "cmd", Flags: program.Flags{Tarnished: true}, Args: []string{}}, nil},

		/*{"for with command and arguments (1)", args{[]string{"for", "match", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Flags: program.Flags{Filters: []string{"match"}}, Args: []string{"a1", "a2"}}, nil},*/
		{"for with command and arguments (2)", args{[]string{"--for", "match", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Flags: program.Flags{Filters: []string{"match"}}, Args: []string{"a1", "a2"}}, nil},
		{"for with command and arguments (3)", args{[]string{"-f", "match", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Flags: program.Flags{Filters: []string{"match"}}, Args: []string{"a1", "a2"}}, nil},

		{"here with command and arguments (1)", args{[]string{"--here", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Flags: program.Flags{Here: true}, Args: []string{"a1", "a2"}}, nil},
		{"here with command and arguments (2)", args{[]string{"-H", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Flags: program.Flags{Here: true}, Args: []string{"a1", "a2"}}, nil},

		{"path with command and arguments (1)", args{[]string{"--path", "P", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P"}}, Args: []string{"a1", "a2"}}, nil},
		{"path with command and arguments (2)", args{[]string{"-P", "P", "cmd", "a1", "a2"}}, Arguments{Command: "cmd", Flags: program.Flags{Path: []string{"P"}}, Args: []string{"a1", "a2"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &Arguments{}
			// TODO: Fix linter warnings
			if err := args.Parse(tt.args.argv); !reflect.DeepEqual(err, tt.wantErr) {
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

func TestCommandArguments_checkPositionalCount(t *testing.T) {
	tests := []struct {
		name string

		options   tDescription
		arguments Arguments

		wantErr string
	}{
		// taking 0 args
		{
			"no arguments",
			tDescription{PosArgsMin: 0, PosArgsMax: 0},
			Arguments{Command: "example", Args: []string{}},
			"",
		},

		// taking 1 arg
		{
			"one argument, too few",
			tDescription{PosArgsMin: 1, PosArgsMax: 1},
			Arguments{Command: "example", Args: []string{}},
			"Wrong number of arguments: 'example' takes exactly 1 argument(s). ",
		},
		{
			"one argument, exactly enough",
			tDescription{PosArgsMin: 1, PosArgsMax: 1},
			Arguments{Command: "example", Args: []string{"world"}},
			"",
		},
		{
			"one argument, too many",
			tDescription{PosArgsMin: 1, PosArgsMax: 1},
			Arguments{Command: "example", Args: []string{"hello", "world"}},
			"Wrong number of arguments: 'example' takes exactly 1 argument(s). ",
		},

		// taking 1 or 2 args
		{
			"1-2 arguments, too few",
			tDescription{PosArgsMin: 1, PosArgsMax: 2},
			Arguments{Command: "example", Args: []string{}},
			"Wrong number of arguments: 'example' takes between 1 and 2 arguments. ",
		},
		{
			"1-2 arguments, enough",
			tDescription{PosArgsMin: 1, PosArgsMax: 2},
			Arguments{Command: "example", Args: []string{"world"}},
			"",
		},
		{
			"1-2 arguments, enough (2)",
			tDescription{PosArgsMin: 1, PosArgsMax: 2},
			Arguments{Command: "example", Args: []string{"hello", "world"}},
			"",
		},
		{
			"1-2 arguments, too many",
			tDescription{PosArgsMin: 1, PosArgsMax: 2},
			Arguments{Command: "example", Args: []string{"hello", "world", "you"}},
			"Wrong number of arguments: 'example' takes between 1 and 2 arguments. ",
		},

		// taking 2 args
		{
			"two arguments, too few",
			tDescription{PosArgsMin: 2, PosArgsMax: 2},
			Arguments{Command: "example", Args: []string{}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},
		{
			"two arguments, too few (2)",
			tDescription{PosArgsMin: 2, PosArgsMax: 2},
			Arguments{Command: "example", Args: []string{"world"}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},
		{
			"two arguments, enough",
			tDescription{PosArgsMin: 2, PosArgsMax: 2},
			Arguments{Command: "example", Args: []string{"hello", "world"}},
			"",
		},
		{
			"two arguments, too many",
			tDescription{PosArgsMin: 2, PosArgsMax: 2},
			Arguments{Command: "example", Args: []string{"hello", "world", "you"}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},

		// at least one argument
		{
			"at least 1 arguments, not enough",
			tDescription{PosArgsMin: 1, PosArgsMax: -1},
			Arguments{Command: "example", Args: []string{}},
			"Wrong number of arguments: 'example' takes at least 1 argument(s). ",
		},
		{
			"at least 1 arguments, enough",
			tDescription{PosArgsMin: 1, PosArgsMax: -1},
			Arguments{Command: "example", Args: []string{"hello"}},
			"",
		},
		{
			"at least 1 arguments, enough (2)",
			tDescription{PosArgsMin: 1, PosArgsMax: -1},
			Arguments{Command: "example", Args: []string{"hello", "cruel", "world"}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &tCommandArguments{
				Description: tt.options,
				Arguments:   tt.arguments,
			}
			err := args.CheckPositionalCount()
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("CommandArguments.checkPositionalCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandArguments_checkFilterArgument(t *testing.T) {
	tests := []struct {
		name      string
		options   tDescription
		arguments Arguments

		wantErr string
	}{
		{
			"for not allowed, for not given",
			tDescription{Requirements: env.Requirement{AllowsFilter: false}},
			Arguments{Command: "example"},
			"",
		},
		{
			"for not allowed, for given",
			tDescription{Requirements: env.Requirement{AllowsFilter: false}},
			Arguments{Command: "example", Flags: program.Flags{Filters: []string{"pattern"}}},
			"Wrong number of arguments: 'example' takes no '--for' argument. ",
		},

		{
			"fuzzy not allowed, fuzzy given",
			tDescription{Requirements: env.Requirement{AllowsFilter: false}},
			Arguments{Command: "example", Flags: program.Flags{NoFuzzyFilter: true, Filters: nil}},
			"Wrong number of arguments: 'example' takes no '--no-fuzzy-filter' argument. ",
		},

		{
			"for allowed, for not given",
			tDescription{Requirements: env.Requirement{AllowsFilter: true}},
			Arguments{Command: "example", Flags: program.Flags{Filters: nil}},
			"",
		},
		{
			"for allowed, for given",
			tDescription{Requirements: env.Requirement{AllowsFilter: true}},
			Arguments{Command: "example", Flags: program.Flags{Filters: []string{"pattern"}}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tCommandArguments{
				Description: tt.options,
				Arguments:   tt.arguments,
			}
			err := args.CheckFilterArgument()
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("CommandArguments.checkFilterArgument() error = %q, wantErr %q", err, tt.wantErr)
			}
		})
	}
}
