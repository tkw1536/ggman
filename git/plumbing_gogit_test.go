package git

import (
	"reflect"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/tkw1536/ggman/testutil"
)

func Test_gogit_IsRepository(t *testing.T) {
	var gg gogit

	t.Run("existing repository is a repository", func(t *testing.T) {
		clonePath, _, cleanup := testutil.NewTestRepo()
		defer cleanup()

		_, isRepo := gg.IsRepository(clonePath)
		if !isRepo {
			t.Error("IsRepository() = false, wanted IsRepository() = true")
		}

	})

	t.Run("empty folder is not repository", func(t *testing.T) {
		clonePath, cleanup := testutil.TempDir()
		defer cleanup()

		_, isRepo := gg.IsRepository(clonePath)
		if isRepo {
			t.Error("IsRepository() = true, wanted IsRepository() = false")
		}
	})

	t.Run("deleted folder is not repository", func(t *testing.T) {
		clonePath, cleanup := testutil.TempDir()
		defer cleanup()

		_, isRepo := gg.IsRepository(clonePath)
		if isRepo {
			t.Error("IsRepository() = true, wanted IsRepository() = false")
		}

	})
}

func Test_gogit_IsRepositoryUnsafe(t *testing.T) {
	var gg gogit

	t.Run("existing repository is a repository", func(t *testing.T) {
		clonePath, _, cleanup := testutil.NewTestRepo()
		defer cleanup()

		isRepo := gg.IsRepositoryUnsafe(clonePath)
		if !isRepo {
			t.Error("IsRepositoryUnsafe() = false, wanted IsRepository() = true")
		}

	})

	t.Run("empty folder is not repository", func(t *testing.T) {
		clonePath, cleanup := testutil.TempDir()
		defer cleanup()

		isRepo := gg.IsRepositoryUnsafe(clonePath)
		if isRepo {
			t.Error("IsRepositoryUnsafe() = true, wanted IsRepository() = false")
		}
	})

	t.Run("deleted folder is not repository", func(t *testing.T) {
		clonePath, cleanup := testutil.TempDir()
		cleanup()

		isRepo := gg.IsRepositoryUnsafe(clonePath)
		if isRepo {
			t.Error("IsRepositoryUnsafe() = true, wanted IsRepository() = false")
		}

	})
}

func Test_gogit_GetHeadRef(t *testing.T) {
	var gg gogit

	// make a temporary repository
	clonePath, repo, cleanup := testutil.NewTestRepo()
	defer cleanup()

	// get the repo object
	ggRepoObject, isRepo := gg.IsRepository(clonePath)
	if !isRepo {
		panic("IsRepository() failed")
	}

	t.Run("head of empty repository is not defined", func(t *testing.T) {
		_, err := gg.GetHeadRef(clonePath, ggRepoObject)
		if err == nil {
			t.Error("GetHeadRef() got err == nil, want err != nil")
		}
	})

	// make a new test commit and check it out on a new hash
	worktree, commitHash := testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})
	if err := worktree.Checkout(&git.CheckoutOptions{
		Hash: commitHash,
	}); err != nil {
		panic(err)
	}

	t.Run("head of commit without branch or tag is a hash", func(t *testing.T) {
		wantRef := commitHash.String()
		ref, err := gg.GetHeadRef(clonePath, ggRepoObject)
		if err != nil {
			t.Error("GetHeadRef() got err != nil, want err == nil")
		}

		if ref != wantRef {
			t.Errorf("GetHeadRef() got ref = %s, want ref = %s", ref, wantRef)
		}
	})

	// checkout a new branch called 'test'
	if err := worktree.Checkout(&git.CheckoutOptions{
		Create: true,
		Hash:   commitHash,
		Branch: "refs/heads/test",
	}); err != nil {
		panic(err)
	}

	t.Run("head of a branch is the branch", func(t *testing.T) {
		wantRef := "test"
		ref, err := gg.GetHeadRef(clonePath, ggRepoObject)
		if err != nil {
			t.Error("GetHeadRef() got err != nil, want err == nil")
		}

		if ref != wantRef {
			t.Errorf("GetHeadRef() got ref = %s, want ref = %s", ref, wantRef)
		}
	})
}

