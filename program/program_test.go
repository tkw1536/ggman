package program

import (
	"bytes"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/internal/testutil"
)

func TestProgram_Main(t *testing.T) {
	root, cleanup := testutil.TempDir()
	defer cleanup()

	// create buffers for input and output
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer

	// create a dummy program
	var program Program
	program.IOStream = ggman.NewIOStream(&stdoutBuffer, &stderrBuffer, nil, 80)

	tests := []struct {
		name      string
		args      []string
		options   Options
		variables env.Variables

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
			wantStdout: "Usage: ggman [--help|-h] [--version|-v] [--for|-f filter] [--here|-H] [--]\nCOMMAND [ARGS...]\n\nggman manages local git repositories.\n\nggman version v0.0.0-unknown\nggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information. \n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -f, --for filter\n      Filter list of repositories to apply COMMAND to by filter. Filter can be\n      a relative or absolute path, or a glob pattern which will be matched\n      against the normalized repository url\n\n   -H, --here\n      Filter the list of repositories to apply COMMAND to only contain the\n      repository in the current directory\n\n   COMMAND [ARGS...]\n      Command to call. One of \"fake\". See individual commands for more help.\n",
			wantCode:   0,
		},

		{
			name:       "display help, don't run command",
			args:       []string{"--help", "fake", "whatever"},
			wantStdout: "Usage: ggman [--help|-h] [--version|-v] [--for|-f filter] [--here|-H] [--]\nCOMMAND [ARGS...]\n\nggman manages local git repositories.\n\nggman version v0.0.0-unknown\nggman is licensed under the terms of the MIT License.\nUse 'ggman license' to view licensing information. \n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -f, --for filter\n      Filter list of repositories to apply COMMAND to by filter. Filter can be\n      a relative or absolute path, or a glob pattern which will be matched\n      against the normalized repository url\n\n   -H, --here\n      Filter the list of repositories to apply COMMAND to only contain the\n      repository in the current directory\n\n   COMMAND [ARGS...]\n      Command to call. One of \"fake\". See individual commands for more help.\n",
			wantCode:   0,
		},

		{
			name:       "display version",
			args:       []string{"--version"},
			wantStdout: "ggman version v0.0.0-unknown, built 1970-01-01 00:00:00 +0000 UTC, using " + runtime.Version() + "\n",
			wantCode:   0,
		},

		{
			name:       "command help",
			args:       []string{"fake", "--help"},
			wantStdout: "Usage: ggman [--help|-h] [--version|-v] [--] fake [--stdout|-o message]\n[--stderr|-e message]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -o, --stdout message\n       (default write to stdout)\n\n   -e, --stderr message\n       (default write to stderr)\n",
			wantCode:   0,
		},

		{
			name:       "command help",
			args:       []string{"--", "fake", "--help"},
			wantStdout: "Usage: ggman [--help|-h] [--version|-v] [--] fake [--stdout|-o message]\n[--stderr|-e message]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\nCommand Arguments:\n\n   -o, --stdout message\n       (default write to stdout)\n\n   -e, --stderr message\n       (default write to stderr)\n",
			wantCode:   0,
		},

		{
			name:       "not enough arguments for fake",
			args:       []string{"fake"},
			options:    Options{MinArgs: 1, MaxArgs: 2},
			wantStderr: "Wrong number of arguments: 'fake' takes between 1 and 2 arguments. \n",
			wantCode:   4,
		},

		{
			name:       "'fake' with unknown argument (not allowed)",
			args:       []string{"fake", "--argument-not-declared"},
			options:    Options{MinArgs: 0, MaxArgs: -1},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "",
			wantStderr: "Error parsing flags: unknown flag `argument-not-declared'\n",
			wantCode:   4,
		},

		{
			name:       "'fake' with unknown argument (allowed)",
			args:       []string{"fake", "--argument-not-declared"},
			options:    Options{MinArgs: 0, MaxArgs: -1, SkipUnknownFlags: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: --argument-not-declared\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' without filter",
			args:       []string{"fake", "hello", "world"},
			options:    Options{MinArgs: 1, MaxArgs: 2},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name:       "'fake' without filter",
			args:       []string{"fake", "hello", "world"},
			options:    Options{MinArgs: 1, MaxArgs: 2},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name:       "'fake' with needsRoot, but no root",
			args:       []string{"fake", "hello", "world"},
			options:    Options{Environment: env.Requirement{NeedsRoot: true}, MinArgs: 1, MaxArgs: 2},
			wantStderr: "Unable to find GGROOT directory. \n",
			wantCode:   5,
		},
		{
			name:       "'fake' with needsroot and root",
			args:       []string{"fake", "hello", "world"},
			options:    Options{Environment: env.Requirement{NeedsRoot: true}, MinArgs: 1, MaxArgs: 2},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name:       "'fake' with filter but not allowed",
			args:       []string{"--for", "example", "fake", "hello", "world"},
			options:    Options{MinArgs: 1, MaxArgs: 2},
			wantStderr: "Wrong number of arguments: 'fake' takes no 'for' argument. \n",
			wantCode:   4,
		},
		{
			name:       "'fake' with filter but no root",
			args:       []string{"--for", "example", "fake", "hello", "world"},
			options:    Options{Environment: env.Requirement{AllowsFilter: true}, MinArgs: 1, MaxArgs: 2},
			wantStderr: "Unable to find GGROOT directory. \n",
			wantCode:   5,
		},

		{
			name:       "'fake' with filter",
			args:       []string{"--for", "example", "fake", "hello", "world"},
			options:    Options{Environment: env.Requirement{AllowsFilter: true}, MinArgs: 1, MaxArgs: 2},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: example\nGot arguments: hello,world\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with here",
			args:       []string{"--here", "fake", "hello", "world"},
			options:    Options{Environment: env.Requirement{AllowsFilter: true}, MinArgs: 1, MaxArgs: 2},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "",
			wantStderr: "Unable to initialize context: Unable to find current repository: Unable to\nresolve repository \".\"\n",
			wantCode:   5,
		},

		{
			name:       "'fake' with non-global here argument",
			args:       []string{"--", "fake", "--here"},
			options:    Options{MinArgs: 0, MaxArgs: -1, SkipUnknownFlags: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: --here\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with parsed short argument",
			args:       []string{"fake", "-o", "message"},
			options:    Options{MinArgs: 0, MaxArgs: -1, SkipUnknownFlags: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: \nmessage\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with non-parsed short argument",
			args:       []string{"fake", "--", "--s", "message"},
			options:    Options{MinArgs: 0, MaxArgs: -1, SkipUnknownFlags: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: --s,message\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with parsed long argument",
			args:       []string{"fake", "--stdout", "message"},
			options:    Options{MinArgs: 0, MaxArgs: -1, SkipUnknownFlags: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: \nmessage\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with non-parsed long argument",
			args:       []string{"fake", "--", "--stdout", "message"},
			options:    Options{MinArgs: 0, MaxArgs: -1, SkipUnknownFlags: true},
			variables:  env.Variables{GGROOT: root},
			wantStdout: "Got filter: \nGot arguments: --stdout,message\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with failure ",
			args:       []string{"fake", "fail"},
			options:    Options{MinArgs: 1, MaxArgs: 2},
			wantStdout: "Got filter: \nGot arguments: fail\nwrite to stdout\n",
			wantStderr: "write to stderr\nTest Failure\n",
			wantCode:   1,
		},

		{
			name:       "'notexistent' command",
			args:       []string{"notexistent"},
			wantStderr: "Unknown command. Must be one of \"fake\". \n",
			wantCode:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset the buffers
			stdoutBuffer.Reset()
			stderrBuffer.Reset()

			// add a fake command
			fake := &echoCommand{name: "fake", options: tt.options}
			program.commands = nil
			program.Register(fake)

			// run the program
			ret := ggman.AsError(program.Main(env.EnvironmentParameters{
				Variables: tt.variables,
				Plumbing:  git.NewPlumbing(),
				Workdir:   "",
			}, tt.args))

			// check all the error values
			gotCode := uint8(ret.ExitCode)
			gotStdout := stdoutBuffer.String()
			gotStderr := stderrBuffer.String()

			if gotCode != tt.wantCode {
				t.Errorf("Program.Main() code = %v, wantCode %v", gotCode, tt.wantCode)
			}

			if gotStdout != tt.wantStdout {
				t.Errorf("Program.Main() stdout = %q, wantStdout %q", gotStdout, tt.wantStdout)
			}

			if gotStderr != tt.wantStderr {
				t.Errorf("Program.Main() stderr = %q, wantStderr %q", gotStderr, tt.wantStderr)
			}
		})
	}
}

// echoCommand is a fake command that can be used for testing.
type echoCommand struct {
	name      string
	StdoutMsg string `short:"o" long:"stdout" value-name:"message" default:"write to stdout"`
	StderrMsg string `short:"e" long:"stderr" value-name:"message" default:"write to stderr"`
	options   Options
}

func (e echoCommand) Name() string {
	return e.name
}
func (e echoCommand) Options() Options {
	return e.options
}
func (e echoCommand) AfterParse() error { return nil }
func (e echoCommand) Run(context Context) error {
	context.Stdout.Write([]byte("Got filter: " + strings.Join(context.Filters, ",")))
	context.Stdout.Write([]byte("\nGot arguments: " + strings.Join(context.Args, ",")))
	context.Stdout.Write([]byte("\n" + e.StdoutMsg + "\n"))
	context.Stderr.Write([]byte(e.StderrMsg + "\n"))

	if len(context.Args) > 0 && context.Args[0] == "fail" {
		return ggman.Error{ExitCode: ggman.ExitGeneric, Message: "Test Failure"}
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
