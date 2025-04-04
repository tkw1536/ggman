// Package walker provides Walker, Scan and Sweep.
//
//spellchecker:words walker
package walker

//spellchecker:words errors slices strings sync atomic github ggman internal record pkglib sema
import (
	"errors"
	"io/fs"
	"slices"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/tkw1536/ggman/internal/record"
	"github.com/tkw1536/pkglib/sema"
)

// Walker is an object that can recursively operate on all subdirectories of a directory and score those matching a specific criterion.
// The criterion is determined by the Process parameter.
//
// Process also determines if the process can operate on multiple directories concurrently.
// Parameters determine the initial root directory (or directories) to start with, and what level of concurrency the walker may make use of.
//
// Each Walker may be used only once.
// A typical use of a walker looks like:
//
//	w := Walker{/* ... */}
//	if err := w.Walk(); err != nil {
//	  return err
//	}
//	results, scores := w.Results(), w.Scores()
type Walker[S any] struct {
	state atomic.Uint32 // see walker* constants

	visited record.Record // contains visited nodes

	wg        sync.WaitGroup // tracks which scanning processes are
	semaphore sema.Semaphore // tracks how many scanners are alive

	errChan    chan error      // contains error in the buffer
	resultChan chan walkResult // contains results temporarily

	ctxPool sync.Pool // pool for *context[S] objects

	paths  []string
	rPaths []string
	scores []float64

	Params  Params
	Process Process[S]
}

// walker* constants represent the state of the walker.
const (
	walkerUnused   uint32 = iota // walker has not yet been used
	walkerRunning                // Walk function has been called, but not yet returned
	walkerFinished               // Walk has returned
)

// Params are parameters for a walk across a filesystem.
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

// Process determines the behavior of a Walker.
//
// Each process may hold intermediate state of type S.
// Processes should not retain references to VisitContexts (or state) beyond the invocation of each method.
type Process[S any] interface {
	// Visit is called for every node that is being visited.
	// It is the first function called for each node.
	//
	// It receives a context, representing the node being visited.
	//
	// Visit should return three things.
	//
	// Snapshot is an arbitrary object that captures that current state of the process
	// It is maintained throughout the processing of one node, and returned to the parent node (when being processed concurrently)
	//
	// shouldVisitChildren determines if any children of this node should be visited or if the process should stop.
	// When shouldVisitChildren is false, no other functions are called for this node, and the snapshot is returned to the parent (if any) immediately.
	//
	// Err is any error that may occur, and should typically be nil.
	// An error immediately causes iteration on this node to be aborted, and the first error of any node will be returned to the caller of Walk.
	Visit(context WalkContext[S]) (shouldVisitChildren bool, err error)

	// VisitChild is called to determine if and how a child node should be processed.
	//
	// A child entry is valid if it can be recursively processed (i.e. is a directory).
	//
	// When child is valid, it determines how the child should be processed; otherwise action is ignored.
	VisitChild(child fs.DirEntry, valid bool, context WalkContext[S]) (action Step, err error)

	// AfterVisitChild is called after a child has been visited synchronously.
	//
	// It is passed to special values, the returned snapshot (as returned from AfterVisit / Visit) and if the child was processed properly.
	// The child was processed improperly when any of the Process functions on it returned an error, listing a directory failed, or it was already processed before (loop detection). In these cases resultValue is nil.
	AfterVisitChild(child fs.DirEntry, resultValue any, resultOK bool, context WalkContext[S]) (err error)

	// AfterVisit is called after all children have been visited (or scheduled to be visited).
	// It is not called for the case where Visit returns shouldVisitChildren = false.
	//
	// result can be used to mark the current node, see also Visit.
	//
	// The returnValue returned from AfterVisit is passed to parent(s) if any.
	AfterVisit(context WalkContext[S]) (err error)
}

// Step describes how a child node should be processed.
type Step int

const (
	// DoNothing ignores the child node, and continue with the next node.
	DoNothing Step = iota
	// DoSync synchronously processes the child node.
	// Once processing the child node has finished the AfterChild() function will be called.
	DoSync
	// The current node will node wait for.
	DoConcurrent
)

// WalkContext represents the current state of a Walker.
// It may additionally hold a snapshot of the state of type S.
//
// Any instance of WalkContext should not be retained past any method it is passed to.
type WalkContext[S any] interface {
	// Root node this instance of the scan started from
	Root() FS

	// Current node being operated on
	Node() FS

	// Path to the current node
	NodePath() string

	// Path from the root node to this node
	Path() []string

	// Depth of this node, equivalent to len(Path())
	Depth() int

	// Update the snapshot corresponding to the current context
	Snapshot(update func(snapshot S) (value S))

	// Mark the current node as a result with the given priority.
	// May be called multiple times, in which case the node is marked as a result multiple times.
	Mark(priority float64)
}

