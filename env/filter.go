package env

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/tkw1536/ggman/internal/pattern"
	"github.com/tkw1536/ggman/internal/text"
)

// Filter is a predicate that matches repositories inside an environment.
//
// A filter is applied by recursively scanning the root folder for git repositories.
// Each folder that is a repository will be passed to clonePath.
//
// Filter may also optionally implement FilterWithCandidates.
type Filter interface {
	// Matches checks if a repository at clonePath matches this filter.
	Matches(env Env, clonePath string) bool
}

// NoFilter is a special filter that matches every directory
var NoFilter Filter = emptyFilter{}

type emptyFilter struct{}

func (emptyFilter) Matches(env Env, clonePath string) bool {
	return true
}

// FilterWithCandidates is a filter that in addition to being applied normally should also be applied to the provided candidates.
type FilterWithCandidates interface {
	Filter

	// Candidates returns a list of folders that should be added regardless of their location.
	// Paths in the return value may be assumed to exist, but may not be repositories.
	// A FilterWithCandidates with a Candidates() function that returns a zero-length slice is equivalent to a regular filter.
	Candidates() []string
}

// Candidates checks if Filter implements FilterWithCandidates and calls the Candidates() method when applicable.
// When Filter does not implement FilterWithCandidates, returns nil
func Candidates(f Filter) []string {
	cFilter, isCFilter := f.(FilterWithCandidates)
	if !isCFilter {
		return nil
	}

	return cFilter.Candidates()
}

// PathFilter is a filter that always matches the provided paths.
// It implements FilterWithCandidates.
type PathFilter struct {
	// Paths is the list of paths this filter should match.
	// It is the callers responsibility to normalize paths accordingly.
	Paths []string
}

// Matches checks if a repository at clonePath matches this filter.
// Root indicates the root of all repositories.
func (pf PathFilter) Matches(env Env, clonePath string) bool {
	return text.SliceContainsAny(pf.Paths, clonePath)
}

// Candidates returns a list of folders that should be scanned regardless of their location.
func (pf PathFilter) Candidates() []string {
	return pf.Paths
}

// NewPatternFilter returns a new pattern filter with the appropriate value
func NewPatternFilter(value string, fuzzy bool) (pat PatternFilter) {
	pat.fuzzy = fuzzy
	pat.Set(value)
	return
}

// PatternFilter is a Filter that matches both paths and URLs according to a pattern.
// PatternFilter implements FilterValue
type PatternFilter struct {
	value   string
	fuzzy   bool
	pattern pattern.SplitPattern
}

func (pat PatternFilter) String() string {
	return pat.value
}

// Set sets the value of this filter.
//
// This function is untested because NewPatternFilter() is tested.
func (pat *PatternFilter) Set(value string) {
	pat.value = value
	pat.pattern = pattern.NewSplitGlobPattern(value, ComponentsOf, pat.fuzzy)
}

// Matches checks if this filter matches the repository at clonePath.
// The caller may assume that there is a repository at clonePath.
func (pat PatternFilter) Matches(env Env, clonePath string) bool {
	// find the remote url to use
	remote, err := env.Git.GetRemote(clonePath)
	if err != nil {
		return false
	}

	// if there is no remote url (because the repo has been cleanly "init"ed)
	// we use the relative path to the root directory to match.
	if remote == "" {
		root, err := env.absRoot()
		if err != nil { // root not resolved
			return false
		}
		actualClonePath, err := filepath.Abs(clonePath)
		if err != nil { // clonepath not resolved
			return false
		}
		remote, err = filepath.Rel(root, actualClonePath)
		if err != nil { // relative path not resolved
			return false
		}
	}

	return pat.pattern.Match(remote)
}

// MatchesURL checks if this filter matches a url
func (pat PatternFilter) MatchesURL(url URL) bool {
	parts := strings.Join(url.Components(), string(os.PathSeparator))
	return pat.pattern.Match(parts)
}

// DisjunctionFilter represents a filter that joins existing filters using an 'or' clause.
type DisjunctionFilter struct {
	Clauses []Filter
}

// Matches checks if this filter matches any of the filters that were joined.
func (or DisjunctionFilter) Matches(env Env, clonePath string) bool {
	for _, f := range or.Clauses {
		if f.Matches(env, clonePath) {
			return true
		}
	}
	return false
}

// Candidates returns the candidates of this filter
func (or DisjunctionFilter) Candidates() []string {

	// gather a list of candidates
	candidates := make([]string, 0, len(or.Clauses))
	for _, clause := range or.Clauses { // most clauses will have exactly one candidate, hence len(or.Clauses) should be enough to never reallocate
		candidates = append(candidates, Candidates(clause)...)
	}

	// remove duplicates from the result
	return text.RemoveDuplicates(candidates)
}

// TODO: Do we need tests for this?

// StatusFilter filters all elements in Filter by if they are clean or dirty
type StatusFilter struct {
	Filter

	Clean bool
	Dirty bool
}

func (sf StatusFilter) Candidates() []string {
	return Candidates(sf.Filter)
}

func (sf StatusFilter) Matches(env Env, clonePath string) bool {
	// first filter by the filter itself
	if !sf.Filter.Matches(env, clonePath) {
		return false
	}
	// if both or neither are included, this is quick to determine.
	if sf.Dirty == sf.Clean {
		return sf.Dirty
	}

	// check if the repo itself is dirty
	dirty, err := env.Git.IsDirty(clonePath)
	if err != nil {
		return false
	}

	// and ensure that it matches the dirty state!
	return dirty == sf.Dirty
}
