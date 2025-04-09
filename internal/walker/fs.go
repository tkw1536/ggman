//spellchecker:words walker
package walker

//spellchecker:words path filepath github pkglib
import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/tkw1536/pkglib/fsx"
)

//spellchecker:words fsys

// FS represents a file system for use by walker
//
// See NewRealFS for a instantiating a sample implementation.
type FS interface {
	// Path returns the path of this FS.
	// The path should not be normalized.
	Path() string

	// ResolvedPath returns the current path of this FS.
	// This function will be called only once, and may perform (potentially slow) normalization.
	//
	// The return value is used for cycle detection, and also passed to all other functions in this interface.
	ResolvedPath() (string, error)

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
	// path and rpath are the Path() and ResolvedPath() values.
	// Sub is only called when CanSub returns true and a nil error.
	CanSub(path string, entry fs.DirEntry) (bool, error)
	Sub(path, rpath string, entry fs.DirEntry) FS
}

// NewRealFS returns a new filesystem rooted at path.
// followLinks indicates if the filesystem should follow and resolve links.
func NewRealFS(path string, followLinks bool) FS {
	return realFS{
		DirPath:     path,
		FollowLinks: followLinks,
	}
}

// realFS represents the real underlying filesystem, implementing FS.
//
// This struct is untested; tests are done via Scan and Sweep.
type realFS struct {
	// DirPath represents the path of the filesystem being represented.
	DirPath string
	// FollowLinks indicates if the filesystem follows symlinks.
	FollowLinks bool
}

func (rfs realFS) Read(path string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(os.DirFS(path), ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	return entries, nil
}

func (rfs realFS) Path() string {
	return rfs.DirPath
}

func (rfs realFS) ResolvedPath() (string, error) {
	if !rfs.FollowLinks {
		return rfs.DirPath, nil
	}
	path, err := filepath.EvalSymlinks(rfs.DirPath)
	if err != nil {
		return "", fmt.Errorf("failed to evaluate symlinks: %w", err)
	}
	return path, nil
}

func (rfs realFS) CanSub(path string, entry fs.DirEntry) (bool, error) {
	child := filepath.Join(path, entry.Name())
	isDir, err := fsx.IsDirectory(child, rfs.FollowLinks)
	if err != nil {
		return false, fmt.Errorf("failed to check directory: %w", err)
	}
	return isDir, nil
}

func (rfs realFS) Sub(path, rpath string, entry fs.DirEntry) FS {
	child := filepath.Join(path, entry.Name())
	return realFS{
		DirPath:     child,
		FollowLinks: rfs.FollowLinks,
	}
}
