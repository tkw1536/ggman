package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/pkg/errors"
)

// ScanOptions are options for a Scanner.
type ScanOptions struct {
	// Root is the root folder to scan
	Root string

	// Filter is a caller-settable function that determines if a directory should be accepted.
	// The filter function returns a pair of booleans match and bool.
	// match indiciates that path should be returned in the array from Scan().
	// cont indicates if Scan() should continue scanning recursively.
	Filter func(path string) (match, cont bool)

	// When FollowLinks is set, the scanner will follow symbolic links and detect cycles.
	FollowLinks bool

	// MaxParallel is maximum number of directories that will be scanned in parallel.
	// Leave blank for unlimited.
	MaxParallel int

	// BufferSize is an integer that can be used to optimize internal behavior.
	// It should be larger than the average number of expected results.
	// Set to 0 to disable.
	BufferSize int
}

// Scan creates a new Scanner and call the Scan method.
// This function is a convenience alternative to:
//
//  scanner := &Scanner{ScanOptions: options}
//  scanner.Scan()
func Scan(options ScanOptions) ([]string, error) {
	scanner := &Scanner{ScanOptions: options}
	return scanner.Scan()
}

// Scanner is an object that can recursively scan a directory for subdirectories
// and return those matching a filter.
//
// Each Scanner may be used only once.
type Scanner struct {
	ScanOptions

	used   bool
	record Record

	semaphore chan struct{}
	wg        sync.WaitGroup

	resultChan chan string
	errChan    chan error

	doneChan chan struct{}
	results  []string
}

// Scan scans the directory tree and returns all directories for which the match value of Filter returns true.
// They are returned in alphabetical order.
//
// When an error occurs, it continues blocking until all scanning routines have finished and returns an error.
// Each scanner should only be used once.
func (s *Scanner) Scan() ([]string, error) {
	if s.used {
		panic("Scanner.Scan(): Attempted reuse")
	}
	s.used = true

	if s.Filter == nil {
		s.Filter = func(path string) (bool, bool) { return true, true }
	}

	// create a channel for results and the repos themselves
	s.resultChan = make(chan string, s.BufferSize)
	s.results = make([]string, 0, s.BufferSize)
	s.doneChan = make(chan struct{})

	s.semaphore = nil
	if s.MaxParallel != 0 {
		s.semaphore = make(chan struct{}, s.MaxParallel)
	}

	// capture errors, this should return only the first error
	s.errChan = make(chan error, 1)

	s.wg = sync.WaitGroup{}
	s.wg.Add(1)
	go s.scan(s.Root)

	// start receiving results, sort them afterwards.
	// one could sort while inserting, but as the order is random
	// this would result in a lot more copy operations.
	go func() {
		for r := range s.resultChan {
			s.results = append(s.results, r)
		}
		sort.Strings(s.results)
		s.doneChan <- struct{}{}
	}()

	// wait for the scan to finish, then return the results
	s.wg.Wait()
	close(s.resultChan)
	close(s.errChan)

	// wait for result collection to finish
	<-s.doneChan

	// return all the found repositories and any error
	return s.results, <-s.errChan
}

// scan performs a recursive scan of path
// one should do s.wg.Add(1) before each call of scan.
func (s *Scanner) scan(path string) {

	// aquire the semaphore
	if s.semaphore != nil {
		s.semaphore <- struct{}{}
		defer func() {
			<-s.semaphore
		}()
	}
	defer s.wg.Done()

	// when following links is enabled, we need to evaluate any symlinks in the path
	if s.FollowLinks {
		var err error
		path, err = filepath.EvalSymlinks(path)
		if err != nil {
			select {
			case s.errChan <- err:
			default:
			}
			return
		}

		// if we already recorded this folder, return.
		if s.record.Record(path) {
			return
		}
	}

	// execute the filter and act on it
	match, cont := s.Filter(path)
	if match {
		s.resultChan <- path
	}
	if !cont {
		return
	}

	// list all the files and folders in this directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		select {
		case s.errChan <- err:
		default:
		}
		return
	}

	// iterate over all the files in this folder
	// having this parallel just adds extra overhead, so we do not do this
	for _, f := range files {
		name := f.Name()
		cpath := filepath.Join(path, name)

		// check if we have a directory
		isDir, err := IsDirectory(cpath, s.FollowLinks)
		if err != nil {
			select {
			case s.errChan <- err:
			default:
			}
			continue
		}

		// if we don't have a directory, we're done
		if !isDir {
			continue
		}

		// scan the folder recursively
		s.wg.Add(1)
		go s.scan(cpath)
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
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, errors.Wrap(err, "Stat failed")
	}

	return stat.IsDir(), nil
}