func Test_gogit_GetRemotes(t *testing.T) {
	var gg gogit

	// For this test we have three repositories:
	// 'remote' <- 'cloneA' <- 'cloneB'
	// 'cloneA' has an origin remote pointing to 'remote'.
	// 'cloneB' has an origin remote pointing to 'cloneA'.
	// 'cloneB' also gets an upstream remote pointing to 'remote'.
	// GetRemotes() should return all of them.

	// create an initial remote repository, and add a new bogus commit to it.
	remote, repo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})

	// clone the remote repository into 'cloneA'.
	// This will create an origin remote pointing to the remote.
	cloneA, cleanup := testutil.TempDir()
	defer cleanup()
	if _, err := git.PlainClone(cloneA, false, &git.CloneOptions{URL: remote}); err != nil {
		panic(err)
	}

	// clone the 'cloneA' repository into 'cloneB'.
	// This will create an origin remote pointing to 'cloneA'
	cloneB, cleanup := testutil.TempDir()
	defer cleanup()
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

	t.Run("GetRemotes() on a repository with a single remote", func(t *testing.T) {
		ggRepoObject, isRepo := gg.IsRepository(cloneA)
		if !isRepo {
			panic("IsRepository() failed")
		}

		wantRemotes := map[string][]string{
			"origin": {remote},
		}
		remotes, err := gg.GetRemotes(cloneA, ggRepoObject)
		if err != nil {
			t.Error("GetRemotes() got err != nil, want err == nil")
		}

		if !reflect.DeepEqual(remotes, wantRemotes) {
			t.Errorf("GetRemotes() got remotes = %v, want remotes = %v", remotes, wantRemotes)
		}
	})

	t.Run("GetRemotes() on a reposuitory with more than one remote", func(t *testing.T) {
		ggRepoObject, isRepo := gg.IsRepository(cloneB)
		if !isRepo {
			panic("IsRepository() failed")
		}

		wantRemotes := map[string][]string{
			"upstream": {remote},
			"origin":   {cloneA},
		}
		remotes, err := gg.GetRemotes(cloneB, ggRepoObject)
		if err != nil {
			t.Error("GetRemotes() got err != nil, want err == nil")
		}

		if !reflect.DeepEqual(remotes, wantRemotes) {
			t.Errorf("GetRemotes() got remotes = %v, want remotes = %v", remotes, wantRemotes)
		}
	})
}

