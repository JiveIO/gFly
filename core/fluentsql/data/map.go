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

import (
	"bytes"
	"fmt"
)

// Any is a shorthand for Go's verbose interface{} type.
type Any interface{}

// A Map associates unique keys (type string) with values (type Any).
type Map interface {
	// IsNil returns true if the Map is empty
	IsNil() bool

	// Set returns a new map in which key and value are associated.
	// If the key didn't exist before, it's created; otherwise, the
	// associated value is changed.
	// This operation is O(log N) in the number of keys.
	Set(key string, value Any) Map

	// Delete returns a new map with the association for key, if any, removed.
	// This operation is O(log N) in the number of keys.
	Delete(key string) Map

	// Lookup returns the value associated with a key, if any.  If the key
	// exists, the second return value is true; otherwise, false.
	// This operation is O(log N) in the number of keys.
	Lookup(key string) (Any, bool)

	// Size returns the number of key value pairs in the map.
	// This takes O(1) time.
	Size() int

	// ForEach executes a callback on each key value pair in the map.
	ForEach(f func(key string, val Any))

	// Keys returns a slice with all keys in this map.
	// This operation is O(N) in the number of keys.
	Keys() []string

	String() string
}

// Immutable (i.e. persistent) associative array
const childCount = 8
const shiftSize = 3

type tree struct {
	count    int
	hash     uint64 // hash of the key (used for tree balancing)
	key      string
	value    Any
	children [childCount]*tree
}

var nilMap = &tree{}

// Recursively set nilMap's subtrees to point at itself.
// This eliminates all nil pointers in the map structure.
// All map nodes are created by cloning this structure so
// they avoid the problem too.
func init() {
	for i := range nilMap.children {
		nilMap.children[i] = nilMap
	}
}

// NewMap allocates a new, persistent map from strings to values of
// any type.
// This is currently implemented as a path-copying binary tree.
func NewMap() Map {
	return nilMap
}

func (m *tree) IsNil() bool {
	return m == nilMap
}

// clone returns an exact duplicate of a tree node
func (m *tree) clone() *tree {
	var t tree = *m
	return &t
}

// constants for FNV-1a hash algorithm
const (
	offset64 uint64 = 14695981039346656037
	prime64  uint64 = 1099511628211
)

// hashKey returns a hash code for a given string
func hashKey(key string) uint64 {
	hash := offset64
	for _, codepoint := range key {
		hash ^= uint64(codepoint)
		hash *= prime64
	}
	return hash
}

// Set returns a new map similar to this one but with key and value
// associated.  If the key didn't exist, it's created; otherwise, the
// associated value is changed.
func (m *tree) Set(key string, value Any) Map {
	hash := hashKey(key)
	return setLowLevel(m, hash, hash, key, value)
}

func setLowLevel(m *tree, partialHash, hash uint64, key string, value Any) *tree {
	if m.IsNil() { // an empty tree is easy
		m1 := m.clone()
		m1.count = 1
		m1.hash = hash
		m1.key = key
		m1.value = value
		return m1
	}

	if hash != m.hash {
		m1 := m.clone()
		i := partialHash % childCount
		m1.children[i] = setLowLevel(m.children[i], partialHash>>shiftSize, hash, key, value)
		recalculateCount(m1)
		return m1
	}

	// replacing a key's previous value
	m1 := m.clone()
	m1.value = value
	return m1
}

// modifies a map by recalculating its key count based on the counts
// of its subtrees
func recalculateCount(m *tree) {
	count := 0
	for _, t := range m.children {
		count += t.Size()
	}
	m.count = count + 1 // add one to count ourself
}

func (m *tree) Delete(key string) Map {
	hash := hashKey(key)
	newMap, _ := deleteLowLevel(m, hash, hash)
	return newMap
}

func deleteLowLevel(m *tree, partialHash, hash uint64) (*tree, bool) {
	// empty trees are easy
	if m.IsNil() {
		return m, false
	}

	if hash != m.hash {
		i := partialHash % childCount
		child, found := deleteLowLevel(m.children[i], partialHash>>shiftSize, hash)
		if !found {
			return m, false
		}
		newMap := m.clone()
		newMap.children[i] = child
		recalculateCount(newMap)
		return newMap, true // ? this wasn't in the original code
	}

	// we must delete our own node
	if m.isLeaf() { // we have no children
		return nilMap, true
	}

	// find a node to replace us
	i := -1
	size := -1
	for j, t := range m.children {
		if t.Size() > size {
			i = j
			size = t.Size()
		}
	}

	// make chosen leaf smaller
	replacement, child := m.children[i].deleteLeftmost()
	newMap := replacement.clone()
	for j := range m.children {
		if j == i {
			newMap.children[j] = child
		} else {
			newMap.children[j] = m.children[j]
		}
	}
	recalculateCount(newMap)
	return newMap, true
}

// delete the leftmost node in a tree returning the node that
// was deleted and the tree left over after its deletion
func (m *tree) deleteLeftmost() (*tree, *tree) {
	if m.isLeaf() {
		return m, nilMap
	}

	for i, t := range m.children {
		if t == nilMap {
			continue
		}

		deleted, child := t.deleteLeftmost()
		newMap := m.clone()
		newMap.children[i] = child
		recalculateCount(newMap)
		return deleted, newMap
	}
	panic("Tree isn't a leaf but also had no children. How does that happen?")
}

// isLeaf returns true if this is a leaf node
func (m *tree) isLeaf() bool {
	return m.Size() == 1
}

func (m *tree) Lookup(key string) (Any, bool) {
	hash := hashKey(key)
	return lookupLowLevel(m, hash, hash)
}

func lookupLowLevel(m *tree, partialHash, hash uint64) (Any, bool) {
	if m.IsNil() { // an empty tree is easy
		return nil, false
	}

	if hash != m.hash {
		i := partialHash % childCount
		return lookupLowLevel(m.children[i], partialHash>>shiftSize, hash)
	}

	// we found it
	return m.value, true
}

func (m *tree) Size() int {
	return m.count
}

func (m *tree) ForEach(f func(key string, val Any)) {
	if m.IsNil() {
		return
	}

	// ourself
	f(m.key, m.value)

	// children
	for _, t := range m.children {
		if t != nilMap {
			t.ForEach(f)
		}
	}
}

func (m *tree) Keys() []string {
	keys := make([]string, m.Size())
	i := 0
	m.ForEach(func(k string, v Any) {
		keys[i] = k
		i++
	})
	return keys
}

// make it easier to display maps for debugging
func (m *tree) String() string {
	keys := m.Keys()
	buf := bytes.NewBufferString("{")
	for _, key := range keys {
		val, _ := m.Lookup(key)
		_, _ = fmt.Fprintf(buf, "%s: %s, ", key, val)
	}
	_, _ = fmt.Fprintf(buf, "}\n")
	return buf.String()
}
