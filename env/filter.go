package env

//spellchecker:words path filepath strings ggman internal pattern pkglib collection
import (
	"os"
	"path/filepath"
	"strings"

	"go.tkw01536.de/ggman/internal/path"
	"go.tkw01536.de/ggman/internal/pattern"
	"go.tkw01536.de/pkglib/collection"
)

//spellchecker:words worktree nolint

// Filter is a predicate that scores repositories inside an environment.
//
// A filter is applied by recursively scanning the root folder for git repositories.
// Each folder that is a repository will be passed to clonePath.
//
// Filter may also optionally implement FilterWithCandidates.
type Filter interface {
	// Score scores the repository at clonePath against this filter.
	//
	// When it does match, returns a float64 between 0 and 1 (inclusive on both ends),
	// If the filter does not match, returns a negative number such as [FilterDoesNotMatch].
	Score(env Env, clonePath string) float64
}

// FilterDoesNotMatch should is used by a [Filter] to indicate that it does not match.
const FilterDoesNotMatch = float64(-1)

// NoFilter is a special filter that matches every directory with the highest possible score.
var NoFilter Filter = emptyFilter{}

type emptyFilter struct{}

func (emptyFilter) Score(env Env, clonePath string) float64 {
	return 1
}

// FilterWithCandidates is a filter that in addition to being applied normally should also be applied to the provided candidates.
type FilterWithCandidates interface {
	Filter

	// Candidates returns a list of folders that should be scanned regardless of their location.
	// A FilterWithCandidates with a Candidates() function that returns a zero-length slice is equivalent to a regular filter.
	Candidates() []string
}

// When Filter does not implement FilterWithCandidates, returns nil.
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
	// Paths is the list of paths this filter should match with the highest possible score.
	// It is the callers responsibility to normalize paths accordingly.
	Paths []string
}

// Score checks if a repository at clonePath matches this filter, and if so returns 1.
// See Filter.Score.
func (pf PathFilter) Score(env Env, clonePath string) float64 {
	for _, p := range pf.Paths {
		if path.HasChild(p, clonePath) {
			return 1
		}
	}
	return FilterDoesNotMatch
}

// Candidates returns a list of folders that should be scanned regardless of their location.
func (pf PathFilter) Candidates() []string {
	return pf.Paths
}

// NewPatternFilter returns a new pattern filter with the appropriate value.
func NewPatternFilter(value string, fuzzy bool) *PatternFilter {
	var pat PatternFilter
	pat.fuzzy = fuzzy
	pat.Set(value)
	return &pat
}

// PatternFilter implements FilterValue.
type PatternFilter struct {
	value   string
	fuzzy   bool
	pattern pattern.SplitPattern
}

func (pat *PatternFilter) String() string {
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
func (pat *PatternFilter) Score(env Env, clonePath string) float64 {
	// find the remote url to use
	remote, err := env.Git.GetRemote(clonePath, "")
	if err != nil {
		return FilterDoesNotMatch
	}

	// if there is no remote url (because the repo has been cleanly "init"ed)
	// we use the relative path to the root directory to match.
	if remote == "" {
		root, err := env.absRoot()
		if err != nil { // root not resolved
			return FilterDoesNotMatch
		}
		actualClonePath, err := filepath.Abs(clonePath)
		if err != nil { // clone path not resolved
			return FilterDoesNotMatch
		}
		remote, err = filepath.Rel(root, actualClonePath)
		if err != nil { // relative path not resolved
			return FilterDoesNotMatch
		}
	}

	return pat.pattern.Score(remote)
}

// MatchesURL checks if this filter matches a url.
func (pat *PatternFilter) MatchesURL(url URL) bool {
	parts := strings.Join(url.Components(), string(os.PathSeparator))
	return pat.pattern.Score(parts) >= 0
}

// DisjunctionFilter represents a filter that joins existing filters using an 'or' clause.
type DisjunctionFilter struct {
	Clauses []Filter
}

// Matches checks if this filter matches any of the filters that were joined.
// It returns the highest possible score.
func (or DisjunctionFilter) Score(env Env, clonePath string) float64 {
	score := FilterDoesNotMatch
	for _, f := range or.Clauses {
		if fScore := f.Score(env, clonePath); fScore > score {
			score = fScore
		}
	}
	return score
}

// Candidates returns the candidates of this filter.
func (or DisjunctionFilter) Candidates() []string {
	// gather a list of candidates
	candidates := make([]string, 0, len(or.Clauses))
	for _, clause := range or.Clauses { // most clauses will have exactly one candidate, hence len(or.Clauses) should be enough to never reallocate
		candidates = append(candidates, Candidates(clause)...)
	}

	// remove duplicates from the result
	return collection.Deduplicate(candidates)
}

// TODO: Do we need tests for this?

// predicateFilter implements [Filter].
// It applies Predicate to each matching repository, and only
// includes them as dictated by [IncludeTrue] and [IncludeFalse].
type predicateFilter struct {
	Filter Filter

	// Predicate is the predicate to apply.
	Predicate func(env Env, clonePath string) bool

	// IncludeTrue and IncludeFalse determine
	// which values of the predicate should be included.
	IncludeTrue, IncludeFalse bool
}

func (pf predicateFilter) Candidates() []string {
	return Candidates(pf.Filter)
}

func (pf predicateFilter) Score(env Env, clonePath string) float64 {
	// include nothing
	if !pf.IncludeTrue && !pf.IncludeFalse {
		return FilterDoesNotMatch
	}

	// does the underlying filter match?
	score := pf.Filter.Score(env, clonePath)
	if score < 0 {
		return FilterDoesNotMatch
	}

	// both are included, so we don't need to do any more checking
	if pf.IncludeTrue && pf.IncludeFalse {
		return score
	}

	// determine if the repository should be included
	include := pf.Predicate(env, clonePath)
	if pf.IncludeFalse {
		include = !include
	}

	if !include {
		return FilterDoesNotMatch
	}
	return score
}

// NewWorktreeFilter returns a filter that filters by repositories having a dirty or clean working directory.
func NewWorktreeFilter(filter Filter, dirty, clean bool) Filter {
	return predicateFilter{
		Filter: filter,
		Predicate: func(env Env, clonePath string) bool {
			dirty, err := env.Git.IsDirty(clonePath)
			return err == nil && dirty
		},

		IncludeTrue:  dirty,
		IncludeFalse: clean,
	}
}

// NewStatusFilter returns  new Filter that filters by repositories being synced or un-synced with the remote.
func NewStatusFilter(filter Filter, synced, unSynced bool) Filter {
	return predicateFilter{
		Filter: filter,
		Predicate: func(env Env, clonePath string) bool {
			sync, err := env.Git.IsSync(clonePath)
			return err == nil && sync
		},

		IncludeTrue:  synced,
		IncludeFalse: unSynced,
	}
}

// NewTarnishFilter returns new Filter that filters by if they have been tarnished or not.
// A repository is tarnished if it has a dirty working directory, or is un-synced with the remote.
func NewTarnishFilter(filter Filter, tarnished, pristine bool) Filter {
	return predicateFilter{
		Filter: filter,
		Predicate: func(env Env, clonePath string) bool {
			dirty, err := env.Git.IsDirty(clonePath)
			if err != nil {
				return false
			}
			if dirty {
				return true
			}

			synced, err := env.Git.IsSync(clonePath)
			if err != nil {
				return false
			}

			return !synced
		},

		IncludeTrue:  tarnished,
		IncludeFalse: pristine,
	}
}
