package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParseArgs tests the ParseArgs Function
func TestParseArgs(t *testing.T) {
	assert := assert.New(t)

	table := []struct {
		input    []string
		expected *GGArgs
	}{
		// no arguments => parsing fails
		{[]string{}, nil},

		// command without arguments => ok
		{[]string{"cmd"}, &GGArgs{"cmd", "", []string{}}},

		// command with arguments => ok
		{[]string{"cmd", "a1", "a2"}, &GGArgs{"cmd", "", []string{"a1", "a2"}}},

		// only a for => parsing fails
		{[]string{"for"}, nil},

		// for without command => parsing fails
		{[]string{"for", "match"}, nil},

		// for with command => ok
		{[]string{"for", "match", "cmd"}, &GGArgs{"cmd", "match", []string{}}},

		// for with command and arguments => ok
		{[]string{"for", "match", "cmd", "a1", "a2"}, &GGArgs{"cmd", "match", []string{"a1", "a2"}}},
	}

	for _, row := range table {
		res, err := ParseArgs(row.input)
		if row.expected == nil {
			assert.Nil(res)
			assert.NotEqual("", err)
		} else {
			assert.Equal(row.expected, res)
			assert.Equal("", err)
		}
	}
}
