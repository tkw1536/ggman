package walker

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// FS is the filesystem abstraction used by Walker
//
// See NewRealFS for a instantiating a sample implementation.
type FS interface {
	// Path returns the current pathof this FS.
	// This function will be called only once, and may perform (potentially slow) normalization.
	//
	// The return value is used for cycle detection, and also passed to all other functions in this interface.
	Path() (string, error)

	// Read reads the root directory of this filesystem.
	// and returns a list of directory entries sorted by filename.
	//
	// If is roughly equivalent to the ReadDir method of fs.ReadDirFS.
	// Assuming fsys is an internal fs.FS the method might be implemented as:
	//
	//  fs.ReadDir(fs.FS(fsys), ".")
	Read(path string) ([]fs.DirEntry, error)

	// CanSub indicates if the given directory entry can be used as a valid FS.
	//
	// Sub creates a new FS for the provided entry.
	// Sub is only called when CanSub returns true and a nil error.
	CanSub(path string, entry fs.DirEntry) (bool, error)
	Sub(path string, entry fs.DirEntry) FS
}

// NewRealFS returns a new filestsystem rooted at path.
// followLinks indicates if the filesystem should follow and resolve links.
func NewRealFS(path string, followLinks bool) FS {
	return realFS{
		DirPath:     path,
		FollowLinks: followLinks,
	}
}

// realFS represents the real underlying filesystem.
// It implements FS.
//
// This struct is untested; tests are done via Scan and Sweep.
type realFS struct {
	// DirPath represents the path of the filesystem being represented.
	DirPath string
	// FollowLinks indicates if the filesystem follows symlinks.
	FollowLinks bool
}

func (real realFS) Read(path string) ([]fs.DirEntry, error) {
	return fs.ReadDir(os.DirFS(path), ".")
}

func (real realFS) Path() (string, error) {
	if !real.FollowLinks {
		return real.DirPath, nil
	}
	return filepath.EvalSymlinks(real.DirPath)
}

func (real realFS) CanSub(path string, entry fs.DirEntry) (bool, error) {
	child := filepath.Join(path, entry.Name())
	return IsDirectory(child, real.FollowLinks)
}

func (real realFS) Sub(path string, entry fs.DirEntry) FS {
	child := filepath.Join(path, entry.Name())
	return realFS{
		DirPath:     child,
		FollowLinks: real.FollowLinks,
	}
}

// IsDirectory checks if path exists and points to a directory
// When includeLinks is true, a symlink counts as a directory.
func IsDirectory(path string, includeLinks bool) (bool, error) {

	// Stat() returns information about the referenced path by default.
	// In case of a symlink, this means the target of the link.
	//
	// If we allow links, this is exactly what we want.
	// If we don't allow links, we want information about the link itself.
	// We thus need to use LStat().
	var stat os.FileInfo
	var err error
	if includeLinks {
		stat, err = os.Stat(path)
	} else {
		stat, err = os.Lstat(path)
	}

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return false, nil
	case err != nil:
		return false, errors.Wrap(err, "Stat failed")
	}

	return stat.IsDir(), nil
}
