package cmd_test

//spellchecker:words strconv testing ggman internal cmdtest mockenv testutil
import (
	"strconv"
	"testing"

	"go.tkw01536.de/ggman/cmd"
	"go.tkw01536.de/ggman/internal/cmdtest"
	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/ggman/internal/testutil"
)

//spellchecker:words workdir nolint tparallel paralleltest

//nolint:tparallel,paralleltest
func TestCommandPull(t *testing.T) {
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
			"pull repository that has a new commit",
			"",
			[]string{"pull"},

			0,
			"Pulling " + escapedClonePath + "\n",
			"",
		},

		{
			"pull repository that doesn't have new commits",
			"",
			[]string{"pull"},

			0,
			"Pulling " + escapedClonePath + "\nalready up-to-date\n",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(cmd.Pull, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

func TestCommandPull_Overlap(t *testing.T) {
	t.Parallel()

	cmdtest.AssertFlagOverlap(t, cmd.Pull, []string{})
}
