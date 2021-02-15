package env

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tkw1536/ggman/git"
	"github.com/tkw1536/ggman/internal/path"
	"github.com/tkw1536/ggman/internal/testutil"
)

func TestEnv_LoadDefaultRoot(t *testing.T) {

	// nopdir does not have a 'Projects' subdirectory
	nopdir, cleanup := testutil.TempDir()
	pnopdir := filepath.Join(nopdir, "Projects")
	defer cleanup()

	// pdir has a 'Projects' subdirectory
	pdir, cleanup := testutil.TempDir()
	ppdir := filepath.Join(pdir, "Projects")
	defer cleanup()
	if err := os.Mkdir(ppdir, os.ModePerm); err != nil {
		panic(err)
	}

	// nodir doesn't exist
	nodir, cleanup := testutil.TempDir()
	cleanup()

	tests := []struct {
		name     string
		vars     Variables
		wantRoot string
		wantErr  bool
	}{
		{"GGROOT exists", Variables{GGROOT: nopdir}, nopdir, false},
		{"GGROOT not exists", Variables{GGROOT: nodir}, nodir, false},

		{"GGROOT unset, HOME unset", Variables{}, "", true},

		{"GGROOT unset, HOME/Projects exists", Variables{HOME: nopdir}, pnopdir, false},
		{"GGROOT unset, HOME/Projects not exists", Variables{HOME: pdir}, ppdir, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := &Env{
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
	// dfltCanFile is the default CanFile
	var dfltCanFile CanFile
	dfltCanFile.ReadDefault()

	// sampleCanFile is the sample CanFile]
	const canLineContent = "sample@^:$.git"
	var sampleCanFile CanFile = []CanLine{{"", canLineContent}}

	// edir is an empty directory without a canFile
	edir, cleanup := testutil.TempDir()
	noggmanfile := filepath.Join(edir, ".ggman")
	defer cleanup()

	// ddir is a directory with a '.ggman' file
	ddir, cleanup := testutil.TempDir()
	ggmanfile := filepath.Join(ddir, ".ggman")
	ioutil.WriteFile(ggmanfile, []byte(canLineContent), os.ModePerm)
	defer cleanup()

	tests := []struct {
		name        string
		vars        Variables
		wantErr     bool
		wantCanfile CanFile
	}{
		{"loading from existing path", Variables{CANFILE: ggmanfile}, false, sampleCanFile},

		{"loading from home", Variables{HOME: ddir}, false, sampleCanFile},
		{"loading from home because of failure", Variables{CANFILE: noggmanfile, HOME: ddir}, false, sampleCanFile},

		{"loading non-existing path", Variables{CANFILE: noggmanfile}, false, dfltCanFile},
		{"loading non-existing home", Variables{HOME: edir}, false, dfltCanFile},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := &Env{
				Vars: tt.vars,
			}
			if err := env.LoadDefaultCANFILE(); (err != nil) != tt.wantErr {
				t.Errorf("Env.LoadDefaultCANFILE() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(env.CanFile, tt.wantCanfile) {
				t.Errorf("Env.LoadDefaultCANFILE() CanFile = %v, wantCanFile %v", env.CanFile, tt.wantCanfile)
			}
		})
	}
}

func TestEnv_Local(t *testing.T) {
	root, cleanup := testutil.TempDir()
	defer cleanup()

	tests := []struct {
		name string
		want string
	}{
		{"git@github.com/user/repo", filepath.Join(root, "github.com", "user", "repo")},
		{"https://github.com/user/repo", filepath.Join(root, "github.com", "user", "repo")},
		{"ssh://git@github.com/hello/world", filepath.Join(root, "github.com", "hello", "world")},
		{"user@server.com:repo", filepath.Join(root, "server.com", "user", "repo")},
		{"ssh://user@server.com:1234/repo", filepath.Join(root, "server.com", "user", "repo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := Env{
				Root: root,
			}

			if got := env.Local(ParseURL(tt.name)); got != tt.want {
				t.Errorf("Env.Local() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnv_At(t *testing.T) {
	root, cleanup := testutil.TempDir()
	defer cleanup()

	// group/repo contains a repository
	group := filepath.Join(root, "group")
	repo := filepath.Join(group, "repo")
	if err := os.MkdirAll(repo, os.ModePerm); err != nil {
		panic(err)
	}
	if testutil.NewTestRepoAt(repo, "") == nil {
		panic("Failed to create test repository")
	}

	// sub is a path inside the repository
	sub := filepath.Join(repo, "sub")
	if err := os.MkdirAll(sub, os.ModePerm); err != nil {
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
			env := Env{
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

func TestEnv_ScanRepos(t *testing.T) {
	root, cleanup := testutil.TempDir()
	defer cleanup()

	// make a dir with parents and turn it into git
	mkgit := func(s string) {
		pth := filepath.Join(root, s)
		err := os.MkdirAll(pth, os.ModePerm)
		if err != nil {
			panic(err)
		}
		if testutil.NewTestRepoAt(pth, s) == nil {
			panic("NewTestRepoAt() returned nil")
		}
	}

	mkgit(filepath.Join("a", "aa", "aaa"))
	mkgit(filepath.Join("a", "aa", "aab"))
	mkgit(filepath.Join("a", "aa", "aac"))
	mkgit(filepath.Join("a", "ab", "aba"))
	mkgit(filepath.Join("a", "ab", "abb"))
	mkgit(filepath.Join("a", "ab", "abc"))
	mkgit(filepath.Join("a", "ac", "aca"))
	mkgit(filepath.Join("a", "ac", "acb"))
	mkgit(filepath.Join("a", "ac", "acc"))

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
			env := Env{
				Root: root,
				Git:  git.NewGitFromPlumbing(nil, ""),

				Filter: NewPatternFilter(tt.Filter),
			}
			got, err := env.ScanRepos(root)
			wantErr := false
			if (err != nil) != wantErr {
				t.Errorf("Env.ScanRepos() error = %v, wantErr %v", err, wantErr)
				return
			}
			trimAll(got)
			path.ToOSPaths(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Env.ScanRepos() = %v, want %v", got, tt.want)
			}
		})
	}
}
