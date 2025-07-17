//spellchecker:words walker
package walker

//spellchecker:words slices
import (
	"io/fs"
	"slices"
)

// WalkContext represents the current state of a Walker.
// It may additionally hold a snapshot of the state of type S.
//
// No caller should create a WalkContext by hand.
//
// Any instance of WalkContext must not be copied.
// Any instance of WalkContext must not be retained past the function call it was initially created for.
type WalkContext[S any] struct {
	w *Walker[S] // walker this context comes from

	root FS       // the root filesystem
	path []string // path to the current node

	node      FS     // current node
	nodePath  string // (unresolved) path to the current node
	rNodePath string // (resolved) path to the current node

	snapshot S // snapshot data carried for the current node
}

//
// create / delete
//

// it is guaranteed to be empty and have a nil context.
func (w *Walker[S]) getCtx() *WalkContext[S] {
	return w.ctxPool.Get().(*WalkContext[S])
}

// returnCtx returns a context to the walker-specific context pool.
// the context is reset before it is put back.
func (w *Walker[S]) returnCtx(ctx *WalkContext[S]) {
	ctx.w = nil
	ctx.root = nil
	ctx.node = nil
	ctx.nodePath = ""
	ctx.rNodePath = ""
	ctx.path = nil

	var nilSnapshot S
	ctx.snapshot = nilSnapshot

	w.ctxPool.Put(ctx)
}

// newContext initializes a new context from the context-specific pool.
func (w *Walker[S]) newContext(root FS) *WalkContext[S] {
	ctx := w.getCtx()

	ctx.w = w

	ctx.root = root
	ctx.node = root

	return ctx
}

// sub creates a new sub-context for the given.
func (w *WalkContext[S]) sub(entry fs.DirEntry) *WalkContext[S] {
	sub := w.w.getCtx()

	sub.w = w.w
	sub.root = w.root

	// create a new sub-path; which will allocate a new path for the child
	sub.path = slices.Clone(w.path)
	sub.path = append(sub.path, entry.Name())

	// return a new node
	sub.node = w.node.Sub(w.nodePath, w.rNodePath, entry)

	return sub
}

//
// public methods
//

// Root returns the node the current scan was started from.
func (w WalkContext[S]) Root() FS {
	return w.root
}

// Returns the current node being operated on.
func (w WalkContext[S]) Node() FS {
	return w.node
}

// Returns the path to the current node.
func (w WalkContext[S]) NodePath() string {
	return w.nodePath
}

// Path from the root node to this node.
func (w WalkContext[S]) Path() []string {
	return slices.Clone(w.path)
}

// Depth of this node, equivalent to len(Path()).
func (w WalkContext[S]) Depth() int {
	return len(w.path)
}

// Mark the current node as a result with the given priority.
// May be called multiple times, in which case the node is marked as a result multiple times.
func (w WalkContext[S]) Mark(priority float64) {
	w.w.reportResult(w.nodePath, w.rNodePath, priority)
}

// Update the snapshot contained in this WalkContext with the given function.
// when update panic()s, the behavior is undefined.
func (w *WalkContext[S]) Snapshot(update func(snapshot S) S) {
	w.snapshot = update(w.snapshot)
}
