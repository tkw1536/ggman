// Package git contains an implementation of git functionality.
//
// The implementation consists of the Git interface and the Plumbing interface.
//
// The Git interface (and it's default instance Default) provide a usable interface to Git Functionality.
// The Git interface will automatically choose between using a os.Exec() call to a native "git" wrapper, or using a pure golang git implementation.
// This should be used directly by callers.
//
// The Plumbing interface provides more direct control over which interface is used to interact with repositories.
// Calls to a Plumbing typically place assumptions on the caller and require some setup.
// For this reason, implementation of the Plumbing interface are not exported.
package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/ggman/internal/testutil"
)

func TestNewGitFromPlumbing(t *testing.T) {

	// create a temporary file
	dir, cleanup := testutil.TempDir()
	defer cleanup()

	// no path => use a gogit
	if _, isgogit := NewGitFromPlumbing(nil, "").Plumbing().(*gogit); !isgogit {
		t.Errorf("NewGitFromPlumbing(): Expected *gogit")
	}

	// path but no git => use a gogit
	if _, isgogit := NewGitFromPlumbing(nil, dir).Plumbing().(*gogit); !isgogit {
		t.Errorf("NewGitFromPlumbing(): Expected *gogit")
	}

	if err := os.WriteFile(filepath.Join(dir, "git"), nil, os.ModePerm&0x111); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "git.exe"), nil, os.ModePerm&0x111); err != nil {
		panic(err)
	}

	// path with git => gitgit
	if _, isgitgit := NewGitFromPlumbing(nil, dir).Plumbing().(*gitgit); !isgitgit {
		t.Errorf("NewGitFromPlumbing(): Expected *gitgit")
	}

}
