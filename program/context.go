package program

import (
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

// Context represents a context that a command is run inside of
type Context struct {
	ggman.IOStream
	env.Env
	CommandArguments
}

// URLV returns the ith argument, parsed as a URL
func (c Context) URLV(i int) env.URL {
	return c.ParseURL(c.Argv[i])
}
