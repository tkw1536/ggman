// Package pattern provides Pattern
//
//spellchecker:words pattern
package pattern

//spellchecker:words math strings github danwakefield fnmatch lithammer fuzzysearch fuzzy
import (
	"math"
	"strings"

	"github.com/danwakefield/fnmatch"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

//spellchecker:words casefold

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
// When not equal, we immediately return a negative score.
//
// When equal we determine the Levenshtein distance between the pattern and score.
// A higher distance, results in a lower score.
//
// Finally the score is normalized to the range [0, 1] using the length of s.
//
// See also [fuzzy.RankMatchFold].
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
// A string matches if they are equal under Unicode case-folding.
//
// See [strings.EqualFold] for a more detailed description.
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
// Matching is determined using case-insensitive glob matching.
func (p GlobPattern) Score(s string) float64 {
	if fnmatch.Match(string(p), s, fnmatch.FNM_CASEFOLD) {
		return 1
	}
	return -1
}

// NewSplitGlobPattern is a pattern that uses the given splitter for a new SplitPattern.
// If patterns starts with '^' or ends with '$', the fuzzy flag is ignored and MatchAt{Start,End} are set appropriately.
// Each sub-pattern consists of a call to NewGlobPattern.
func NewSplitGlobPattern(pattern string, splitter func(string) []string, fuzzy bool) SplitPattern {
	// check for the special case with ^ and '$'
	forceStart := false
	if len(pattern) > 0 && pattern[0] == '^' {
		forceStart = true
		pattern = pattern[1:]
	}
	forceEnd := false
	if len(pattern) > 0 && pattern[len(pattern)-1] == '$' {
		forceEnd = true
		pattern = pattern[:len(pattern)-1]
	}

	// disable fuzzy matching when ^ or $ are set
	if forceStart || forceEnd {
		fuzzy = false
	}

	// do the splitting
	globs := splitter(pattern)

	patterns := make([]Pattern, len(globs))
	for i, glob := range globs {
		patterns[i] = NewGlobPattern(glob, fuzzy)
	}

	return SplitPattern{
		Split:    splitter,
		Patterns: patterns,

		MatchAtStart: forceStart,
		MatchAtEnd:   forceEnd,
	}
}

// SplitPattern is a pattern that splits an input string and matches each string according to a sub-pattern.
// To compute the overall score, the pattern scores are averaged.
// It can optionally force matches at the start or end (or both) of the string.
//
// Finally, SplitPattern prioritizes matches at the edge of the input string.
// This means that a score at the start or the end of the string will always score higher than one in the middle.
type SplitPattern struct {
	// Split splits the input string into components
	Split func(s string) []string

	// Patterns are the patterns to score components with
	// A contiguous sequence of patterns of at least length 1 must score non-negatively in order for the predicate to apply.
	// The combined score is the average of all sub scores.
	//
	// The SplitPattern with an empty list of patterns always matches with the maximum score.
	Patterns []Pattern

	// MatchAtStart and MatchAtEnd determine if the SplitPattern must match at the start or the end of the string.
	MatchAtStart bool
	MatchAtEnd   bool
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

	// special cases for matching at start or at the end.
	// we need to match at most one candidate, and don't need to iterate.
	switch {
	case sp.MatchAtStart && sp.MatchAtEnd:
		// exact match, require that parts and patterns match!
		if len(parts) != len(sp.Patterns) {
			return -1
		}
		fallthrough // check position 0
	case sp.MatchAtStart:
		// match only at position 0
		score := sp.scoreFrom(parts, 0, last)
		if score > 0 {
			return score
		}
	case sp.MatchAtEnd:
		// match only at the last possible position
		score := sp.scoreFrom(parts, last, last)
		if score > 0 {
			return score
		}
	default:
		// check all possible positions
		for start := last; start >= 0; start-- {
			score := sp.scoreFrom(parts, start, last)
			if score > 0 {
				return score
			}
		}
	}

	// no match found
	return -1
}

// scoreFrom returns the score for this SplitPattern starting at start.
// It ignores sp.MatchAtStart and sp.MatchAtEnd.
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
