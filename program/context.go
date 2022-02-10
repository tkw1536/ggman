package program

import (
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/program/stream"
)

// Context represents a context that a command is run inside of
type Context struct {
	stream.IOStream
	CommandArguments
	runtime Runtime
}

// Runtime represents the runtime of this command
// TODO: type parameter
type Runtime interface{}

// Runtime returns the runtime belonging to this context
// TODO: type parameter
func (c *Context) Runtime() Runtime {
	return c.runtime
}

// URLV returns the ith argument, parsed as a URL.
// TODO: figure out how to make this a type parameter
//
// It is a convenience wrapper for:
//  c.ParseURL(c.Args[i])
// This function is untested.
func (c Context) URLV(i int) env.URL {
	return env.ParseURL(c.Args[i])
}
