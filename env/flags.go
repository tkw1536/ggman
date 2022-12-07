package env

import (
	"bufio"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/tkw1536/ggman/internal/walker"
	"github.com/tkw1536/goprogram/exit"
)

// Flags represents a set of filter flags used for the ggman goprogram.
type Flags struct {
	For           []string `short:"f" long:"for" value-name:"FILTER" description:"filter list of repositories by FILTER. FILTER can be a relative or absolute path, or a glob pattern which will be matched against the normalized repository url"`
	FromFile      []string `short:"i" long:"from-file" value-name:"FILE" description:"filter list of repositories to only those matching filters in FILE. FILE should contain one filter per line, with common comment chars being ignored"`
	NoFuzzyFilter bool     `short:"n" long:"no-fuzzy-filter" description:"disable fuzzy matching for filters"`

	Here bool     `short:"H" long:"here" description:"filter list of repositories to only contain those that are in the current directory or subtree. alias for \"-p .\""`
	Path []string `short:"P" long:"path" description:"filter list of repositories to only contain those that are in or under the specified path. may be used multiple times"`

	Dirty bool `short:"d" long:"dirty" description:"filter list of repositories to only contain repositories with uncommited changes"`
	Clean bool `short:"c" long:"clean" description:"filter list of repositories to only contain repositories without uncommited changes"`

	Synced   bool `short:"s" long:"synced" description:"filter list of repositories to only contain those which are up-to-date with remote"`
	UnSynced bool `short:"u" long:"unsynced" description:"filter list of repositories to only contain those not up-to-date with remote"`

	Tarnished bool `short:"t" long:"tarnished" description:"filter list of repositories to only contain those that are dirty or unsynced"`
	Pristine  bool `short:"p" long:"pristine" description:"filter list of repositories to only contain those that are clean and synced"`
}

var errNotADirectory = exit.Error{
	ExitCode: ExitInvalidRepo,
	Message:  "not a directory: %q",
}

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

	// add a WorktreeFilter filter if requested
	if flags.Dirty || flags.Clean {
		filter = WorktreeFilter{
			Filter: filter,

			Dirty: flags.Dirty,
			Clean: flags.Clean,
		}
	}

	if flags.Synced || flags.UnSynced {
		filter = StatusFilter{
			Filter: filter,

			Synced:   flags.Synced,
			UnSynced: flags.UnSynced,
		}
	}

	if flags.Tarnished || flags.Pristine {
		filter = TarnishFilter{
			Filter: filter,

			Tarnished: flags.Tarnished,
			Pristine:  flags.Pristine,
		}
	}

	return
}

// NewForFilter creates a new 'for' filter for this environment.
//
// A 'for' filter may be either:
//   - a (relative or absolute) path to the root of a repository (see env.AtRoot)
//   - a repository url or pattern (see NewPatternFilter)
func (env Env) NewForFilter(filter string, fuzzy bool) Filter {
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
func (env Env) NewFromFileFilter(p string, fuzzy bool) (filters []Filter, err error) {
	// resolve the path
	path, err := env.Abs(p)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to resolve path %q", p)
	}

	// open the file
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to open path %q", p)
	}
	defer file.Close()

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
		return nil, errors.Wrapf(err, "Unable to read file %q", p)
	}

	return filters, nil
}

// ResolvePathFilter resolves and validates p for use within a PathFilter.
//
// p must be an existing absolute absolute or relative path pointing to:
//   - a repository directory (see env.At)
//   - a (possibly nested) directory containing repositories
func (env Env) ResolvePathFilter(p string) (path string, err error) {
	// a repository directly
	path, _, err = env.At(p)
	if err == nil {
		return
	}

	// sub-repositories contained in a path
	path, err = env.Abs(p)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to resolve path %q", p)
	}

	// must be a directory!
	if ok, err := walker.IsDirectory(path, true); err != nil || !ok {
		return "", errNotADirectory.WithMessageF(p)
	}

	return
}
