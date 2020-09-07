package env

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/tkw1536/ggman/git"
)

// reposInternal is the internal implementation of the Repos function
func reposInternal(git git.Git, paths *[]string, rootpath string, rootpattern string, pattern string, allowSymlinks bool) {
	// check if the current directory is a git repository
	// and if so, return a list containing this repo
	if git.IsRepositoryQuick(rootpath) {
		// fmt.Printf("matching %q against %q\n", pattern, rootpattern) // for debugging
		if MatchesString(pattern, rootpattern) {
			*paths = append(*paths, rootpath)
		}
		return
	}

	// read all the folders in this directory
	// but bail out if an error occurs
	files, err := ioutil.ReadDir(rootpath)
	if err != nil {
		return
	}

	// iterate over all the subdirectories
	// and recursively call this function on all sub-directories
	for _, f := range files {
		// get name and path of the file
		name := f.Name()
		fn := path.Join(rootpath, name)

		// resolve symlink (if needed)
		wasLink, _ := resolveSymlinks(fn, &f)

		// and if it is a directory, recurse
		if f.IsDir() {
			reposInternal(git, paths, fn, path.Join(rootpattern, name), pattern, !wasLink)
		}
	}
}

// resolveSymlinks resolves the symlinks of a fileinfo
func resolveSymlinks(fn string, info *os.FileInfo) (wasLink bool, err error) {
	newinfo := *info

	// if we have a symlink, resolve the symlink
	if newinfo.Mode()&os.ModeSymlink != 0 {
		wasLink = true

		// resolve the directory
		base := filepath.Dir(fn)

		// resolve the new path
		fn, err = filepath.EvalSymlinks(fn)
		if err != nil {
			return
		}

		// if filepath is not absolute, make it absolute relative to the base
		if !filepath.IsAbs(fn) {
			fn = filepath.Join(base, fn)
		}

		newinfo, err = os.Stat(fn)
		if err == nil {
			*info = newinfo
		}
	} else {
		wasLink = false
	}

	// and return if we had an error
	return
}
