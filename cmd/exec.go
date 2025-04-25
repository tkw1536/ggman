package cmd

//spellchecker:words errors exec essio shellescape github ggman goprogram exit parser pkglib sema status stream
import (
	"errors"
	"fmt"
	"os/exec"

	"al.essio.dev/pkg/shellescape"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/goprogram/exit"
	"github.com/tkw1536/goprogram/parser"
	"github.com/tkw1536/pkglib/sema"
	"github.com/tkw1536/pkglib/status"
	"github.com/tkw1536/pkglib/stream"
)

//spellchecker:words positionals compat nolint wrapcheck

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
		Exe  string   `description:"program to execute"           positional-arg-name:"EXE" required:"1-1"`
		Args []string `description:"arguments to pass to program" positional-arg-name:"ARG"`
	} `positional-args:"true"`

	Parallel int  `default:"1"                                                                                  description:"number of commands to run in parallel, 0 for no limit" long:"parallel" short:"p"`
	Simulate bool `description:"instead of actually running a command, print a bash script that would run them" long:"simulate"                                                     short:"s"`
	NoRepo   bool `description:"do not print name of repos command is being run in"                             long:"no-repo"                                                      short:"n"`
	Quiet    bool `description:"do not provide input or output streams to the command being run"                long:"quiet"                                                        short:"q"`
	Force    bool `description:"continue execution even if an executable returns a non-zero exit code"          long:"force"                                                        short:"f"`
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

var (
	errExecFatal              = exit.NewErrorWithCode("", exit.ExitGeneric)
	errExecParallelNegative   = exit.NewErrorWithCode("argument for `--parallel` must be non-negative", exit.ExitCommandArguments)
	errExecNoParallelSimulate = exit.NewErrorWithCode("`--simulate` expects `--parallel` to be 1", exit.ExitCommandArguments)
)

func (e exe) AfterParse() error {
	if e.Parallel < 0 {
		return errExecParallelNegative
	}
	return nil
}

func (e exe) Run(context ggman.Context) error {
	if e.Simulate {
		return e.runSimulate(context)
	}
	return e.runReal(context)
}

// runReal implements ggman exec for simulate = False.
func (e exe) runReal(context ggman.Context) (err error) {
	repos := context.Environment.Repos(true)

	statusIO := e.Parallel != 1 && !e.Quiet

	var st *status.Status
	if statusIO {
		st = status.NewWithCompat(context.Stdout, 0)
		st.Start()
		defer st.Stop()
	}

	// schedule each command to be run in parallel by using a semaphore!
	err = sema.Schedule(func(i uint64) error {
		repo := repos[i]

		io := context.IOStream
		if statusIO {
			line := st.OpenLine(repo+": ", "")
			defer func() {
				errClose := line.Close()
				if errClose == nil {
					return
				}
				if err == nil {
					err = errClose
				}
			}()
			io = io.Streams(line, line, nil, 0).NonInteractive()
		}

		if !e.NoRepo && !statusIO {
			if _, err := io.EPrintln(repo); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		}

		return e.runRepo(io, repo)
	}, uint64(len(repos)), sema.Concurrency{
		Limit: e.Parallel,
		Force: e.Force,
	})
	if err != nil {
		return fmt.Errorf("process reported error: %w", err)
	}
	return nil
}

func (e exe) runRepo(io stream.IOStream, repo string) error {
	cmd := exec.Command(e.Positionals.Exe, e.Positionals.Args...) /* #nosec G204 -- by design */
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
	var exitError *exec.ExitError
	if errors.As(err, &exitError) {
		return exit.NewErrorWithCode(err.Error(), exit.Code(exitError.ExitCode()))
	}

	return fmt.Errorf("%w%w", errExecFatal, err)
}

// runSimulate runs the --simulate flag.
func (e exe) runSimulate(context ggman.Context) (err error) {
	if e.Parallel != 1 {
		return fmt.Errorf("%w, but got %d", errExecNoParallelSimulate, e.Parallel)
	}

	// print header of the bash script
	if _, err := context.Println("#!/bin/bash"); err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	if !e.Force {
		if _, err := context.Println("set -e"); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}
	if _, err := context.Println(""); err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}

	exec := shellescape.QuoteCommand(append([]string{e.Positionals.Exe}, e.Positionals.Args...))

	// iterate over each repository
	// then print each of the commands to be run!
	for _, repo := range context.Environment.Repos(true) {
		if _, err := context.Printf("cd %s\n", shellescape.Quote(repo)); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if !e.NoRepo {
			if _, err := context.Printf("echo %s\n", shellescape.Quote(repo)); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		}

		if _, err := context.Println(exec); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if _, err := context.Println(""); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}

	return err
}

//spellchecker:words nosec
