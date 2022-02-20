package ggman

import (
	"github.com/pkg/errors"
	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/internal/walker"
	"github.com/tkw1536/ggman/program/exit"
)

// NewRuntime makes a new runtime for ggman
func NewRuntime(params env.Parameters, context Context) (*env.Env, error) {
	// create a new environment
	env, err := env.NewEnv(context.Description.Requirements, params)
	if err != nil {
		return nil, err
	}

	// make a filter
	f, err := makeFilter(env, context.Args)
	if err != nil {
		return nil, err
	}
	env.Filter = f

	return env, nil

}

var errNotADirectory = exit.Error{
	ExitCode: env.ExitInvalidRepo,
	Message:  "Not a directory: %q",
}

func makeFilter(e *env.Env, c Arguments) (env.Filter, error) {
	// generate pattern filters for the "--for" arguments
	clauses := make([]env.Filter, len(c.Flags.Filters))
	for i, pat := range c.Flags.Filters {

		// check if 'pat' represents the root of a repository
		if repo, err := e.AtRoot(pat); err == nil && repo != "" {
			clauses[i] = env.PathFilter{Paths: []string{repo}}
			continue
		}

		// create a normal pattern filter
		clauses[i] = env.NewPatternFilter(pat, !c.Flags.NoFuzzyFilter)
	}

	// here filter: alias for --path .
	if c.Flags.Here {
		c.Flags.Path = append(c.Flags.Path, ".")
	}

	// for each of the candidate paths, add a path filter
	pf := env.PathFilter{Paths: make([]string, len(c.Flags.Path))}
	for i, p := range c.Flags.Path {
		var err error
		pf.Paths[i], _, err = e.At(p) // try to use the current repository first.
		if err != nil {
			// filter sub-repositories under this repo!
			pf.Paths[i], err = e.Abs(p)
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
	if c.Flags.Dirty || c.Flags.Clean {
		dj = env.WorktreeFilter{
			Filter: dj,

			Dirty: c.Flags.Dirty,
			Clean: c.Flags.Clean,
		}
	}

	if c.Flags.Synced || c.Flags.UnSynced {
		dj = env.StatusFilter{
			Filter: dj,

			Synced:   c.Flags.Synced,
			UnSynced: c.Flags.UnSynced,
		}
	}

	if c.Flags.Tarnished || c.Flags.Pristine {
		dj = env.TarnishFilter{
			Filter: dj,

			Tarnished: c.Flags.Tarnished,
			Pristine:  c.Flags.Pristine,
		}
	}

	return dj, nil
}
