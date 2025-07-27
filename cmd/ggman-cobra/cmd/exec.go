package cmd

//spellchecker:words errors exec essio shellescape ggman goprogram exit pkglib sema status stream
import (
	"errors"
	"fmt"
	"os/exec"

	"al.essio.dev/pkg/shellescape"
	"github.com/spf13/cobra"
	"go.tkw01536.de/ggman"
	"go.tkw01536.de/ggman/env"
	"go.tkw01536.de/goprogram/exit"
	"go.tkw01536.de/pkglib/sema"
	"go.tkw01536.de/pkglib/status"
	"go.tkw01536.de/pkglib/stream"
)

func NewExecCommand() *cobra.Command {
	impl := new(exe)

	cmd := &cobra.Command{
		Use:   "exec EXE [ARGS...]",
		Short: "execute a command for all repositories",
		Long: `Exec executes an external command for every repository known to ggman.

Each program is run with a working directory set to the root of the provided repository.
Each program is inherits standard input, output and error streams from the ggman process.

Exec prints the path to the repository the command is being run in to standard error.
By default, 'ggman exec' exits with the exit code as soon as the first program that does not return code 0.
If all programs return code 0, 'ggman exec' also exits with code 0.`,
		Args: cobra.MinimumNArgs(1),

		PreRunE: PreRunE(impl),
		RunE:    impl.Exec,
	}

	flags := cmd.Flags()
	flags.IntVarP(&impl.Parallel, "parallel", "p", 1, "number of commands to run in parallel, 0 for no limit")
	flags.BoolVarP(&impl.Simulate, "simulate", "s", false, "instead of actually running a command, print a bash script that would run them")
	flags.BoolVarP(&impl.NoRepo, "no-repo", "n", false, "do not print name of repos command is being run in")
	flags.BoolVarP(&impl.Quiet, "quiet", "q", false, "do not provide input or output streams to the command being run")
	flags.BoolVarP(&impl.Force, "force", "f", false, "continue execution even if an executable returns a non-zero exit code")

	return cmd
}

type exe struct {
	Positionals struct {
		Exe  string
		Args []string
	}

	Parallel int
	Simulate bool
	NoRepo   bool
	Quiet    bool
	Force    bool
}

func (exe) Description() ggman.Description {
	return ggman.Description{
		Command:     "exec",
		Description: "execute a command for all repositories",

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

func (e *exe) AfterParse(cmd *cobra.Command, args []string) error {
	if e.Parallel < 0 {
		return errExecParallelNegative
	}

	e.Positionals.Exe = args[0]
	e.Positionals.Args = args[1:]

	return nil
}

func (e *exe) Exec(cmd *cobra.Command, args []string) error {
	environment, err := ggman.GetEnv(cmd)
	if err != nil {
		return err
	}

	if e.Simulate {
		return e.execSimulate(cmd, environment)
	}
	return e.execReal(cmd, environment)
}

// execReal implements ggman exec for simulate = False.
func (e *exe) execReal(cmd *cobra.Command, environment env.Env) (err error) {

	repos := environment.Repos(true)

	statusIO := e.Parallel != 1 && !e.Quiet

	var st *status.Status
	if statusIO {
		st = status.NewWithCompat(cmd.OutOrStdout(), 0)
		st.Start()
		defer st.Stop()
	}

	// schedule each command to be run in parallel by using a semaphore!
	err = sema.Schedule(func(i uint64) error {
		repo := repos[i]

		io := streamFromCommand(cmd)
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

		return e.execRepo(io, repo)
	}, uint64(len(repos)), sema.Concurrency{
		Limit: e.Parallel,
		Force: e.Force,
	})
	if err != nil {
		return fmt.Errorf("process reported error: %w", err)
	}
	return nil
}

func (e *exe) execRepo(io stream.IOStream, repo string) error {
	exe := exec.Command(e.Positionals.Exe, e.Positionals.Args...) /* #nosec G204 -- by design */
	exe.Dir = repo

	// setup standard output / input, using either the environment
	// or be quiet
	if !e.Quiet {
		exe.Stdin = io.Stdin
		exe.Stdout = io.Stdout
		exe.Stderr = io.Stderr
	}

	// run the actual command, and return if the command was oK!
	err := exe.Run()
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
func (e exe) execSimulate(cmd *cobra.Command, environment env.Env) (err error) {
	if e.Parallel != 1 {
		return fmt.Errorf("%w, but got %d", errExecNoParallelSimulate, e.Parallel)
	}

	// print header of the bash script
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), "#!/bin/bash"); err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}
	if !e.Force {
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), "set -e"); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}
	if _, err := fmt.Fprintln(cmd.OutOrStdout(), ""); err != nil {
		return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
	}

	exec := shellescape.QuoteCommand(append([]string{e.Positionals.Exe}, e.Positionals.Args...))

	// iterate over each repository
	// then print each of the commands to be run!
	for _, repo := range environment.Repos(true) {
		if _, err := fmt.Fprintf(cmd.OutOrStdout(), "cd %s\n", shellescape.Quote(repo)); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if !e.NoRepo {
			if _, err := fmt.Fprintf(cmd.OutOrStdout(), "echo %s\n", shellescape.Quote(repo)); err != nil {
				return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
			}
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), exec); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
		if _, err := fmt.Fprintln(cmd.OutOrStdout(), ""); err != nil {
			return fmt.Errorf("%w: %w", ggman.ErrGenericOutput, err)
		}
	}

	return err
}

//spellchecker:words nosec
