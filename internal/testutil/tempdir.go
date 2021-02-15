package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// TempDir creates a new temporary directory to be used during testing.
// If something goes wrong, calls panic().
//
// The caller is expected to call cleanup() to remove the temporary directory.
// A typical invocation would be something like:
//
//  dir, cleanup = TempDir()
//  defer cleanup()
//
func TempDir() (path string, cleanup func()) {
	// This function is more or less a thin wrapper around ioutil.TempDir.
	// The reason it exists is because it saves a lot of boilerplate, like checking err != nil.
	path, err := ioutil.TempDir("", "")
	if err != nil {
		panic(err)
	}
	cleanup = func() {
		if err := os.RemoveAll(path); err != nil {
			panic(err)
		}
	}
	path, err = filepath.EvalSymlinks(path)
	if err != nil {
		cleanup()
		panic(err)
	}
	return
}
