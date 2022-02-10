package program

import (
	"github.com/tkw1536/ggman/program/stream"
)

// Context represents a context that a command is run inside of
type Context[Runtime any, Parameters any, Requirements any] struct {
	stream.IOStream
	CommandArguments[Runtime, Parameters, Requirements]
	runtime Runtime
}

// Runtime returns the runtime belonging to this context
func (c Context[Runtime, Parameters, Requirements]) Runtime() Runtime {
	return c.runtime
}
