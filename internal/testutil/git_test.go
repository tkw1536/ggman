//spellchecker:words testutil
package testutil_test

//spellchecker:words testing github ggman internal testutil
import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"go.tkw01536.de/ggman/internal/testutil"
)

func TestNewTestRepo(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	_, repo := testutil.NewTestRepo(t)

	_, hash := testutil.CommitTestFiles(repo)

	// check that the repository has 'hash' checked out.
	head, err := repo.Head()
	if err != nil {
		t.Errorf("CommitTestFiles: No HEAD")
	}

	if head.Hash() != hash {
		t.Errorf("CommitTestFiles: returned hash not checked out")
	}
}

func TestCreateTrackingBranch(t *testing.T) {
	t.Parallel()

	_, repo := testutil.NewTestRepo(t)
	testutil.CommitTestFiles(repo)

	testutil.CreateTrackingBranch(repo, "origin", "feature", "main")

	// check that the branch was created
	branch, err := repo.Branch("feature")
	if err != nil {
		t.Errorf("CreateTrackingBranch: branch not created: %v", err)
	}

	// check that the branch tracks the correct remote
	if branch.Remote != "origin" {
		t.Errorf("CreateTrackingBranch: expected remote 'origin', got %q", branch.Remote)
	}

	// check that the branch merges with the correct remote branch
	if branch.Merge.Short() != "main" {
		t.Errorf("CreateTrackingBranch: expected merge 'main', got %q", branch.Merge.Short())
	}
}
