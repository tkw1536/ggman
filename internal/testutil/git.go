//spellchecker:words testutil
package testutil

//spellchecker:words path sync atomic testing time github config plumbing object pkglib testlib
import (
	"fmt"
	"os"
	"path"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"go.tkw01536.de/pkglib/testlib"
)

//spellchecker:words worktree nosec Storer

// NewTestRepo creates a new empty repository for testing at an unspecified path.
//
// It returns a pair of path the repository has been created at and a git repository.
// It is the callers responsibility to delete the test repository once it is no longer needed.
//
// If something goes wrong, the function calls panic().
func NewTestRepo(t *testing.T) (clonePath string, repo *git.Repository) {
	t.Helper()

	// first create a new temporary directory to put the git repository in
	clonePath = testlib.TempDirAbs(t)

	// then create a test repo there
	repo = NewTestRepoAt(clonePath, "")
	if repo == nil {
		panic("NewTestRepo: Repository not created")
	}

	return
}

// NewTestRepo creates a new empty repository for testing at the specified path.
// If remote is a non-empty string, also creates a new remote called "origin" pointing to the given remote.
//
// It returns a reference to the underlying git repository.
// It is the callers responsibility to delete the test repository once it is no longer needed.
// If something goes wrong, the function returns nil.
//
// The 'remote' part of this function is untested.
func NewTestRepoAt(clonePath, remote string) (repo *git.Repository) {
	repo, err := git.PlainInit(clonePath, false)
	if err != nil {
		return nil
	}
	if remote != "" {
		if _, err := repo.CreateRemote(&config.RemoteConfig{
			Name: "origin",
			URLs: []string{remote},
		}); err != nil {
			return nil
		}
	}
	return repo
}

const (
	// CommitMessage is the message to be used for the commit.
	CommitMessage = "CommitTestFiles() commit"

	// AuthorName is the name to be used for authors of test git commit-likes.
	AuthorName = "Jane Doe"

	// AuthorEmail is the email to be used for email of the author of test git commit-likes.
	AuthorEmail = "jane.doe@example.com"

	// FileNamePrefix is the prefix of the name of each dummy file.
	FileNamePrefix = "dummy-"

	// FileContents is the content of each dummy file.
	FileContents = ""
)

// counter for test file names.
var dummyFileCounter atomic.Uint64

// CommitTestFiles makes a new commit in the repository repo.
// The commit will contain new dummy files and content each time it is called.
// The commit will appear to have been authored from a bogus author and have a bogus commit message.
//
// # The function returns the worktree of the repository and the commit hash produced
//
// The files will be written out to disk.
// If an error occurs, panic() is called.
func CommitTestFiles(repo *git.Repository) (*git.Worktree, plumbing.Hash) {
	// get the worktree of the repository
	// and the root directory
	worktree, err := repo.Worktree()
	if err != nil {
		panic(err)
	}
	root := worktree.Filesystem.Root()

	// write the file to disk and add it to the staging area
	fileName := fmt.Sprintf("%s%d", FileNamePrefix, dummyFileCounter.Add(1))
	if err := os.WriteFile(path.Join(root, fileName), []byte(FileContents), os.ModePerm /* #nosec G306 -- fine for testing */); err != nil {
		panic(err)
	}
	if _, err := worktree.Add(fileName); err != nil {
		panic(err)
	}

	// make the commit
	commit, err := worktree.Commit(CommitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  AuthorName,
			Email: AuthorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		panic(err)
	}

	return worktree, commit
}

// CreateTrackingBranch creates a new branch named 'localBranch' that tracks the remote branch 'remoteBranch' of the given remote.
// If something goes wrong, the function calls panic().
func CreateTrackingBranch(repo *git.Repository, remoteName, localBranch, remoteBranch string) {
	headRef, err := repo.Head()
	if err != nil {
		panic(fmt.Sprintf("failed to get HEAD: %v", err))
	}
	headHash := headRef.Hash()

	// create the new branch reference
	newBranch := plumbing.NewHashReference(
		plumbing.NewBranchReferenceName(localBranch),
		headHash,
	)
	if err := repo.Storer.SetReference(newBranch); err != nil {
		panic(fmt.Sprintf("failed to create branch ref: %v", err))
	}

	// setup the branch config
	cfg, err := repo.Config()
	if err != nil {
		panic(fmt.Sprintf("failed to read repo config: %v", err))
	}
	if cfg.Branches == nil {
		cfg.Branches = make(map[string]*config.Branch)
	}
	cfg.Branches[localBranch] = &config.Branch{
		Name:   localBranch,
		Remote: remoteName,
		Merge:  plumbing.NewBranchReferenceName(remoteBranch),
	}

	if err := repo.Storer.SetConfig(cfg); err != nil {
		panic(fmt.Sprintf("failed to write repo config: %v", err))
	}
}
