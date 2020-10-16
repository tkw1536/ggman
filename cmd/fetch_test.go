package cmd

import (
	"testing"

	"github.com/tkw1536/ggman/testutil"
	"github.com/tkw1536/ggman/testutil/mockenv"
)

func TestCommandFetch(t *testing.T) {
	mock, cleanup := mockenv.NewMockEnv()
	defer cleanup()

	// install git repo and make an extra commit
	repo, _ := mock.Register("https://github.com/hello/world.git")
	clonePath := mock.Install("https://github.com/hello/world.git", "hello", "world")
	testutil.CommitTestFiles(repo, nil)

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"fetch repository that has a new commit",
			"",
			[]string{"fetch"},

			0,
			"Fetching \"" + clonePath + "\"\nEnumerating objects: 3, done.\nCounting objects:  33% (1/3)\rCounting objects:  66% (2/3)\rCounting objects: 100% (3/3)\rCounting objects: 100% (3/3), done.\nCompressing objects:  50% (1/2)\rCompressing objects: 100% (2/2)\rCompressing objects: 100% (2/2), done.\nTotal 3 (delta 0), reused 0 (delta 0), pack-reused 0\n",
			"",
		},

		{
			"fetch repository that doesn't have new commits",
			"",
			[]string{"fetch"},

			0,
			"Fetching \"" + clonePath + "\"\nalready up-to-date\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(Fetch, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, stdout, tt.wantStdout, "Stdout")
			mock.AssertOutput(t, stderr, tt.wantStderr, "Stderr")
		})
	}
}
