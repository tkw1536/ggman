//spellchecker:words mockenv
package mockenv_test

//spellchecker:words reflect testing github config ggman internal gggit mockenv testutil pkglib stream testlib
import (
	"reflect"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	gggit "go.tkw01536.de/ggman/internal/git"
	"go.tkw01536.de/ggman/internal/mockenv"
	"go.tkw01536.de/ggman/internal/testutil"
	"go.tkw01536.de/pkglib/stream"
	"go.tkw01536.de/pkglib/testlib"
)

//spellchecker:words gogit tparallel paralleltest

func TestDevPlumbing_Forward(t *testing.T) {
	t.Parallel()

	mp := &mockenv.DevPlumbing{URLMap: make(map[string]string)}
	mp.URLMap["forward-a"] = "backward-a"
	mp.URLMap["forward-b"] = "backward-b"

	tests := []struct {
		name string
		url  string
		want string
	}{
		{"forward-a", "forward-a", "backward-a"},
		{"forward-b", "forward-b", "backward-b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := mp.Forward(tt.url); got != tt.want {
				t.Errorf("DevPlumbing.Forward() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("does-not-exist", func(t *testing.T) {
		t.Parallel()

		if panics, _ := testlib.DoesPanic(func() { mp.Forward("does-not-exist") }); !panics {
			t.Errorf("Expected DevPlumbing.Forward() to panic")
		}
	})
}

func TestDevPlumbing_Backward(t *testing.T) {
	t.Parallel()

	mp := &mockenv.DevPlumbing{URLMap: make(map[string]string)}
	mp.URLMap["forward-a"] = "backward-a"
	mp.URLMap["forward-b"] = "backward-b"

	tests := []struct {
		name string
		url  string
		want string
	}{
		{"backward-a", "backward-a", "forward-a"},
		{"backward-b", "backward-b", "forward-b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := mp.Backward(tt.url); got != tt.want {
				t.Errorf("DevPlumbing.Backward() = %v, want %v", got, tt.want)
			}
		})
	}

	t.Run("does-not-exist", func(t *testing.T) {
		t.Parallel()

		if panics, _ := testlib.DoesPanic(func() { mp.Backward("does-not-exist") }); !panics {
			t.Errorf("Expected DevPlumbing.Backward() to panic")
		}
	})
}

func Test_DevPlumbing_GetRemotes(t *testing.T) {
	t.Parallel()

	// This test has been adapted from Test_goGit_GetRemotes.

	mp := &mockenv.DevPlumbing{
		Plumbing: gggit.NewPlumbing(),
		URLMap:   make(map[string]string),
	}

	// For this test we have three repositories:
	// 'remote' <- 'cloneA' <- 'cloneB'
	// 'cloneA' has an origin remote pointing to 'remote'.
	// 'cloneB' has an origin remote pointing to 'cloneA'.
	// 'cloneB' also gets an upstream remote pointing to 'remote'.

	// In the mapping we set up:

	// 'git://server.com:2222/cloneA.git' -> cloneA
	// 'git://server.com:2222/remote.git' -> remote

	// GetRemotes() should return the mapped remotes.

	// create an initial remote repository, and add a new bogus commit to it.
	remote, repo := testutil.NewTestRepo(t)
	testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})

	// clone the remote repository into 'cloneA'.
	// This will create an origin remote pointing to the remote.
	cloneA := testlib.TempDirAbs(t)
	if _, err := git.PlainClone(cloneA, false, &git.CloneOptions{URL: remote}); err != nil {
		panic(err)
	}

	// clone the 'cloneA' repository into 'cloneB'.
	// This will create an origin remote pointing to 'cloneA'
	cloneB := testlib.TempDirAbs(t)
	repo, err := git.PlainClone(cloneB, false, &git.CloneOptions{URL: cloneA})
	if err != nil {
		panic(err)
	}

	// create an 'upstream' remote to point to 'remote'.
	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{remote},
	}); err != nil {
		panic(err)
	}

	mappedCloneA := "git://server.com:2222/cloneA.git"
	mp.URLMap[mappedCloneA] = cloneA

	mappedRemote := "git://server.com:2222/remote.git"
	mp.URLMap[mappedRemote] = remote

	t.Run("GetRemotes() on a repository with a single remote", func(t *testing.T) {
		t.Parallel()

		ggRepoObject, isRepo := mp.IsRepository(t.Context(), cloneA)
		if !isRepo {
			panic("IsRepository() failed")
		}

		wantRemotes := map[string][]string{
			"origin": {mappedRemote},
		}
		remotes, err := mp.GetRemotes(t.Context(), cloneA, ggRepoObject)
		if err != nil {
			t.Error("GetRemotes() got err != nil, want err == nil")
		}

		if !reflect.DeepEqual(remotes, wantRemotes) {
			t.Errorf("GetRemotes() got remotes = %v, want remotes = %v", remotes, wantRemotes)
		}
	})

	t.Run("GetRemotes() on a repository with more than one remote", func(t *testing.T) {
		t.Parallel()

		ggRepoObject, isRepo := mp.IsRepository(t.Context(), cloneB)
		if !isRepo {
			panic("IsRepository() failed")
		}

		wantRemotes := map[string][]string{
			"upstream": {mappedRemote},
			"origin":   {mappedCloneA},
		}
		remotes, err := mp.GetRemotes(t.Context(), cloneB, ggRepoObject)
		if err != nil {
			t.Error("GetRemotes() got err != nil, want err == nil")
		}

		if !reflect.DeepEqual(remotes, wantRemotes) {
			t.Errorf("GetRemotes() got remotes = %v, want remotes = %v", remotes, wantRemotes)
		}
	})
}

