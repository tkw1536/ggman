//spellchecker:words main
package main

//spellchecker:words context runtime debug ggman pkglib exit
import (
	"context"
	"fmt"
	"os"
	"runtime/debug"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/env"
	"go.tkw01536.de/pkglib/exit"
)

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
			exit.ExitPanic.Return()
		}
	}()

	// build the parameters
	params := env.Parameters{
		Variables: env.ReadVariables(),
		Plumbing:  nil,
		Workdir:   "",
	}

	// and run the command
	cmd := cmd.NewCommand(context.Background(), params)
	if err := cmd.Execute(); err != nil {
		code, _ := exit.CodeFromError(err)
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), err)
		code.Return()
	}
}
