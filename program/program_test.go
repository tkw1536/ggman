package program

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/testutil"
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
			wantStderr: "Unable to parse arguments: unknown flag: --this-flag-doesnt-exist\n",
			wantCode:   3,
		},

		{
			name:       "display help",
			args:       []string{"--help"},
			wantStdout: "ggman version v0.0.0-unknown\n\nUsage: ggman [--help|-h] [--version|-v] [--for|-f filter] [--] COMMAND [ARGS...]\n\n   -h, --help\n       Print this usage dialog and exit.\n\n   -v, --version\n       Print version message and exit.\n\n   -f, --for filter\n       Filter the list of repositories to apply command to by filter.\n\n   COMMAND [ARGS...]\n       Command to call. One of 'fake'. See individual commands for more help.\n\nggman is licensed under the terms of the MIT License. Use 'ggman license' to\nview licensing information.\n",
			wantCode:   0,
		},

		{
			name:       "display help, don't run command",
			args:       []string{"--help", "fake", "whatever"},
			wantStdout: "ggman version v0.0.0-unknown\n\nUsage: ggman [--help|-h] [--version|-v] [--for|-f filter] [--] COMMAND [ARGS...]\n\n   -h, --help\n       Print this usage dialog and exit.\n\n   -v, --version\n       Print version message and exit.\n\n   -f, --for filter\n       Filter the list of repositories to apply command to by filter.\n\n   COMMAND [ARGS...]\n       Command to call. One of 'fake'. See individual commands for more help.\n\nggman is licensed under the terms of the MIT License. Use 'ggman license' to\nview licensing information.\n",
			wantCode:   0,
		},

		{
			name:       "display version",
			args:       []string{"--version"},
			wantStdout: "ggman version v0.0.0-unknown, built 1970-01-01 00:00:00 +0000 UTC\n",
			wantCode:   0,
		},

		{
			name:       "command help",
			args:       []string{"fake", "--help"},
			wantStdout: "Usage: ggman [global arguments] [--] fake [--help|-h]\n\n   -h, --help\n       Print this usage message and exit.\n\n   global arguments\n       Global arguments for ggman. See ggman --help for more information.\n",
			wantCode:   0,
		},

		{
			name:       "command help",
			args:       []string{"--", "fake", "--help"},
			wantStdout: "Usage: ggman [global arguments] [--] fake [--help|-h]\n\n   -h, --help\n       Print this usage message and exit.\n\n   global arguments\n       Global arguments for ggman. See ggman --help for more information.\n",
			wantCode:   0,
		},

		{
			name:       "",
			args:       []string{"--", "fake", "--help"},
			wantStdout: "Usage: ggman [global arguments] [--] fake [--help|-h]\n\n   -h, --help\n       Print this usage message and exit.\n\n   global arguments\n       Global arguments for ggman. See ggman --help for more information.\n",
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
			wantStderr: "Unknown command. Must be one of 'fake'. \n",
			wantCode:   2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset the buffers
			stdoutBuffer.Reset()
			stderrBuffer.Reset()

			// add a fake command
			fake := &programFakeCommandT{name: "fake", options: tt.options}
			program.commands = nil
			program.Register(fake)

			// run the program
			ret := ggman.AsError(program.Main(tt.variables, nil, "", tt.args))

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

// programFakeCommandT is a fake command that can be used for testing.
type programFakeCommandT struct {
	name    string
	options Options
}

func (p programFakeCommandT) Name() string {
	return p.name
}
func (p programFakeCommandT) Options(*pflag.FlagSet) Options {
	return p.options
}
func (programFakeCommandT) AfterParse() error { return nil }
func (programFakeCommandT) Run(context Context) error {
	context.Stdout.Write([]byte("Got filter: " + context.Filter.String()))
	context.Stdout.Write([]byte("\nGot arguments: " + strings.Join(context.Args, ",")))
	context.Stdout.Write([]byte("\nwrite to stdout\n"))
	context.Stderr.Write([]byte("write to stderr\n"))

	if len(context.Args) > 0 && context.Args[0] == "fail" {
		return ggman.Error{ExitCode: ggman.ExitGeneric, Message: "Test Failure"}
	}

	return nil
}
