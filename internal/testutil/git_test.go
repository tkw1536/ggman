//spellchecker:words testutil
package testutil

//spellchecker:words testing github
import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
)

func TestNewTestRepo(t *testing.T) {
	dir, repo := NewTestRepo(t)

	if s, err := os.Stat(dir); err != nil || !s.IsDir() {
		t.Errorf("NewTestRepo: Directory was not created. ")
	}

	if repo == nil {
		t.Errorf("NewTestRepo: Repository was not returned. ")
	}

	if _, err := git.PlainOpen(dir); err != nil {
		t.Errorf("NewTestRepo: Repository was not created. ")
	}
}

func TestCommitTestFiles(t *testing.T) {
	_, repo := NewTestRepo(t)

	_, hash := CommitTestFiles(repo, nil)

	// check that the repository has 'hash' checked out.
	head, err := repo.Head()
	if err != nil {
		t.Errorf("CommitTestFiles: No HEAD")
	}

	if head.Hash() != hash {
		t.Errorf("CommitTestFiles: returned hash not checked out")
	}
}
