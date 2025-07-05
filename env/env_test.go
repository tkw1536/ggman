package env_test

//spellchecker:words path filepath reflect testing github ggman internal testutil pkglib testlib
import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/env"
	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/internal/path"
	"github.com/tkw1536/ggman/internal/testutil"
	"go.tkw01536.de/pkglib/testlib"
)

//spellchecker:words GGNORM GGROOT CANFILE worktree

func TestEnv_LoadDefaultRoot(t *testing.T) {
	t.Parallel()

	// noProjectsDir does not have a 'Projects' subdirectory
	noProjectsDir := testlib.TempDirAbs(t)
	missingProjectsDir := filepath.Join(noProjectsDir, "Projects")

	// withProjectsDir has a 'Projects' subdirectory
	withProjectsDir := testlib.TempDirAbs(t)
	existingProjectsDir := filepath.Join(withProjectsDir, "Projects")

	if err := os.Mkdir(existingProjectsDir, 0750); err != nil {
		panic(err)
	}

	// noExistsDir doesn't exist
	noExistsDir := filepath.Join(testlib.TempDirAbs(t), "noExist")

	tests := []struct {
		name     string
		vars     env.Variables
		wantRoot string
		wantErr  bool
	}{
		{"GGROOT exists", env.Variables{GGROOT: noProjectsDir}, noProjectsDir, false},
		{"GGROOT not exists", env.Variables{GGROOT: noExistsDir}, noExistsDir, false},

		{"GGROOT unset, HOME unset", env.Variables{}, "", true},

		{"GGROOT unset, HOME/Projects exists", env.Variables{HOME: noProjectsDir}, missingProjectsDir, false},
		{"GGROOT unset, HOME/Projects not exists", env.Variables{HOME: withProjectsDir}, existingProjectsDir, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := env.Env{
				Vars: tt.vars,
			}
			if err := env.LoadDefaultRoot(); (err != nil) != tt.wantErr {
				t.Errorf("Env.LoadDefaultRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
			if env.Root != tt.wantRoot {
				t.Errorf("Env.LoadDefaultRoot() root = %v, wantRoot %v", env.Root, tt.wantRoot)
			}
		})
	}
}

func TestEnv_LoadDefaultCANFILE(t *testing.T) {
	t.Parallel()

	// defaultCanFile is the default CanFile
	var defaultCanFile env.CanFile
	defaultCanFile.ReadDefault()

	// sampleCanFile is the sample CanFile]
	const canLineContent = "sample@^:$.git"
	var sampleCanFile env.CanFile = []env.CanLine{{"", canLineContent}}

	// emptyDir is an empty directory without a canFile
	emptyDir := testlib.TempDirAbs(t)
	noGGMANFile := filepath.Join(emptyDir, ".ggman")

	// dDir is a directory with a '.ggman' file
	dDir := testlib.TempDirAbs(t)
	GGMANFile := filepath.Join(dDir, ".ggman")
	if err := os.WriteFile(GGMANFile, []byte(canLineContent), 0600); err != nil {
		panic(err)
	}

	tests := []struct {
		name        string
		vars        env.Variables
		wantErr     bool
		wantCanfile env.CanFile
	}{
		{"loading from existing path", env.Variables{CANFILE: GGMANFile}, false, sampleCanFile},

		{"loading from home", env.Variables{HOME: dDir}, false, sampleCanFile},
		{"loading from home because of failure", env.Variables{CANFILE: noGGMANFile, HOME: dDir}, false, sampleCanFile},

		{"loading non-existing path", env.Variables{CANFILE: noGGMANFile}, false, defaultCanFile},
		{"loading non-existing home", env.Variables{HOME: emptyDir}, false, defaultCanFile},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := env.Env{
				Vars: tt.vars,
			}
			if _, err := env.LoadDefaultCANFILE(); (err != nil) != tt.wantErr {
				t.Errorf("Env.LoadDefaultCANFILE() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(env.CanFile, tt.wantCanfile) {
				t.Errorf("Env.LoadDefaultCANFILE() CanFile = %v, wantCanFile %v", env.CanFile, tt.wantCanfile)
			}
		})
	}
}

