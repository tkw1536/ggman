package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/cmd"
	"github.com/tkw1536/ggman/program"
)

// subcommands is a list of all supported subcommands.
var subcommands = []program.Command{
	cmd.Root,

	cmd.Ls,
	cmd.Lsr,

	cmd.Where,
	cmd.Canon,
	cmd.Comps,

	cmd.Fetch,
	cmd.Pull,

	cmd.Fix,

	cmd.Clone,
	cmd.Link,

	cmd.License,

	cmd.Here,

	cmd.Web,
	cmd.URL,

	cmd.FindBranch,
}

func main() {

	// recover from calls to panic(), and exit the program appropriatly.
	// This has to be in the main() function because any of the libary functions might be broken.
	// For this reason, no ggman functions are used here; just stuff from the main package.
	defer func() {
		if err := recover(); err != nil {
			fmt.Fprintf(os.Stderr, fatalPanicMessage, err)
			debug.PrintStack()
			ggman.ExitPanic.Return()
		}
	}()

	// Create the 'ggman' program and register all the subcommands
	// Then execute the program and handle the exit code.
	cmd := &program.Program{IOStream: ggman.NewEnvIOStream()}
	for _, c := range subcommands {
		cmd.Register(c)
	}

	err := ggman.AsError(cmd.Main(os.Args[1:]))
	err.Return()
}

const fatalPanicMessage = `Fatal Error: Panic

The ggman program panicked and had to abort execution. This is usually
indicative of a bug. If this occurs repeatedly you might want to consider
filing an issue in the issue tracker at
https://github.com/tkw1536/ggman/issues. Below is debug information that might
help the developer to track down what happened. 

panic: %v
`
