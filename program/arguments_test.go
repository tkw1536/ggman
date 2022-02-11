package program_test

import (
	"reflect"
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

var errParseArgsNeedOneArgument = program.ErrParseArgsNeedOneArgument
var errParseMissingForArgument = &flags.Error{Type: flags.ErrExpectedArgument, Message: "expected argument for flag `-f, --for'"}

func TestArguments_Parse(t *testing.T) {
	type args struct {
		argv []string
	}
	tests := []struct {
		name       string
		args       args
		wantParsed tArguments
		wantErr    error
	}{
		{"no arguments", args{[]string{}}, tArguments{}, errParseArgsNeedOneArgument},
		{"command without arguments", args{[]string{"cmd"}}, tArguments{Command: "cmd", Pos: []string{}}, nil},

		{"help with command (2)", args{[]string{"--help", "cmd"}}, tArguments{Universals: program.Universals{Help: true}, Pos: []string{"cmd"}}, nil},
		{"help with command (3)", args{[]string{"-h", "cmd"}}, tArguments{Universals: program.Universals{Help: true}, Pos: []string{"cmd"}}, nil},

		{"help without command (2)", args{[]string{"--help"}}, tArguments{Universals: program.Universals{Help: true}, Pos: []string{}}, nil},
		{"help without command (3)", args{[]string{"-h"}}, tArguments{Universals: program.Universals{Help: true}, Pos: []string{}}, nil},

		{"version with command (2)", args{[]string{"--version", "cmd"}}, tArguments{Universals: program.Universals{Version: true}, Pos: []string{"cmd"}}, nil},
		{"version with command (3)", args{[]string{"-v", "cmd"}}, tArguments{Universals: program.Universals{Version: true}, Pos: []string{"cmd"}}, nil},

		{"version without command (2)", args{[]string{"--version"}}, tArguments{Universals: program.Universals{Version: true}, Pos: []string{}}, nil},
		{"version without command (3)", args{[]string{"-v"}}, tArguments{Universals: program.Universals{Version: true}, Pos: []string{}}, nil},

		{"command with arguments", args{[]string{"cmd", "a1", "a2"}}, tArguments{Command: "cmd", Pos: []string{"a1", "a2"}}, nil},

		{"command with help (1)", args{[]string{"cmd", "help", "a1"}}, tArguments{Command: "cmd", Pos: []string{"help", "a1"}}, nil},
		{"command with help (2)", args{[]string{"cmd", "--help", "a1"}}, tArguments{Command: "cmd", Pos: []string{"--help", "a1"}}, nil},
		{"command with help (3)", args{[]string{"cmd", "-h", "a1"}}, tArguments{Command: "cmd", Pos: []string{"-h", "a1"}}, nil},

		{"command with version (1)", args{[]string{"cmd", "version", "a1"}}, tArguments{Command: "cmd", Pos: []string{"version", "a1"}}, nil},
		{"command with version (2)", args{[]string{"cmd", "--version", "a1"}}, tArguments{Command: "cmd", Pos: []string{"--version", "a1"}}, nil},
		{"command with version (3)", args{[]string{"cmd", "-v", "a1"}}, tArguments{Command: "cmd", Pos: []string{"-v", "a1"}}, nil},

		{"only a for (2)", args{[]string{"--for"}}, tArguments{}, errParseMissingForArgument},
		{"only a for (3)", args{[]string{"-f"}}, tArguments{}, errParseMissingForArgument},

		{"only a here (1)", args{[]string{"--here"}}, tArguments{}, errParseArgsNeedOneArgument},
		{"only a here (2)", args{[]string{"-H"}}, tArguments{}, errParseArgsNeedOneArgument},

		{"only a path (1)", args{[]string{"--path", "p"}}, tArguments{}, errParseArgsNeedOneArgument},
		{"only a path (2)", args{[]string{"-P", "p"}}, tArguments{}, errParseArgsNeedOneArgument},

		{"for without command (2)", args{[]string{"--for", "match"}}, tArguments{}, errParseArgsNeedOneArgument},
		{"for without command (3)", args{[]string{"-f", "match"}}, tArguments{}, errParseArgsNeedOneArgument},

		{"for with command (2)", args{[]string{"--for", "match", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Filters: []string{"match"}}, Pos: []string{}}, nil},
		{"for with command (3)", args{[]string{"-f", "match", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Filters: []string{"match"}}, Pos: []string{}}, nil},

		{"here with command (1)", args{[]string{"--here", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Here: true}, Pos: []string{}}, nil},
		{"here with command (2)", args{[]string{"-H", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Here: true}, Pos: []string{}}, nil},

		{"path with command (1)", args{[]string{"--path", "P", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P"}}, Pos: []string{}}, nil},
		{"path with command (2)", args{[]string{"-P", "P", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P"}}, Pos: []string{}}, nil},

		{"multiple paths with command (1)", args{[]string{"--path", "P1", "--path", "P2", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P1", "P2"}}, Pos: []string{}}, nil},
		{"multiple paths with command (2)", args{[]string{"-P", "P1", "--path", "P2", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P1", "P2"}}, Pos: []string{}}, nil},

		{"path + here with command (1)", args{[]string{"--path", "P", "--here", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P"}, Here: true}, Pos: []string{}}, nil},
		{"path + here with command (2)", args{[]string{"--path", "P", "-H", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P"}, Here: true}, Pos: []string{}}, nil},
		{"path + here with command (3)", args{[]string{"-P", "P", "--here", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P"}, Here: true}, Pos: []string{}}, nil},
		{"path + here with command (4)", args{[]string{"-P", "P", "-H", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P"}, Here: true}, Pos: []string{}}, nil},

		{"dirty with command (1)", args{[]string{"--dirty", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Dirty: true}, Pos: []string{}}, nil},
		{"dirty with command (2)", args{[]string{"-d", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Dirty: true}, Pos: []string{}}, nil},
		{"clean with command (1)", args{[]string{"--clean", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Clean: true}, Pos: []string{}}, nil},
		{"clean with command (2)", args{[]string{"-c", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Clean: true}, Pos: []string{}}, nil},

		{"synced with command (1)", args{[]string{"--synced", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Synced: true}, Pos: []string{}}, nil},
		{"synced with command (2)", args{[]string{"-s", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Synced: true}, Pos: []string{}}, nil},
		{"unsynced with command (1)", args{[]string{"--unsynced", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{UnSynced: true}, Pos: []string{}}, nil},
		{"unsynced with command (2)", args{[]string{"-u", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{UnSynced: true}, Pos: []string{}}, nil},

		{"pristine with command (1)", args{[]string{"--pristine", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Pristine: true}, Pos: []string{}}, nil},
		{"pristine with command (2)", args{[]string{"-p", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Pristine: true}, Pos: []string{}}, nil},
		{"pristine with command (1)", args{[]string{"--tarnished", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Tarnished: true}, Pos: []string{}}, nil},
		{"pristine with command (2)", args{[]string{"-t", "cmd"}}, tArguments{Command: "cmd", Flags: tFlags{Tarnished: true}, Pos: []string{}}, nil},

		{"for with command and arguments (2)", args{[]string{"--for", "match", "cmd", "a1", "a2"}}, tArguments{Command: "cmd", Flags: tFlags{Filters: []string{"match"}}, Pos: []string{"a1", "a2"}}, nil},
		{"for with command and arguments (3)", args{[]string{"-f", "match", "cmd", "a1", "a2"}}, tArguments{Command: "cmd", Flags: tFlags{Filters: []string{"match"}}, Pos: []string{"a1", "a2"}}, nil},

		{"here with command and arguments (1)", args{[]string{"--here", "cmd", "a1", "a2"}}, tArguments{Command: "cmd", Flags: tFlags{Here: true}, Pos: []string{"a1", "a2"}}, nil},
		{"here with command and arguments (2)", args{[]string{"-H", "cmd", "a1", "a2"}}, tArguments{Command: "cmd", Flags: tFlags{Here: true}, Pos: []string{"a1", "a2"}}, nil},

		{"path with command and arguments (1)", args{[]string{"--path", "P", "cmd", "a1", "a2"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P"}}, Pos: []string{"a1", "a2"}}, nil},
		{"path with command and arguments (2)", args{[]string{"-P", "P", "cmd", "a1", "a2"}}, tArguments{Command: "cmd", Flags: tFlags{Path: []string{"P"}}, Pos: []string{"a1", "a2"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &tArguments{}
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

func TestContext_checkPositionalCount(t *testing.T) {
	tests := []struct {
		name string

		options   tDescription
		arguments tArguments

		wantErr string
	}{
		// taking 0 args
		{
			"no arguments",
			tDescription{PosArgsMin: 0, PosArgsMax: 0},
			tArguments{Command: "example", Pos: []string{}},
			"",
		},

		// taking 1 arg
		{
			"one argument, too few",
			tDescription{PosArgsMin: 1, PosArgsMax: 1},
			tArguments{Command: "example", Pos: []string{}},
			"Wrong number of arguments: 'example' takes exactly 1 argument(s). ",
		},
		{
			"one argument, exactly enough",
			tDescription{PosArgsMin: 1, PosArgsMax: 1},
			tArguments{Command: "example", Pos: []string{"world"}},
			"",
		},
		{
			"one argument, too many",
			tDescription{PosArgsMin: 1, PosArgsMax: 1},
			tArguments{Command: "example", Pos: []string{"hello", "world"}},
			"Wrong number of arguments: 'example' takes exactly 1 argument(s). ",
		},

		// taking 1 or 2 args
		{
			"1-2 arguments, too few",
			tDescription{PosArgsMin: 1, PosArgsMax: 2},
			tArguments{Command: "example", Pos: []string{}},
			"Wrong number of arguments: 'example' takes between 1 and 2 arguments. ",
		},
		{
			"1-2 arguments, enough",
			tDescription{PosArgsMin: 1, PosArgsMax: 2},
			tArguments{Command: "example", Pos: []string{"world"}},
			"",
		},
		{
			"1-2 arguments, enough (2)",
			tDescription{PosArgsMin: 1, PosArgsMax: 2},
			tArguments{Command: "example", Pos: []string{"hello", "world"}},
			"",
		},
		{
			"1-2 arguments, too many",
			tDescription{PosArgsMin: 1, PosArgsMax: 2},
			tArguments{Command: "example", Pos: []string{"hello", "world", "you"}},
			"Wrong number of arguments: 'example' takes between 1 and 2 arguments. ",
		},

		// taking 2 args
		{
			"two arguments, too few",
			tDescription{PosArgsMin: 2, PosArgsMax: 2},
			tArguments{Command: "example", Pos: []string{}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},
		{
			"two arguments, too few (2)",
			tDescription{PosArgsMin: 2, PosArgsMax: 2},
			tArguments{Command: "example", Pos: []string{"world"}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},
		{
			"two arguments, enough",
			tDescription{PosArgsMin: 2, PosArgsMax: 2},
			tArguments{Command: "example", Pos: []string{"hello", "world"}},
			"",
		},
		{
			"two arguments, too many",
			tDescription{PosArgsMin: 2, PosArgsMax: 2},
			tArguments{Command: "example", Pos: []string{"hello", "world", "you"}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},

		// at least one argument
		{
			"at least 1 arguments, not enough",
			tDescription{PosArgsMin: 1, PosArgsMax: -1},
			tArguments{Command: "example", Pos: []string{}},
			"Wrong number of arguments: 'example' takes at least 1 argument(s). ",
		},
		{
			"at least 1 arguments, enough",
			tDescription{PosArgsMin: 1, PosArgsMax: -1},
			tArguments{Command: "example", Pos: []string{"hello"}},
			"",
		},
		{
			"at least 1 arguments, enough (2)",
			tDescription{PosArgsMin: 1, PosArgsMax: -1},
			tArguments{Command: "example", Pos: []string{"hello", "cruel", "world"}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &tContext{
				Args:        tt.arguments,
				Description: tt.options,
			}
			err := context.CheckPositionalCount()
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("Context.checkPositionalCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommandArguments_checkFilterArgument(t *testing.T) {
	tests := []struct {
		name      string
		options   tDescription
		arguments tArguments

		wantErr string
	}{
		{
			"for not allowed, for not given",
			tDescription{Requirements: env.Requirement{AllowsFilter: false}},
			tArguments{Command: "example"},
			"",
		},
		{
			"for not allowed, for given",
			tDescription{Requirements: env.Requirement{AllowsFilter: false}},
			tArguments{Command: "example", Flags: tFlags{Filters: []string{"pattern"}}},
			"Wrong number of arguments: 'example' takes no '--for' argument. ",
		},

		{
			"fuzzy not allowed, fuzzy given",
			tDescription{Requirements: env.Requirement{AllowsFilter: false}},
			tArguments{Command: "example", Flags: tFlags{NoFuzzyFilter: true, Filters: nil}},
			"Wrong number of arguments: 'example' takes no '--no-fuzzy-filter' argument. ",
		},

		{
			"for allowed, for not given",
			tDescription{Requirements: env.Requirement{AllowsFilter: true}},
			tArguments{Command: "example", Flags: tFlags{Filters: nil}},
			"",
		},
		{
			"for allowed, for given",
			tDescription{Requirements: env.Requirement{AllowsFilter: true}},
			tArguments{Command: "example", Flags: tFlags{Filters: []string{"pattern"}}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := &tContext{
				Args:        tt.arguments,
				Description: tt.options,
			}
			err := context.CheckFilterArgument()
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("Context.checkFilterArgument() error = %q, wantErr %q", err, tt.wantErr)
			}
		})
	}
}
