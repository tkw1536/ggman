package program

import (
	"github.com/jessevdk/go-flags"
	"github.com/tkw1536/ggman/program/stream"
)

// Context represents an execution environment for a command.
// it takes the same type parameters as a command and program.
type Context[E any, P any, F any, R Requirement[F]] struct {
	// IOStream describes the input and output the command reads from and writes to.
	stream.IOStream

	// Args holds arguments passed to this command.
	Args Arguments[F]

	// Description is the description of the command being invoked
	Description Description[F, R]

	// Environment holds the environment for this command.
	Environment E

	// commandParser holds a parser that parsed argument passed to the command itself
	commandParser *flags.Parser
}

// Arguments represent a set of command-independent arguments passed to a command.
// These should be further parsed into CommandArguments using the appropriate Parse() method.
//
// Command line argument are annotated using syntax provided by "github.com/jessevdk/go-flags".
type Arguments[F any] struct {
	Universals Universals `group:"universals"`
	Flags      F          `group:"flags"`

	Command string   // command to run
	Pos     []string // positional arguments
}

// Universals holds flags added to every executable.
//
// Command line arguments are annotated using syntax provided by "github.com/jessevdk/go-flags".
type Universals struct {
	Help    bool `short:"h" long:"help" description:"Print a help message and exit"`
	Version bool `short:"v" long:"version" description:"Print a version message and exit"`
}
