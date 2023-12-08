package gfly

import (
	"errors"
	"fmt"
	gstrings "github.com/savsgio/gotils/strings"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/valyala/bytebufferpool"
)

const (
	root nodeType = iota
	static
	param
	wildcard
)

const (
	errSetHandler         = "a handler is already registered for path '%s'"
	errSetWildcardHandler = "a wildcard handler is already registered for path '%s'"
	errWildPathConflict   = "'%s' in new path '%s' conflicts with existing wild path '%s' in existing prefix '%s'"
	errWildcardConflict   = "'%s' in new path '%s' conflicts with existing wildcard '%s' in existing prefix '%s'"
	errWildcardSlash      = "no / before wildcard in path '%s'"
	errWildcardNotAtEnd   = "wildcard routes are only allowed at the end of the path in path '%s'"
)

type radixError struct {
	msg    string
	params []interface{}
}

func (err radixError) Error() string {
	return fmt.Sprintf(err.msg, err.params...)
}

func newRadixError(msg string, params ...interface{}) radixError {
	return radixError{msg, params}
}

type nodeType uint8

type nodeWildcard struct {
	path     string
	paramKey string
	handler  IHandler
}

type node struct {
	nType nodeType

	path         string
	tsr          bool
	handler      IHandler
	hasWildChild bool
	children     []*node
	wildcard     *nodeWildcard

	paramKeys  []string
	paramRegex *regexp.Regexp
}

type wildPath struct {
	path  string
	keys  []string
	start int
	end   int
	pType nodeType

	pattern string
	regex   *regexp.Regexp
}

