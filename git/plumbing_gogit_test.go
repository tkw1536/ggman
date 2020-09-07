package git

import (
	"reflect"
	"testing"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/tkw1536/ggman"
	"github.com/tkw1536/ggman/testutil"
)

func Test_gogit_IsRepository(t *testing.T) {
	var gg gogit

	// for this test, we make three directories for testing:
	// - a folder with an existing repository in it
	// - an empty folder
	// - a deleted folder
	// Only in the first of these IsRepository() should return true on.

	// make a folder with an empty repository
	existingRepo, _, cleanup := testutil.NewTestRepo()
	defer cleanup()

	// make an empty folder
	emptyFolder, cleanup := testutil.TempDir()
	defer cleanup()

	// create a new folder that is deleted
	deletedFolder, cleanup := testutil.TempDir()
	cleanup()

	type args struct {
		localPath string
	}
	tests := []struct {
		name       string
		args       args
		wantIsRepo bool
	}{
		{"existing repository is a repository", args{existingRepo}, true},
		{"empty folder is not repository", args{emptyFolder}, false},
		{"deleted folder is not repository", args{deletedFolder}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotIsRepo := gg.IsRepository(tt.args.localPath)
			if gotIsRepo != tt.wantIsRepo {
				t.Errorf("gogit.IsRepository() gotIsRepo = %v, want %v", gotIsRepo, tt.wantIsRepo)
			}
		})
	}
}

func Test_gogit_IsRepositoryUnsafe(t *testing.T) {
	var gg gogit

	// This test behaves like the IsRepository() test.
	// We again make three directories for testing:
	// - a folder with an existing repository in it
	// - an empty folder
	// - a deleted folder
	// Only in the first of these IsRepositoryUnsafe() should return true on.

	// make a folder with an empty repository
	existingRepo, _, cleanup := testutil.NewTestRepo()
	defer cleanup()

	// make an empty folder
	emptyFolder, cleanup := testutil.TempDir()
	defer cleanup()

	// create a new folder that is deleted
	deletedFolder, cleanup := testutil.TempDir()
	cleanup()

	type args struct {
		localPath string
	}
	tests := []struct {
		name       string
		args       args
		wantIsRepo bool
	}{
		{"existing repository is a repository", args{existingRepo}, true},
		{"empty folder is not repository", args{emptyFolder}, false},
		{"deleted folder is not repository", args{deletedFolder}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsRepo := gg.IsRepositoryUnsafe(tt.args.localPath)
			if gotIsRepo != tt.wantIsRepo {
				t.Errorf("gogit.IsRepositoryUnsafe() gotIsRepo = %v, want %v", gotIsRepo, tt.wantIsRepo)
			}
		})
	}
}

func Test_gogit_GetHeadRefA(t *testing.T) {
	var gg gogit

	// for this test we make three repositories:
	// - one with an empty repository
	// - one with a checked out branch 'test'
	// - one with a checked out hash

	// make an empty repository
	emptyRepo, _, cleanup := testutil.NewTestRepo()
	defer cleanup()

	// make a new repository and checkout a new branch 'test'
	branchTestCheckout, repo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	worktree, commit := testutil.CommitTestFiles(repo, nil)
	if err := worktree.Checkout(&git.CheckoutOptions{
		Hash:   commit,
		Branch: plumbing.NewBranchReferenceName("test"),
		Create: true,
	}); err != nil {
		panic(err)
	}

	// make a new repository and checkout a hash
	hashCheckout, repo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	worktree, commit = testutil.CommitTestFiles(repo, nil)
	if err := worktree.Checkout(&git.CheckoutOptions{
		Hash: commit,
	}); err != nil {
		panic(err)
	}

	type args struct {
		clonePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"head of empty repository is not defined", args{emptyRepo}, "", true},
		{"head of a branch is the branch", args{branchTestCheckout}, "test", false},
		{"head of commit without branch or tag is a hash", args{hashCheckout}, commit.String(), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// get the repo object
			ggRepoObject, isRepo := gg.IsRepository(tt.args.clonePath)
			if !isRepo {
				panic("IsRepository() failed")
			}

			got, err := gg.GetHeadRef(tt.args.clonePath, ggRepoObject)
			if (err != nil) != tt.wantErr {
				t.Errorf("gogit.GetHeadRef() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("gogit.GetHeadRef() = %v, want %v", got, tt.want)
			}
		})
	}
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

		err := gg.Clone(ggman.NewEnvIOStream(), remote, clone)
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

		err := gg.Clone(ggman.NewEnvIOStream(), remote, clone, "--branch", "main")
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
		err := gg.Fetch(ggman.NewEnvIOStream(), clone, ggRepoObject)
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
		err := gg.Fetch(ggman.NewEnvIOStream(), clone, ggRepoObject)
		if err != nil {
			t.Error("Fetch() returned err != nil, want err = nil")
		}
	})
}

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
		err := gg.Pull(ggman.NewEnvIOStream(), clone, ggRepoObject)
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
		err := gg.Pull(ggman.NewEnvIOStream(), clone, ggRepoObject)
		if err != nil {
			t.Error("Pull() returned err != nil, want err = nil")
		}
	})
}

func Test_gogit_ContainsBranch(t *testing.T) {
	var gg gogit

	// In this test we only have a single repository.
	// We create two branches 'branchA' and 'branchB'
	clone, repo, cleanup := testutil.NewTestRepo()
	defer cleanup()
	repo.CreateBranch(&config.Branch{Name: "branchA"})
	repo.CreateBranch(&config.Branch{Name: "branchB"})

	type args struct {
		branch string
	}
	tests := []struct {
		name         string
		args         args
		wantContains bool
		wantErr      bool
	}{
		{"branchA exists", args{"branchA"}, true, false},
		{"branchB exists", args{"branchB"}, true, false},
		{"branchC does not exist", args{"branchC"}, false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ggRepoObject, isRepo := gg.IsRepository(clone)
			if !isRepo {
				panic("IsRepository() failed")
			}

			gotContains, err := gg.ContainsBranch(clone, ggRepoObject, tt.args.branch)
			if (err != nil) != tt.wantErr {
				t.Errorf("gogit.ContainsBranch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotContains != tt.wantContains {
				t.Errorf("gogit.ContainsBranch() = %v, want %v", gotContains, tt.wantContains)
			}
		})
	}
}
