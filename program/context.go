package program

import (
	"github.com/tkw1536/ggman/program/stream"
)

// Context represents a context that a command is run inside of
type Context[Runtime any, Parameters any, Flags any, Requirements Requirement[Flags]] struct {
	stream.IOStream
	Args CommandArguments[Runtime, Parameters, Flags, Requirements]

	runtime Runtime
}

// Runtime returns the runtime belonging to this context
func (c Context[Runtime, Parameters, Flags, Requirements]) Runtime() Runtime {
	return c.runtime
}
