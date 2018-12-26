package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestComponents tests the Components Function
func TestComponents(t *testing.T) {
	assert := assert.New(t)

	// the two expected repositories
	gHWC := []string{"github.com", "hello", "world"}
	sURC := []string{"server.com", "user", "repository"}

	table := []struct {
		s        string
		expected []string
	}{
		// git@github.com/user/repo
		{"git@github.com/hello/world.git", gHWC},
		{"git@github.com:hello/world", gHWC},
		{"git@github.com:hello/world.git", gHWC},
		{"git@github.com:hello/world/", gHWC},
		{"git@github.com:hello/world//", gHWC},

		// ssh://git@github.com/hello/world
		{"ssh://git@github.com/hello/world.git", gHWC},
		{"ssh://git@github.com/hello/world", gHWC},
		{"ssh://git@github.com/hello/world/", gHWC},
		{"ssh://git@github.com/hello/world//", gHWC},

		// https://github.com/user/repo
		{"https://github.com:hello/world", gHWC},
		{"https://github.com/hello/world.git", gHWC},
		{"https://github.com:hello/world/", gHWC},
		{"https://github.com:hello/world//", gHWC},

		// user@server.com
		{"user@server.com:repository", sURC},
		{"user@server.com:repository/", sURC},
		{"user@server.com:repository//", sURC},
		{"user@server.com:repository.git", sURC},

		// ssh://user@server.com
		{"ssh://user@server.com/repository", sURC},
		{"ssh://user@server.com/repository/", sURC},
		{"ssh://user@server.com/repository//", sURC},
		{"ssh://user@server.com/repository.git", sURC},

		// ssh://user@server.com:1234
		{"ssh://user@server.com:1234/repository", sURC},
		{"ssh://user@server.com:1234/repository/", sURC},
		{"ssh://user@server.com:1234/repository//", sURC},
		{"ssh://user@server.com:1234/repository.git", sURC},
	}

	for _, row := range table {
		res, err := Components(row.s)
		if row.expected == nil {
			assert.Nil(res)
			assert.Error(err)
		} else {
			assert.Equal(res, row.expected, "Components of "+row.s)
			assert.NoError(err)
		}
	}
}