// Walk begins recursively walking the directory tree starting at the roots defined in Config.
//
// Walk must be called at most once for each Walker and will panic() if called multiple times.
//
// This function is untested because the tests for Scan and Sweep suffice.
func (w *Walker[S]) Walk() error {
	// state of the walker
	if !w.state.CompareAndSwap(walkerUnused, walkerRunning) {
		panic("Walker.Walk: Attempted reuse")
	}
	defer w.state.Store(walkerFinished)

	// setup a pool for new contexts
	w.ctxPool.New = func() any {
		return new(context[S])
	}

	// configure concurrency
	w.semaphore = sema.New(w.Params.MaxParallel)

	// create channels for result & error
	w.resultChan = make(chan walkResult, w.Params.BufferSize)
	w.errChan = make(chan error, 1)

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
		slices.SortFunc(results, walkResult.compare)

		// store the real (textual) results
		w.paths = make([]string, len(results))
		w.rPaths = make([]string, len(results))
		w.scores = make([]float64, len(results))
		for i, r := range results {
			w.paths[i] = r.NodePath
			w.rPaths[i] = r.NodeRPath
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

// walkRoot starts a walk through the provided root.
func (w *Walker[S]) walkRoot(root FS) {
	w.semaphore.Lock()
	defer w.semaphore.Unlock()

	ctx := w.newContext(root)
	defer w.returnCtx(ctx)

	w.walk(true, ctx)
}

// walk walks recursively through the provided context.
func (w *Walker[S]) walk(sync bool, ctx *context[S]) (ok bool) {
	defer w.wg.Done()

	if !sync {
		w.semaphore.Lock()
		defer w.semaphore.Unlock()
	}

	// get the (normalized) path
	path, err := ctx.node.ResolvedPath()
	if err != nil {
		w.reportError(err)
		return false
	}
	ctx.rNodePath = path
	ctx.nodePath = ctx.node.Path()

	// bail out if we already visited this node!
	if w.visited.Record(path) {
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
			// work asynchronously and discard the parent!
			go func(cContext *context[S]) {
				defer w.returnCtx(cContext)
				w.walk(false, cContext)
			}(ctx.sub(entry))
		case action == DoSync:
			// run the child processing!
			ok, value := func(cContext *context[S]) (bool, any) {
				defer w.returnCtx(cContext)

				ok := w.walk(true, cContext)
				return ok, cContext.snapshot
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

	// we have finished all (synchronous) operations
	if err := w.Process.AfterVisit(ctx); err != nil {
		w.reportError(err)
		return false
	}
	return true
}

// reportResults reports the node with the given path and resolved paths.
// might block until a slot in the results is available.
func (w *Walker[S]) reportResult(path, rpath string, score float64) {
	w.resultChan <- walkResult{NodePath: path, NodeRPath: rpath, Score: score}
}

// When another error has already occurred, does nothing.
func (w *Walker[S]) reportError(err error) {
	select {
	case w.errChan <- err:
	default:
	}
}

// DEPRECATED.
func (w *Walker[S]) Results() []string {
	return w.Paths(true)
}

// Paths returns the path of all nodes which have been marked as a result.
//
// When resolved is true, returns the normalized (resolved) paths; else the non-normalized versions are returned.
// Directories are returned in sorted order; sorted first ascending by priority then by lexicographically by resolved node path.
// Each call to result returns a new copy of the results.
//
// Paths expects the Scan() function to have returned, and will panic if this is not the case.
func (w *Walker[S]) Paths(resolved bool) []string {
	if w.state.Load() != walkerFinished {
		panic("Walker.Paths: Results() called before Walk() returned")
	}

	if resolved {
		return slices.Clone(w.rPaths)
	} else {
		return slices.Clone(w.paths)
	}
}

// Scores returns the scores which have been marked as a result.
// They are returned in the same order as Results()
//
// Results expects the Scan() function to have returned, and will panic if this is not the case.
func (w *Walker[S]) Scores() []float64 {
	if w.state.Load() != walkerFinished {
		panic("Walker.Walk: Scores() called before Walk() returned")
	}

	return slices.Clone(w.scores)
}

var ErrUnknownAction = errors.New("Process.BeforeChild: Unknown action")

// walkResult represents an internal result of the walk function.
type walkResult struct {
	NodePath  string
	NodeRPath string
	Score     float64
}

// Compare compares w with v and returns -1 when w should occur before v, 1 if v before w, and 0 if they are equal
//
// Sorting first occurs descending by Score, then ascending by lexicographic order on Node.
func (w walkResult) LessThan(v walkResult) bool {
	switch {
	case w.Score < v.Score:
		return false
	case w.Score > v.Score:
		return true
	case w.NodeRPath < v.NodeRPath:
		return true
	default:
		return false
	}
}

// compare implements a strict weak ordering to compare w with v in the context of sorting a result list.
// It returns 1 if w should occur after v, -1 if w before v, and 0 if they are equal.
func (w walkResult) compare(v walkResult) int {
	cmp := v.Score - w.Score
	switch {
	case cmp == 0:
		return strings.Compare(w.NodeRPath, v.NodeRPath)
	case cmp < 0:
		return -1
	case cmp > 0:
		return 1
	}
	panic("never reached")
}