func TestEnv_Local_Exact(t *testing.T) {
	t.Parallel()

	root := testlib.TempDirAbs(t)

	// make the 'HELLO' directory, to ensure that it already exists
	if err := os.MkdirAll(filepath.Join(root, "server.com", "HELLO"), os.ModePerm|os.ModeDir); err != nil {
		panic(err)
	}

	tests := []struct {
		name   string
		GGNORM string
		want   string
	}{
		// smart
		{"git@github.com/user/repo", "smart", filepath.Join(root, "github.com", "user", "repo")},
		{"https://github.com/user/repo", "smart", filepath.Join(root, "github.com", "user", "repo")},
		{"ssh://git@github.com/hello/world", "smart", filepath.Join(root, "github.com", "hello", "world")},
		{"user@server.com:repo", "smart", filepath.Join(root, "server.com", "user", "repo")},
		{"ssh://user@server.com:1234/repo", "smart", filepath.Join(root, "server.com", "user", "repo")},

		{"ssh://server.com/hello/world", "smart", filepath.Join(root, "server.com", "HELLO", "world")}, // using existing case

		// exact
		{"git@github.com/user/repo", "exact", filepath.Join(root, "github.com", "user", "repo")},
		{"https://github.com/user/repo", "exact", filepath.Join(root, "github.com", "user", "repo")},
		{"ssh://git@github.com/hello/world", "exact", filepath.Join(root, "github.com", "hello", "world")},
		{"user@server.com:repo", "exact", filepath.Join(root, "server.com", "user", "repo")},
		{"ssh://user@server.com:1234/repo", "exact", filepath.Join(root, "server.com", "user", "repo")},

		{"ssh://server.com/hello/world", "exact", filepath.Join(root, "server.com", "hello", "world")}, // don't use existing case

		// fold
		{"git@github.com/user/repo", "fold", filepath.Join(root, "github.com", "user", "repo")},
		{"https://github.com/user/repo", "fold", filepath.Join(root, "github.com", "user", "repo")},
		{"ssh://git@github.com/hello/world", "fold", filepath.Join(root, "github.com", "hello", "world")},
		{"user@server.com:repo", "fold", filepath.Join(root, "server.com", "user", "repo")},
		{"ssh://user@server.com:1234/repo", "fold", filepath.Join(root, "server.com", "user", "repo")},

		{"ssh://server.com/hello/world", "fold", filepath.Join(root, "server.com", "HELLO", "world")}, // using existing case
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := env.Env{
				Root: root,
				Vars: env.Variables{
					GGNORM: tt.GGNORM,
				},
			}

			got, gotErr := e.Local(env.ParseURL(tt.name))

			if gotErr != nil {
				t.Errorf("Env.Local() err = %v, want err = nil", gotErr)
			}

			if got != tt.want {
				t.Errorf("Env.Local() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_At(t *testing.T) {
	t.Parallel()

	root := testlib.TempDirAbs(t)

	// group/repo contains a repository
	group := filepath.Join(root, "group")
	repo := filepath.Join(group, "repo")
	if err := os.MkdirAll(repo, 0750); err != nil {
		panic(err)
	}
	if testutil.NewTestRepoAt(repo, "") == nil {
		panic("Failed to create test repository")
	}

	// sub is a path inside the repository
	sub := filepath.Join(repo, "sub")
	if err := os.MkdirAll(sub, 0750); err != nil {
		panic(err)
	}

	tests := []struct {
		name         string
		path         string
		wantRepo     string
		wantWorktree string
		wantErr      bool
	}{
		{"no repository at root", root, "", "", true},
		{"no repository in group root", group, "", "", true},
		{"repository in repo", repo, repo, ".", false},
		{"repository in repo/sub", sub, repo, "sub", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := env.Env{
				Git:  git.NewGitFromPlumbing(nil, ""),
				Root: root,
			}
			gotRepo, gotWorktree, err := env.At(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.At() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("Env.At() gotRepo = %v, want %v", gotRepo, tt.wantRepo)
			}
			if gotWorktree != tt.wantWorktree {
				t.Errorf("Env.At() gotWorktree = %v, want %v", gotWorktree, tt.wantWorktree)
			}
		})
	}
}

func TestEnv_AtRoot(t *testing.T) {
	t.Parallel()

	root := testlib.TempDirAbs(t)

	// group/repo contains a repository
	group := filepath.Join(root, "group")
	repo := filepath.Join(group, "repo")
	if err := os.MkdirAll(repo, 0750); err != nil {
		panic(err)
	}
	if testutil.NewTestRepoAt(repo, "") == nil {
		panic("Failed to create test repository")
	}

	// sub is a path inside the repository
	sub := filepath.Join(repo, "sub")
	if err := os.MkdirAll(sub, 0750); err != nil {
		panic(err)
	}

	tests := []struct {
		name     string
		path     string
		wantRepo string
		wantErr  bool
	}{
		{"no repository root at root", root, "", false},
		{"no repository root in group root", group, "", false},
		{"repository root in repo", repo, repo, false},
		{"no repository root in repo/sub", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := env.Env{
				Git:  git.NewGitFromPlumbing(nil, ""),
				Root: root,
			}
			gotRepo, err := env.AtRoot(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Env.AtRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("Env.AtRoot() gotRepo = %v, want %v", gotRepo, tt.wantRepo)
			}
		})
	}
}

func TestEnv_ScanRepos(t *testing.T) {
	t.Parallel()

	root := testlib.TempDirAbs(t)

	// make a dir with parents and turn it into git
	makeGit := func(s string) {
		pth := filepath.Join(root, s)
		err := os.MkdirAll(pth, 0750)
		if err != nil {
			panic(err)
		}
		if testutil.NewTestRepoAt(pth, s) == nil {
			panic("NewTestRepoAt() returned nil")
		}
	}

	makeGit(filepath.Join("a", "aa", "aaa"))
	makeGit(filepath.Join("a", "aa", "aab"))
	makeGit(filepath.Join("a", "aa", "aac"))
	makeGit(filepath.Join("a", "ab", "aba"))
	makeGit(filepath.Join("a", "ab", "abb"))
	makeGit(filepath.Join("a", "ab", "abc"))
	makeGit(filepath.Join("a", "ac", "aca"))
	makeGit(filepath.Join("a", "ac", "acb"))
	makeGit(filepath.Join("a", "ac", "acc"))

	// utility to remove root from all the paths
	trimPath := func(path string) string {
		t, err := filepath.Rel(root, path)
		if err != nil {
			return path
		}
		return t
	}
	trimAll := func(paths []string) {
		for idx := range paths {
			paths[idx] = trimPath(paths[idx])
		}
	}

	tests := []struct {
		name   string
		Filter string
		want   []string
	}{
		{
			"all repos", "", []string{
				"a/aa/aaa",
				"a/aa/aab",
				"a/aa/aac",
				"a/ab/aba",
				"a/ab/abb",
				"a/ab/abc",
				"a/ac/aca",
				"a/ac/acb",
				"a/ac/acc",
			},
		},

		{
			"'aa' only", "aa", []string{
				"a/aa/aaa",
				"a/aa/aab",
				"a/aa/aac",
			},
		},

		{
			"'a/a*/*a' only", "a/a*/*a", []string{
				"a/aa/aaa",
				"a/ab/aba",
				"a/ac/aca",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := env.Env{
				Root: root,
				Git:  git.NewGitFromPlumbing(nil, ""),

				Filter: env.NewPatternFilter(tt.Filter, false),
			}
			got, err := env.ScanRepos(root, true)
			wantErr := false
			if (err != nil) != wantErr {
				t.Errorf("Env.ScanRepos() error = %v, wantErr %v", err, wantErr)
				return
			}
			trimAll(got)
			testutil.ToOSPaths(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.ScanRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_ScanRepos_fuzzy(t *testing.T) {
	t.Parallel()

	root := testlib.TempDirAbs(t)

	// make a dir with parents and turn it into git
	makeGit := func(s string) {
		pth := filepath.Join(root, s)
		err := os.MkdirAll(pth, 0750)
		if err != nil {
			panic(err)
		}
		if testutil.NewTestRepoAt(pth, s) == nil {
			panic("NewTestRepoAt() returned nil")
		}
	}

	makeGit("abc") // matches the filter 'bc' with a score of 0.66, but lexicographically first
	makeGit("bc")  // matches the filter 'bc' with a score of 1, but lexicographically last

	// utility to remove root from all the paths
	trimPath := func(path string) string {
		t, err := filepath.Rel(root, path)
		if err != nil {
			return path
		}
		return t
	}
	trimAll := func(paths []string) {
		for idx := range paths {
			paths[idx] = trimPath(paths[idx])
		}
	}

	tests := []struct {
		name   string
		Filter string
		want   []string
	}{
		{
			"all repos", "", []string{
				"abc",
				"bc",
			},
		},
		{
			"filter 'bc', sorted by priority", "bc", []string{
				"bc",
				"abc",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			env := env.Env{
				Root: root,
				Git:  git.NewGitFromPlumbing(nil, ""),

				Filter: env.NewPatternFilter(tt.Filter, true),
			}
			got, err := env.ScanRepos(root, true)
			wantErr := false
			if (err != nil) != wantErr {
				t.Errorf("Env.ScanRepos() error = %v, wantErr %v", err, wantErr)
				return
			}
			trimAll(got)
			testutil.ToOSPaths(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.ScanRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_Normalization(t *testing.T) {
	t.Parallel()

	tests := []struct {
		GGNORM string
		want   path.Normalization
	}{
		{"fold", path.FoldNorm},
		{"smart", path.FoldPreferExactNorm},
		{"exact", path.NoNorm},
		{"", path.FoldPreferExactNorm},
		{"this-norm-doesn't-exist", path.FoldPreferExactNorm},
	}
	for _, tt := range tests {
		t.Run(tt.GGNORM, func(t *testing.T) {
			t.Parallel()

			env := env.Env{
				Vars: env.Variables{
					GGNORM: tt.GGNORM,
				},
			}
			if got := env.Normalization(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.Normalization() = %v, want %v", got, tt.want)
			}
		})
	}
}