func Test_gogit_GetCanonicalRemote(t *testing.T) {
	var gg gogit

	// For this test we have three repositories:
	// 'remote' <- 'cloneA' <- 'cloneB'
	// Each repository has remotes pointing to the previous ones.
	// we then ask 'cloneA' and 'cloneB' for their canonical remotes.
	// These should return 'remote' and 'cloneA' respectively.

	// create an initial remote repository, and add a new bogus commit to it.
	remote, repo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})

	// clone the remote repository into 'cloneA'.
	// This will create an origin remote pointing to the remote.
	cloneA, cleanup := testutil.TempDir()
	defer cleanup()
	if _, err := git.PlainClone(cloneA, false, &git.CloneOptions{URL: remote}); err != nil {
		panic(err)
	}

	// clone the 'cloneA' repository into 'cloneB'.
	// This will create an origin remote pointing to 'cloneA'
	cloneB, cleanup := testutil.TempDir()
	defer cleanup()
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

	t.Run("GetCanonicalRemote() on a repository with a single remote", func(t *testing.T) {
		ggRepoObject, isRepo := gg.IsRepository(cloneA)
		if !isRepo {
			panic("IsRepository() failed")
		}

		wantName := "origin"
		wantURLs := []string{remote}

		name, urls, err := gg.GetCanonicalRemote(cloneA, ggRepoObject)
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
		ggRepoObject, isRepo := gg.IsRepository(cloneB)
		if !isRepo {
			panic("IsRepository() failed")
		}

		wantName := "origin"
		wantURLs := []string{cloneA}

		name, urls, err := gg.GetCanonicalRemote(cloneB, ggRepoObject)
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

func Test_gogit_SetRemoteURLs(t *testing.T) {
	var gg gogit

	// for this test we have two repositories:
	// 'remote' <- 'clone'
	// Clone is a clone of remote.
	// We then try to set the remote url of the origin remote to bogus values.
	// This should succeed as long as the number of urls stays the same.

	// create an initial remote repository, and add a new bogus commit to it.
	remote, repo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})

	// clone the remote repository into 'clone'
	clone, cleanup := testutil.TempDir()
	defer cleanup()
	repo, err := git.PlainClone(clone, false, &git.CloneOptions{URL: remote})
	if err != nil {
		panic(err)
	}

	// get a repo object
	ggRepoObject, isRepo := gg.IsRepository(clone)
	if !isRepo {
		panic("IsRepository() failed")
	}

	t.Run("setting existing remote with correct length", func(t *testing.T) {
		urls := []string{"https://example.com"}

		err := gg.SetRemoteURLs(clone, ggRepoObject, "origin", urls)
		if err != nil {
			t.Error("SetRemoteURLs() got err != nil, want err = nil")
		}

		cfg, err := repo.Remote("origin")
		if err != nil {
			panic(err)
		}
		gotURLs := cfg.Config().URLs

		if !reflect.DeepEqual(gotURLs, urls) {
			t.Errorf("SetRemoteURLs() set urls = %v, want urls = %v", gotURLs, urls)
		}
	})

	t.Run("setting existing remote with incorrect length", func(t *testing.T) {
		urls := []string{"https://example.com", "https://example2.com"}

		err := gg.SetRemoteURLs(clone, ggRepoObject, "origin", urls)
		if err == nil {
			t.Error("SetRemoteURLs() got err = nil, want err != nil")
		}

		cfg, err := repo.Remote("origin")
		if err != nil {
			panic(err)
		}
		gotURLs := cfg.Config().URLs

		if reflect.DeepEqual(gotURLs, urls) {
			t.Errorf("SetRemoteURLs() set urls = %v, did not want urls = %v", gotURLs, urls)
		}
	})

	t.Run("setting non-existent remote", func(t *testing.T) {
		urls := []string{"https://example.com", "https://example2.com"}

		err := gg.SetRemoteURLs(clone, ggRepoObject, "upstream", urls)
		if err == nil {
			t.Error("SetRemoteURLs() got err = nil, want err != nil")
		}
	})
}

func Test_gogit_Clone(t *testing.T) {
	var gg gogit

	// create an initial remote repository, and add a new bogus commit to it.
	remote, repo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	testutil.CommitTestFiles(repo, map[string]string{"commit1.txt": "I was added in commit 1. "})

	t.Run("cloning a repository", func(t *testing.T) {
		clone, cleanup := testutil.TempDir()
		defer cleanup()

		err := gg.Clone(remote, clone)
		if err != nil {
			t.Error("Clone() got err != nil, want err = nil")
		}

		if _, err := git.PlainOpen(clone); err != nil {
			t.Error("Clone() did not clone repository")
		}
	})

	t.Run("cloning a repository with arguments is not supported", func(t *testing.T) {
		clone, cleanup := testutil.TempDir()
		defer cleanup()

		err := gg.Clone(remote, clone, "--branch", "main")
		if err != ErrArgumentsUnsupported {
			t.Error("Clone() got err != ErrArgumentsUnsupported, want err = ErrArgumentsUnsupported")
		}

	})
}

