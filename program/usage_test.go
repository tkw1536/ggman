package program

import (
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/env"
)

func TestProgram_UsagePage(t *testing.T) {
	var program Program
	program.Register(fakeCommand("a"))
	program.Register(fakeCommand("c"))
	program.Register(fakeCommand("b"))

	got := program.UsagePage().String()
	want := "Usage: ggman [--help|-h] [--version|-v] [--for|-f filter] [--no-fuzzy-filter|-n] [--here|-H] [--dirty|-d] [--clean|-c] [--synced|-s] [--unsynced|-u] [--tarnished|-t] [--pristine|-p] [--] COMMAND [ARGS...]\n\nggman manages local git repositories.\n\nggman version v0.0.0-unknown\nggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information.\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -f, --for filter\n      Filter list of repositories to apply COMMAND to by filter. Filter can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url\n\n   -n, --no-fuzzy-filter\n      Disable fuzzy matching for filters\n\n   -H, --here\n      Filter the list of repositories to apply COMMAND to only contain the repository in the current directory\n\n   -d, --dirty\n      List only repositories with uncommited changes\n\n   -c, --clean\n      List only repositories without uncommited changes\n\n   -s, --synced\n      List only repositories which are up-to-date with remote\n\n   -u, --unsynced\n      List only repositories not up-to-date with remote\n\n   -t, --tarnished\n      List only repositories which are dirty or unsynced\n\n   -p, --pristine\n      List only repositories which are clean and synced\n\n   COMMAND [ARGS...]\n      Command to call. One of \"a\", \"b\", \"c\". See individual commands for more help."
	if got != want {
		t.Errorf("Program.UsagePage() = %q, want %q", got, want)
	}
}

// fakeCommand represents a dummy command with the given name
type fakeCommand string

var _ Command = (fakeCommand)("")

func (fakeCommand) BeforeRegister(program *Program) {}
func (f fakeCommand) Description() Description      { return Description{Name: string(f)} }
func (fakeCommand) AfterParse() error               { panic("fakeCommand: not implemented") }
func (fakeCommand) Run(Context) error               { panic("fakeCommand: not implemented") }

func TestCommandArguments_UsagePage(t *testing.T) {

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
			"Usage: ggman [--help|-h] [--version|-v] [--for|-f filter] [--no-fuzzy-filter|-n] [--here|-H] [--dirty|-d] [--clean|-c] [--synced|-s] [--unsynced|-u] [--tarnished|-t] [--pristine|-p] [--] a [--bool|-b random] [--int dummy]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -f, --for filter\n      Filter list of repositories to apply COMMAND to by filter. Filter can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url\n\n   -n, --no-fuzzy-filter\n      Disable fuzzy matching for filters\n\n   -H, --here\n      Filter the list of repositories to apply COMMAND to only contain the repository in the current directory\n\n   -d, --dirty\n      List only repositories with uncommited changes\n\n   -c, --clean\n      List only repositories without uncommited changes\n\n   -s, --synced\n      List only repositories which are up-to-date with remote\n\n   -u, --unsynced\n      List only repositories not up-to-date with remote\n\n   -t, --tarnished\n      List only repositories which are dirty or unsynced\n\n   -p, --pristine\n      List only repositories which are clean and synced\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   \n      usage",
		},

		{
			"command without args and not allowing filter",
			fields{Environment: env.Requirement{}, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   \n      usage",
		},

		{
			"command with max finite args",
			fields{Environment: env.Requirement{}, MaxArgs: 4, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy] [--] [META [META [META [META]]]]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   [META [META [META [META]]]]\n      usage",
		},

		{
			"command with finite args",
			fields{Environment: env.Requirement{}, MinArgs: 1, MaxArgs: 2, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy] [--] META [META]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META]\n      usage",
		},

		{
			"command with infinite args",
			fields{Environment: env.Requirement{}, MinArgs: 1, MaxArgs: -1, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy] [--] META [META ...]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META ...]\n      usage",
		},

		{
			"command with description",
			fields{Description: "A fake command", Environment: env.Requirement{}, MinArgs: 1, MaxArgs: -1, UsageDescription: "usage", Metavar: "META"},
			args{"a"},
			"Usage: ggman [--help|-h] [--version|-v] [--] a [--bool|-b random] [--int dummy] [--] META [META ...]\n\nA fake command\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -b, --bool random\n      a random boolean argument with short\n\n   --int dummy\n      a dummy integer flag (default 12)\n\n   META [META ...]\n      usage",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := &CommandArguments{
				Arguments: Arguments{
					Command: tt.args.command,
				},

				parser: parser,

				description: Description{
					Description:       tt.fields.Description,
					Environment:       tt.fields.Environment,
					PosArgsMin:        tt.fields.MinArgs,
					PosArgsMax:        tt.fields.MaxArgs,
					PosArgName:        tt.fields.Metavar,
					PosArgDescription: tt.fields.UsageDescription,
				},
			}
			if gotUsage := args.UsagePage().String(); gotUsage != tt.wantUsage {
				t.Errorf("CommandArguments.UsagePage() = %q, want %q", gotUsage, tt.wantUsage)
			}
		})
	}
}
