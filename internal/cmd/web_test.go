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

//spellchecker:words workdir reclone godoc tparallel paralleltest worktree

func TestCommandURL(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	clonePath := mock.Clone(t.Context(), "git@github.com/hello/world.git", "hello", "world")

	subClonePath := filepath.Join(clonePath, "sub")
	if err := os.MkdirAll(subClonePath, 0750); err != nil {
		panic(err)
	}

	nonRepoPath := filepath.Join(clonePath, "..", "..", "example.com", "other")
	if err := os.MkdirAll(nonRepoPath, 0750); err != nil {
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
			"Open url at root",
			clonePath,
			[]string{"url"},
			0,
			"https://github.com/hello/world\n",
			"",
		},

		{
			"Open url for specific remote",
			clonePath,
			[]string{"url", "--remote", "origin"},
			0,
			"https://github.com/hello/world\n",
			"",
		},

		{
			"Print clone url at root",
			clonePath,
			[]string{"url", "--clone"},
			0,
			"git clone https://github.com/hello/world.git\n",
			"",
		},

		{
			"Print reclone url at root",
			clonePath,
			[]string{"url", "--reclone"},
			0,
			"git clone git@github.com/hello/world.git\n",
			"",
		},

		{
			"Open url at faked root",
			nonRepoPath,
			[]string{"url", "--force-repo-here"},
			0,
			"https://example.com/other\n",
			"",
		},

		{
			"Open url with custom base",
			clonePath,
			[]string{"url", "https://mybase"},
			0,
			"https://mybase/hello/world\n",
			"",
		},

		{
			"Open url with custom and prefix base",
			clonePath,
			[]string{"url", "https://mybase/", "--prefix"},
			0,
			"https://mybase/github.com/hello/world\n",
			"",
		},

		{
			"Open url with predefined base with prefix",
			clonePath,
			[]string{"url", "godoc"},
			0,
			"https://pkg.go.dev/github.com/hello/world\n",
			"",
		},

		{
			"Open url with predefined base without prefix",
			clonePath,
			[]string{"url", "travis"},
			0,
			"https://travis-ci.com/hello/world\n",
			"",
		},

		{
			"Open url with tree at root",
			clonePath,
			[]string{"url", "--tree"},
			0,
			"https://github.com/hello/world/tree/master/\n",
			"",
		},

		{
			"Do not print clone url with tree at root",
			clonePath,
			[]string{"url", "--clone", "--tree"},
			4,
			"",
			`incompatible flags for "ggman web": "clone" and "tree"` + "\n",
		},

		{
			"Do not print reclone url with tree at root",
			clonePath,
			[]string{"url", "--reclone", "--tree"},
			4,
			"",
			`incompatible flags for "ggman web": "reclone" and "tree"` + "\n",
		},

		{
			"Open url at faked root with tree",
			nonRepoPath,
			[]string{"url", "--force-repo-here", "--tree"},
			0,
			"https://example.com/other\n",
			"",
		},

		{
			"Open url with branch at root",
			clonePath,
			[]string{"url", "--branch"},
			0,
			"https://github.com/hello/world/tree/master\n",
			"",
		},

		{
			"Print clone url with branch at root",
			clonePath,
			[]string{"url", "--clone", "--branch"},
			0,
			"git clone https://github.com/hello/world.git --branch master\n",
			"",
		},

		{
			"Print reclone url with branch at root",
			clonePath,
			[]string{"url", "--reclone", "--branch"},
			0,
			"git clone git@github.com/hello/world.git --branch master\n",
			"",
		},

		{
			"Open url at faked root with branch",
			nonRepoPath,
			[]string{"url", "--force-repo-here", "--branch"},
			0,
			"https://example.com/other\n",
			"",
		},

		{
			"Open url at subpath",
			subClonePath,
			[]string{"url"},
			0,
			"https://github.com/hello/world\n",
			"",
		},

		{
			"Print clone url at subpath",
			subClonePath,
			[]string{"url", "--clone"},
			0,
			"git clone https://github.com/hello/world.git\n",
			"",
		},

		{
			"Print reclone url at subpath",
			subClonePath,
			[]string{"url", "--reclone"},
			0,
			"git clone git@github.com/hello/world.git\n",
			"",
		},

		{
			"Open url with tree at subpath",
			subClonePath,
			[]string{"url", "--tree"},
			0,
			"https://github.com/hello/world/tree/master/sub\n",
			"",
		},

		{
			"Open url with branch at subpath",
			subClonePath,
			[]string{"url", "--branch"},
			0,
			"https://github.com/hello/world/tree/master\n",
			"",
		},

		{
			"Print clone url with branch at subpath",
			subClonePath,
			[]string{"url", "--clone", "--branch"},
			0,
			"git clone https://github.com/hello/world.git --branch master\n",
			"",
		},

		{
			"Print reclone url with branch at subpath",
			subClonePath,
			[]string{"url", "--reclone", "--branch"},
			0,
			"git clone git@github.com/hello/world.git --branch master\n",
			"",
		},

		{
			"List all bases",
			clonePath,
			[]string{"url", "--list-bases"},
			0,
			"circle: https://app.circleci.com/pipelines/github\ngodoc: https://pkg.go.dev/\nlocalgodoc: http://localhost:6060/pkg/\ntravis: https://travis-ci.com\n",
			"",
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

//nolint:tparallel,paralleltest
func TestCommandURL_MultipleRemotes(t *testing.T) {
	t.Parallel()
	// The sub-tests need to be run sequentially, because they modify the repository
	// using checkout.
	// This cannot happen concurrently.

	mock := mockenv.NewMockEnv(t)

	const (
		mainRemote = "git@example.com:main.git"
		mainURL    = "https://example.com/main"

		forkRemote = "git@example.com:fork.git"
		forkURL    = "https://example.com/fork"
	)

	mock.Register(mainRemote)
	_, featureRemotes := mock.Register(forkRemote)
	clonePath := mock.Install(t.Context(), mainRemote, "github.com", "user", "repo")

	// Open the cloned repository to manipulate it
	repo, err := git.PlainOpen(clonePath)
	if err != nil {
		panic(err)
	}

	// Add a second remote "upstream" pointing to the feature
	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "fork",
		URLs: []string{featureRemotes[0]},
	}); err != nil {
		panic(err)
	}

	// Get the worktree to create and checkout branches
	wt, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	// Create three branches:
	// - main => pointing to the main remote
	// - fork => pointing to the fork remote
	// - no_remote => no remote at all

	testutil.CommitTestFiles(repo)
	testutil.CreateTrackingBranch(repo, "origin", "main", "main")
	testutil.CreateTrackingBranch(repo, "fork", "fork", "fork")
	if err := wt.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("no_remote"),
		Create: true,
	}); err != nil {
		panic(err)
	}

	// helper function so we can checkout specific branches
	checkout := func(branch string) {
		if err := wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branch),
		}); err != nil {
			panic(err)
		}
	}

	t.Run("main branch", func(t *testing.T) {
		checkout("main")

		if err := wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName("main"),
		}); err != nil {
			panic(err)
		}

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, clonePath, "", "url")
		if code != 0 {
			t.Errorf("main branch: Code = %d, wantCode = 0", code)
		}
		mock.AssertOutput(t, "main branch: Stdout", stdout, mainURL+"\n")
		mock.AssertOutput(t, "main branch: Stderr", stderr, "")
	})

	t.Run("main branch with --remote", func(t *testing.T) {
		checkout("main")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, clonePath, "", "url", "--remote", "fork")
		if code != 0 {
			t.Errorf("main branch with --remote: Code = %d, wantCode = 0", code)
		}
		mock.AssertOutput(t, "main branch with --remote: Stdout", stdout, forkURL+"\n")
		mock.AssertOutput(t, "main branch with --remote: Stderr", stderr, "")
	})

	t.Run("fork branch", func(t *testing.T) {
		checkout("fork")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, clonePath, "", "url")
		if code != 0 {
			t.Errorf("feature branch: Code = %d, wantCode = 0", code)
		}
		mock.AssertOutput(t, "feature branch: Stdout", stdout, forkURL+"\n")
		mock.AssertOutput(t, "feature branch: Stderr", stderr, "")
	})

	t.Run("fork branch with --remote", func(t *testing.T) {
		checkout("fork")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, clonePath, "", "url", "--remote", "origin")
		if code != 0 {
			t.Errorf("feature branch with --remote: Code = %d, wantCode = 0", code)
		}
		mock.AssertOutput(t, "feature branch with --remote: Stdout", stdout, mainURL+"\n")
		mock.AssertOutput(t, "feature branch with --remote: Stderr", stderr, "")
	})

	t.Run("branch without remote", func(t *testing.T) {
		checkout("no_remote")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, clonePath, "", "url")
		if code != 0 {
			t.Errorf("feature branch with --remote: Code = %d, wantCode = 0", code)
		}
		mock.AssertOutput(t, "feature branch with --remote: Stdout", stdout, mainURL+"\n")
		mock.AssertOutput(t, "feature branch with --remote: Stderr", stderr, "")
	})

	t.Run("branch without remote with --remote", func(t *testing.T) {
		checkout("no_remote")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, clonePath, "", "url", "--remote", "fork")
		if code != 0 {
			t.Errorf("feature branch with --remote: Code = %d, wantCode = 0", code)
		}
		mock.AssertOutput(t, "feature branch with --remote: Stdout", stdout, forkURL+"\n")
		mock.AssertOutput(t, "feature branch with --remote: Stderr", stderr, "")
	})
}
