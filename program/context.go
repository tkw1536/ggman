package program

import (
	"github.com/pkg/errors"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/walker"
)

// Context represents a context that a command is run inside of
type Context struct {
	ggman.IOStream
	Env env.Env
	CommandArguments
}

// init initializes this context by setting up the environment according to the arguments
func (c *Context) init() (err error) {
	c.Env.Filter, err = c.makeFilter()
	return
}

var errNotADirectory = ggman.Error{
	ExitCode: ggman.ExitInvalidRepo,
	Message:  "Not a directory: %q",
}

// makeFilter creates a new filter for this context.
// It should only be used during initialization.
func (c Context) makeFilter() (env.Filter, error) {
	// generate pattern filters for the "--for" arguments
	clauses := make([]env.Filter, len(c.Filters))
	for i, pat := range c.Filters {

		// check if 'pat' represents the root of a repository
		if repo, err := c.Env.AtRoot(pat); err == nil && repo != "" {
			clauses[i] = env.PathFilter{Paths: []string{repo}}
			continue
		}

		// create a normal pattern filter
		clauses[i] = env.NewPatternFilter(pat, !c.NoFuzzyFilter)
	}

	// here filter: alias for --path .
	if c.Here {
		c.Path = append(c.Path, ".")
	}

	// for each of the candidate paths, add a path filter
	pf := env.PathFilter{Paths: make([]string, len(c.Path))}
	for i, p := range c.Path {
		var err error
		pf.Paths[i], _, err = c.Env.At(p) // try to use the current repository first.
		if err != nil {
			// filter sub-repositories under this repo!
			pf.Paths[i], err = c.Env.Abs(p)
			if err != nil {
				return nil, errors.Wrapf(err, "Unable to resolve path: %q", p)
			}

			// make sure it is actually a directory!
			if ok, err := walker.IsDirectory(pf.Paths[i], true); err != nil || !ok {
				return nil, errNotADirectory.WithMessageF(p)
			}
		}
	}

	if len(pf.Paths) > 0 {
		clauses = append(clauses, pf)
	}

	// only set the filter when we actually have something to filter by
	var dj env.Filter = env.DisjunctionFilter{Clauses: clauses}
	if len(clauses) == 0 {
		dj = env.NoFilter
	}

	// add a WorktreeFilter filter if requested
	if c.Dirty || c.Clean {
		dj = env.WorktreeFilter{
			Filter: dj,

			Dirty: c.Dirty,
			Clean: c.Clean,
		}
	}

	if c.Synced || c.UnSynced {
		dj = env.StatusFilter{
			Filter: dj,

			Synced:   c.Synced,
			UnSynced: c.UnSynced,
		}
	}

	if c.Tarnished || c.Pristine {
		dj = env.TarnishFilter{
			Filter: dj,

			Tarnished: c.Tarnished,
			Pristine:  c.Pristine,
		}
	}

	return dj, nil
}

// URLV returns the ith argument, parsed as a URL.
//
// It is a convenience wrapper for:
//  c.ParseURL(c.Args[i])
// This function is untested.
func (c Context) URLV(i int) env.URL {
	return env.ParseURL(c.Args[i])
}
