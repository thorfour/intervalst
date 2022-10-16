// Package interval implements a generic interval tree with left leaning red-black (LLRB) algorithm.
package interval

import "sync"

// SearchTree is a generic type representing the Interval Search Tree
// where V is a generic value type, and T is a generic interval key type.
type SearchTree[V, T any] struct {
	mu   sync.RWMutex // used to serialize read and write operations
	root *node[V, T]
	cmp  CmpFunc[T]
}

// NewSearchTree returns an initialized interval search tree.
// The cmp parameter is used for comparing total order of the interval key type T
// when inserting or looking up an interval in the tree.
// For more details on cmp, see the CmpFunc type.
func NewSearchTree[V, T any](cmp CmpFunc[T]) *SearchTree[V, T] {
	if cmp == nil {
		panic("NewSearchTree: comparison function cmp cannot be nil")
	}
	return &SearchTree[V, T]{
		cmp: cmp,
	}
}

func (st *SearchTree[V, T]) rotateLeft(h *node[V, T]) *node[V, T] {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	x.maxEnd = h.maxEnd
	h.color = red

	st.updateMaxEnd(h)
	return x
}

func (st *SearchTree[V, T]) rotateRight(h *node[V, T]) *node[V, T] {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	x.maxEnd = h.maxEnd
	h.color = red

	st.updateMaxEnd(h)
	return x
}

func (st *SearchTree[V, T]) updateMaxEnd(h *node[V, T]) {
	h.maxEnd = h.interval.end
	if h.left != nil && st.cmp.gt(h.left.maxEnd, h.maxEnd) {
		h.maxEnd = h.left.maxEnd
	}

	if h.right != nil && st.cmp.gt(h.right.maxEnd, h.maxEnd) {
		h.maxEnd = h.right.maxEnd
	}
}

func (st *SearchTree[V, T]) balanceNode(h *node[V, T]) *node[V, T] {
	if isRed(h.right) && !isRed(h.left) {
		h = st.rotateLeft(h)
	}

	if isRed(h.left) && isRed(h.left.left) {
		h = st.rotateRight(h)
	}

	if isRed(h.left) && isRed(h.right) {
		flipColors(h)
	}

	return h
}

func (st *SearchTree[V, T]) moveRedLeft(h *node[V, T]) *node[V, T] {
	flipColors(h)
	if h.right != nil && isRed(h.right.left) {
		h.right = st.rotateRight(h.right)
		h = st.rotateLeft(h)
		flipColors(h)
	}
	return h
}

func (st *SearchTree[V, T]) moveRedRight(h *node[V, T]) *node[V, T] {
	flipColors(h)
	if h.left != nil && isRed(h.left.left) {
		h = st.rotateRight(h)
		flipColors(h)
	}
	return h
}

func (st *SearchTree[V, T]) fixUp(h *node[V, T]) *node[V, T] {
	st.updateMaxEnd(h)

	return st.balanceNode(h)
}