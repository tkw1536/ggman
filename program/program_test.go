package program

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/tkw1536/ggman/internal/testutil"
	"github.com/tkw1536/ggman/internal/text"
	"github.com/tkw1536/ggman/program/exit"
	"github.com/tkw1536/ggman/program/meta"
	"github.com/tkw1536/ggman/program/stream"
)

// This file contains dummy implementations of everything required to assemble a program.
// It is reused across the test suite, however there is no versioning guarantee.
// It may change in a future revision of the test suite.

// Environment of each command is a single string value.
// Parameters to initialize each command is also a string value.
type tEnvironment string
type tParameters string

func makeEchoCommand(name string) iCommand {
	return tCommand{
		desc: iDescription{
			Command: name,
			Positional: meta.Positional{
				Max: -1,
			},
			Requirements: func(flag meta.Flag) bool { return true },
		},

		beforeRegister: func() error { return nil },
		afterParse:     func() error { return nil },
		run: func(context iContext) error {
			context.Printf("%v\n", context.Args.Pos)
			return nil
		},
	}
}

// makeProgram creates a new program and registers an echo command with it.
func makeProgram() iProgram {
	return iProgram{
		NewEnvironment: iNewEnvironment,
		Info: meta.Info{
			BuildVersion: "42.0.0",
			BuildTime:    time.Unix(0, 0).UTC(),

			Executable:  "exe",
			Description: "something something dark side",
		},
	}
}

// iNewEnvivoment implements new environment for parameters
func iNewEnvironment(params tParameters, context iContext) (tEnvironment, error) {
	return tEnvironment(string(params)), nil
}

// tFlags holds a set of dummy global flags.
type tFlags struct {
	GlobalOne string `short:"a" long:"global-one"`
	GlobalTwo string `short:"b" long:"global-two"`
}

// tRequirements is the implementation of the AllowsFlag function
type tRequirements func(flag meta.Flag) bool

func (t tRequirements) AllowsFlag(flag meta.Flag) bool { return t(flag) }
func (t tRequirements) Validate(args Arguments[tFlags]) error {
	return ValidateAllowedFlags[tFlags](t, args)
}

// instiantiated types for the test suite
type iProgram = Program[tEnvironment, tParameters, tFlags, tRequirements]
type iCommand = Command[tEnvironment, tParameters, tFlags, tRequirements]
type iContext = Context[tEnvironment, tParameters, tFlags, tRequirements]
type iArguments = Arguments[tFlags]
type iDescription = Description[tFlags, tRequirements]

// tCommand represents a sample test suite command.
// It runs the associated private functions, or prints an info message to stdout.
type tCommand struct {
	StdoutMsg string `short:"o" long:"stdout" value-name:"message" default:"write to stdout"`
	StderrMsg string `short:"e" long:"stderr" value-name:"message" default:"write to stderr"`

	beforeRegister func() error
	desc           iDescription
	afterParse     func() error
	run            func(context iContext) error
}

func (t tCommand) BeforeRegister(program *iProgram) {
	if t.beforeRegister == nil {
		fmt.Println("BeforeRegister()")
		return
	}
	t.beforeRegister()
}
func (t tCommand) Description() iDescription {
	return t.desc
}
func (t tCommand) AfterParse() error {
	if t.afterParse == nil {
		fmt.Println("AfterParse()")
		return nil
	}
	return t.afterParse()
}
func (t tCommand) Run(ctx iContext) error {
	if t.run == nil {
		fmt.Println("Run()")
		return nil
	}
	return t.run(ctx)
}