//nolint:tparallel,paralleltest
func Test_DevPlumbing_SetRemoteURLs(t *testing.T) {
	t.Parallel()

	// This test has been adapted from Test_gogit_SetRemoteURLs.

	mp := &mockenv.DevPlumbing{
		Plumbing: gggit.NewPlumbing(),
		URLMap:   make(map[string]string),
	}

	// for this test we have two repositories:
	// 'remote' <- 'clone'
	// Clone is a clone of remote.
	// We then try to set the remote url of the origin remote to bogus values.
	// We furthermore map these bogus urls using the URLMap.
	// This should succeed as long as the number of urls stays the same.

	mp.URLMap["https://example.com"] = "https://real.example.com"
	mp.URLMap["https://example2.com"] = "https://real.example2.com"

	// create an initial remote repository, and add a new bogus commit to it.
	remote, repo := testutil.NewTestRepo(t)
	testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})

	// clone the remote repository into 'clone'
	clone := testlib.TempDirAbs(t)
	repo, err := git.PlainClone(clone, false, &git.CloneOptions{URL: remote})
	if err != nil {
		panic(err)
	}

	// get a repo object
	ggRepoObject, isRepo := mp.IsRepository(t.Context(), clone)
	if !isRepo {
		panic("IsRepository() failed")
	}

	t.Run("setting existing remote with correct length", func(t *testing.T) {
		urls := []string{"https://example.com"}
		wantURLs := []string{"https://real.example.com"}

		err := mp.SetRemoteURLs(t.Context(), clone, ggRepoObject, "origin", urls)
		if err != nil {
			t.Error("SetRemoteURLs() got err != nil, want err = nil")
		}

		cfg, err := repo.Remote("origin")
		if err != nil {
			panic(err)
		}
		gotURLs := cfg.Config().URLs

		if !reflect.DeepEqual(gotURLs, wantURLs) {
			t.Errorf("SetRemoteURLs() set urls = %v, want urls = %v", gotURLs, wantURLs)
		}
	})

	t.Run("setting existing remote with incorrect length", func(t *testing.T) {
		urls := []string{"https://example.com", "https://example2.com"}
		wantURLs := []string{"https://real.example.com", "https://real.example2.com"}

		err := mp.SetRemoteURLs(t.Context(), clone, ggRepoObject, "origin", urls)
		if err == nil {
			t.Error("SetRemoteURLs() got err = nil, want err != nil")
		}

		cfg, err := repo.Remote("origin")
		if err != nil {
			panic(err)
		}
		gotURLs := cfg.Config().URLs

		if reflect.DeepEqual(gotURLs, wantURLs) {
			t.Errorf("SetRemoteURLs() set urls = %v, did not want urls = %v", gotURLs, wantURLs)
		}
	})

	t.Run("setting non-existent remote", func(t *testing.T) {
		urls := []string{"https://example.com", "https://example2.com"}

		err := mp.SetRemoteURLs(t.Context(), clone, ggRepoObject, "upstream", urls)
		if err == nil {
			t.Error("SetRemoteURLs() got err = nil, want err != nil")
		}
	})
}

