package walker

import (
	"io/fs"
	"sync"
)

// WalkContext represents the current state of a Walk.
//
// Any instance of WalkContext should not be retained past any callback it is passed in.
type WalkContext interface {
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
	Snapshot(update func(snapshot any) (value any))

	// Mark the current node as a result with the given priority.
	// May be called multiple times, in which case the node is marked as a result multiple times.
	Mark(prio float64)
}

var walkContextPool = sync.Pool{
	New: func() any {
		return new(walkContext)
	},
}

// walkContext implements WalkContext
type walkContext struct {
	w *Walker

	root FS

	node     FS
	nodePath string

	path     []string
	snapshot any
}

func (w *Walker) newContext(root FS) *walkContext {
	ctx := walkContextPool.Get().(*walkContext)

	ctx.w = w

	ctx.root = root

	ctx.node = root
	ctx.nodePath = "" // never used!

	ctx.path = nil
	ctx.snapshot = nil

	return ctx
}

func (w walkContext) sub(entry fs.DirEntry) *walkContext {
	sub := walkContextPool.Get().(*walkContext)

	sub.w = w.w
	sub.root = w.root

	// create a new sub-path; which will allocate a new path for the child
	sub.path = make([]string, len(w.path)+1)
	copy(sub.path, w.path)
	sub.path[len(w.path)] = entry.Name()

	sub.node = w.node.Sub(w.nodePath, entry)
	sub.snapshot = nil

	return sub
}

func (w walkContext) Root() FS {
	return w.root
}

func (w walkContext) Node() FS {
	return w.node
}

func (w walkContext) NodePath() string {
	return w.nodePath
}

func (w walkContext) Path() []string {
	path := make([]string, len(w.path))
	copy(path, w.path)
	return path
}

func (w walkContext) Depth() int {
	return len(w.path)
}

func (w walkContext) Mark(prio float64) {
	w.w.reportResult(w.nodePath, prio)
}

func (w *walkContext) Snapshot(update func(snapshot any) any) {
	w.snapshot = update(w.snapshot)
}
