package program_test

import (
	"testing"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program"
)

// fakeCommand represents a dummy command with the given name
type fakeCommand string

var _ tCommand = (fakeCommand)("")

func (fakeCommand) BeforeRegister(program *tProgram) {}
func (f fakeCommand) Description() tDescription      { return tDescription{Name: string(f)} }
func (fakeCommand) AfterParse() error                { panic("fakeCommand: not implemented") }
func (fakeCommand) Run(tContext) error               { panic("fakeCommand: not implemented") }

var testInfo = program.Info{
	BuildVersion: "42.0.0",
	BuildTime:    time.Unix(0, 0).UTC(),

	MainName:    "exe",
	Description: "something something dark side",
}

func TestProgram_MainUsage(t *testing.T) {
	program := tProgram{
		Info: testInfo,
	}

	program.Register(fakeCommand("a"))
	program.Register(fakeCommand("c"))
	program.Register(fakeCommand("b"))

	got := program.MainUsage().String()
	want := "Usage: exe [--help|-h] [--version|-v] [--for|-f filter] [--no-fuzzy-filter|-n] [--here|-H] [--path|-P] [--dirty|-d] [--clean|-c] [--synced|-s] [--unsynced|-u] [--tarnished|-t] [--pristine|-p] [--] COMMAND [ARGS...]\n\nsomething something dark side\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -f, --for filter\n      Filter list of repositories to apply COMMAND to by filter. Filter can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url\n\n   -n, --no-fuzzy-filter\n      Disable fuzzy matching for filters\n\n   -H, --here\n      Filter the list of repositories to apply COMMAND to only contain repository in the current directory or subtree. Alias for '-p .'\n\n   -P, --path\n      Filter the list of repositories to apply COMMAND to only contain repositories in or under the specified path. May be used multiple times\n\n   -d, --dirty\n      List only repositories with uncommited changes\n\n   -c, --clean\n      List only repositories without uncommited changes\n\n   -s, --synced\n      List only repositories which are up-to-date with remote\n\n   -u, --unsynced\n      List only repositories not up-to-date with remote\n\n   -t, --tarnished\n      List only repositories which are dirty or unsynced\n\n   -p, --pristine\n      List only repositories which are clean and synced\n\n   COMMAND [ARGS...]\n      Command to call. One of \"a\", \"b\", \"c\". See individual commands for more help."
	if got != want {
		t.Errorf("Program.UsagePage() = %q, want %q", got, want)
	}
}

func TestProgram_CommandUsage(t *testing.T) {

	program := tProgram{
		Info: testInfo,
	}

	parser := flags.NewParser(&struct {
		Boolean bool `short:"b" value-name:"random" long:"bool" description:"a random boolean argument with short"`
		Int     int  `long:"int" value-name:"dummy" description:"a dummy integer flag" default:"12"`
	}{}, flags.Default)

	type fields struct {
		Description      string
		Environment      env.Requirement
		MinArgs          int
		MaxArgs          int
		Metavar          string
		UsageDescription string
	}
	type args struct {
		command string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantUsage string
	}{
		{
			"command without args and allowing filter",
			fields{Environment: env.Requirement{AllowsFilter: true}, UsageDescription: "usage"},
			args{"a"},
			"Usage: exe [--help|-h] [--version|-v] [--for|-f filter] [--no-fuzzy-filter|-n] [--here|-H] [--path|-P] [--dirty|-d] [--clean|-c] [--synced|-s] [--unsynced|-u] [--tarnished|-t] [--pristine|-p] [--] a [--bool|-b random] [--int dummy]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -f, --for filter\n      Filter list of repositories to apply COMMAND to by filter. Filter can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url\n\n   -n, --no-fuzzy-filter\n      Disable fuzzy matching for filters\n\n   -H, --here\n      Filter the list of repositories to apply COMMAND to only contain repository in the current directory or subtree. Alias for '-p .'\n\n   -P, --path\n      Filter the list of repositories to apply COMMAND to only contain repositories in or under the specified path. May be used multiple times\n\n   -d, --dirty\n      List only repositories with uncommited changes\n\n   -c, --clean\n      List only repositories without uncommited changes\n\n   -s, --synced\n      List only repositories which are up-to-date with remote\n\n   -u, --unsynced\n      List only repositories not up-to-date with remote\n\n   -t, --tarnished\n      List only repositories which are dirty or unsynced\n\n   -p, --pristine\n      List only repositories which are clean and synced\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   \n      usage",
		},

		{
			"command without args and not allowing filter",
			fields{Environment: env.Requirement{}, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: exe [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   \n      usage",
		},

		{
			"command with max finite args",
			fields{Environment: env.Requirement{}, MaxArgs: 4, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: exe [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy] [--] [META [META [META [META]]]]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   [META [META [META [META]]]]\n      usage",
		},

		{
			"command with finite args",
			fields{Environment: env.Requirement{}, MinArgs: 1, MaxArgs: 2, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: exe [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy] [--] META [META]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META]\n      usage",
		},

		{
			"command with infinite args",
			fields{Environment: env.Requirement{}, MinArgs: 1, MaxArgs: -1, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: exe [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy] [--] META [META ...]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META ...]\n      usage",
		},

		{
			"command with description",
			fields{Description: "A fake command", Environment: env.Requirement{}, MinArgs: 1, MaxArgs: -1, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: exe [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy] [--] META [META ...]\n\nA fake command\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META ...]\n      usage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tCommandArguments{
				Arguments: Arguments{
					Command: tt.args.command,
				},

				Parser: parser, // TODO: Fix public / private issue

				Description: tDescription{
					Description:       tt.fields.Description,
					Requirements:      tt.fields.Environment,
					PosArgsMin:        tt.fields.MinArgs,
					PosArgsMax:        tt.fields.MaxArgs,
					PosArgName:        tt.fields.Metavar,
					PosArgDescription: tt.fields.UsageDescription,
				},
			}
			if gotUsage := program.CommandUsage(args).String(); gotUsage != tt.wantUsage {
				t.Errorf("CommandArguments.UsagePage() = %q, want %q", gotUsage, tt.wantUsage)
			}
		})
	}
}