func Test_DevPlumbing_GetCanonicalRemote(t *testing.T) {
	t.Parallel()

	// This test has been adapted from Test_gogit_GetCanonicalRemote.

	mp := &mockenv.DevPlumbing{
		Plumbing: gggit.NewPlumbing(),
		URLMap:   make(map[string]string),
	}

	// For this test we have three repositories:
	// 'remote' <- 'cloneA' <- 'cloneB'
	// Each repository has remotes pointing to the previous ones.
	// we also set up a mapping for remote and cloneA.

	// we then ask 'cloneA' and 'cloneB' for their canonical remotes.
	// These should return the mapped versions of 'remote' and 'cloneA' respectively.

	// create an initial remote repository, and add a new bogus commit to it.
	remote, repo := testutil.NewTestRepo(t)
	testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})

	// clone the remote repository into 'cloneA'.
	// This will create an origin remote pointing to the remote.
	cloneA := testlib.TempDirAbs(t)
	if _, err := git.PlainClone(cloneA, false, &git.CloneOptions{URL: remote}); err != nil {
		panic(err)
	}

	// clone the 'cloneA' repository into 'cloneB'.
	// This will create an origin remote pointing to 'cloneA'
	cloneB := testlib.TempDirAbs(t)
	repo, err := git.PlainClone(cloneB, false, &git.CloneOptions{URL: cloneA})
	if err != nil {
		panic(err)
	}

	// create an 'upstream' remote to point to 'remote'.
	if _, err := repo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{remote},
	}); err != nil {
		panic(err)
	}

	mappedCloneA := "git://server.com:2222/cloneA.git"
	mp.URLMap[mappedCloneA] = cloneA

	mappedRemote := "git://server.com:2222/remote.git"
	mp.URLMap[mappedRemote] = remote

	t.Run("GetCanonicalRemote() on a repository with a single remote", func(t *testing.T) {
		t.Parallel()

		ggRepoObject, isRepo := mp.IsRepository(t.Context(), cloneA)
		if !isRepo {
			panic("IsRepository() failed")
		}

		wantName := "origin"
		wantURLs := []string{mappedRemote}

		name, urls, err := mp.GetCanonicalRemote(t.Context(), cloneA, ggRepoObject)
		if err != nil {
			t.Error("GetCanonicalRemote() got err != nil, want err == nil")
		}

		if name != wantName {
			t.Errorf("GetCanonicalRemote() got name = %v, want name = %v", name, wantName)
		}
		if !reflect.DeepEqual(urls, wantURLs) {
			t.Errorf("GetCanonicalRemote() got urls = %v, want urls = %v", urls, wantURLs)
		}
	})

	t.Run("GetCanonicalRemote() on a repository with more than a single remote", func(t *testing.T) {
		t.Parallel()

		ggRepoObject, isRepo := mp.IsRepository(t.Context(), cloneB)
		if !isRepo {
			panic("IsRepository() failed")
		}

		wantName := "origin"
		wantURLs := []string{mappedCloneA}

		name, urls, err := mp.GetCanonicalRemote(t.Context(), cloneB, ggRepoObject)
		if err != nil {
			t.Error("GetCanonicalRemote() got err != nil, want err == nil")
		}

		if name != wantName {
			t.Errorf("GetCanonicalRemote() got name = %v, want name = %v", name, wantName)
		}
		if !reflect.DeepEqual(urls, wantURLs) {
			t.Errorf("GetCanonicalRemote() got urls = %v, want urls = %v", urls, wantURLs)
		}
	})
}

func Test_DevPlumbing_Clone(t *testing.T) {
	t.Parallel()

	// This test has been adapted from Test_gogit_Clone.
	mp := &mockenv.DevPlumbing{
		Plumbing: gggit.NewPlumbing(),
		URLMap:   make(map[string]string),
	}

	// create an initial remote repository, and add a new bogus commit to it.
	remote, repo := testutil.NewTestRepo(t)
	testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})

	mappedRemote := "https://example.com/example.git"
	mp.URLMap[mappedRemote] = remote

	t.Run("cloning a repository", func(t *testing.T) {
		t.Parallel()

		clone := testlib.TempDirAbs(t)

		err := mp.Clone(t.Context(), stream.FromNil(), mappedRemote, clone)
		if err != nil {
			t.Error("Clone() got err != nil, want err = nil")
		}

		if _, err := git.PlainOpen(clone); err != nil {
			t.Error("Clone() did not clone repository")
		}
	})
}
