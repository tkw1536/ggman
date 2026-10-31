package cmd_test

//spellchecker:words strconv strings testing ggman internal mockenv testutil
import (
	"strconv"
	"strings"
	"testing"

	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/ggman/internal/testutil"
)

//spellchecker:words workdir tparallel paralleltest

//nolint:tparallel,paralleltest
func TestCommandPull(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// install git repo and make an extra commit
	repo, _ := mock.Register("https://gitforge.example/hello/world.git")
	clonePath := mock.Install(t.Context(), "https://gitforge.example/hello/world.git", "hello", "world")
	testutil.CommitTestFiles(repo)

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
			"Pulling " + escapedClonePath + "\nPulled 1 repositories, 0 failed\n",
			"",
		},

		{
			"pull repository that doesn't have new commits",
			"",
			[]string{"pull"},

			0,
			"Pulling " + escapedClonePath + "\nalready up-to-date\nPulled 1 repositories, 0 failed\n",
			"",
		},

		{
			"reject negative parallel",
			"",
			[]string{"pull", "--parallel", "-1"},

			4,
			"",
			"argument for \"--parallel\" must be non-negative\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

//nolint:tparallel,paralleltest
func TestCommandPull_parallel(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	repo1, _ := mock.Register("https://gitforge.example/hello/world.git")
	mock.Install(t.Context(), "https://gitforge.example/hello/world.git", "hello", "world")
	testutil.CommitTestFiles(repo1)

	repo2, _ := mock.Register("https://githost.example/hello/world.git")
	mock.Install(t.Context(), "https://githost.example/hello/world.git", "githost.example", "hello", "world")
	testutil.CommitTestFiles(repo2)

	for _, parallel := range []string{"0", "2"} {
		t.Run("parallel="+parallel, func(t *testing.T) {
			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "pull", "--parallel", parallel)
			if code != 0 {
				t.Errorf("Code = %d, wantCode = 0", code)
			}
			if stderr != "" {
				t.Errorf("Stderr = %q, want empty", stderr)
			}
			if !strings.Contains(stdout, "Pulled 2 repositories, 0 failed\n") {
				t.Errorf("Stdout = %q, want summary for 2 repositories", stdout)
			}
		})
	}
}
