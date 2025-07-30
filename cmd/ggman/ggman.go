//spellchecker:words main
package main

//spellchecker:words context signal runtime debug syscall ggman internal pkglib exit
import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

//spellchecker:words workdir

const fatalPanicMessage = `Fatal Error: Panic

The ggman program panicked and had to abort execution. This is usually
indicative of a bug. If this occurs repeatedly you might want to consider
filing an issue in the issue tracker at:

https://github.com/tkw1536/ggman/issues

Below is debug information that might help the developer track down what
happened.

panic: %v
`

func main() {
	// recover from calls to panic(), and exit the program appropriately.
	// This has to be in the main() function because any of the library functions might be broken.
	// For this reason, as few ggman functions as possible are used here; just stuff from the top-level ggman package.
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, fatalPanicMessage, err)
			debug.PrintStack()
			env.ExitPanic.Return()
		}
	}()

	// build the parameters
	params := env.Parameters{
		Variables: env.ReadVariables(),
		Plumbing:  nil,
		Workdir:   "",
	}

	// create a context that can be cancelled with Ctrl+C
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// and run the command
	cmd := cmd.NewCommand(ctx, params)
	if err := cmd.Execute(); err != nil {
		code, _ := exit.CodeFromError(err, env.ExitGeneric)
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
		code.Return()
	}
}
