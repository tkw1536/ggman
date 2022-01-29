// Package walker provides Walker, Scan and Sweep.
package walker

import (
	"errors"
	"io/fs"
	"sort"
	"sync"
	"sync/atomic"

	"github.com/tkw1536/ggman/internal/record"
	"github.com/tkw1536/ggman/internal/sema"
)

// Walker is an object that can recursively scan a directory for subdirectories
// and return those determined by a used.
//
// A cursor uses concurrent operations when possible, as determiend by Callbacks
//
// Each Walker may be used only once; it may not in and of itself be used concurrently
type Walker struct {
	state uint32 // 0 => initial, 1 => in Walk(), 2 => done

	record record.Record // contains visited nodes

	wg        sync.WaitGroup // tracks which scanning processes are
	semaphore sync.Locker    // tracks how many scanners are alive

	errChan    chan error      // contains error in the buffer
	resultChan chan walkResult // contains results temporarily

	results []string  // results
	scores  []float64 // scores

	Params  Params
	Process Process
}

// Params are parameters for a walk across a filesystem
type Params struct {
	// Root is the root filesystem to begin the walk at
	Root FS

	// ExtraRoots are extra root folders to walk on top of root.
	// These may be nil.
	ExtraRoots []FS

	// MaxParallel is maximum number of nodes that will be scanned in parallel.
	// Zero or negative values are treated as no limit.
	MaxParallel int

	// BufferSize is an integer that can be used to optimize internal behavior.
	// It should be larger than the average number of expected results.
	// Set to 0 to disable.
	BufferSize int
}

// Process represents a walking process of a walker
//
// Processes should not retain references to VisitContexts beyond the invocation of the function
type Process interface {
	// Visit is called for every node that is being visited.
	// It is the first function called for each node.
	//
	// It receives several parameters:
	//
	// The node being visited and the appropriate context.
	//
	// A function that can be used to mark this node as a result.
	// prio is a prioritisation of a node that is used for sorting; see the Result() method for details.
	//
	// Visit should return three things.
	//
	// Snapshot is an arbitrary object that captures that current state of the process
	// It is maintained throughout the processing of one node, and returned to the parent node (when being processed concurrently)
	//
	// shouldVisitChildren determines if any children of this node should be visited or if the process should stop.
	// When shouldVisitChildren is false, no other functions are called for this node, and the snapshot is returned to the parent (if any) immediatly.
	//
	// Err is any error that may occur, and should typically be nil.
	// An error immediatly cause iteration on this node to be aborted, and the first error of any node will be returned to the caller of Walk.
	Visit(context WalkContext) (shouldVisitChildren bool, err error)

	// VisitChild is called to determine if and how a child node should be processed.
	//
	// A child entry is valid if it can be recursivly processed (i.e. is a directory).
	//
	// When child is valid, it determines how the child should be processed; otherwise action is ignored.
	VisitChild(child fs.DirEntry, valid bool, context WalkContext) (action Step, err error)

	// AfterVisitChild is called after a child has been visited syncronously.
	//
	// It is passed to special values, the returned snapshot (as returned from AfterVisit / Visit) and if the child was processed properly.
	// The child was processed improperly when any of the Process functions on it returned an error, listing a directory failed, or it was already processed before (loop detection). In these cases resultValue is nil.
	AfterVisitChild(child fs.DirEntry, resultValue interface{}, resultOK bool, context WalkContext) (err error)

	// AfterVisit is called after all children have been visited (or scheduled to be visited).
	// It is not called for the case where Visit returns shouldVisitChildren = false.
	//
	// result can be used to mark the current node, see also Visit.
	//
	// The returnValue returned from AfterVisit is passed to parent(s) if any.
	AfterVisit(context WalkContext) (err error)
}

// Step describes how a child node should be processed
type Step int

const (
	// DoNothing ignores the child node, and continue with the next node.
	DoNothing Step = iota
	// DoSync syncronously processes the child node.
	// Once processing the child node has finished the AfterChild() function will be called.
	DoSync
	// DoConcurrent queues the child node to be processed concurrently.
	// The current node will node wait for
	DoConcurrent
)

// Walk begins recursively walking the directory tree starting at the roots defined in Config.
//
// Walk must be called at most once for each Walker and will panic() if called multiple times.
//
// This function is untested because the tests for Scan and Sweep suffice.
func (w *Walker) Walk() error {
	// state of the walker
	if !atomic.CompareAndSwapUint32(&w.state, 0, 1) {
		panic("Walker.Walk(): Attempted reuse")
	}
	defer atomic.StoreUint32(&w.state, 2)

	// configure concurrency
	w.semaphore = sema.NewSemaphore(w.Params.MaxParallel)

	// create channels for result & error
	w.resultChan = make(chan walkResult, w.Params.BufferSize)
	w.errChan = make(chan error, 1)

	w.wg = sync.WaitGroup{}

	// scan the root
	w.wg.Add(1)
	go w.walkRoot(w.Params.Root)

	// scan all the extra roots
	w.wg.Add(len(w.Params.ExtraRoots))
	for _, root := range w.Params.ExtraRoots {
		go w.walkRoot(root)
	}

	// start another goroutine to begin receiving results
	// then sort these and mark everything as finished
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)

		results := make([]walkResult, 0)
		for r := range w.resultChan {
			results = append(results, r)
		}

		// sort the slice
		sort.Slice(results, func(i, j int) bool {
			return results[i].LessThan(results[j])
		})

		// store the real (textual) results
		w.results = make([]string, len(results))
		w.scores = make([]float64, len(results))
		for i, r := range results {
			w.results[i] = r.Node
			w.scores[i] = r.Score
		}

	}()

	// wait for the scan to finish, then return the results
	w.wg.Wait()
	close(w.resultChan)
	close(w.errChan)

	// wait for result collection to finish
	<-doneChan

	// return all the found repositories and any error
	return <-w.errChan
}

