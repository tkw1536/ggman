package cmd_test

//spellchecker:words path filepath testing github config plumbing ggman internal mockenv testutil
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/ggman/internal/testutil"
)

//spellchecker:words workdir GGROOT tparallel paralleltest nolint worktree

//nolint:tparallel,paralleltest
func TestCommandRelocate(t *testing.T) {
	t.Parallel()

	symlink := func(oldName, newName string) {
		err := os.Symlink(oldName, newName)
		if err != nil {
			panic(err)
		}
	}

	mock := mockenv.NewMockEnv(t)

	mock.Clone(t.Context(), "https://gitforge.example/right/directory.git", "gitforge.example", "right", "directory")
	mock.Clone(t.Context(), "https://gitforge.example/correct/directory.git", "gitforge.example", "incorrect", "directory")

	// link in an external repository in the right place
	external1 := mock.Clone(t.Context(), "https://gitforge.example/right/external1.git", "..", "external-path-1")
	symlink(external1, mock.Resolve(filepath.Join("gitforge.example", "right", "external1")))

	// link in an external repository in the right place
	external2 := mock.Clone(t.Context(), "https://gitforge.example/right/external2.git", "..", "external-path-2")
	symlink(external2, mock.Resolve(filepath.Join("gitforge.example", "right", "wrong-external")))

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"relocate with simulate",
			"",
			[]string{"relocate", "--simulate"},

			0,
			"mkdir -p `${GGROOT gitforge.example right}`\nmv `${GGROOT gitforge.example right wrong-external}` `${GGROOT gitforge.example right external2}`\nmkdir -p `${GGROOT gitforge.example correct}`\nmv `${GGROOT gitforge.example incorrect directory}` `${GGROOT gitforge.example correct directory}`\n",

			"",
		},

		{
			"relocate without simulate",
			"",
			[]string{"relocate"},

			0,
			"mkdir -p `${GGROOT gitforge.example right}`\nmv `${GGROOT gitforge.example right wrong-external}` `${GGROOT gitforge.example right external2}`\nmkdir -p `${GGROOT gitforge.example correct}`\nmv `${GGROOT gitforge.example incorrect directory}` `${GGROOT gitforge.example correct directory}`\n",

			"",
		},

		{
			"nothing to relocate",
			"",
			[]string{"relocate"},

			0,
			"",

			"",
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

func TestCommandRelocate_existsRepo(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// clone the same repository twice
	mock.Register("https://gitforge.example/right/directory.git")
	mock.Install(t.Context(), "https://gitforge.example/right/directory.git", "gitforge.example", "right", "directory")
	mock.Install(t.Context(), "https://gitforge.example/right/directory.git", "gitforge.example", "right", "other")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"relocate with simulate",
			"",
			[]string{"relocate", "--simulate"},

			0,
			"mkdir -p `${GGROOT gitforge.example right}`\nmv `${GGROOT gitforge.example right other}` `${GGROOT gitforge.example right directory}`\n",

			"",
		},

		{
			"relocate without simulate",
			"",
			[]string{"relocate"},

			1,
			"mkdir -p `${GGROOT gitforge.example right}`\nmv `${GGROOT gitforge.example right other}` `${GGROOT gitforge.example right directory}`\n",

			"failed to move repository: repository already exists at \"${GGROOT gitforge.example right directory}\"\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

func TestCommandRelocate_existsPath(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	// clone the same repository twice
	mock.Clone(t.Context(), "https://gitforge.example/right/directory.git", "gitforge.example", "wrong", "directory")

	if err := os.MkdirAll(mock.Resolve("gitforge.example", "right", "directory"), os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"relocate with simulate",
			"",
			[]string{"relocate", "--simulate"},

			0,
			"mkdir -p `${GGROOT gitforge.example right}`\nmv `${GGROOT gitforge.example wrong directory}` `${GGROOT gitforge.example right directory}`\n",

			"",
		},

		{
			"relocate without simulate",
			"",
			[]string{"relocate"},

			1,
			"mkdir -p `${GGROOT gitforge.example right}`\nmv `${GGROOT gitforge.example wrong directory}` `${GGROOT gitforge.example right directory}`\n",

			"\"${GGROOT gitforge.example right directory}\": failed to move repository: path already exists\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, tt.workdir, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			mock.AssertOutput(t, "Stdout", stdout, tt.wantStdout)
			mock.AssertOutput(t, "Stderr", stderr, tt.wantStderr)
		})
	}
}

