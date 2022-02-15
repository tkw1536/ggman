package program

import (
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/program/meta"
)

func TestProgram_MainUsage(t *testing.T) {
	program := makeProgram()
	program.Info = ttInfo

	program.Register(makeEchoCommand("a"))
	program.Register(makeEchoCommand("c"))
	program.Register(makeEchoCommand("b"))

	got := program.MainUsage().String()
	want := "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] COMMAND [ARGS...]\n\nsomething something dark side\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\n   COMMAND [ARGS...]\n      Command to call. One of \"a\", \"b\", \"c\". See individual commands for more help."
	if got != want {
		t.Errorf("Program.MainUsage() = %q, want %q", got, want)
	}
}

func TestProgram_CommandUsage(t *testing.T) {

	program := iProgram{
		Info: ttInfo,
	}

	// define requirements to allow only the Global1 (or any) arguments
	reqOne := ttRequirements(func(flag meta.Flag) bool {
		return flag.FieldName == "Global1"
	})

	// define requirements to allow anything
	reqAny := ttRequirements(func(flag meta.Flag) bool { return true })

	parser := flags.NewParser(&struct {
		Boolean bool `short:"b" value-name:"random" long:"bool" description:"a random boolean argument with short"`
		Int     int  `long:"int" value-name:"dummy" description:"a dummy integer flag" default:"12"`
	}{}, flags.Default)

	type args struct {
		Command     string
		Description string
		Requirement ttRequirements
		Positional  meta.Positional
	}
	tests := []struct {
		name      string
		args      args
		wantUsage string
	}{
		{
			"command without args and allowing all globals",
			args{Command: "cmd", Requirement: reqAny, Positional: meta.Positional{}},
			"Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--] cmd [--bool|-b random] [--int dummy]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)",
		},

		{
			"command without args and allowing only global1",
			args{Command: "cmd", Requirement: reqOne, Positional: meta.Positional{Description: "usage", Value: "META"}},
			"Usage: exe [--help|-h] [--version|-v] [--] cmd [--bool|-b random] [--int dummy]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   \n      usage",
		},

		{
			"command with max finite args",
			args{Command: "cmd", Requirement: reqOne, Positional: meta.Positional{Max: 4, Description: "usage", Value: "META"}},
			"Usage: exe [--help|-h] [--version|-v] [--] cmd [--bool|-b random] [--int dummy] [--] [META [META [META [META]]]]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   [META [META [META [META]]]]\n      usage",
		},

		{
			"command with finite args",
			args{Command: "cmd", Requirement: reqOne, Positional: meta.Positional{Min: 1, Max: 2, Description: "usage", Value: "META"}},
			"Usage: exe [--help|-h] [--version|-v] [--] cmd [--bool|-b random] [--int dummy] [--] META [META]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META]\n      usage",
		},

		{
			"command with infinite args",
			args{Command: "cmd", Requirement: reqOne, Positional: meta.Positional{Min: 1, Max: -1, Description: "usage", Value: "META"}},
			"Usage: exe [--help|-h] [--version|-v] [--] cmd [--bool|-b random] [--int dummy] [--] META [META ...]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META ...]\n      usage",
		},

		{
			"command with description",
			args{Command: "cmd", Description: "A fake command", Requirement: reqOne, Positional: meta.Positional{Min: 1, Max: -1, Description: "usage", Value: "META"}},
			"Usage: exe [--help|-h] [--version|-v] [--] cmd [--bool|-b random] [--int dummy] [--] META [META ...]\n\nA fake command\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META ...]\n      usage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context := iContext{
				Args: iArguments{
					Command: tt.args.Command,
				},

				parser: parser, // TODO: Fix public / private issue

				Description: iDescription{
					Command:      tt.args.Command,
					Description:  tt.args.Description,
					Positional:   tt.args.Positional,
					Requirements: tt.args.Requirement,
				},
			}
			if gotUsage := program.CommandUsage(context).String(); gotUsage != tt.wantUsage {
				t.Errorf("Program.CommandUsage() = %q\n\n, want %q", gotUsage, tt.wantUsage)
			}
		})
	}
}
