package testutil

import (
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// NewTestRepo creates a new (empty) repository at the provided path
// clonePath is the path where the repository is located on disk.
// repo is a reference to the created repository.
// cleanup is a function that removes the repository from disk.
// If something goes wrong, the function calls panic()
//
// This function is intended to be called like:
//   clonePath, repo, cleanup := NewTestRepo()
//   defer cleanup()
//
func NewTestRepo() (clonePath string, repo *git.Repository, cleanup func()) {

	// first create a new temporary directory to put the git repository in
	clonePath, cleanup = TempDir()

	// then actually do a git PlainInit
	repo = NewTestRepoAt(clonePath)
	if repo == nil {
		cleanup()
		panic("NewTestRepoAt(): Repository not created")
	}

	return
}

// NewTestRepoAt creates a new repository at the provided path.
// When an error occurs, returns nil.
//
// This function is untested.
func NewTestRepoAt(clonePath string) (repo *git.Repository) {
	repo, err := git.PlainInit(clonePath, false)
	if err != nil {
		return nil
	}
	return repo
}

const commitMessage = "CommitTestFiles() commit"

// AuthorName is the name to be used for authors of test git commit-likes
const AuthorName = "Jane Doe"

// AuthorEmail is the email to be used for email of the author of test git commit-likes
const AuthorEmail = "jane.doe@example.com"

// CommitTestFiles makes a new commit in the repository repo.
// The commit will contain files with the names and content of the contained map.
// When the map is nil, a default dummy file will be used instead.
// The commit will appear to have been authored from a bogus author and have a bogus commit message.
//
// The function returns the worktree of the repository and the commit hash produced
//
// The files will be written out to disk.
// If an error occurs, panic() is called.
//
func CommitTestFiles(repo *git.Repository, files map[string]string) (*git.Worktree, plumbing.Hash) {
	// get the worktree of the repository
	// and the root directory
	worktree, err := repo.Worktree()
	if err != nil {
		panic(err)
	}
	root := worktree.Filesystem.Root()

	if files == nil {
		files = map[string]string{"dummy.txt": "I am a dummy file. "}
	}

	// write each file to disk and add it to the staging area
	for file, content := range files {
		if err := ioutil.WriteFile(path.Join(root, file), []byte(content), os.ModePerm); err != nil {
			panic(err)
		}
		if _, err := worktree.Add(file); err != nil {
			panic(err)
		}
	}

	// make the commit
	commit, err := worktree.Commit(commitMessage, &git.CommitOptions{
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
