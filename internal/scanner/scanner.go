// Package scanner provides Scan.
package scanner

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/pkg/errors"
	"github.com/tkw1536/ggman/internal/record"
)

// Options are options for a Scanner.
type Options struct {
	// Root is the root folder to scan
	Root string

	// ExtraRoots are extra root folders to scan on top of root.
	// These may be nil.
	ExtraRoots []string

	// Visit is a caller-settable function that is called once for each directory that is being scanned.
	// It determines function returns a pair of booleans match and bool.
	//
	// match indiciates that path should be returned in the array from Scan().
	// cont indicates if Scan() should continue scanning recursively.
	//
	// Visit may be nil.
	// In such a case, it is assumed to return the pair (true, true) for every invocation.
	Visit func(path string, context VisitContext) (match, cont bool)

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

// VisitContext represents the context of the Visit function
type VisitContext struct {
	// Root is the root this scan started from
	Root string

	// Depth determines the depth this scan was started from
	Depth int
}

// next returns a new ScanVisitContext that can be used for the next level of a scan
func (context VisitContext) next() (next VisitContext) {
	next.Root = context.Root
	next.Depth = context.Depth + 1
	return
}

// Scan creates a new Scanner, calls the Scan method, and returns a pair of results and error.
// This function is a convenience alternative to:
//
//  scanner := &Scanner{ScanOptions: options}
//  err := scanner.Scan()
//  results := scanner.Results()
func Scan(options Options) ([]string, error) {
	scanner := &Scanner{Options: options}

	// we cannot directly return here
	// Results() MUST be called after Scan()
	err := scanner.Scan()
	results := scanner.Results()

	return results, err
}

// Scanner is an object that can recursively scan a directory for subdirectories
// and return those matching a filter.
//
// Each Scanner may be used only once.
// Scanner is not safe for access by multiple goroutines.
type Scanner struct {
	used uint32 // 0 => not used, 1 => used; first to guarantee alignment

	Options

	record record.Record

	semaphore chan struct{}
	wg        sync.WaitGroup

	resultChan chan string
	errChan    chan error

	doneChan chan struct{}

	resultsSorted bool
	results       []string
}

// Scan scans the directory tree.
//
// When an error occurs, it continues blocking until all scanning routines have finished and returns an error.
// Each scanner should only be used once.
func (s *Scanner) Scan() error {
	if !atomic.CompareAndSwapUint32(&s.used, 0, 1) {
		panic("Scanner.Scan(): Attempted reuse")
	}

	if s.Visit == nil {
		s.Visit = func(path string, context VisitContext) (bool, bool) { return true, true }
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
	go s.scan(s.Root, VisitContext{
		Root:  s.Root,
		Depth: 0,
	})

	// scan all the extra roots
	s.wg.Add(len(s.ExtraRoots))
	for _, root := range s.ExtraRoots {
		go s.scan(root, VisitContext{
			Root:  root,
			Depth: 0,
		})
	}

	// start receiving results, and storing results
	go func() {
		for r := range s.resultChan {
			s.results = append(s.results, r)
		}
		s.doneChan <- struct{}{}
	}()

	// wait for the scan to finish, then return the results
	s.wg.Wait()
	close(s.resultChan)
	close(s.errChan)

	// wait for result collection to finish
	<-s.doneChan

	// return all the found repositories and any error
	return <-s.errChan
}

// Results returns all directories for which the match value of Visit function returns true.
// Directories are returned in sorted order.
//
// Results expects the Scan() function to have returned, but performs no checks that this is actually the case.
// When the Scan() function has not returned, the return value of this function is not defined.
func (s *Scanner) Results() []string {
	if !s.resultsSorted {
		sort.Strings(s.results)
		s.resultsSorted = true
	}
	return s.results
}

// scan performs a recursive scan of path
// one should do s.wg.Add(1) before each call of scan.
func (s *Scanner) scan(path string, context VisitContext) {

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
	}

	// if we already recorded this folder, return.
	if s.record.Record(path) {
		return
	}

	// execute the filter and act on it
	match, cont := s.Visit(path, context)
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

	// find the next context to visit
	nextContext := context.next()

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
		go s.scan(cpath, nextContext)
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