// walkRoot starts a walk through the provided root
func (w *Walker) walkRoot(root FS) {
	w.semaphore.Lock()
	defer w.semaphore.Unlock()

	ctx := w.newContext(root)
	defer walkContextPool.Put(ctx)

	w.walk(true, ctx)
}

// walk walks recursively through the provided context
func (w *Walker) walk(sync bool, ctx *walkContext) (ok bool) {
	defer w.wg.Done()

	if !sync {
		w.semaphore.Lock()
		defer w.semaphore.Unlock()
	}

	// get the (normalized) path
	path, err := ctx.node.Path()
	if err != nil {
		w.reportError(err)
		return false
	}
	ctx.nodePath = path

	// bail out if we already visited this node!
	if w.record.Record(path) {
		return true
	}

	shouldVisitChildren, err := w.Process.Visit(ctx)
	if err != nil {
		w.reportError(err)
		return false
	}
	if !shouldVisitChildren {
		return false
	}

	// list all the entries and folders in this directory
	entries, err := ctx.node.Read(path)
	if err != nil {
		w.reportError(err)
		return false
	}

	// iterate over all the entries and figure out what to do!
	w.wg.Add(len(entries))
	for _, entry := range entries {
		// check if we have a valid child!
		valid, err := ctx.node.CanSub(path, entry)
		if err != nil {
			w.reportError(err)
			continue
		}

		var action Step
		action, err = w.Process.VisitChild(entry, valid, ctx)
		if err != nil {
			w.reportError(err)
			return false
		}

		switch {
		case action == DoNothing || !valid:
			w.wg.Done()
		case action == DoConcurrent:
			// work asyncronously and discard the parent!
			go func(cctx *walkContext) {
				defer walkContextPool.Put(cctx)
				w.walk(false, cctx)
			}(ctx.sub(entry))
		case action == DoSync:
			// run the child processing!
			ok, value := func(cctx *walkContext) (bool, interface{}) {
				defer walkContextPool.Put(cctx)

				ok := w.walk(true, cctx)
				return ok, cctx.snapshot
			}(ctx.sub(entry))

			if err := w.Process.AfterVisitChild(entry, value, ok, ctx); err != nil {
				w.reportError(err)
				return false
			}
		default:
			w.reportError(ErrUnknownAction)
			return false
		}
	}

	// we have finished all (syncronous) operations
	if err := w.Process.AfterVisit(ctx); err != nil {
		w.reportError(err)
		return false
	}
	return true
}

// reportResults reports the given node as a result.
// might block until a slot in the results is available.
func (w *Walker) reportResult(node string, score float64) {
	w.resultChan <- walkResult{Node: node, Score: score}
}

// reportErrors reports the provided error to the caller of Walk()
// When another error has already occured, does nothing
func (w *Walker) reportError(err error) {
	select {
	case w.errChan <- err:
	default:
	}
}

// Results returns all nodes which have been marked as a result.
// Directories are returned in sorted order; sorted first ascending by priority then by lexiographically by node.
// Each call to result returns a new copy of the results.
//
// Results expects the Scan() function to have returned, and will panic if this is not the case.
func (w *Walker) Results() []string {
	if atomic.LoadUint32(&w.state) != 2 {
		panic("Walker.Walk(): Results() called before Walk() returned")
	}

	results := make([]string, len(w.results))
	copy(results, w.results)
	return results
}

// Scores returns the scores which have been marked as a result.
// They are returned in the same order as Results()
//
// Results expects the Scan() function to have returned, and will panic if this is not the case.
func (w *Walker) Scores() []float64 {
	if atomic.LoadUint32(&w.state) != 2 {
		panic("Walker.Walk(): Scores() called before Walk() returned")
	}

	scores := make([]float64, len(w.scores))
	copy(scores, w.scores)
	return scores
}

var ErrUnknownAction = errors.New("Process.BeforeChild(): Unknown action")

// walkResult represents an internal result of the wlak function
type walkResult struct {
	Node  string
	Score float64
}

// LessThan returns true if w should occur before v when sorting a slice of walkResults
//
// Sorting first occurs descending by Score, then ascending by lexiographic order on Node.
func (w walkResult) LessThan(v walkResult) bool {
	switch {
	case w.Score < v.Score:
		return false
	case w.Score > v.Score:
		return true
	case w.Node < v.Node:
		return true
	default:
		return false
	}
}