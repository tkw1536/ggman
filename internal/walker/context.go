package walker

import (
	"io/fs"
	"sync"
)

// ProcessContext represents the current state of the filesystem
//
// An instance of ProcessContext should not be retained by callbacks after invocation.
type WalkContext interface {
	// Root returns the root filesystem the scan of this node started from
	Root() FS

	// Node returns the current node being operated on
	Node() FS

	// NodePath returns the path to the current node
	NodePath() string

	// Path returns the path from the root node to this node.
	Path() []string

	// Depth returns the depth of this node
	Depth() int

	// Snapshot updates the snapshot based on the update function.
	Snapshot(update func(snapshot interface{}) (value interface{}))

	// Mark marks the current node as a result with the given priority
	Mark(prio int)
}

var walkContextPool = sync.Pool{
	New: func() interface{} {
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
	snapshot interface{}
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

func (w walkContext) Mark(prio int) {
	w.w.reportResult(w.nodePath, prio)
}

func (w *walkContext) Snapshot(update func(snapshot interface{}) interface{}) {
	w.snapshot = update(w.snapshot)
}
