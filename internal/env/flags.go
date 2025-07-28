package env

//spellchecker:words bufio strings pkglib exit
import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go.tkw01536.de/pkglib/exit"
	"go.tkw01536.de/pkglib/fsx"
)

//spellchecker:words ggman unsynced worktree

// Flags represents a set of filter flags used for the ggman program.
type Flags struct {
	For           []string
	FromFile      []string
	NoFuzzyFilter bool

	Here bool
	Path []string

	Dirty bool
	Clean bool

	Synced   bool
	UnSynced bool

	Tarnished bool
	Pristine  bool
}

var errNotADirectory = exit.NewErrorWithCode("not a directory", ExitInvalidRepo)

// NewFilter creates a new filter corresponding to the given Flags and Environment.
func NewFilter(flags Flags, env *Env) (filter Filter, err error) {
	// generate pattern filters for the "--for" arguments
	clauses := make([]Filter, len(flags.For))
	for i, pat := range flags.For {
		clauses[i] = env.NewForFilter(pat, !flags.NoFuzzyFilter)
	}

	// read filters from file
	for _, p := range flags.FromFile {
		filters, err := env.NewFromFileFilter(p, !flags.NoFuzzyFilter)
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, filters...)
	}

	// here filter: alias for --path .
	if flags.Here {
		flags.Path = append(flags.Path, ".")
	}

	// for each of the candidate paths, add a path filter
	pf := PathFilter{Paths: make([]string, len(flags.Path))}
	for i, p := range flags.Path {
		pf.Paths[i], err = env.ResolvePathFilter(p)
		if err != nil {
			return nil, err
		}
	}

	if len(pf.Paths) > 0 {
		clauses = append(clauses, pf)
	}

	// only set the filter when we actually have something to filter by
	filter = DisjunctionFilter{Clauses: clauses}
	if len(clauses) == 0 {
		filter = NoFilter
	}

	// setup some additional filters
	if flags.Dirty || flags.Clean {
		filter = NewWorktreeFilter(filter, flags.Dirty, flags.Clean)
	}
	if flags.Synced || flags.UnSynced {
		filter = NewStatusFilter(filter, flags.Synced, flags.UnSynced)
	}
	if flags.Tarnished || flags.Pristine {
		filter = NewTarnishFilter(filter, flags.Tarnished, flags.Pristine)
	}

	return
}

// NewForFilter creates a new 'for' filter for this environment.
//
// A 'for' filter may be either:
//   - a (relative or absolute) path to the root of a repository (see env.AtRoot)
//   - a repository url or pattern (see NewPatternFilter)
func (env *Env) NewForFilter(filter string, fuzzy bool) Filter {
	// check if 'pat' represents the root of a repository
	if repo, err := env.AtRoot(filter); err == nil && repo != "" {
		return PathFilter{Paths: []string{repo}}
	}

	// create a normal pattern filter
	return NewPatternFilter(filter, fuzzy)
}

// NewFromFileFilter creates a list of filters from the file at path.
//
// To create a filter, p is opened and each (whitespace-trimmed) line is passed to env.NewForFilter.
// Blank lines, or those starting with ';', '//' or '#' are ignored.
func (env *Env) NewFromFileFilter(p string, fuzzy bool) (filters []Filter, err error) {
	// resolve the path
	path, err := env.Abs(p)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve path %q: %w", p, err)
	}

	// open the file
	file, oErr := os.Open(path) /* #nosec G304 -- explicitly passed as a parameter */
	if oErr != nil {
		return nil, fmt.Errorf("unable to open path %q: %w", p, oErr)
	}
	defer func() {
		eClose := file.Close()
		if eClose == nil {
			return
		}
		eClose = fmt.Errorf("unable to close path %q: %w", p, eClose)

		if err == nil {
			err = eClose
		}
	}()

	// read each line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// ignore blank or comment lines
		if line == "" || line[0] == ';' || line[0] == '#' || strings.HasPrefix(line, "//") {
			continue
		}
		filters = append(filters, env.NewForFilter(line, fuzzy))
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("unable to read file %q: %w", p, err)
	}

	return filters, nil
}

// ResolvePathFilter resolves and validates p for use within a PathFilter.
//
// p must be an existing absolute absolute or relative path pointing to:
//   - a repository directory (see env.At)
//   - a (possibly nested) directory containing repositories
func (env *Env) ResolvePathFilter(p string) (path string, err error) {
	// a repository directly
	path, _, err = env.At(p)
	if err == nil {
		return
	}

	// sub-repositories contained in a path
	path, err = env.Abs(p)
	if err != nil {
		return "", fmt.Errorf("unable to resolve path %q: %w", p, err)
	}

	// must be a directory!
	if ok, err := fsx.IsDirectory(path, true); err != nil || !ok {
		return "", fmt.Errorf("%q %w", p, errNotADirectory)
	}

	return
}

//spellchecker:words nosec
