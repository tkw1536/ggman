package cmd_test

//spellchecker:words strconv testing ggman internal mockenv testutil
import (
	"strconv"
	"testing"

	"go.tkw01536.de/ggman/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/ggman/internal/testutil"
)

//spellchecker:words workdir nolint tparallel paralleltest

//nolint:tparallel,paralleltest
func TestCommandFetch(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// install git repo and make an extra commit
	repo := mock.Register("https://github.com/hello/world.git")
	clonePath := mock.Install("https://github.com/hello/world.git", "hello", "world")
	testutil.CommitTestFiles(repo, nil)

	escapedClonePath := strconv.Quote(clonePath)

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
			"Fetching " + escapedClonePath + "\n",
			"",
		},

		{
			"fetch repository that doesn't have new commits",
			"",
			[]string{"fetch"},

			0,
			"Fetching " + escapedClonePath + "\nalready up-to-date\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(t, cmd.NewCommand, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}