func Test_gogit_Fetch(t *testing.T) {
	var gg gogit

	// In this test we have three repositories:
	// 'upstream' with commits 'commitA' and 'commitB2'
	// 'origin' with commits 'commitA' and 'commitB1'
	// 'clone' with commits 'commitA'
	//
	// 'clone' has two remotes, upstream and origin which point to the respective repositories.
	// After fetching, the clone should become aware of both commits.

	// create an initial upstream repository, and add a new bogus commit to it.
	upstream, upstreamRepo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	_, commitA := testutil.CommitTestFiles(upstreamRepo, map[string]string{"commita.txt": "Commit A"})

	// clone upstream@commitA to the remote
	remote, cleanup := testutil.TempDir()
	defer cleanup()
	remoteRepo, err := git.PlainClone(remote, false, &git.CloneOptions{URL: upstream})
	if err != nil {
		panic(err)
	}

	// clone remote to the local clone
	clone, cleanup := testutil.TempDir()
	defer cleanup()
	cloneRepo, err := git.PlainClone(clone, false, &git.CloneOptions{URL: remote})
	if err != nil {
		panic(err)
	}

	// create an 'upstream' remote to point to 'remote'.
	if _, err := cloneRepo.CreateRemote(&config.RemoteConfig{
		Name: "upstream",
		URLs: []string{upstream},
	}); err != nil {
		panic(err)
	}

	// make distinct commits to the upstream and remote repo
	_, commitB1 := testutil.CommitTestFiles(upstreamRepo, map[string]string{"commitb1.txt": "Commit B1"})
	_, commitB2 := testutil.CommitTestFiles(remoteRepo, map[string]string{"commitb2.txt": "Commit B2"})

	// get a repo object
	ggRepoObject, isRepo := gg.IsRepository(clone)
	if !isRepo {
		panic("IsRepository() failed")
	}

	t.Run("fetching fetches all remotes", func(t *testing.T) {
		err := gg.Fetch(clone, ggRepoObject)
		if err != nil {
			t.Error("Fetch() returned err != nil, want err = nil")
		}

		head, err := cloneRepo.Head()
		if err != nil {
			panic(err)
		}

		if head.Hash() != commitA {
			t.Error("Fetch() updated HEAD")
		}

		gotOrigin, err := cloneRepo.Reference("refs/remotes/origin/master", true)
		if err != nil {
			panic(err)
		}
		if gotOrigin.Hash() != commitB2 {
			t.Error("Fetch() did not fetch origin properly")
		}

		gotUpstream, err := cloneRepo.Reference("refs/remotes/upstream/master", true)
		if err != nil {
			panic(err)
		}
		if gotUpstream.Hash() != commitB1 {
			t.Error("Fetch() did not fetch upstream properly")
		}
	})

	t.Run("fetching an up-to-date repo returns no error", func(t *testing.T) {
		err := gg.Fetch(clone, ggRepoObject)
		if err != nil {
			t.Error("Fetch() returned err != nil, want err = nil")
		}
	})
}

// TODO: More git testing

func Test_gogit_Pull(t *testing.T) {
	var gg gogit

	// In this test we have two repositories:
	// 'origin' with commits 'commitA' and 'commitB'
	// 'clone' with commit 'commitA'
	// after pulling clone should have updated to commitB.

	// create the upstream repository
	origin, originRepo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	testutil.CommitTestFiles(originRepo, map[string]string{"commita.txt": "Commit A"})

	// clone remote to the local clone
	clone, cleanup := testutil.TempDir()
	defer cleanup()
	cloneRepo, err := git.PlainClone(clone, false, &git.CloneOptions{URL: origin})
	if err != nil {
		panic(err)
	}

	// create a second commmit in the remote repo
	_, commitB := testutil.CommitTestFiles(originRepo, map[string]string{"commitb.txt": "Commit B"})

	// get a repo object
	ggRepoObject, isRepo := gg.IsRepository(clone)
	if !isRepo {
		panic("IsRepository() failed")
	}

	t.Run("pulling pulls a repository", func(t *testing.T) {
		err := gg.Pull(clone, ggRepoObject)
		if err != nil {
			t.Error("Pull() returned err != nil, want err = nil")
		}

		head, err := cloneRepo.Head()
		if err != nil {
			panic(err)
		}

		if head.Hash() != commitB {
			t.Error("Pull() did not update HEAD")
		}
	})

	t.Run("pulling an up-to-date repo returns no error", func(t *testing.T) {
		err := gg.Pull(clone, ggRepoObject)
		if err != nil {
			t.Error("Pull() returned err != nil, want err = nil")
		}
	})
}
