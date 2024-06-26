//spellchecker:words walker
package walker

//spellchecker:words path filepath github pkglib
import (
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

func (real realFS) Read(path string) ([]fs.DirEntry, error) {
	return fs.ReadDir(os.DirFS(path), ".")
}

func (real realFS) Path() string {
	return real.DirPath
}

func (real realFS) ResolvedPath() (string, error) {
	if !real.FollowLinks {
		return real.DirPath, nil
	}
	return filepath.EvalSymlinks(real.DirPath)
}

func (real realFS) CanSub(path string, entry fs.DirEntry) (bool, error) {
	child := filepath.Join(path, entry.Name())
	return fsx.IsDirectory(child, real.FollowLinks)
}

func (real realFS) Sub(path, rpath string, entry fs.DirEntry) FS {
	child := filepath.Join(path, entry.Name())
	return realFS{
		DirPath:     child,
		FollowLinks: real.FollowLinks,
	}
}
