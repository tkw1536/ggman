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

// init initializes this context by setting up the environment according to the arguments
func (c *Context) init() error {

	// generate pattern filters for the "--for" arguments
	clauses := make([]env.Filter, len(c.Filters))
	for i, pat := range c.Filters {

		// check if 'pat' represents the root of a repository
		if repo, err := c.AtRoot(pat); err == nil && repo != "" {
			clauses[i] = env.PathFilter{Paths: []string{repo}}
			continue
		}

		// create a normal pattern filter
		clauses[i] = env.NewPatternFilter(pat, !c.NoFuzzyFilter) // TODO: Make fuzzyness optional
	}

	// generate a 'here' filter for the current repository
	if c.Here {
		repo, _, err := c.At(".")
		if err != nil {
			return errors.Wrap(err, "Unable to find current repository")
		}

		clauses = append(clauses, env.PathFilter{Paths: []string{repo}})
	}

	// only set the filter when we actually have something to filter by
	if len(clauses) != 0 {
		c.Filter = env.DisjunctionFilter{Clauses: clauses}
	}
	return nil
}

// URLV returns the ith argument, parsed as a URL.
//
// It is a convenience wrapper for:
//  c.ParseURL(c.Args[i])
// This function is untested.
func (c Context) URLV(i int) env.URL {
	return env.ParseURL(c.Args[i])
}
