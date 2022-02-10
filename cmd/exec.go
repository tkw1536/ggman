package cmd

import (
	"os/exec"

	"github.com/alessio/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/sema"
	"github.com/tkw1536/ggman/internal/stream"
	"github.com/tkw1536/ggman/program"
	"github.com/tkw1536/ggman/program/exit"
)

// Exec is the 'ggman exec' command.
//
// Exec excecutes an external command for every repository known to ggman.
//
// Each program is run with a working directory set to the root of the provided repository.
// Each program is inherits standard input, output and error streams from the ggman process.
//
// Exec prints the path to the repository the command is being run in to standard error.
// By default, 'ggman exec' exits with the exit code as soon as the first program that does not return code 0.
// If all programs return code 0, 'ggman exec' also exits with code 0.
//
//   --simulate
// Instead of actually running a command, print a bash script that would run them.
//   --parallel
// Number of commands to run in parallel, 0 for no limit
//   --no-repo
// Do not print name of repos command is being run in.
//   --quiet
// Do not provide input or output streams to the command being run.
//   --force
// Continue execution of programs, even if one returns a non-zero exit code.
// 'exec' will still return code 0 as the final exit code.
var Exec program.Command = &exe{}

type exe struct {
	Parallel int  `short:"p" long:"parallel" default:"1" description:"number of commands to run in parallel, 0 for no limit"`
	Simulate bool `short:"s" long:"simulate" description:"Instead of actually running a command, print a bash script that would run them."`
	NoRepo   bool `short:"n" long:"no-repo" description:"Do not print name of repos command is being run in"`
	Quiet    bool `short:"q" long:"quiet" description:"Do not provide input or output streams to the command being run"`
	Force    bool `short:"f" long:"force" description:"Continue execution even if an executable returns a non-zero exit code"`
}

func (*exe) BeforeRegister(program *program.Program) {}

func (*exe) Description() program.Description {
	return program.Description{
		Name:        "exec",
		Description: "Execute a command for all repositories",

		SkipUnknownOptions: true,
		PosArgsMin:         1,
		PosArgsMax:         -1,
		PosArgName:         "ARGS",

		Environment: env.Requirement{
			AllowsFilter: true,
			NeedsRoot:    true,
		},
	}
}

var ErrExecParalllelNegative = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "argument for --parallel must be non-negative",
}

func (e *exe) AfterParse() error {
	if e.Parallel < 0 {
		return ErrExecParalllelNegative
	}
	return nil
}

func (e *exe) Run(context program.Context) error {
	if e.Simulate {
		return e.runSimulate(context)
	}
	return e.runReal(context)
}

// runReal implements ggman exec for simulate = False
func (e *exe) runReal(context program.Context) error {
	repos := ggman.C2E(context).Repos()

	// schedule each command to be run in parallel by using a semaphore!
	return sema.Schedule(func(i int) error {
		repo := repos[i]

		if !e.NoRepo {
			context.EPrintln(repo)
		}

		return e.runRepo(context, repo)
	}, len(repos), sema.Concurrency{
		Limit: e.Parallel,
		Force: e.Force,
	})
}

var ErrExecFatal = exit.Error{
	ExitCode: exit.ExitGeneric,
}

func (e *exe) runRepo(context program.Context, repo string) error {
	cmd := exec.Command(context.Args[0], context.Args[1:]...)
	cmd.Dir = repo

	// setup standard output / input, using either the environment
	// or be quiet
	if !e.Quiet {
		cmd.Stdin = context.Stdin
		cmd.Stdout = context.Stdout
		cmd.Stderr = context.Stderr
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
			ExitCode: exit.ExitCode(ee.ExitCode()),
			Message:  ee.Error(),
		}
	}

	return ErrExecFatal.WithMessage(err.Error())
}

var ErrExecNoParallelSimulate = exit.Error{
	ExitCode: exit.ExitCommandArguments,
	Message:  "--simulate expects --parallel to be 1, but got %d",
}

// runSimulate runs the --simulate flag
func (e *exe) runSimulate(context program.Context) (err error) {
	if e.Parallel != 1 {
		return ErrExecNoParallelSimulate.WithMessageF(e.Parallel)
	}

	// print header of the bash script
	context.Println("#!/bin/bash")
	if !e.Force {
		context.Println("set -e")
	}
	context.Println("")

	exec := shellescape.QuoteCommand(context.Args)

	// iterate over each repository
	// then print each of the commands to be run!
	repos := ggman.C2E(context).Repos()
	for _, repo := range repos {
		context.Printf("cd %s\n", shellescape.Quote(repo))
		if !e.NoRepo {
			context.Printf("echo %s\n", shellescape.Quote(repo))
		}

		context.Println(exec)
		context.Println("")
	}

	return err
}