//nolint:paralleltest
func TestCommandRelocate_multipleRemotes(t *testing.T) {
	const (
		originRemote = "https://gitforge.example/origin/repo.git"
		forkRemote   = "https://gitforge.example/fork/repo.git"
	)

	// setup sets up a new mock environment.
	// It then creates a new repository with two branches:
	// - 'origin_branch' => tracking the origin remote
	// - 'fork_branch' => tracking the fork remote
	// The given repository is installed to the given install path.
	// Finally, the given branch is checked out.
	setup := func(t *testing.T, installPath []string, branch string) *mockenv.MockEnv {
		t.Helper()

		//
		mock := mockenv.NewMockEnv(t)
		mock.Register(originRemote)
		_, forkURLs := mock.Register(forkRemote)

		repoPath := mock.Install(t.Context(), originRemote, installPath...)

		// Open the cloned repository to manipulate it
		repo, err := git.PlainOpen(repoPath)
		if err != nil {
			panic(err)
		}

		// Add a second remote "fork" pointing to the fork remote
		if _, err := repo.CreateRemote(&config.RemoteConfig{
			Name: "fork",
			URLs: []string{forkURLs[0]},
		}); err != nil {
			panic(err)
		}

		// Get the worktree to create and checkout branches
		wt, err := repo.Worktree()
		if err != nil {
			panic(err)
		}

		// Create two branches:
		// - origin_branch => tracking origin
		// - fork_branch => tracking fork
		testutil.CommitTestFiles(repo)
		testutil.CreateTrackingBranch(repo, "origin", "origin_branch", "main")
		testutil.CreateTrackingBranch(repo, "fork", "fork_branch", "main")

		// Checkout the specified branch
		if err := wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branch),
		}); err != nil {
			panic(err)
		}

		return mock
	}

	t.Run("at canonical path for checked out branch (origin)", func(t *testing.T) {
		// Repository is at gitforge.example/origin/repo (canonical for origin remote)
		// and origin_branch is checked out (tracking origin remote)
		mock := setup(t, []string{"gitforge.example", "origin", "repo"}, "origin_branch")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "relocate", "--simulate")
		if code != 0 {
			t.Errorf("Code = %d, wantCode = 0", code)
		}
		// No relocation needed - already at canonical path for origin remote
		mock.AssertOutput(t, "Stdout", stdout, "")
		mock.AssertOutput(t, "Stderr", stderr, "")
	})

	t.Run("at canonical path for checked out branch (origin) with only current remote", func(t *testing.T) {
		// Repository is at gitforge.example/origin/repo (canonical for origin remote)
		// and origin_branch is checked out (tracking origin remote)
		mock := setup(t, []string{"gitforge.example", "origin", "repo"}, "origin_branch")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "relocate", "--simulate", "--only-current-remote")
		if code != 0 {
			t.Errorf("Code = %d, wantCode = 0", code)
		}
		// No relocation needed - already at canonical path for origin remote
		mock.AssertOutput(t, "Stdout", stdout, "")
		mock.AssertOutput(t, "Stderr", stderr, "")
	})

	t.Run("at canonical path for non-checked out branch (fork)", func(t *testing.T) {
		// Repository is at gitforge.example/fork/repo (canonical for fork remote)
		// but origin_branch is checked out (tracking origin remote)
		mock := setup(t, []string{"gitforge.example", "fork", "repo"}, "origin_branch")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "relocate", "--simulate")
		if code != 0 {
			t.Errorf("Code = %d, wantCode = 0", code)
		}
		// No relocation needed - at canonical path for fork remote (even though origin_branch is checked out)
		mock.AssertOutput(t, "Stdout", stdout, "")
		mock.AssertOutput(t, "Stderr", stderr, "")
	})

	t.Run("at canonical path for non-checked out branch (fork) with only current remote", func(t *testing.T) {
		// Repository is at gitforge.example/fork/repo (canonical for fork remote)
		// but origin_branch is checked out (tracking origin remote)
		mock := setup(t, []string{"gitforge.example", "fork", "repo"}, "origin_branch")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "relocate", "--simulate", "--only-current-remote")
		if code != 0 {
			t.Errorf("Code = %d, wantCode = 0", code)
		}

		// Should relocate to origin remote path (canonical remote of current branch)
		mock.AssertOutput(t, "Stdout", stdout, "mkdir -p `${GGROOT gitforge.example origin}`\nmv `${GGROOT gitforge.example fork repo}` `${GGROOT gitforge.example origin repo}`\n")
		mock.AssertOutput(t, "Stderr", stderr, "")
	})

	t.Run("not at canonical path for any remote", func(t *testing.T) {
		// Repository is at gitforge.example/wrong/repo (not canonical for either remote)
		// fork_branch is checked out (tracking fork remote)
		mock := setup(t, []string{"gitforge.example", "wrong", "repo"}, "fork_branch")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "relocate", "--simulate")
		if code != 0 {
			t.Errorf("Code = %d, wantCode = 0", code)
		}
		// Should relocate to fork remote path (canonical remote of current branch)
		mock.AssertOutput(t, "Stdout", stdout, "mkdir -p `${GGROOT gitforge.example fork}`\nmv `${GGROOT gitforge.example wrong repo}` `${GGROOT gitforge.example fork repo}`\n")
		mock.AssertOutput(t, "Stderr", stderr, "")
	})

	t.Run("not at canonical path for any remote and only-current-remote set", func(t *testing.T) {
		// Repository is at gitforge.example/wrong/repo (not canonical for either remote)
		// fork_branch is checked out (tracking fork remote)
		mock := setup(t, []string{"gitforge.example", "wrong", "repo"}, "fork_branch")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "relocate", "--simulate", "--only-current-remote")
		if code != 0 {
			t.Errorf("Code = %d, wantCode = 0", code)
		}
		// Should relocate to fork remote path (canonical remote of current branch)
		mock.AssertOutput(t, "Stdout", stdout, "mkdir -p `${GGROOT gitforge.example fork}`\nmv `${GGROOT gitforge.example wrong repo}` `${GGROOT gitforge.example fork repo}`\n")
		mock.AssertOutput(t, "Stderr", stderr, "")
	})
}
