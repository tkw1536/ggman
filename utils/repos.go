package utils

import (
	"io/ioutil"
	"os"
	"path"
)

// Repos collects all git repositories in a given root folder
func Repos(root string, pattern string) (paths []string) {
	reposInternal(&paths, root, "", pattern)
	return
}

// reposInternal is the internal implementation of the Repos function
func reposInternal(paths *[]string, rootpath string, rootpattern string, pattern string) {
	// check if the current directory is a git repository
	// and if so, return a list containing this repo
	if isGitRepoRoot(rootpath) {
		// fmt.Printf("matching %q against %q\n", pattern, rootpattern) // for debugging
		if Matches(pattern, rootpattern) {
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
		if f.IsDir() {
			name := f.Name()
			reposInternal(paths, path.Join(rootpath, name), path.Join(rootpattern, name), pattern)
		}
	}
}

// isGitRepoRoot checks if a folder is the root of a git repository
func isGitRepoRoot(folder string) (isGit bool) {
	gitPath := path.Join(folder, ".git")
	s, err := os.Stat(gitPath)

	return !os.IsNotExist(err) && s.Mode().IsDir()
}