func TestProgram_Main(t *testing.T) {
	root := testutil.TempDirAbs(t)
	if err := os.Mkdir(filepath.Join(root, "real"), os.ModeDir&os.ModePerm); err != nil {
		panic(err)
	}

	wrapLength := 80

	// create buffers for input and output
	var stdoutBuffer bytes.Buffer
	var stderrBuffer bytes.Buffer
	stream := stream.NewIOStream(&stdoutBuffer, &stderrBuffer, nil, wrapLength)

	// define requirements to allow only the Global1 (or any) arguments
	reqOne := tRequirements(func(flag meta.Flag) bool {
		return flag.FieldName == "Global1"
	})

	// define requirements to allow anything
	reqAny := tRequirements(func(flag meta.Flag) bool { return true })

	tests := []struct {
		name       string
		args       []string
		desc       iDescription
		parameters tParameters

		// alias to register (if any)
		alias Alias

		// should the output and error test case data be wrapped automatically
		wrapError bool
		wrapOut   bool

		wantStdout string
		wantStderr string
		wantCode   uint8
	}{

		{
			name:       "no arguments",
			args:       []string{},
			wantStderr: "Unable to parse arguments: Need at least one argument.\n",
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
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--]\nCOMMAND [ARGS...]\n\nsomething something dark side\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\n   COMMAND [ARGS...]\n      Command to call. One of \"fake\". See individual commands for more help.\n",
			wantCode:   0,
		},

		{
			name:       "display help, don't run command",
			args:       []string{"--help", "fake", "whatever"},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--]\nCOMMAND [ARGS...]\n\nsomething something dark side\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\n   COMMAND [ARGS...]\n      Command to call. One of \"fake\". See individual commands for more help.\n",
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
			name:       "command help (1)",
			args:       []string{"fake", "--help"},
			desc:       iDescription{Requirements: reqAny},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--]\nfake [--stdout|-o message] [--stderr|-e message]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   -o, --stdout message\n       (default write to stdout)\n\n   -e, --stderr message\n       (default write to stderr)\n",
			wantCode:   0,
		},

		{
			name:       "command help (2)",
			args:       []string{"--", "fake", "--help"},
			desc:       iDescription{Requirements: reqAny},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--]\nfake [--stdout|-o message] [--stderr|-e message]\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   -o, --stdout message\n       (default write to stdout)\n\n   -e, --stderr message\n       (default write to stderr)\n",
			wantCode:   0,
		},

		{
			name: "alias help (1)",

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			args: []string{"alias", "--help"},
			desc: iDescription{Requirements: reqAny},

			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--]\nalias [--] [ARG ...]\n\nAlias for `exe fake`. See `exe fake --help` for detailed help page about fake.\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake`.\n",
			wantCode:   0,
		},

		{
			name: "alias help (2)",

			alias: Alias{
				Name:    "alias",
				Command: "fake",
			},

			args:       []string{"--", "alias", "--help"},
			desc:       iDescription{Requirements: reqAny},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--]\nalias [--] [ARG ...]\n\nAlias for `exe fake`. See `exe fake --help` for detailed help page about fake.\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake`.\n",
			wantCode:   0,
		},

		{
			name: "alias help (3)",

			alias: Alias{
				Name:        "alias",
				Command:     "fake",
				Args:        []string{"something", "else"},
				Description: "Some useful alias",
			},

			args:       []string{"alias", "--help"},
			desc:       iDescription{Requirements: reqAny},
			wantStdout: "Usage: exe [--help|-h] [--version|-v] [--global-one|-a] [--global-two|-b] [--]\nalias [--] [ARG ...]\n\nSome useful alias\n\nAlias for `exe fake something else`. See `exe fake --help` for detailed help\npage about fake.\n\nGlobal Arguments:\n\n   -h, --help\n      Print a help message and exit\n\n   -v, --version\n      Print a version message and exit\n\n   -a, --global-one\n      \n\n   -b, --global-two\n      \n\nCommand Arguments:\n\n   [ARG ...]\n      Arguments to pass after `exe fake something else`.\n",
			wantCode:   0,
		},

		{
			name:       "not enough arguments for fake",
			args:       []string{"fake"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 1, Max: 2}},
			wantStderr: "Wrong number of positional arguments for fake: Between 1 and 2 argument(s)\nrequired\n",
			wantCode:   4,
		},

		{
			name:       "'fake' with unknown argument (not allowed)",
			args:       []string{"fake", "--argument-not-declared"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1, IncludeUnknown: false}},
			wantStdout: "",
			wantStderr: "Error parsing flags: unknown flag `argument-not-declared'\n",
			wantCode:   4,
		},

		{
			name:       "'fake' with unknown argument (allowed)",
			args:       []string{"fake", "--argument-not-declared"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1, IncludeUnknown: true}},
			wantStdout: "Got Flags: { }\nGot Pos: [--argument-not-declared]\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' without global",
			args:       []string{"fake", "hello", "world"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 1, Max: 2}},
			wantStdout: "Got Flags: { }\nGot Pos: [hello world]\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name:       "'fake' with global (1)",
			args:       []string{"-a", "real", "fake", "hello", "world"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 1, Max: 2}},
			wantStdout: "Got Flags: {real }\nGot Pos: [hello world]\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name:       "'fake' with global (2)",
			args:       []string{"--global-two", "real", "fake", "hello", "world"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 1, Max: 2}},
			wantStdout: "Got Flags: { real}\nGot Pos: [hello world]\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
		{
			name:       "'fake' with disallowed global",
			args:       []string{"--global-one", "notallowed", "fake", "hello", "world"},
			desc:       iDescription{Requirements: reqOne, Positional: meta.Positional{Min: 1, Max: 2}},
			wantStderr: "Wrong number of arguments: 'fake' takes no '--global-one' argument.\n",
			wantCode:   4,
		},

		{
			name:       "'fake' with allowed and disallowed global",
			args:       []string{"--global-one", "one", "--global-two", "two", "fake", "hello", "world"},
			desc:       iDescription{Requirements: reqOne, Positional: meta.Positional{Min: 1, Max: 2}},
			wantStderr: "Wrong number of arguments: 'fake' takes no '--global-one' argument.\n",
			wantCode:   4,
		},

		{
			name:       "'fake' with non-global argument with identical name",
			args:       []string{"--", "fake", "--global-one"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1, IncludeUnknown: true}},
			wantStdout: "Got Flags: { }\nGot Pos: [--global-one]\nwrite to stdout\n",
			wantStderr: "write to stderr\n", //
			wantCode:   0,
		},

		{
			name:       "'fake' with parsed short argument",
			args:       []string{"fake", "-o", "message"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1, IncludeUnknown: true}},
			wantStdout: "Got Flags: { }\nGot Pos: []\nmessage\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with non-parsed short argument",
			args:       []string{"fake", "--", "--s", "message"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1, IncludeUnknown: true}},
			wantStdout: "Got Flags: { }\nGot Pos: [--s message]\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with parsed long argument",
			args:       []string{"fake", "--stdout", "message"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1, IncludeUnknown: true}},
			wantStdout: "Got Flags: { }\nGot Pos: []\nmessage\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with non-parsed long argument",
			args:       []string{"fake", "--", "--stdout", "message"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1, IncludeUnknown: true}},
			wantStdout: "Got Flags: { }\nGot Pos: [--stdout message]\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},

		{
			name:       "'fake' with failure ",
			args:       []string{"fake", "fail"},
			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 1, Max: 2}},
			wantStdout: "Got Flags: { }\nGot Pos: [fail]\nwrite to stdout\n",
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

			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1}},
			wantStdout: "Got Flags: { }\nGot Pos: [hello world]\nwrite to stdout\n",
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

			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1}},
			wantStdout: "Got Flags: { }\nGot Pos: [hello world]\nwrite to stdout\n",
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

			desc:       iDescription{Requirements: reqAny, Positional: meta.Positional{Min: 0, Max: -1}},
			wantStdout: "Got Flags: { }\nGot Pos: [hello world]\nwrite to stdout\n",
			wantStderr: "write to stderr\n",
			wantCode:   0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// reset the buffers
			stdoutBuffer.Reset()
			stderrBuffer.Reset()

			program := makeProgram()

			tt.desc.Command = "fake"
			fake := &tCommand{
				desc:           tt.desc,
				beforeRegister: func() error { return nil },
				afterParse:     func() error { return nil },
			}

			fake.run = func(context iContext) error {
				context.Printf("Got Flags: %s\n", context.Args.Flags)
				context.Printf("Got Pos: %v\n", context.Args.Pos)

				context.Println(fake.StdoutMsg)
				context.EPrintln(fake.StderrMsg)

				// fail when requested
				if len(context.Args.Pos) > 0 && context.Args.Pos[0] == "fail" {
					return exit.Error{ExitCode: exit.ExitGeneric, Message: "Test Failure"}
				}

				return nil
			}
			program.Register(fake)

			if tt.alias.Name != "" {
				program.RegisterAlias(tt.alias)
			}

			// run the program
			ret := exit.AsError(program.Main(stream, tt.parameters, tt.args))

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
