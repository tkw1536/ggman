package cmd

//spellchecker:words exec github alessio shellescape ggman goprogram exit parser pkglib sema status stream
import (
	"os/exec"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
	"github.com/tkw1536/pkglib/sema"
	"github.com/tkw1536/pkglib/status"
	"github.com/tkw1536/pkglib/stream"
)

//spellchecker:words positionals compat

// Exec is the 'ggman exec' command.
//
// Exec executes an external command for every repository known to ggman.
//
// Each program is run with a working directory set to the root of the provided repository.
// Each program is inherits standard input, output and error streams from the ggman process.
//
// Exec prints the path to the repository the command is being run in to standard error.
// By default, 'ggman exec' exits with the exit code as soon as the first program that does not return code 0.
// If all programs return code 0, 'ggman exec' also exits with code 0.
//
//	--simulate
//
// Instead of actually running a command, print a bash script that would run them.
//
//	--parallel
//
// Number of commands to run in parallel, 0 for no limit.
// Output will be shown on different status lines, except when parallel == 1.
//
//	--no-repo
//
// Do not print name of repos command is being run in.
//
//	--quiet
//
// Do not provide input or output streams to the command being run.
//
//	--force
//
// Continue execution of programs, even if one returns a non-zero exit code.
// 'exec' will still return code 0 as the final exit code.
var Exec ggman.Command = exe{}

type exe struct {
	Positionals struct {
		Exe  string   `positional-arg-name:"EXE" required:"1-1" description:"program to execute"`
		Args []string `positional-arg-name:"ARG"  description:"arguments to pass to program"`
	} `positional-args:"true"`

	Parallel int  `short:"p" long:"parallel" default:"1" description:"number of commands to run in parallel, 0 for no limit"`
	Simulate bool `short:"s" long:"simulate" description:"instead of actually running a command, print a bash script that would run them"`
	NoRepo   bool `short:"n" long:"no-repo" description:"do not print name of repos command is being run in"`
	Quiet    bool `short:"q" long:"quiet" description:"do not provide input or output streams to the command being run"`
	Force    bool `short:"f" long:"force" description:"continue execution even if an executable returns a non-zero exit code"`
}

func (exe) Description() ggman.Description {
	return ggman.Description{
		Command:     "exec",
		Description: "execute a command for all repositories",
		ParserConfig: parser.Config{
			IncludeUnknown: true,
		},

		Requirements: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

var ErrExecParallelNegative = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "argument for `--parallel` must be non-negative",
}

func (e exe) AfterParse() error {
	if e.Parallel < 0 {
		return ErrExecParallelNegative
	}
	return nil
}

func (e exe) Run(context ggman.Context) error {
	if e.Simulate {
		return e.runSimulate(context)
	}
	return e.runReal(context)
}

// runReal implements ggman exec for simulate = False
func (e exe) runReal(context ggman.Context) error {
	repos := context.Environment.Repos(true)

	statusIO := e.Parallel != 1 && !e.Quiet

	var st *status.Status
	if statusIO {
		st = status.NewWithCompat(context.IOStream.Stdout, 0)
		st.Start()
		defer st.Stop()
	}

	// schedule each command to be run in parallel by using a semaphore!
	return sema.Schedule(func(i int) error {
		repo := repos[i]

		io := context.IOStream
		if statusIO {
			line := st.OpenLine(repo+": ", "")
			defer line.Close()
			io = io.Streams(line, line, nil, 0).NonInteractive()
		}

		if !e.NoRepo && !statusIO {
			io.EPrintln(repo)
		}

		return e.runRepo(io, repo)
	}, len(repos), sema.Concurrency{
		Limit: e.Parallel,
		Force: e.Force,
	})
}

var ErrExecFatal = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (e exe) runRepo(io stream.IOStream, repo string) error {
	cmd := exec.Command(e.Positionals.Exe, e.Positionals.Args...)
	cmd.Dir = repo

	// setup standard output / input, using either the environment
	// or be quiet
	if !e.Quiet {
		cmd.Stdin = io.Stdin
		cmd.Stdout = io.Stdout
		cmd.Stderr = io.Stderr
	} else {
		cmd.Stdin = stream.Null
		cmd.Stdout = stream.Null
		cmd.Stderr = stream.Null
	}

	// run the actual command, and return if the command was oK!
	err := cmd.Run()
	if err == nil {
		return nil
	}

	// when something went wrong intercept ExitErrors
	// but actually return other error properly!
	if ee, ok := err.(*exec.ExitError); ok {
		return exit.Error{
			ExitCode: exit.ExitCode(ee.ExitCode()), // #nosec G115 exit status guaranteed to fit into uint8
			Message:  ee.Error(),
		}
	}

	return ErrExecFatal.WithMessage(err.Error())
}

var ErrExecNoParallelSimulate = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "`--simulate` expects `--parallel` to be 1, but got %d",
}

// runSimulate runs the --simulate flag
func (e exe) runSimulate(context ggman.Context) (err error) {
	if e.Parallel != 1 {
		return ErrExecNoParallelSimulate.WithMessageF(e.Parallel)
	}

	// print header of the bash script
	context.Println("#!/bin/bash")
	if !e.Force {
		context.Println("set -e")
	}
	context.Println("")

	exec := shellescape.QuoteCommand(append([]string{e.Positionals.Exe}, e.Positionals.Args...))

	// iterate over each repository
	// then print each of the commands to be run!
	for _, repo := range context.Environment.Repos(true) {
		context.Printf("cd %s\n", shellescape.Quote(repo))
		if !e.NoRepo {
			context.Printf("echo %s\n", shellescape.Quote(repo))
		}

		context.Println(exec)
		context.Println("")
	}

	return err
}

// spellchecker:words nosec
