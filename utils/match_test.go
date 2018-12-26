package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMatches tests the Matches Function
func TestMatches(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		s       string
		pattern string
		matches bool
	}{
		// matching the empty pattern
		{"", "", true},

		// matching one-component parts of a/b/c
		{"a/b/c", "a", true},
		{"a/b/c", "b", true},
		{"a/b/c", "c", true},
		{"a/b/c", "d", false},

		// matching constant sub-paths
		{"a/b/c/d/e/f", "b/c", true},
		{"a/b/c/d/e/f", "f/g", false},

		// matching sub-paths
		{"a/b/c/d/e/f", "b/*/d", true},
		{"a/b/c/d/e/f", "b/*/c", false},
	}

	for _, row := range table {
		got := Matches(row.pattern, row.s)
		assert.Equal(row.matches, got, "Matching "+row.s+" against "+row.pattern)
	}
}
