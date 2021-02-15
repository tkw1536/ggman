package testutil

import (
	"os"
	"testing"
)

func TestTempDir(t *testing.T) {
	dir, cleanup := TempDir()
	defer cleanup()

	if s, err := os.Stat(dir); err != nil || !s.IsDir() {
		t.Errorf("TempDir(): Directory was not created. ")
	}

	cleanup()
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("TempDir(): Cleanup did not remove directory")
	}
}
