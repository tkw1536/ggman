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

//spellchecker:words path filepath testing github pkglib testlib
import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tkw1536/pkglib/testlib"
)

//spellchecker:words gogit gitgit

func TestNewGitFromPlumbing(t *testing.T) {

	// create a temporary file
	dir := testlib.TempDirAbs(t)

	// no path => use a gogit
	if _, isGogit := NewGitFromPlumbing(nil, "").Plumbing().(*gogit); !isGogit {
		t.Errorf("NewGitFromPlumbing: Expected *gogit")
	}

	// path but no git => use a gogit
	if _, isGogit := NewGitFromPlumbing(nil, dir).Plumbing().(*gogit); !isGogit {
		t.Errorf("NewGitFromPlumbing: Expected *gogit")
	}

	if err := os.WriteFile(filepath.Join(dir, "git"), nil, os.ModePerm&0x111); err != nil {
		panic(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "git.exe"), nil, os.ModePerm&0x111); err != nil {
		panic(err)
	}

	// path with git => gitgit
	if _, isGitgit := NewGitFromPlumbing(nil, dir).Plumbing().(*gitgit); !isGitgit {
		t.Errorf("NewGitFromPlumbing: Expected *gitgit")
	}

}
