package env

import (
	"path/filepath"
	"strings"

	"github.com/danwakefield/fnmatch"
)

// Filter is a filter for an environment.
// It satisifies the pflag Value interface.
// A filter should only be created using the NewFilter method or by using the NoFilter variable.
type Filter struct {
	filter string

	pattern []pattern
}

func (f Filter) String() string {
	return f.filter
}

// Set sets the value of this filter
func (f *Filter) Set(value string) error {
	components := ParseURL(value).Components()
	pattern := make([]pattern, len(components))
	for i, c := range components {
		pattern[i] = newPattern(c)
	}

	f.filter = value
	f.pattern = pattern

	return nil
}

// Type returns the type of this filter
func (f Filter) Type() string {
	return "filter"
}

// IsEmpty checks if the filter f includes all repositories
func (f Filter) IsEmpty() bool {
	return f.filter == "" || f.filter == "*"
}

// NoFilter is the absence of a Filter.
var NoFilter Filter

// NewFilter creates a new filter from a string
func NewFilter(input string) Filter {
	f := &Filter{}
	f.Set(input)
	return *f
}

// Matches checks if this filter matches the repository at clonePath.
// The caller may assume that there is a repository at clonePath.
func (f Filter) Matches(root, clonePath string) bool {
	if f.IsEmpty() { // this is neccessary with the current implementation
		return true
	}

	relpath, err := filepath.Rel(root, clonePath)
	if err != nil {
		return false
	}

	return ParseURL(relpath).MatchesFilter(f)
}

// Matches checks if a URL matches a given filter.
func (url URL) Matches(pattern string) bool {
	filter := NewFilter(pattern)
	if filter.IsEmpty() {
		return true
	}
	return url.MatchesFilter(filter)
}

// MatchesFilter checks if a filter matches a pattern.
func (url URL) MatchesFilter(filter Filter) bool {
	components := url.Components()

	last := len(components) - len(filter.pattern)
outer:
	for i := 0; i <= last; i++ {
		for j, pattern := range filter.pattern {
			if !pattern.Match(components[i+j]) {
				continue outer
			}
		}

		return true
	}

	return false
}

// pattern represents a single component of a filter.
type pattern interface {
	Match(s string) bool
}

const patternExpensive = "*?\\["

// newPattern makes a new pattern from a string
func newPattern(s string) pattern {
	if !strings.ContainsAny(s, patternExpensive) {
		return cheappattern(s)
	}
	return exppattern(s)
}

// cheappattern is a pattern that does not contain any special characters
// and can be evaluated by using string comparison operators.
type cheappattern string

func (p cheappattern) Match(s string) bool {
	return strings.EqualFold(string(p), s)
}

// exppattern is a pattern that contains special characters that require it to be evaluated using
// a call to fnamtch.
type exppattern string

func (p exppattern) Match(s string) bool {
	return fnmatch.Match(string(p), s, fnmatch.FNM_CASEFOLD)
}
