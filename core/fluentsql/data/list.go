// Package data
// Fully persistent data structures. A persistent data structure is a data
// structure that always preserves the previous version of itself when
// it is modified. Such data structures are effectively immutable,
// as their operations do not update the structure in-place, but instead
// always yield a new structure.
//
// Persistent
// data structures typically share structure among themselves.  This allows
// operations to avoid copying the entire data structure.
package data

// List is a persistent list of possibly heterogeneous values.
type List interface {
	// IsNil returns true if the list is empty
	IsNil() bool

	// Cons returns a new list with val as the head
	Cons(val Any) List

	// Head returns the first element of the list;
	// panics if the list is empty
	Head() Any

	// Tail returns a list with all elements except the head;
	// panics if the list is empty
	Tail() List

	// Size returns the list's length.  This takes O(1) time.
	Size() int

	// ForEach executes a callback for each value in the list.
	ForEach(f func(Any))

	// Reverse returns a list whose elements are in the opposite order as
	// the original list.
	Reverse() List
}

// Immutable (i.e. persistent) list
type list struct {
	depth int // the number of nodes after, and including, this one
	value Any
	tail  *list
}

// An empty list shared by all lists
var nilList = &list{}

// NewList returns a new, empty list.  The result is a singly linked
// list implementation.  All lists share an empty tail, so allocating
// empty lists is efficient in time and memory.
func NewList() List {
	return nilList
}

func (l *list) IsNil() bool {
	return l == nilList
}

func (l *list) Size() int {
	return l.depth
}

func (l *list) Cons(val Any) List {
	var xs list
	xs.depth = l.depth + 1
	xs.value = val
	xs.tail = l
	return &xs
}

func (l *list) Head() Any {
	if l.IsNil() {
		panic("Called Head() on an empty list")
	}

	return l.value
}

func (l *list) Tail() List {
	if l.IsNil() {
		panic("Called Tail() on an empty list")
	}

	return l.tail
}

// ForEach executes a callback for each value in the list
func (l *list) ForEach(f func(Any)) {
	if l.IsNil() {
		return
	}
	f(l.Head())
	l.Tail().ForEach(f)
}

// Reverse returns a list with elements in opposite order as this list
func (l *list) Reverse() List {
	reversed := NewList()
	l.ForEach(func(v Any) { reversed = reversed.Cons(v) })
	return reversed
}
