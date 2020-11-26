package util

import (
	"strings"

	"github.com/danwakefield/fnmatch"
)

// Pattern is a predicate on strings.
type Pattern interface {
	Match(s string) bool
}

const globPatternBlacklist = "*?\\["

// NewGlobPattern creates a new case-insensitive pattern that matches a string based on a 'glob'-like pattern.
// This function essentially behaves like GlobPattern, however uses performance optimizations to avoid calls to fnmatch.
// As a special case, a GlobPattern with the empty string matches any string.
func NewGlobPattern(s string) Pattern {
	if s == "" {
		return AnyStringPattern{}
	}
	// when no special characters are contained, we can use an equality pattern
	if !strings.ContainsAny(s, globPatternBlacklist) {
		return EqualityFoldPattern(s)
	}
	return GlobPattern(s)
}

// AnyStringPattern matches any string
type AnyStringPattern struct{}

// Match checks if a string matches this pattern.
// It always returns true.
func (AnyStringPattern) Match(s string) bool {
	return true
}

// EqualityFoldPattern is a pattern that matches strings based on equality.
type EqualityFoldPattern string

// Match checks if a string matches this pattern.
// A string matches an EqualityFoldPattern if they are equal under Unicode case-folding.
// See strings.EqualFold for a more detailed description.
func (p EqualityFoldPattern) Match(s string) bool {
	return strings.EqualFold(string(p), s)
}

// GlobPattern represents a pattern that matches based on a string based on a 'glob'-like pattern.
type GlobPattern string

// Match checks if a string matches this pattern.
// Matching is determined using case-insenstive glob matching.
func (p GlobPattern) Match(s string) bool {
	return fnmatch.Match(string(p), s, fnmatch.FNM_CASEFOLD)
}

// NewSplitGlobPattern is a pattern that uses the given splitter for a new SplitPattern that consists of GlobPatterns.
func NewSplitGlobPattern(pattern string, splitter func(string) []string) SplitPattern {
	globs := splitter(pattern)

	patterns := make([]Pattern, len(globs))
	for i, glob := range globs {
		patterns[i] = NewGlobPattern(glob)
	}

	return SplitPattern{
		Split:    splitter,
		Patterns: patterns,
	}
}

// SplitPattern is a pattern that splits an input string and matches each string according to a subpattern.
type SplitPattern struct {
	// Split splits the input string
	Split func(s string) []string

	// Patterns are the patterns to match components with
	// A contiguous sequence of patterns of at least length 1 must be matched in order for the predicate to apply.
	// The SplitPattern with an empty list of patterns always matches.
	Patterns []Pattern
}

// Match checks if s matches this SplitPattern
func (sp SplitPattern) Match(s string) bool {
	// when we have no patterns, we can return true right away!
	if len(sp.Patterns) == 0 {
		return true
	}

	parts := sp.Split(s)
	last := len(parts) - len(sp.Patterns)
outer:
	for i := 0; i <= last; i++ {
		for j, pattern := range sp.Patterns {
			if !pattern.Match(parts[i+j]) {
				continue outer
			}
		}

		return true
	}
	return false
}
