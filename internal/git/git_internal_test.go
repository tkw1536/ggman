// Package git contains an implementation of git functionality.
package git

//spellchecker:words path filepath testing pkglib testlib
import (
	"os"
	"path/filepath"
	"testing"

	"go.tkw01536.de/pkglib/testlib"
)

//spellchecker:words gogit gitgit

func TestNewGitFromPlumbing(t *testing.T) {
	t.Parallel()

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
