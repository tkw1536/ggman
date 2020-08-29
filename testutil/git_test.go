package testutil

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestNewTestRepo(t *testing.T) {
	dir, repo, cleanup := NewTestRepo()
	defer cleanup()

	if s, err := os.Stat(dir); err != nil || !s.IsDir() {
		t.Errorf("NewTestRepo(): Directory was not created. ")
	}

	if repo == nil {
		t.Errorf("NewTestRepo(): Repository was not returned. ")
	}

	if _, err := git.PlainOpen(dir); err != nil {
		t.Errorf("NewTestRepo(): Repository was not created. ")
	}

	cleanup()
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("NewTestRepo(): Cleanup did not remove directory")
	}
}

func TestCommitTestFiles(t *testing.T) {
	_, repo, cleanup := NewTestRepo()
	defer cleanup()

	_, hash := CommitTestFiles(repo, nil)

	// check that the repository has 'hash' checked out.
	head, err := repo.Head()
	if err != nil {
		t.Errorf("CommitTestFiles(): No HEAD")
	}

	if head.Hash() != hash {
		t.Errorf("CommitTestFiles(): returned hash not checked out")
	}
}
