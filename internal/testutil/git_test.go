//spellchecker:words testutil
package testutil_test

//spellchecker:words testing github ggman internal testutil
import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/tkw1536/ggman/internal/testutil"
)

func TestNewTestRepo(t *testing.T) {
	dir, repo := testutil.NewTestRepo(t)

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
	_, repo := testutil.NewTestRepo(t)

	_, hash := testutil.CommitTestFiles(repo, nil)

	// check that the repository has 'hash' checked out.
	head, err := repo.Head()
	if err != nil {
		t.Errorf("CommitTestFiles: No HEAD")
	}

	if head.Hash() != hash {
		t.Errorf("CommitTestFiles: returned hash not checked out")
	}
}
