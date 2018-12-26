package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTrimSuffixWhile tests the TrimWhile Function
func TestTrimSuffixWhile(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		s        string
		suffix   string
		expected string
	}{
		{"abcd", "d", "abc"},       // trim a single character
		{"abc", "d", "abc"},        // trim a non-existing character
		{"abcddd", "d", "abc"},     // trim a repeated character
		{"abc def", "", "abc def"}, // trim off the empty string
	}

	for _, row := range table {
		assert.Equal(row.expected, TrimSuffixWhile(row.s, row.suffix))
	}
}

// TestTrimPrefixWhile tests the TrimWhile Function
func TestTrimPrefixWhile(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		s        string
		prefix   string
		expected string
	}{
		{"abcd", "a", "bcd"},       // trim a single character
		{"bcd", "a", "bcd"},        // trim a non-existing character
		{"aaabcd", "a", "bcd"},     // trim a repeated character
		{"abc def", "", "abc def"}, // trim off the empty string
	}

	for _, row := range table {
		assert.Equal(row.expected, TrimPrefixWhile(row.s, row.prefix))
	}
}