func panicf(s string, args ...interface{}) {
	panic(fmt.Sprintf(s, args...))
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func bufferRemoveString(buf *bytebufferpool.ByteBuffer, s string) {
	buf.B = buf.B[:len(buf.B)-len(s)]
}

// func isIndexEqual(a, b string) bool {
// 	ra, _ := utf8.DecodeRuneInString(a)
// 	rb, _ := utf8.DecodeRuneInString(b)

// 	return unicode.ToLower(ra) == unicode.ToLower(rb)
// }

// longestCommonPrefix finds the longest common prefix.
// This also implies that the common prefix contains no ':' or '*'
// since the existing key can't contain those chars.
func longestCommonPrefix(a, b string) int {
	i := 0
	max := min(utf8.RuneCountInString(a), utf8.RuneCountInString(b))

	for i < max {
		ra, sizeA := utf8.DecodeRuneInString(a)
		rb, sizeB := utf8.DecodeRuneInString(b)

		a = a[sizeA:]
		b = b[sizeB:]

		if ra != rb {
			return i
		}

		i += sizeA
	}

	return i
}

// segmentEndIndex returns the index where the segment ends from the given path
func segmentEndIndex(path string, includeTSR bool) int {
	end := 0
	for end < len(path) && path[end] != '/' {
		end++
	}

	if includeTSR && path[end:] == "/" {
		end++
	}

	return end
}

// findWildPath search for a wild path segment and check the name for invalid characters.
// Returns -1 as index, if no param/wildcard was found.
func findWildPath(path, fullPath string) *wildPath {
	// Find start
	for start, c := range []byte(path) {
		// A wildcard starts with ':' (param) or '*' (wildcard)
		if c != '{' {
			continue
		}

		withRegex := false
		keys := 0

		// Find end and check for invalid characters
		for end, c := range []byte(path[start+1:]) {
			switch c {
			case '}':
				if keys > 0 {
					keys--
					continue
				}

				end := start + end + 2
				wp := &wildPath{
					path:  path[start:end],
					keys:  []string{path[start+1 : end-1]},
					start: start,
					end:   end,
					pType: param,
				}

				if len(path) > end && path[end] == '{' {
					panic("the wildcards must be separated by at least 1 char")
				}

				sn := strings.SplitN(wp.keys[0], ":", 2)
				if len(sn) > 1 {
					wp.keys = []string{sn[0]}
					pattern := sn[1]

					if pattern == "*" {
						wp.pattern = pattern
						wp.pType = wildcard
					} else {
						wp.pattern = "(" + pattern + ")"
						wp.regex = regexp.MustCompile(wp.pattern)
					}
				} else if path[len(path)-1] != '/' {
					wp.pattern = "(.*)"
				}

				if wp.keys[0] == "" {
					panicf("wildcards must be named with a non-empty name in path '%s'", fullPath)
				}

				segEnd := end + segmentEndIndex(path[end:], true)
				path = path[end:segEnd]

				if path == "/" {
					// Last segment, so include the TSR
					path = ""
					wp.end++
				}

				if len(path) > 0 {
					// Rebuild the wildpath with the prefix
					wp2 := findWildPath(path, fullPath)
					if wp2 != nil {
						prefix := path[:wp2.start]

						wp.end += wp2.end
						wp.path += prefix + wp2.path
						wp.pattern += prefix + wp2.pattern
						wp.keys = append(wp.keys, wp2.keys...)
					} else {
						wp.path += path
						wp.pattern += path
						wp.end += len(path)
					}

					wp.regex = regexp.MustCompile(wp.pattern)
				}

				return wp

			case ':':
				withRegex = true

			case '{':
				if !withRegex && keys == 0 {
					panic("the char '{' is not allowed in the param name")
				}

				keys++
			}
		}
	}

	return nil
}

// Tree is a routes storage
type Tree struct {
	root *node

	// If enabled, the node handler could be updated
	Mutable bool
}

// NewTree returns an empty routes storage
func NewTree() *Tree {
	return &Tree{
		root: &node{
			nType: root,
		},
	}
}

// Add adds a node with the given handle to the path.
//
// WARNING: Not concurrency-safe!
func (t *Tree) Add(path string, handler IHandler) {
	if !strings.HasPrefix(path, "/") {
		panicf("path must begin with '/' in path '%s'", path)
	} else if handler == nil {
		panic("nil handler")
	}

	fullPath := path

	i := longestCommonPrefix(path, t.root.path)
	if i > 0 {
		if len(t.root.path) > i {
			t.root.split(i)
		}

		path = path[i:]
	}

	n, err := t.root.add(path, fullPath, handler)
	if err != nil {
		var radixErr radixError

		if errors.As(err, &radixErr) && t.Mutable && !n.tsr {
			switch radixErr.msg {
			case errSetHandler:
				n.handler = handler
				return
			case errSetWildcardHandler:
				n.wildcard.handler = handler
				return
			}
		}

		panic(err)
	}

	if t.root.path == "" {
		t.root = t.root.children[0]
		t.root.nType = root
	}

	// Reorder the nodes
	t.root.sort()
}

// Get returns the handle registered with the given path (key). The values of
// param/wildcard are saved as ctx.UserValue.
// If no handle can be found, a TSR (trailing slash redirect) recommendation is
// made if a handle exists with an extra (without the) trailing slash for the
// given path.
func (t *Tree) Get(path string, ctx *Ctx) (IHandler, bool) {
	if len(path) > len(t.root.path) {
		if path[:len(t.root.path)] != t.root.path {
			return nil, false
		}

		path = path[len(t.root.path):]

		return t.root.getFromChild(path, ctx)

	} else if path == t.root.path {
		switch {
		case t.root.tsr:
			return nil, true
		case t.root.handler != nil:
			return t.root.handler, false
		case t.root.wildcard != nil:
			if ctx != nil {
				ctx.Root().SetUserValue(t.root.wildcard.paramKey, "")
			}

			return t.root.wildcard.handler, false
		}
	}

	return nil, false
}

// FindCaseInsensitivePath makes a case-insensitive lookup of the given path
// and tries to find a handler.
// It can optionally also fix trailing slashes.
// It returns the case-corrected path and a bool indicating whether the lookup
// was successful.
func (t *Tree) FindCaseInsensitivePath(path string, fixTrailingSlash bool, buf *bytebufferpool.ByteBuffer) bool {
	found, tsr := t.root.find(path, buf)

	if !found || (tsr && !fixTrailingSlash) {
		buf.Reset()

		return false
	}

	return true
}

func newNode(path string) *node {
	return &node{
		nType: static,
		path:  path,
	}
}

// conflict raises a panic with some details
func (n *nodeWildcard) conflict(path, fullPath string) error {
	prefix := fullPath[:strings.LastIndex(fullPath, path)] + n.path

	return newRadixError(errWildcardConflict, path, fullPath, n.path, prefix)
}

// wildPathConflict raises a panic with some details
func (n *node) wildPathConflict(path, fullPath string) error {
	pathSeg := strings.SplitN(path, "/", 2)[0]
	prefix := fullPath[:strings.LastIndex(fullPath, path)] + n.path

	return newRadixError(errWildPathConflict, pathSeg, fullPath, n.path, prefix)
}

// clone clones the current node in a new pointer
func (n *node) clone() *node {
	cloneNode := new(node)
	cloneNode.nType = n.nType
	cloneNode.path = n.path
	cloneNode.tsr = n.tsr
	cloneNode.handler = n.handler

	if len(n.children) > 0 {
		cloneNode.children = make([]*node, len(n.children))

		for i, child := range n.children {
			cloneNode.children[i] = child.clone()
		}
	}

	if n.wildcard != nil {
		cloneNode.wildcard = &nodeWildcard{
			path:     n.wildcard.path,
			paramKey: n.wildcard.paramKey,
			handler:  n.wildcard.handler,
		}
	}

	if len(n.paramKeys) > 0 {
		cloneNode.paramKeys = make([]string, len(n.paramKeys))
		copy(cloneNode.paramKeys, n.paramKeys)
	}

	cloneNode.paramRegex = n.paramRegex

	return cloneNode
}

func (n *node) split(i int) {
	cloneChild := n.clone()
	cloneChild.nType = static
	cloneChild.path = cloneChild.path[i:]
	cloneChild.paramKeys = nil
	cloneChild.paramRegex = nil

	n.path = n.path[:i]
	n.handler = nil
	n.tsr = false
	n.wildcard = nil
	n.children = append(n.children[:0], cloneChild)
}

func (n *node) findEndIndexAndValues(path string) (int, []string) {
	index := n.paramRegex.FindStringSubmatchIndex(path)
	if len(index) == 0 || index[0] != 0 {
		return -1, nil
	}

	end := index[1]

	index = index[2:]
	values := make([]string, len(index)/2)

	i := 0
	for j := range index {
		if (j+1)%2 != 0 {
			continue
		}

		values[i] = gstrings.Copy(path[index[j-1]:index[j]])

		i++
	}

	return end, values
}

func (n *node) setHandle(handler IHandler, fullPath string) (*node, error) {
	if n.handler != nil || n.tsr {
		return n, newRadixError(errSetHandler, fullPath)
	}

	n.handler = handler
	foundTSR := false

	// Set TSR in method
	for i := range n.children {
		child := n.children[i]

		if child.path != "/" {
			continue
		}

		child.tsr = true
		foundTSR = true

		break
	}

	if n.path != "/" && !foundTSR {
		if strings.HasSuffix(n.path, "/") {
			n.split(len(n.path) - 1)
			n.tsr = true
		} else {
			childTSR := newNode("/")
			childTSR.tsr = true
			n.children = append(n.children, childTSR)
		}
	}

	return n, nil
}

func (n *node) insert(path, fullPath string, handler IHandler) (*node, error) {
	end := segmentEndIndex(path, true)
	child := newNode(path)

	wp := findWildPath(path, fullPath)
	if wp != nil {
		j := end
		if wp.start > 0 {
			j = wp.start
		}

		child.path = path[:j]

		if wp.start > 0 {
			n.children = append(n.children, child)

			return child.insert(path[j:], fullPath, handler)
		}

		switch wp.pType {
		case param:
			n.hasWildChild = true

			child.nType = wp.pType
			child.paramKeys = wp.keys
			child.paramRegex = wp.regex
		case wildcard:
			if len(path) == end && n.path[len(n.path)-1] != '/' {
				return nil, newRadixError(errWildcardSlash, fullPath)
			} else if len(path) != end {
				return nil, newRadixError(errWildcardNotAtEnd, fullPath)
			}

			if n.path != "/" && n.path[len(n.path)-1] == '/' {
				n.split(len(n.path) - 1)
				n.tsr = true

				n = n.children[0]
			}

			if n.wildcard != nil {
				if n.wildcard.path == path {
					return n, newRadixError(errSetWildcardHandler, fullPath)
				}

				return nil, n.wildcard.conflict(path, fullPath)
			}

			n.wildcard = &nodeWildcard{
				path:     wp.path,
				paramKey: wp.keys[0],
				handler:  handler,
			}

			return n, nil
		}

		path = path[wp.end:]

		if len(path) > 0 {
			n.children = append(n.children, child)

			return child.insert(path, fullPath, handler)
		}
	}

	child.handler = handler
	n.children = append(n.children, child)

	switch {
	case child.path == "/":
		// Add TSR when split a edge and the remain path to insert is "/"
		n.tsr = true
	case strings.HasSuffix(child.path, "/"):
		child.split(len(child.path) - 1)
		child.tsr = true
	default:
		childTSR := newNode("/")
		childTSR.tsr = true
		child.children = append(child.children, childTSR)
	}

	return child, nil
}

// add adds the handler to node for the given path
func (n *node) add(path, fullPath string, handler IHandler) (*node, error) {
	if path == "" {
		return n.setHandle(handler, fullPath)
	}

	for _, child := range n.children {
		i := longestCommonPrefix(path, child.path)
		if i == 0 {
			continue
		}

		switch child.nType {
		case static:
			if len(child.path) > i {
				child.split(i)
			}

			if len(path) > i {
				return child.add(path[i:], fullPath, handler)
			}
		case param:
			wp := findWildPath(path, fullPath)

			isParam := wp.start == 0 && wp.pType == param
			hasHandler := child.handler != nil || handler == nil

			if len(path) == wp.end && isParam && hasHandler {
				// The current segment is a param and it's duplicated
				if child.path == path {
					return child, newRadixError(errSetHandler, fullPath)
				}

				return nil, child.wildPathConflict(path, fullPath)
			}

			if len(path) > i {
				if child.path == wp.path {
					return child.add(path[i:], fullPath, handler)
				}

				return n.insert(path, fullPath, handler)
			}
		}

		if path == "/" {
			n.tsr = true
		}

		return child.setHandle(handler, fullPath)
	}

	return n.insert(path, fullPath, handler)
}

func (n *node) getFromChild(path string, ctx *Ctx) (IHandler, bool) {
	for _, child := range n.children {
		switch child.nType {
		case static:

			// Checks if the first byte is equal
			// It's faster than compare strings
			if path[0] != child.path[0] {
				continue
			}

			if len(path) > len(child.path) {
				if path[:len(child.path)] != child.path {
					continue
				}

				h, tsr := child.getFromChild(path[len(child.path):], ctx)
				if h != nil || tsr {
					return h, tsr
				}
			} else if path == child.path {
				switch {
				case child.tsr:
					return nil, true
				case child.handler != nil:
					return child.handler, false
				case child.wildcard != nil:
					if ctx != nil {
						ctx.Root().SetUserValue(child.wildcard.paramKey, "")
					}

					return child.wildcard.handler, false
				}

				return nil, false
			}

		case param:
			end := segmentEndIndex(path, false)
			values := []string{gstrings.Copy(path[:end])}

			if child.paramRegex != nil {
				end, values = child.findEndIndexAndValues(path[:end])
				if end == -1 {
					continue
				}
			}

			if len(path) > end {
				h, tsr := child.getFromChild(path[end:], ctx)
				if tsr {
					return nil, tsr
				} else if h != nil {
					if ctx != nil {
						for i, key := range child.paramKeys {
							ctx.Root().SetUserValue(key, values[i])
						}
					}

					return h, false
				}

			} else if len(path) == end {
				switch {
				case child.tsr:
					return nil, true
				case child.handler == nil:
					// try another child
					continue
				case ctx != nil:
					for i, key := range child.paramKeys {
						ctx.Root().SetUserValue(key, values[i])
					}
				}

				return child.handler, false
			}

		default:
			panic("invalid node type")
		}
	}

	if n.wildcard != nil {
		if ctx != nil {
			ctx.Root().SetUserValue(n.wildcard.paramKey, gstrings.Copy(path))
		}

		return n.wildcard.handler, false
	}

	return nil, false
}

func (n *node) find(path string, buf *bytebufferpool.ByteBuffer) (bool, bool) {
	if len(path) > len(n.path) {
		if !strings.EqualFold(path[:len(n.path)], n.path) {
			return false, false
		}

		path = path[len(n.path):]
		_, err := buf.WriteString(n.path)
		if err != nil {
			return false, false
		}

		found, tsr := n.findFromChild(path, buf)
		if found {
			return found, tsr
		}

		bufferRemoveString(buf, n.path)

	} else if strings.EqualFold(path, n.path) {
		_, err := buf.WriteString(n.path)
		if err != nil {
			return false, false
		}

		if n.tsr {
			if n.path == "/" {
				bufferRemoveString(buf, n.path)
			} else {
				err := buf.WriteByte('/')
				if err != nil {
					return false, false
				}
			}

			return true, true
		}

		if n.handler != nil {
			return true, false
		} else {
			bufferRemoveString(buf, n.path)
		}
	}

	return false, false
}

func (n *node) findFromChild(path string, buf *bytebufferpool.ByteBuffer) (bool, bool) {
	for _, child := range n.children {
		switch child.nType {
		case static:
			found, tsr := child.find(path, buf)
			if found {
				return found, tsr
			}

		case param:
			end := segmentEndIndex(path, false)

			if child.paramRegex != nil {
				end, _ = child.findEndIndexAndValues(path[:end])
				if end == -1 {
					continue
				}
			}

			_, err := buf.WriteString(path[:end])
			if err != nil {
				return false, false
			}

			if len(path) > end {
				found, tsr := child.findFromChild(path[end:], buf)
				if found {
					return found, tsr
				}

			} else if len(path) == end {
				if child.tsr {
					err := buf.WriteByte('/')
					if err != nil {
						return false, false
					}

					return true, true
				}

				if child.handler != nil {
					return true, false
				}
			}

			bufferRemoveString(buf, path[:end])

		default:
			panic("invalid node type")
		}
	}

	if n.wildcard != nil {
		_, err := buf.WriteString(path)
		if err != nil {
			return false, false
		}

		return true, false
	}

	return false, false
}

// sort sorts the current node and their children
func (n *node) sort() {
	for _, child := range n.children {
		child.sort()
	}

	sort.Sort(n)
}

// Len returns the total number of children the node has
func (n *node) Len() int {
	return len(n.children)
}

// Swap swaps the order of children nodes
func (n *node) Swap(i, j int) {
	n.children[i], n.children[j] = n.children[j], n.children[i]
}

// Less checks if the node 'i' has less priority than the node 'j'
func (n *node) Less(i, j int) bool {
	if n.children[i].nType < n.children[j].nType {
		return true
	} else if n.children[i].nType > n.children[j].nType {
		return false
	}

	return len(n.children[i].children) > len(n.children[j].children)
}
