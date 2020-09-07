package program

import (
	"testing"

	"github.com/tkw1536/ggman/env"
)

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
			Arguments{Command: "example", Argv: []string{}},
			"",
		},

		// taking 1 arg
		{
			"one argument, too few",
			Options{MinArgs: 1, MaxArgs: 1},
			Arguments{Command: "example", Argv: []string{}},
			"Wrong number of arguments: 'example' takes exactly 1 argument(s). ",
		},
		{
			"one argument, exactly enough",
			Options{MinArgs: 1, MaxArgs: 1},
			Arguments{Command: "example", Argv: []string{"world"}},
			"",
		},
		{
			"one argument, too many",
			Options{MinArgs: 1, MaxArgs: 1},
			Arguments{Command: "example", Argv: []string{"hello", "world"}},
			"Wrong number of arguments: 'example' takes exactly 1 argument(s). ",
		},

		// taking 1 or 2 args
		{
			"1-2 arguments, too few",
			Options{MinArgs: 1, MaxArgs: 2},
			Arguments{Command: "example", Argv: []string{}},
			"Wrong number of arguments: 'example' takes between 1 and 2 arguments. ",
		},
		{
			"1-2 arguments, enough",
			Options{MinArgs: 1, MaxArgs: 2},
			Arguments{Command: "example", Argv: []string{"world"}},
			"",
		},
		{
			"1-2 arguments, enough (2)",
			Options{MinArgs: 1, MaxArgs: 2},
			Arguments{Command: "example", Argv: []string{"hello", "world"}},
			"",
		},
		{
			"1-2 arguments, too many",
			Options{MinArgs: 1, MaxArgs: 2},
			Arguments{Command: "example", Argv: []string{"hello", "world", "you"}},
			"Wrong number of arguments: 'example' takes between 1 and 2 arguments. ",
		},

		// taking 2 args
		{
			"two arguments, too few",
			Options{MinArgs: 2, MaxArgs: 2},
			Arguments{Command: "example", Argv: []string{}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},
		{
			"two arguments, too few (2)",
			Options{MinArgs: 2, MaxArgs: 2},
			Arguments{Command: "example", Argv: []string{"world"}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},
		{
			"two arguments, enough",
			Options{MinArgs: 2, MaxArgs: 2},
			Arguments{Command: "example", Argv: []string{"hello", "world"}},
			"",
		},
		{
			"two arguments, too many",
			Options{MinArgs: 2, MaxArgs: 2},
			Arguments{Command: "example", Argv: []string{"hello", "world", "you"}},
			"Wrong number of arguments: 'example' takes exactly 2 argument(s). ",
		},

		// at least one argument
		{
			"at least 1 arguments, not enough",
			Options{MinArgs: 1, MaxArgs: -1},
			Arguments{Command: "example", Argv: []string{}},
			"Wrong number of arguments: 'example' takes at least 1 argument(s). ",
		},
		{
			"at least 1 arguments, enough",
			Options{MinArgs: 1, MaxArgs: -1},
			Arguments{Command: "example", Argv: []string{"hello"}},
			"",
		},
		{
			"at least 1 arguments, enough (2)",
			Options{MinArgs: 1, MaxArgs: -1},
			Arguments{Command: "example", Argv: []string{"hello", "cruel", "world"}},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &CommandArguments{
				Options:   tt.options,
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

func TestCommandArguments_checkForArgument(t *testing.T) {
	tests := []struct {
		name      string
		options   Options
		arguments Arguments

		wantErr string
	}{
		{
			"for not allowed, for not given",
			Options{Environment: env.Requirement{AllowsFilter: false}},
			Arguments{Command: "example", For: ""},
			"",
		},
		{
			"for not allowed, for given",
			Options{Environment: env.Requirement{AllowsFilter: false}},
			Arguments{Command: "example", For: "pattern"},
			"Wrong number of arguments: 'example' takes no 'for' argument. ",
		},

		{
			"for allowed, for not given",
			Options{Environment: env.Requirement{AllowsFilter: true}},
			Arguments{Command: "example", For: ""},
			"",
		},
		{
			"for allowed, for given",
			Options{Environment: env.Requirement{AllowsFilter: true}},
			Arguments{Command: "example", For: "pattern"},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := CommandArguments{
				Options:   tt.options,
				Arguments: tt.arguments,
			}
			err := args.checkForArgument()
			gotErr := ""
			if err != nil {
				gotErr = err.Error()
			}
			if gotErr != tt.wantErr {
				t.Errorf("CommandArguments.checkForArgument() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
