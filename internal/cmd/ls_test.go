package cmd_test

//spellchecker:words encoding json path filepath slices testing essio shellescape github config plumbing ggman internal mockenv testutil
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"al.essio.dev/pkg/shellescape"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"go.tkw01536.de/ggman/internal/cmd"
	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/ggman/internal/testutil"
)

//spellchecker:words workdir GGROOT wrld tparallel paralleltest worktree

var testInputFile = `
; this and the following lines are ignored
# gitlab.com/hello/world
` + filepath.Join("server.com", "user", "repo") + `

// blank lines too

https://github.com/hello/world.git

`

func TestCommandLs(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	ghHelloWorld := mock.Clone(t.Context(), "https://github.com/hello/world.git", "github.com", "hello", "world")
	serverRepo := mock.Clone(t.Context(), "user@server.com/repo", "server.com", "user", "repo")
	glHelloWorld := mock.Clone(t.Context(), "https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	inputFile := mock.Resolve("file.txt")
	if err := os.WriteFile(inputFile, []byte(testInputFile), 0600); err != nil {
		panic(err)
	}

	// make glHelloWorldDirty
	if err := os.WriteFile(filepath.Join(glHelloWorld, "dirty"), []byte{}, 0600); err != nil {
		panic(err)
	}

	glHelloDir := filepath.Join(glHelloWorld, "..")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{

		{
			"list all repositories",
			"",
			[]string{"ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"list dirty and clean repositories",
			"",
			[]string{"--dirty", "--clean", "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"list dirty repositories only",
			"",
			[]string{"--dirty", "ls"},

			0,
			"${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list clean repositories only",
			"",
			[]string{"--clean", "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"list all repositories with exit code",
			"",
			[]string{"ls", "--exit-code"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"list all repositories with one",
			"",
			[]string{"ls", "--one"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},

		{
			"list all repositories with specific count",
			"",
			[]string{"ls", "--count", "2"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list all repositories with higher than available count",
			"",
			[]string{"ls", "--count", "5"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n${GGROOT server.com user repo}\n",

			"",
		},

		{
			"don't support both one and count at the same time",
			"",
			[]string{"ls", "--one", "--count", "2"},

			4,
			"",

			`only one of "--one" and "--count" may be provided` + "\n",
		},

		{
			"don't support negative limit",
			"",
			[]string{"ls", "--count", "-1"},

			4,
			"",

			`"--count" may not be negative` + "\n",
		},

		{
			"list only hello/world repositories",
			"",
			[]string{"--for", "hello/world", "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list only clean hello/world repositories",
			"",
			[]string{"--for", "hello/world", "--clean", "ls"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},

		{
			"list repositories fuzzy",
			"",
			[]string{"--for", "wrld", "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list repositories with start flags",
			"",
			[]string{"--for", "^github.com", "ls"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},

		{
			"list repositories with scores",
			"",
			[]string{"--for", "wrld", "ls", "--scores"},

			0,
			"0.900000 ${GGROOT github.com hello world}\n0.900000 ${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list repositories non-fuzzy",
			"",
			[]string{"--no-fuzzy-filter", "--for", "wrld", "ls"},

			0,
			"",

			"",
		},

		{
			"list non-existing repositories",
			"",
			[]string{"--for", "does/not/exist", "ls"},

			0,
			"",

			"",
		},

		{
			"list non-existing repositories with exit code",
			"",
			[]string{"--for", "does/not/exist", "ls", "--exit-code"},

			1,
			"",

			"",
		},

		{
			"list only current repository (github.com hello world)",
			ghHelloWorld,
			[]string{"--here", "ls"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},

		{
			"list only current repository (server.com user repo)",
			serverRepo,
			[]string{"--here", "ls"},

			0,
			"${GGROOT server.com user repo}\n",

			"",
		},
		{
			"list only current repository (gitlab.com hello world)",
			glHelloWorld,
			[]string{"--here", "ls"},

			0,
			"${GGROOT gitlab.com hello world}\n",

			"",
		},
		{
			"list an absolute path",
			serverRepo,
			[]string{"--for", ghHelloWorld, "ls"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},

		{
			"list an absolute path with --path",
			serverRepo,
			[]string{"--path", ghHelloWorld, "ls"},

			0,
			"${GGROOT github.com hello world}\n",

			"",
		},

		{
			"list a relative path",
			glHelloDir,
			[]string{"--for", filepath.Join(".", "world"), "ls"},

			0,
			"${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list a relative path with --path",
			glHelloDir,
			[]string{"--path", filepath.Join(".", "world"), "ls"},

			0,
			"${GGROOT gitlab.com hello world}\n",

			"",
		},

		{
			"list multiple paths with --path",
			glHelloDir,
			[]string{"--path", filepath.Join(".", "world"), "--path", ghHelloWorld, "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT gitlab.com hello world}\n",

			"",
		},
		{
			"list from input file",
			"",
			[]string{"--from-file", inputFile, "ls"},

			0,
			"${GGROOT github.com hello world}\n${GGROOT server.com user repo}\n",

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

//nolint:paralleltest
func TestCommandLsPriorities(t *testing.T) {
	mock := mockenv.NewMockEnv(t)

	// Create a repository called "needle_in_haystack" with a single remote
	// This repo should NOT be matched first when searching for "needle"
	mock.Clone(t.Context(), "https://github.com/example/needle_in_haystack.git", "github.com", "example", "needle_in_haystack")

	// Create a repository called "needle" with two branches pointing to two remotes:
	// - "needle_match" branch => remote that contains "needle" exactly
	// - "needle_no_match" branch => remote that does NOT contain "needle"
	needleMatchRemote := "https://github.com/example/needle.git"
	needleNoMatchRemote := "https://github.com/fork/repo.git"

	mock.Register(needleMatchRemote)
	_, needleNoMatchRemoteURLs := mock.Register(needleNoMatchRemote)
	needlePath := mock.Install(t.Context(), needleMatchRemote, "github.com", "example", "needle")

	// Open the cloned repository to manipulate it
	repo, err := git.PlainOpen(needlePath)
	if err != nil {
		panic(err)
	}

	// Add a second remote "fork" pointing to a remote without "needle" in the name
	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "fork",
		URLs: []string{needleNoMatchRemoteURLs[0]},
	}); err != nil {
		panic(err)
	}

	// Get the worktree to create and checkout branches
	wt, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	// Create two branches:
	// - needle_match => pointing to origin (which has "needle" in the URL)
	// - needle_no_match => pointing to fork (which does NOT have "needle" in the URL)
	testutil.CommitTestFiles(repo)
	testutil.CreateTrackingBranch(repo, "origin", "needle_match", "main")
	testutil.CreateTrackingBranch(repo, "fork", "needle_no_match", "main")

	// helper function so we can checkout specific branches
	checkout := func(branch string) {
		if err := wt.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(branch),
		}); err != nil {
			panic(err)
		}
	}

	// Test: When searching for "needle", the "needle" repo should always come first,
	// regardless of which branch is checked out, because it has a remote that matches "needle" exactly

	t.Run("needle_match branch checked out", func(t *testing.T) {
		checkout("needle_match")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "--for", "needle", "ls")
		if code != 0 {
			t.Errorf("Code = %d, wantCode = 0", code)
		}
		// needle should come first because it has an exact match
		mock.AssertOutput(t, "Stdout", stdout, "${GGROOT github.com example needle}\n${GGROOT github.com example needle_in_haystack}\n")
		mock.AssertOutput(t, "Stderr", stderr, "")
	})

	t.Run("needle_no_match branch checked out", func(t *testing.T) {
		checkout("needle_no_match")

		code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "--for", "needle", "ls")
		if code != 0 {
			t.Errorf("Code = %d, wantCode = 0", code)
		}
		// needle should still come first because it has a remote (origin) that matches "needle" exactly
		// even though the currently checked out branch points to a remote without "needle"
		mock.AssertOutput(t, "Stdout", stdout, "${GGROOT github.com example needle}\n${GGROOT github.com example needle_in_haystack}\n")
		mock.AssertOutput(t, "Stderr", stderr, "")
	})
}

func TestCommandLsRemote(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	mock.Clone(t.Context(), "https://github.com/hello/world.git", "github.com", "hello", "world")
	mock.Clone(t.Context(), "user@server.com/repo", "server.com", "user", "repo")
	mock.Clone(t.Context(), "https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{
		{
			"list remotes of all repositories",
			"",
			[]string{"ls", "--remote"},

			0,
			"https://github.com/hello/world.git\nhttps://gitlab.com/hello/world.git\nuser@server.com/repo\n",

			"",
		},

		{
			"list canonical remotes of all repositories",
			"",
			[]string{"ls", "--remote", "--canonical"},

			0,
			"git@github.com:hello/world.git\ngit@gitlab.com:hello/world.git\ngit@server.com:user/repo.git\n",

			"",
		},

		{
			"list remotes only hello/world repositories",
			"",
			[]string{"--for", "hello/world", "ls", "--remote"},

			0,
			"https://github.com/hello/world.git\nhttps://gitlab.com/hello/world.git\n",

			"",
		},

		{
			"list canonical remotes only hello/world repositories",
			"",
			[]string{"--for", "hello/world", "ls", "--remote", "--canonical"},

			0,
			"git@github.com:hello/world.git\ngit@gitlab.com:hello/world.git\n",

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

func TestLsCommandRelative(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	ghHelloWorld := mock.Clone(t.Context(), "https://github.com/hello/world.git", "github.com", "hello", "world")
	serverRepo := mock.Clone(t.Context(), "user@server.com/repo", "server.com", "user", "repo")
	glHelloWorld := mock.Clone(t.Context(), "https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	inputFile := mock.Resolve("file.txt")
	if err := os.WriteFile(inputFile, []byte(testInputFile), 0600); err != nil {
		panic(err)
	}

	// make glHelloWorldDirty
	if err := os.WriteFile(filepath.Join(glHelloWorld, "dirty"), []byte{}, 0600); err != nil {
		panic(err)
	}

	glHelloDir := filepath.Join(glHelloWorld, "..")

	// Relative paths for expected output
	ghHelloWorldRel := filepath.Join("github.com", "hello", "world")
	glHelloWorldRel := filepath.Join("gitlab.com", "hello", "world")
	serverRepoRel := filepath.Join("server.com", "user", "repo")

	tests := []struct {
		name    string
		workdir string
		args    []string

		wantCode   uint8
		wantStdout string
		wantStderr string
	}{

		{
			"list all repositories",
			"",
			[]string{"ls", "--relative"},

			0,
			fmt.Sprintf("%s\n%s\n%s\n", ghHelloWorldRel, glHelloWorldRel, serverRepoRel),

			"",
		},

		{
			"list dirty and clean repositories",
			"",
			[]string{"--dirty", "--clean", "ls", "--relative"},

			0,
			fmt.Sprintf("%s\n%s\n%s\n", ghHelloWorldRel, glHelloWorldRel, serverRepoRel),

			"",
		},

		{
			"list dirty repositories only",
			"",
			[]string{"--dirty", "ls", "--relative"},

			0,
			glHelloWorldRel + "\n",

			"",
		},

		{
			"list clean repositories only",
			"",
			[]string{"--clean", "ls", "--relative"},

			0,
			fmt.Sprintf("%s\n%s\n", ghHelloWorldRel, serverRepoRel),

			"",
		},

		{
			"list all repositories with exit code",
			"",
			[]string{"ls", "--relative", "--exit-code"},

			0,
			fmt.Sprintf("%s\n%s\n%s\n", ghHelloWorldRel, glHelloWorldRel, serverRepoRel),

			"",
		},

		{
			"list all repositories with one",
			"",
			[]string{"ls", "--relative", "--one"},

			0,
			ghHelloWorldRel + "\n",

			"",
		},

		{
			"list all repositories with specific count",
			"",
			[]string{"ls", "--relative", "--count", "2"},

			0,
			fmt.Sprintf("%s\n%s\n", ghHelloWorldRel, glHelloWorldRel),

			"",
		},

		{
			"list all repositories with higher than available count",
			"",
			[]string{"ls", "--relative", "--count", "5"},

			0,
			fmt.Sprintf("%s\n%s\n%s\n", ghHelloWorldRel, glHelloWorldRel, serverRepoRel),

			"",
		},

		{
			"don't support both one and count at the same time",
			"",
			[]string{"ls", "--relative", "--one", "--count", "2"},

			4,
			"",

			`only one of "--one" and "--count" may be provided` + "\n",
		},

		{
			"list only hello/world repositories",
			"",
			[]string{"--for", "hello/world", "ls", "--relative"},

			0,
			fmt.Sprintf("%s\n%s\n", ghHelloWorldRel, glHelloWorldRel),

			"",
		},

		{
			"list only clean hello/world repositories",
			"",
			[]string{"--for", "hello/world", "--clean", "ls", "--relative"},

			0,
			ghHelloWorldRel + "\n",

			"",
		},

		{
			"list repositories fuzzy",
			"",
			[]string{"--for", "wrld", "ls", "--relative"},

			0,
			fmt.Sprintf("%s\n%s\n", ghHelloWorldRel, glHelloWorldRel),

			"",
		},

		{
			"list repositories with start flags",
			"",
			[]string{"--for", "^github.com", "ls", "--relative"},

			0,
			ghHelloWorldRel + "\n",

			"",
		},

		{
			"list repositories with scores",
			"",
			[]string{"--for", "wrld", "ls", "--relative", "--scores"},

			0,
			fmt.Sprintf("0.900000 %s\n0.900000 %s\n", ghHelloWorldRel, glHelloWorldRel),

			"",
		},

		{
			"list repositories non-fuzzy",
			"",
			[]string{"--no-fuzzy-filter", "--for", "wrld", "ls", "--relative"},

			0,
			"",

			"",
		},

		{
			"list non-existing repositories",
			"",
			[]string{"--for", "does/not/exist", "ls", "--relative"},

			0,
			"",

			"",
		},

		{
			"list non-existing repositories with exit code",
			"",
			[]string{"--for", "does/not/exist", "ls", "--relative", "--exit-code"},

			1,
			"",

			"",
		},

		{
			"list only current repository (github.com hello world)",
			ghHelloWorld,
			[]string{"--here", "ls", "--relative"},

			0,
			ghHelloWorldRel + "\n",

			"",
		},

		{
			"list only current repository (server.com user repo)",
			serverRepo,
			[]string{"--here", "ls", "--relative"},

			0,
			serverRepoRel + "\n",

			"",
		},
		{
			"list only current repository (gitlab.com hello world)",
			glHelloWorld,
			[]string{"--here", "ls", "--relative"},

			0,
			glHelloWorldRel + "\n",

			"",
		},
		{
			"list an absolute path",
			serverRepo,
			[]string{"--for", ghHelloWorld, "ls", "--relative"},

			0,
			ghHelloWorldRel + "\n",

			"",
		},

		{
			"list an absolute path with --path",
			serverRepo,
			[]string{"--path", ghHelloWorld, "ls", "--relative"},

			0,
			ghHelloWorldRel + "\n",

			"",
		},

		{
			"list a relative path",
			glHelloDir,
			[]string{"--for", filepath.Join(".", "world"), "ls", "--relative"},

			0,
			glHelloWorldRel + "\n",

			"",
		},

		{
			"list a relative path with --path",
			glHelloDir,
			[]string{"--path", filepath.Join(".", "world"), "ls", "--relative"},

			0,
			glHelloWorldRel + "\n",

			"",
		},

		{
			"list multiple paths with --path",
			glHelloDir,
			[]string{"--path", filepath.Join(".", "world"), "--path", ghHelloWorld, "ls", "--relative"},

			0,
			fmt.Sprintf("%s\n%s\n", ghHelloWorldRel, glHelloWorldRel),

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

func TestLsCommandJSON(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	ghHelloWorld := mock.Clone(t.Context(), "https://github.com/hello/world.git", "github.com", "hello", "world")
	serverRepo := mock.Clone(t.Context(), "user@server.com/repo", "server.com", "user", "repo")
	glHelloWorld := mock.Clone(t.Context(), "https://gitlab.com/hello/world.git", "gitlab.com", "hello", "world")

	// make glHelloWorld dirty
	if err := os.WriteFile(filepath.Join(glHelloWorld, "dirty"), []byte{}, 0600); err != nil {
		panic(err)
	}

	tests := []struct {
		name     string
		args     []string
		wantCode uint8
		want     []cmd.Repo
	}{
		{
			"list all repositories with paths",
			[]string{"ls", "--json"},
			0,
			[]cmd.Repo{
				{Path: ghHelloWorld, Score: 1},
				{Path: glHelloWorld, Score: 1},
				{Path: serverRepo, Score: 1},
			},
		},
		{
			"list all repositories with relative paths",
			[]string{"ls", "--json", "--relative"},
			0,
			[]cmd.Repo{
				{Path: ghHelloWorld, Relative: filepath.Join("github.com", "hello", "world"), Score: 1},
				{Path: glHelloWorld, Relative: filepath.Join("gitlab.com", "hello", "world"), Score: 1},
				{Path: serverRepo, Relative: filepath.Join("server.com", "user", "repo"), Score: 1},
			},
		},
		{
			"list all repositories with remotes",
			[]string{"ls", "--json", "--remote"},
			0,
			[]cmd.Repo{
				{Path: ghHelloWorld, Remote: "https://github.com/hello/world.git", Score: 1},
				{Path: glHelloWorld, Remote: "https://gitlab.com/hello/world.git", Score: 1},
				{Path: serverRepo, Remote: "user@server.com/repo", Score: 1},
			},
		},
		{
			"list dirty repositories only",
			[]string{"--dirty", "ls", "--json"},
			0,
			[]cmd.Repo{
				{Path: glHelloWorld, Score: 1},
			},
		},
		{
			"list with limit",
			[]string{"ls", "--json", "--count", "2"},
			0,
			[]cmd.Repo{
				{Path: ghHelloWorld, Score: 1},
				{Path: glHelloWorld, Score: 1},
			},
		},
		{
			"list all repositories with canonical remotes",
			[]string{"ls", "--json", "--remote", "--canonical"},
			0,
			[]cmd.Repo{
				{Path: ghHelloWorld, Remote: "https://github.com/hello/world.git", Canonical: "git@github.com:hello/world.git", Score: 1},
				{Path: glHelloWorld, Remote: "https://gitlab.com/hello/world.git", Canonical: "git@gitlab.com:hello/world.git", Score: 1},
				{Path: serverRepo, Remote: "user@server.com/repo", Canonical: "git@server.com:user/repo.git", Score: 1},
			},
		},
		{
			"list all repositories with remote, canonical, and relative",
			[]string{"ls", "--json", "--remote", "--canonical", "--relative"},
			0,
			[]cmd.Repo{
				{Path: ghHelloWorld, Relative: filepath.Join("github.com", "hello", "world"), Remote: "https://github.com/hello/world.git", Canonical: "git@github.com:hello/world.git", Score: 1},
				{Path: glHelloWorld, Relative: filepath.Join("gitlab.com", "hello", "world"), Remote: "https://gitlab.com/hello/world.git", Canonical: "git@gitlab.com:hello/world.git", Score: 1},
				{Path: serverRepo, Relative: filepath.Join("server.com", "user", "repo"), Remote: "user@server.com/repo", Canonical: "git@server.com:user/repo.git", Score: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("Code = %d, wantCode = %d", code, tt.wantCode)
			}
			if stderr != "" {
				t.Errorf("Stderr = %q, want empty", stderr)
			}

			var got []cmd.Repo
			if err := json.Unmarshal([]byte(stdout), &got); err != nil {
				t.Fatalf("failed to unmarshal JSON output: %v", err)
			}

			if !slices.Equal(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandLsExport(t *testing.T) {
	t.Parallel()

	mock := mockenv.NewMockEnv(t)

	mock.Clone(t.Context(), "https://github.com/hello/world.git", "github.com", "hello", "world")
	mock.Clone(t.Context(), "user@server.com:user/repo", "server.com", "user", "repo")
	mock.Clone(t.Context(), "https://gitlab.com/org/project.git", "gitlab.com", "org", "project")

	// Build expected bash script using shellescape.Quote and filepath.Join
	ghPath := filepath.Join("github.com", "hello", "world")
	ghURL := "https://github.com/hello/world.git"

	serverPath := filepath.Join("server.com", "user", "repo")
	serverURL := "user@server.com:user/repo"

	glPath := filepath.Join("gitlab.com", "org", "project")
	glURL := "https://gitlab.com/org/project.git"

	wantStdout := "#!/bin/bash\nset -e\n\n# Generated by ggman export\n" +
		fmt.Sprintf("mkdir -p %s\n", shellescape.Quote(ghPath)) +
		fmt.Sprintf("git clone %s %s\n", shellescape.Quote(ghURL), shellescape.Quote(ghPath)) +
		fmt.Sprintf("mkdir -p %s\n", shellescape.Quote(glPath)) +
		fmt.Sprintf("git clone %s %s\n", shellescape.Quote(glURL), shellescape.Quote(glPath)) +
		fmt.Sprintf("mkdir -p %s\n", shellescape.Quote(serverPath)) +
		fmt.Sprintf("git clone %s %s\n", shellescape.Quote(serverURL), shellescape.Quote(serverPath))

	code, stdout, stderr := mock.Run(t, nil, cmd.NewCommand, "", "", "ls", "--export")
	if code != 0 {
		t.Errorf("Code = %d, wantCode = 0", code)
	}
	if stdout != wantStdout {
		t.Errorf("Stdout mismatch:\ngot:\n%s\nwant:\n%s", stdout, wantStdout)
	}
	if stderr != "" {
		t.Errorf("Stderr = %q, want empty", stderr)
	}
}
