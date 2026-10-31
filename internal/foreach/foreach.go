// Package foreach runs an operation across a list of repository paths,
// optionally in parallel.
//
//spellchecker:words foreach
package foreach

//spellchecker:words context sync atomic ggman internal pkglib exit sema status stream
import (
	"context"
	"fmt"
	"sync/atomic"

	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/sema"
	"go.tkw01536.de/pkglib/status"
	"go.tkw01536.de/pkglib/stream"
)

//spellchecker:words compat

// ErrParallelNegative is returned when parallel is negative.
var ErrParallelNegative = exit.NewErrorWithCode(`argument for "--parallel" must be non-negative`, env.ExitCommandArguments)

var errOutput = exit.NewErrorWithCode("unknown output error", env.ExitGeneric)

// Repos runs run once for every repository path in repos.
//
// When parallel is 1, repositories are processed sequentially and each is announced
// as `label "path"` on stdout. When parallel is not 1, repositories are processed
// concurrently (0 means unlimited) using status lines for per-repo output.
// Operation failures are counted and do not stop remaining repositories.
//
// ok and failed count successful and failed repository operations respectively.
// err is a fatal I/O or argument error; repository operation failures are not returned as err.
func Repos(
	ctx context.Context,
	io stream.IOStream,
	repos []string,
	parallel int,
	label string,
	run func(ctx context.Context, io stream.IOStream, repo string) error,
) (ok, failed int, err error) {
	if parallel < 0 {
		return 0, 0, ErrParallelNegative
	}

	if len(repos) == 0 {
		return 0, 0, nil
	}

	if parallel == 1 {
		return reposSequential(ctx, io, repos, label, run)
	}
	return reposParallel(ctx, io, repos, parallel, run)
}

func reposSequential(
	ctx context.Context,
	io stream.IOStream,
	repos []string,
	label string,
	run func(ctx context.Context, io stream.IOStream, repo string) error,
) (ok, failed int, err error) {
	for _, repo := range repos {
		if _, err := fmt.Fprintf(io.Stdout, "%s %q\n", label, repo); err != nil {
			return ok, failed, fmt.Errorf("%w: %w", errOutput, err)
		}

		if e := run(ctx, io, repo); e != nil {
			if _, err := fmt.Fprintln(io.Stderr, e); err != nil {
				return ok, failed, fmt.Errorf("%w: %w", errOutput, err)
			}
			failed++
			continue
		}
		ok++
	}

	return ok, failed, nil
}

func reposParallel(
	ctx context.Context,
	io stream.IOStream,
	repos []string,
	parallel int,
	run func(ctx context.Context, io stream.IOStream, repo string) error,
) (okCount, failedCount int, err error) {
	st := status.NewWithCompat(io.Stdout, 0)
	st.Start()
	defer st.Stop()

	var okAtomic, failedAtomic atomic.Int64

	err = sema.Schedule(func(i uint64) error {
		repo := repos[i]

		line := st.OpenLine(repo+": ", "")
		defer func() {
			errClose := line.Close()
			if errClose == nil {
				return
			}
			if err == nil {
				err = fmt.Errorf("%w: %w", errOutput, errClose)
			}
		}()

		repoIO := io.Streams(line, line, nil, 0).NonInteractive()
		if e := run(ctx, repoIO, repo); e != nil {
			if _, printErr := fmt.Fprintln(repoIO.Stderr, e); printErr != nil {
				return fmt.Errorf("%w: %w", errOutput, printErr)
			}
			failedAtomic.Add(1)
			return nil
		}
		okAtomic.Add(1)
		return nil
	}, uint64(len(repos)), sema.Concurrency{
		Limit: parallel,
		Force: true,
	})
	if err != nil {
		return int(okAtomic.Load()), int(failedAtomic.Load()), err
	}

	return int(okAtomic.Load()), int(failedAtomic.Load()), nil
}
