package program

import (
	"github.com/tkw1536/ggman/program/stream"
)

// Context represents a context that a command is run inside of
type Context[Runtime any, Requirements any] struct {
	stream.IOStream
	CommandArguments[Runtime, Requirements]
	runtime Runtime
}

// Runtime returns the runtime belonging to this context
func (c Context[Runtime, Requirements]) Runtime() Runtime {
	return c.runtime
}
