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

// URLV returns the ith argument, parsed as a URL.
//
// It is a convenience wrapper for:
//  c.ParseURL(c.Args[i])
// This function is untested.
func (c Context) URLV(i int) env.URL {
	// TODO: Consider making this work similar to env.At(), except returning the url.
	// i.e.: If this resolves to an existing path, and that path contains a repository, then it should return the url.
	return env.ParseURL(c.Args[i])
}
