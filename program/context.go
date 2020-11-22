package program

import (
	"github.com/pkg/errors"

	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
)

// Context represents a context that a command is run inside of
type Context struct {
	ggman.IOStream
	env.Env
	CommandArguments
}

// init initializes this context and sets up additional variables
func (c *Context) init() error {
	if !c.filterHere {
		return nil
	}

	if !c.options.Environment.AllowsFilter {
		return errors.New("--here provided, but not allowed")
	}

	// find the current repository
	repo, _, err := c.At(".")
	if err != nil {
		return errors.Wrap(err, "Unable to find current repository")
	}

	// get the remote
	url, err := c.Git.GetRemote(repo)
	if err != nil {
		return errors.Wrap(err, "Unable to identify current repository")
	}
	if url == "" {
		return errors.New("Current repository does not have a url and can not be filtered")
	}

	// and set it
	return c.Env.Filter.Set(url)
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
