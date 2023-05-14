// Package pattern provides Pattern
package pattern

import (
	"math"
	"strings"

	"github.com/danwakefield/fnmatch"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// Pattern is a float predicate on strings.
//
// A pattern must return a float between 0 and 1 (inclusive) when it matches and a negative number when not.
// A higher score indicates a higher match.
type Pattern interface {
	Score(s string) float64
}

const globPatternBlacklist = "*?\\["

// NewGlobPattern creates a new case-insensitive pattern that scores strings either using a glob pattern (when it contains special characters) or using fuzzy or exact case-folded equality.
// This function essentially behaves like GlobPattern, however uses performance optimizations to avoid calls to fnmatch.
// As a special case, a GlobPattern with the empty string scores any string with the highest possible score.
func NewGlobPattern(s string, fuzzy bool) Pattern {
	if s == "" {
		return AnyStringPattern{}
	}
	// when no special characters are contained, we can use an equality pattern
	if !strings.ContainsAny(s, globPatternBlacklist) {
		if fuzzy {
			return FuzzyFoldPattern(s)
		}
		return EqualityFoldPattern(s)
	}
	return GlobPattern(s)
}

// AnyStringPattern scores any string with the highest possible score.
type AnyStringPattern struct{}

// Score scores s against this string pattern.
// Returns 1.
func (AnyStringPattern) Score(s string) float64 {
	return 1
}

// FuzzyFoldPattern is a pattern that matches strings based on fuzzy equality.
type FuzzyFoldPattern string

// Score scores a string against this pattern.
//
// To determine the score, we first check if they are reasonably equal.
// When not equal, we immediatly return a negative score.
//
// When equal we determine the Levenshtein distance between the pattern and score.
// A higher distance, results in a lower score.
//
// Finally the score is normalized to the range [0, 1] using the length of s.
//
// See also "github.com/lithammer/fuzzysearch/fuzzy".RankMatchFold.
func (p FuzzyFoldPattern) Score(s string) float64 {
	score := float64(fuzzy.RankMatchFold(string(p), s))
	if score == -1 {
		return -1
	}
	return 1 - (score / float64(len(s)))
}

// EqualityFoldPattern is a pattern that scores strings -1 or 1 based on equality.
type EqualityFoldPattern string

// Score scores a string against this pattern.
// A string matches an EqualityFoldPattern if they are equal under Unicode case-folding.
//
// See strings.EqualFold for a more detailed description.
// Returns 1 when the string matches, and -1 when not.
func (p EqualityFoldPattern) Score(s string) float64 {
	if strings.EqualFold(string(p), s) {
		return 1
	}
	return -1
}

// GlobPattern represents a pattern that scores a string based on a 'glob'-like pattern.
type GlobPattern string

// Score checks if a string matches this pattern.
// When a string matches, returns a score of 1, else -1.
// Matching is determined using case-insenstive glob matching.
func (p GlobPattern) Score(s string) float64 {
	if fnmatch.Match(string(p), s, fnmatch.FNM_CASEFOLD) {
		return 1
	}
	return -1
}

// NewSplitGlobPattern is a pattern that uses the given splitter for a new SplitPattern.
// Each subpattern consists of a call to NewGlobPattern.
func NewSplitGlobPattern(pattern string, splitter func(string) []string, fuzzy bool) SplitPattern {
	globs := splitter(pattern)

	patterns := make([]Pattern, len(globs))
	for i, glob := range globs {
		patterns[i] = NewGlobPattern(glob, fuzzy)
	}

	return SplitPattern{
		Split:    splitter,
		Patterns: patterns,
	}
}

// SplitPattern is a pattern that splits an input string and matches each string according to a sub-pattern.
// To compute the overall score, the pattern scores are averaged.
//
// Finally, SplitPattern prioritizes matches at the edge of the input string.
// This means that a score at the start or the end of the string will always score higher than one in the middle.
type SplitPattern struct {
	// Split splits the input string
	Split func(s string) []string

	// Patterns are the patterns to score components with
	// A contiguous sequence of patterns of at least length 1 must score non-negatively in order for the predicate to apply.
	// The combined score is the average of all sub scores.
	//
	// The SplitPattern with an empty list of patterns always matches with the maximum score.
	Patterns []Pattern
}

// Score scores a string.
//
// It first splits the function according to the split function.
// Then it matches each part according to a corresponding sub-pattern.
// It adds the scores and normalizes according to the number of parts matched, and where the sequence of these matches starts.
func (sp SplitPattern) Score(s string) (score float64) {
	// quick path: no patterns, so we can immediately return!
	if len(sp.Patterns) == 0 {
		return 1
	}

	// split the string into parts
	parts := sp.Split(s)

	// find the last possible place where a match can start
	last := len(parts) - len(sp.Patterns)
	if last < 0 {
		// no possible match!
		return -1
	}

	// find the last match
	for start := last; start >= 0; start-- {
		score := sp.scoreFrom(parts, start, last)
		if score > 0 {
			return score
		}
	}

	// no match found
	return -1
}

// scoreFrom returns the score for this SplitPattern starting at start.
func (sp SplitPattern) scoreFrom(parts []string, start int, last int) (score float64) {
	// compute the average score for each pattern
	for i, pat := range sp.Patterns {
		partial := pat.Score(parts[start+i])
		if partial < 0 {
			return -1 // no match
		}

		score += partial
	}
	score /= float64(len(sp.Patterns))

	// Normalize for where the pattern starts.
	//
	// Consider for example the sub-patterns ['hello', 'world'] to be matched against:
	//
	// ['i-have-three-levels', 'hello', 'world', 'stuff']
	// ['i-have-three-levels', 'stuff', 'hello', 'world']
	// ['i-have-two-levels', 'hello', 'world']
	//
	// We want the latter two scores to be higher (as they are matched near the edges).
	// We use the following algorithm:
	//
	// Suppose the match is nth-from-the-edge (i.e. at position n or ending at last - n).
	// Then we place the score into the interval (1/2^(n+1) ... 1/2^(n)).
	n := math.Min(float64(start), float64(last-start))
	m := math.Pow(2, -(n + 1)) // 1/2^(n+1)
	score = m * (score + 1)    // align (0....1) to the interval (m,2m)

	return
}
