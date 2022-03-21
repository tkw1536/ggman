package walker

import (
	"io/fs"

	"golang.org/x/exp/slices"
)

// context implements WalkerContext
type context[S any] struct {
	w *Walker[S] // walker this context comes from

	root FS       // the root filesystem
	path []string // path to the current node

	node     FS     // current node
	nodePath string // name of the current node

	snapshot S // snapshot data carried for the current node
}

//
// create / delete
//

// getCtx fetches a new (uninitialized) context from the walker-specific context pool.
// it is guaranteed to be empty and have a nil context
func (w *Walker[S]) getCtx() *context[S] {
	return w.ctxPool.Get().(*context[S])
}

// returnCtx returns a context to the walker-specific context pool.
// the context is reset before it is put back.
func (w *Walker[S]) returnCtx(ctx *context[S]) {
	ctx.w = nil
	ctx.root = nil
	ctx.node = nil
	ctx.nodePath = ""
	ctx.path = nil

	var nilSnapshot S
	ctx.snapshot = nilSnapshot

	w.ctxPool.Put(ctx)
}

// newContext initializes a new context from the context-specific pool
func (w *Walker[S]) newContext(root FS) *context[S] {
	ctx := w.getCtx()

	ctx.w = w

	ctx.root = root
	ctx.node = root

	return ctx
}

// sub creates a new sub-context for the given
func (w *context[S]) sub(entry fs.DirEntry) *context[S] {
	sub := w.w.getCtx()

	sub.w = w.w
	sub.root = w.root

	// create a new sub-path; which will allocate a new path for the child
	sub.path = slices.Clone(w.path)
	sub.path = append(w.path, entry.Name())

	// return a new node
	sub.node = w.node.Sub(w.nodePath, entry)

	return sub
}

func (w context[S]) Root() FS {
	return w.root
}

func (w context[S]) Node() FS {
	return w.node
}

func (w context[S]) NodePath() string {
	return w.nodePath
}

func (w context[S]) Path() []string {
	return slices.Clone(w.path)
}

func (w context[S]) Depth() int {
	return len(w.path)
}

func (w context[S]) Mark(prio float64) {
	w.w.reportResult(w.nodePath, prio)
}

func (w *context[S]) Snapshot(update func(snapshot S) S) {
	w.snapshot = update(w.snapshot)
}
