package program

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/internal/testutil"
	"github.com/tkw1536/ggman/internal/text"
	"github.com/tkw1536/ggman/program/exit"
	"github.com/tkw1536/ggman/program/stream"
)

// TODO: Fix broken tests (after type parameters)

func TestProgram_Main(t *testing.T) {
	root := testutil.TempDirAbs(t)
	if err := os.Mkdir(filepath.Join(root, "real"), os.ModeDir&os.ModePerm); err != nil {
		panic(err)
	}

	// create buffers for input and output
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	// create a dummy program
	program := Program{
		Initalizer: func(params env.EnvironmentParameters, cmdargs CommandArguments) (Runtime, error) {
			return nil, nil
		},
		Info: testInfo,
	}

	wrapLength := 80
	stream := stream.NewIOStream(&stdoutBuffer, &stderrBuffer, nil, wrapLength)

	tests := []struct {
		name      string
		args      []string
		options   Description
		variables env.Variables
		workdir   string

		alias Alias

		// wrap the "want" variables automatically?
		wrapError bool
		wrapOut   bool

		wantStdout string
		wantStderr string
		wantCode   uint8
	}{
		{
			name:       "no arguments",
			args:       []string{},
			wantStderr: "Unable to parse arguments: Need at least one argument. Use `ggman license` to\nview licensing information.\n",
			wantCode:   3,
		},

		{
			name:       "unknown general args",
			args:       []string{"--this-flag-doesnt-exist", "--", "fake"},
			wantStderr: "Unable to parse arguments: unknown flag `this-flag-doesnt-exist'\n",
			wantCode:   3,
		},

		{
			name:       "display help",
			args:       []string{"--help"},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--for|-f filter] [--no-fuzzy-filter|-n]\n[--here|-H] [--path|-P] [--dirty|-d] [--clean|-c] [--synced|-s] [--unsynced|-u]\n[--tarnished|-t] [--pristine|-p] [--] COMMAND [ARGS...]\n\nsomething something dark side\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -f, --for filter\n      Filter list of repositories to apply COMMAND to by filter. Filter can be\n      a relative or absolute path, or a glob pattern which will be matched\n      against the normalized repository url\n\n   -n, --no-fuzzy-filter\n      Disable fuzzy matching for filters\n\n   -H, --here\n      Filter the list of repositories to apply COMMAND to only contain\n      repository in the current directory or subtree. Alias for '-p .'\n\n   -P, --path\n      Filter the list of repositories to apply COMMAND to only contain\n      repositories in or under the specified path. May be used multiple times\n\n   -d, --dirty\n      List only repositories with uncommited changes\n\n   -c, --clean\n      List only repositories without uncommited changes\n\n   -s, --synced\n      List only repositories which are up-to-date with remote\n\n   -u, --unsynced\n      List only repositories not up-to-date with remote\n\n   -t, --tarnished\n      List only repositories which are dirty or unsynced\n\n   -p, --pristine\n      List only repositories which are clean and synced\n\n   COMMAND [ARGS...]\n      Command to call. One of \"fake\". See individual commands for more help.\n",
			wantCode:   0,
		},

		{
			name:       "display help, don't run command",
			args:       []string{"--help", "fake", "whatever"},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--for|-f filter] [--no-fuzzy-filter|-n]\n[--here|-H] [--path|-P] [--dirty|-d] [--clean|-c] [--synced|-s] [--unsynced|-u]\n[--tarnished|-t] [--pristine|-p] [--] COMMAND [ARGS...]\n\nsomething something dark side\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -f, --for filter\n      Filter list of repositories to apply COMMAND to by filter. Filter can be\n      a relative or absolute path, or a glob pattern which will be matched\n      against the normalized repository url\n\n   -n, --no-fuzzy-filter\n      Disable fuzzy matching for filters\n\n   -H, --here\n      Filter the list of repositories to apply COMMAND to only contain\n      repository in the current directory or subtree. Alias for '-p .'\n\n   -P, --path\n      Filter the list of repositories to apply COMMAND to only contain\n      repositories in or under the specified path. May be used multiple times\n\n   -d, --dirty\n      List only repositories with uncommited changes\n\n   -c, --clean\n      List only repositories without uncommited changes\n\n   -s, --synced\n      List only repositories which are up-to-date with remote\n\n   -u, --unsynced\n      List only repositories not up-to-date with remote\n\n   -t, --tarnished\n      List only repositories which are dirty or unsynced\n\n   -p, --pristine\n      List only repositories which are clean and synced\n\n   COMMAND [ARGS...]\n      Command to call. One of \"fake\". See individual commands for more help.\n",
			wantCode:   0,
		},

		{
			name:       "display version",
			args:       []string{"--version"},
			wantStdout: "exe version 42.0.0, built 1970-01-01 00:00:00 +0000 UTC, using " + runtime.Version() + "\n",
			wrapOut:    true,
			wantCode:   0,
		},

		{
			name:       "command help",
			args:       []string{"fake", "--help"},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--] fake [--stdout|-o message]\n[--stderr|-e message]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -o, --stdout message\n       (default write to stdout)\n\n   -e, --stderr message\n       (default write to stderr)\n",
			wantCode:   0,
		},

		{
			name: "alias help",

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			args:       []string{"alias", "--help"},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--] alias [--] [ARG ...]\n\nAlias for `exe fake`. See `exe fake --help` for detailed help page about fake.\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake`.\n",
			wantCode:   0,
		},

		{
			name:       "command help",
			args:       []string{"--", "fake", "--help"},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--] fake [--stdout|-o message]\n[--stderr|-e message]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -o, --stdout message\n       (default write to stdout)\n\n   -e, --stderr message\n       (default write to stderr)\n",
			wantCode:   0,
		},

		{
			name: "alias help",

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			args:       []string{"--", "alias", "--help"},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--] alias [--] [ARG ...]\n\nAlias for `exe fake`. See `exe fake --help` for detailed help page about fake.\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake`.\n",
			wantCode:   0,
		},

		{
			name: "long alias help",

			alias: Alias{
				Name:        "alias",
				Command:     "fake",
				Args:        []string{"something", "else"},
				Description: "Some useful alias",
			},

			args:       []string{"alias", "--help"},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--] alias [--] [ARG ...]\n\nSome useful alias\n\nAlias for `exe fake something else`. See `exe fake --help` for detailed help\npage about fake.\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake something else`.\n",
			wantCode:   0,
		},

		{
			name:       "not enough arguments for fake",
			args:       []string{"fake"},
			options:    Description{PosArgsMin: 1, PosArgsMax: 2},
			wantStderr: "Wrong number of arguments: 'fake' takes between 1 and 2 arguments.\n",
			wantCode:   4,
		},

		{
			name:       "'fake' with unknown argument (not allowed)",
			args:       []string{"fake", "--argument-not-declared"},
			options:    Description{PosArgsMin: 0, PosArgsMax: -1},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "",
			wantStderr: "Error parsing flags: unknown flag `argument-not-declared'\n",
			wantCode:   4,
		},

		{
			name:       "'fake' with unknown argument (allowed)",
			args:       []string{"fake", "--argument-not-declared"},
			options:    Description{PosArgsMin: 0, PosArgsMax: -1, SkipUnknownOptions: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: --argument-not-declared\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' without filter",
			args:       []string{"fake", "hello", "world"},
			options:    Description{PosArgsMin: 1, PosArgsMax: 2},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name:       "'fake' without filter",
			args:       []string{"fake", "hello", "world"},
			options:    Description{PosArgsMin: 1, PosArgsMax: 2},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		/*
			{
				name:       "'fake' with needsRoot, but no root",
				args:       []string{"fake", "hello", "world"},
				options:    Description{Environment: env.Requirement{NeedsRoot: true}, PosArgsMin: 1, PosArgsMax: 2},
				wantStderr: "Unable to find GGROOT directory.\n",
				wantCode:   5,
			},
		*/
		{
			name:       "'fake' with needsroot and root",
			args:       []string{"fake", "hello", "world"},
			options:    Description{Environment: env.Requirement{NeedsRoot: true}, PosArgsMin: 1, PosArgsMax: 2},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name:       "'fake' with filter but not allowed",
			args:       []string{"--for", "example", "fake", "hello", "world"},
			options:    Description{PosArgsMin: 1, PosArgsMax: 2},
			wantStderr: "Wrong number of arguments: 'fake' takes no '--for' argument.\n",
			wantCode:   4,
		},
		/*
			{
				name:       "'fake' with filter but no root",
				args:       []string{"--for", "example", "fake", "hello", "world"},
				options:    Description{Environment: env.Requirement{AllowsFilter: true}, PosArgsMin: 1, PosArgsMax: 2},
				wantStderr: "Unable to find GGROOT directory.\n",
				wantCode:   5,
			},
		*/

		{
			name:       "'fake' with filter",
			args:       []string{"--for", "example", "fake", "hello", "world"},
			options:    Description{Environment: env.Requirement{AllowsFilter: true}, PosArgsMin: 1, PosArgsMax: 2},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: example\nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		/*
			{
				name:       "'fake' with here (not working)",
				args:       []string{"--here", "fake", "hello", "world"},
				options:    Description{Environment: env.Requirement{AllowsFilter: true}, PosArgsMin: 1, PosArgsMax: 2},
				variables:  env.Variables{GGROOT: root},
				workdir:    filepath.Join(root, "doesnotexist"),
				wantStdout: "",
				wantStderr: "Unable to initialize context: Not a directory: \".\"\n",
				wantCode:   5,
			},
		*/

		{
			name:       "'fake' with path (working)",
			args:       []string{"--path", "real", "fake", "hello", "world"},
			options:    Description{Environment: env.Requirement{AllowsFilter: true}, PosArgsMin: 1, PosArgsMax: 2},
			variables:  env.Variables{GGROOT: root},
			workdir:    root,
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with non-global here argument",
			args:       []string{"--", "fake", "--here"},
			options:    Description{PosArgsMin: 0, PosArgsMax: -1, SkipUnknownOptions: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: --here\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with parsed short argument",
			args:       []string{"fake", "-o", "message"},
			options:    Description{PosArgsMin: 0, PosArgsMax: -1, SkipUnknownOptions: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: \nmessage\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with non-parsed short argument",
			args:       []string{"fake", "--", "--s", "message"},
			options:    Description{PosArgsMin: 0, PosArgsMax: -1, SkipUnknownOptions: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: --s,message\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with parsed long argument",
			args:       []string{"fake", "--stdout", "message"},
			options:    Description{PosArgsMin: 0, PosArgsMax: -1, SkipUnknownOptions: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: \nmessage\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with non-parsed long argument",
			args:       []string{"fake", "--", "--stdout", "message"},
			options:    Description{PosArgsMin: 0, PosArgsMax: -1, SkipUnknownOptions: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: --stdout,message\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with failure ",
			args:       []string{"fake", "fail"},
			options:    Description{PosArgsMin: 1, PosArgsMax: 2},
			wantStdout: "Got filter: \nGot arguments: fail\nwrite to stdout\n",
			wantStderr: "write to stderr\nTest Failure\n",
			wantCode:   1,
		},

		{
			name:       "'notexistent' command",
			args:       []string{"notexistent"},
			wantStderr: "Unknown command. Must be one of \"fake\".\n",
			wantCode:   2,
		},

		{
			name: "'notexistent' command (with alias)",

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			args:       []string{"notexistent"},
			wantStderr: "Unknown command. Must be one of \"fake\".\n",
			wantCode:   2,
		},

		{
			name: "'alias' without args",
			args: []string{"alias", "hello", "world"},

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			options:    Description{PosArgsMin: 0, PosArgsMax: -1},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name: "'alias' with args",
			args: []string{"alias", "world"},

			alias: Alias{
				Name:    "alias",
				Command: "fake",
				Args:    []string{"hello"},
			},

			options:    Description{PosArgsMin: 0, PosArgsMax: -1},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name: "recursive 'fake' alias ",
			args: []string{"fake", "world"},

			alias: Alias{
				Name:    "fake",
				Command: "fake",
				Args:    []string{"hello"},
			},

			options:    Description{PosArgsMin: 0, PosArgsMax: -1},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset the buffers
			stdoutBuffer.Reset()
			stderrBuffer.Reset()

			// add a fake command
			fake := &echoCommand{name: "fake", description: tt.options}
			program.commands = nil
			program.Register(fake)

			program.aliases = nil
			if tt.alias.Name != "" {
				program.RegisterAlias(tt.alias)
			}

			// run the program
			ret := exit.AsError(program.Main(stream, env.EnvironmentParameters{
				Variables: tt.variables,
				Plumbing:  git.NewPlumbing(),
				Workdir:   tt.workdir,
			}, tt.args))

			// check all the error values
			gotCode := uint8(ret.ExitCode)
			gotStdout := stdoutBuffer.String()
			gotStderr := stderrBuffer.String()

			if gotCode != tt.wantCode {
				t.Errorf("Program.Main() code = %v, wantCode %v", gotCode, tt.wantCode)
			}

			// wrap if requested
			if tt.wrapOut {
				tt.wantStdout = text.WrapString(wrapLength, tt.wantStdout)
			}

			if gotStdout != tt.wantStdout {
				t.Errorf("Program.Main() stdout = %q, wantStdout %q", gotStdout, tt.wantStdout)
			}

			// wrap if requested
			if tt.wrapError {
				tt.wantStderr = text.WrapString(wrapLength, tt.wantStderr)
			}

			if gotStderr != tt.wantStderr {
				t.Errorf("Program.Main() stderr = %q, wantStderr %q", gotStderr, tt.wantStderr)
			}
		})
	}
}

// echoCommand is a fake command that can be used for testing.
type echoCommand struct {
	name        string
	StdoutMsg   string `short:"o" long:"stdout" value-name:"message" default:"write to stdout"`
	StderrMsg   string `short:"e" long:"stderr" value-name:"message" default:"write to stderr"`
	description Description
}

func (e echoCommand) BeforeRegister(program *Program) {}
func (e echoCommand) Description() Description {
	e.description.Name = e.name
	return e.description
}
func (e echoCommand) AfterParse() error { return nil }
func (e echoCommand) Run(context Context) error {
	context.Stdout.Write([]byte("Got filter: " + strings.Join(context.Filters, ",")))
	context.Stdout.Write([]byte("\nGot arguments: " + strings.Join(context.Args, ",")))
	context.Stdout.Write([]byte("\n" + e.StdoutMsg + "\n"))
	context.Stderr.Write([]byte(e.StderrMsg + "\n"))

	if len(context.Args) > 0 && context.Args[0] == "fail" {
		return exit.Error{ExitCode: exit.ExitGeneric, Message: "Test Failure"}
	}

	return nil
}

func TestProgram_Commands(t *testing.T) {
	var program Program
	program.Register(fakeCommand("a"))
	program.Register(fakeCommand("c"))
	program.Register(fakeCommand("b"))

	got := program.Commands()
	want := []string{"a", "b", "c"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Program.Commands() = %v, want = %v", got, want)
	}
}

func TestProgram_FmtCommands(t *testing.T) {
	var program Program
	program.Register(fakeCommand("a"))
	program.Register(fakeCommand("c"))
	program.Register(fakeCommand("b"))

	got := program.FmtCommands()
	want := `"a", "b", "c"`

	if got != want {
		t.Errorf("Program.FmtCommands() = %v, want = %v", got, want)
	}
}
