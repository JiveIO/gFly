package gfly

import (
	"app/core/log"
	"app/core/utils"
	"fmt"
	"strings"

	"github.com/savsgio/gotils/bytes"
	"github.com/savsgio/gotils/strconv"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

// MethodWild wild HTTP method
const MethodWild = "*"

var (
	questionMark = byte('?')

	// matchedRoutePathParam is the param name under which the path of the matched
	// route is stored, if Router.SaveMatchedRoutePath is set.
	matchedRoutePathParam = fmt.Sprintf("__matchedRoutePath::%s__", bytes.Rand(make([]byte, 15)))
)

// Router is a RequestHandler which can be used to dispatch requests to different
// handler functions via configurable routes
type Router struct {
	trees              []*Tree
	treeMutable        bool
	customMethodsIndex map[string]int
	registeredPaths    map[string][]string

	// If enabled, adds the matched route path onto the ctx.UserValue context
	// before invoking the handler.
	// The matched route path is only added to handlers of routes that were
	// registered when this option was enabled.
	SaveMatchedRoutePath bool

	// Enables automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	// For example if /foo/ is requested but a route only exists for /foo, the
	// client is redirected to /foo with http status code 301 for GET requests
	// and 308 for all other request methods.
	RedirectTrailingSlash bool

	// If enabled, the router tries to fix the current request path, if no
	// handle is registered for it.
	// First superfluous path elements like ../ or // are removed.
	// Afterwards the router does a case-insensitive lookup of the cleaned path.
	// If a handle can be found for this route, the router makes a redirection
	// to the corrected path with status code 301 for GET requests and 308 for
	// all other request methods.
	// For example /FOO and /..//Foo could be redirected to /foo.
	// RedirectTrailingSlash is independent of this option.
	RedirectFixedPath bool

	// If enabled, the router checks if another method is allowed for the
	// current route, if the current request can not be routed.
	// If this is the case, the request is answered with 'Method Not Allowed'
	// and HTTP status code 405.
	// If no other Method is allowed, the request is delegated to the NotFound
	// handler.
	HandleMethodNotAllowed bool

	// If enabled, the router automatically replies to OPTIONS requests.
	// Custom OPTIONS handlers take priority over automatic replies.
	HandleOPTIONS bool

	// An optional RequestHandler that is called on automatic OPTIONS requests.
	// The handler is only called if HandleOPTIONS is true and no OPTIONS
	// handler for the specific path was set.
	// The "Allowed" header is set before calling the handler.
	GlobalOPTIONS RequestHandler

	// Configurable RequestHandler which is called when no matching route is
	// found. If it is not set, default NotFound is used.
	NotFound RequestHandler

	// Configurable RequestHandler which is called when a request
	// cannot be routed and HandleMethodNotAllowed is true.
	// If it is not set, ctx.Error with fasthttp.StatusMethodNotAllowed is used.
	// The "Allow" header with allowed request methods is set before the handler
	// is called.
	MethodNotAllowed RequestHandler

	// Function to handle panics recovered from http handlers.
	// It should be used to generate a error page and return the http error code
	// 500 (Internal Server Error).
	// The handler can be used to keep your server from crashing because of
	// unrecovered panics.
	PanicHandler func(*Ctx, interface{})

	// Cached value of global (*) allowed methods
	globalAllowed string
}

// NewRouter returns a new router.
// Path auto-correction, including trailing slashes, is enabled by default.
func NewRouter() *Router {
	return &Router{
		trees:                  make([]*Tree, 10),
		customMethodsIndex:     make(map[string]int),
		registeredPaths:        make(map[string][]string),
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		PanicHandler: func(ctx *Ctx, data interface{}) {
			log.Errorf("Router Panic Handling %v", data)
		},
		GlobalOPTIONS: func(ctx *Ctx) error {
			// Set CORs headers
			ctx.root.Response.Header.Set(HeaderAccessControlAllowOrigin, "*")
			ctx.root.Response.Header.Set(HeaderAccessControlAllowMethods, "PUT, POST, GET, DELETE, OPTIONS, PATCH")
			ctx.root.Response.Header.Set(HeaderAccessControlAllowHeaders, "Authorization, Content-Type, x-requested-with, origin, true-client-ip, X-Correlation-ID")

			return nil
		},
	}
}

// Group returns a new group.
// Path auto-correction, including trailing slashes, is enabled by default.
func (r *Router) Group(path string) *Group {
	validatePath(path)

	if path != "/" && strings.HasSuffix(path, "/") {
		panic("group path must not end with a trailing slash")
	}

	return &Group{
		router: r,
		prefix: path,
	}
}

func (r *Router) saveMatchedRoutePath(path string, handler IHandler) IHandler {
	return &saveMatchedRoutePathHandler{
		path:    path,
		handler: handler,
	}
}

// saveMatchedRoutePathHandler Default handler
type saveMatchedRoutePathHandler struct {
	Endpoint
	path    string
	handler IHandler
}

func (e *saveMatchedRoutePathHandler) Handle(c *Ctx) error {
	c.root.SetUserValue(matchedRoutePathParam, e.path)
	return e.handler.Handle(c)
}

func (r *Router) methodIndexOf(method string) int {
	switch method {
	case fasthttp.MethodGet:
		return 0
	case fasthttp.MethodHead:
		return 1
	case fasthttp.MethodPost:
		return 2
	case fasthttp.MethodPut:
		return 3
	case fasthttp.MethodPatch:
		return 4
	case fasthttp.MethodDelete:
		return 5
	case fasthttp.MethodConnect:
		return 6
	case fasthttp.MethodOptions:
		return 7
	case fasthttp.MethodTrace:
		return 8
	case MethodWild:
		return 9
	}

	if i, ok := r.customMethodsIndex[method]; ok {
		return i
	}

	return -1
}

// Mutable allows updating the route handler
//
// # It's disabled by default
//
// WARNING: Use with care. It could generate unexpected behaviours
func (r *Router) Mutable(v bool) {
	r.treeMutable = v

	for i := range r.trees {
		tree := r.trees[i]

		if tree != nil {
			tree.Mutable = v
		}
	}
}

// List returns all registered routes grouped by method
func (r *Router) List() map[string][]string {
	return r.registeredPaths
}

// GET is a shortcut for router.Handle(fasthttp.MethodGet, path, handler)
func (r *Router) GET(path string, handler IHandler) {
	r.Handle(fasthttp.MethodGet, path, handler)
}

// HEAD is a shortcut for router.Handle(fasthttp.MethodHead, path, handler)
func (r *Router) HEAD(path string, handler IHandler) {
	r.Handle(fasthttp.MethodHead, path, handler)
}

// POST is a shortcut for router.Handle(fasthttp.MethodPost, path, handler)
func (r *Router) POST(path string, handler IHandler) {
	r.Handle(fasthttp.MethodPost, path, handler)
}

// PUT is a shortcut for router.Handle(fasthttp.MethodPut, path, handler)
func (r *Router) PUT(path string, handler IHandler) {
	r.Handle(fasthttp.MethodPut, path, handler)
}

// PATCH is a shortcut for router.Handle(fasthttp.MethodPatch, path, handler)
func (r *Router) PATCH(path string, handler IHandler) {
	r.Handle(fasthttp.MethodPatch, path, handler)
}

// DELETE is a shortcut for router.Handle(fasthttp.MethodDelete, path, handler)
func (r *Router) DELETE(path string, handler IHandler) {
	r.Handle(fasthttp.MethodDelete, path, handler)
}

// CONNECT is a shortcut for router.Handle(fasthttp.MethodConnect, path, handler)
func (r *Router) CONNECT(path string, handler IHandler) {
	r.Handle(fasthttp.MethodConnect, path, handler)
}

// OPTIONS is a shortcut for router.Handle(fasthttp.MethodOptions, path, handler)
func (r *Router) OPTIONS(path string, handler IHandler) {
	r.Handle(fasthttp.MethodOptions, path, handler)
}

// TRACE is a shortcut for router.Handle(fasthttp.MethodTrace, path, handler)
func (r *Router) TRACE(path string, handler IHandler) {
	r.Handle(fasthttp.MethodTrace, path, handler)
}

// ServeFiles serves files from the given file system root.
// The path must end with "/{filepath:*}", files are then served from the local
// path /defined/root/dir/{filepath:*}.
// For example if root is "/etc" and {filepath:*} is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a fasthttp.FSHandler is used, therefore fasthttp.NotFound is used instead
// Use:
//
//	router.ServeFiles("/src/{filepath:*}", "./")
func (r *Router) ServeFiles(path, rootPath string) {
	r.ServeFilesCustom(path, &fasthttp.FS{
		Root:               rootPath,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: true,
		AcceptByteRange:    true,
	})
}

// ServeFilesCustom serves files from the given file system settings.
// The path must end with "/{filepath:*}", files are then served from the local
// path /defined/root/dir/{filepath:*}.
// For example if root is "/etc" and {filepath:*} is "passwd", the local file
// "/etc/passwd" would be served.
// Internally a fasthttp.FSHandler is used, therefore http.NotFound is used instead
// of the Router's NotFound handler.
// Use:
//
//	router.ServeFilesCustom("/src/{filepath:*}", *customFS)
func (r *Router) ServeFilesCustom(path string, fs *fasthttp.FS) {
	suffix := "/{filepath:*}"

	if !strings.HasSuffix(path, suffix) {
		panic("path must end with " + suffix + " in path '" + path + "'")
	}

	prefix := path[:len(path)-len(suffix)]
	stripSlashes := strings.Count(prefix, "/")

	if fs.PathRewrite == nil && stripSlashes > 0 {
		fs.PathRewrite = fasthttp.NewPathSlashesStripper(stripSlashes)
	}
	handler := &serveFilesCustomHandler{
		fileHandler: fs.NewRequestHandler(),
	}

	r.GET(path, handler)
}

// ServeFilesCustomEndpoint Default handler
type serveFilesCustomHandler struct {
	Endpoint
	fileHandler fasthttp.RequestHandler
}

func (e *serveFilesCustomHandler) Handle(c *Ctx) error {
	e.fileHandler(c.root)

	return nil
}

// Handle registers a new request handler with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *Router) Handle(method, path string, handler IHandler) {
	switch {
	case method == "":
		panic("method must not be empty")
	case handler == nil:
		panic("handler must not be nil")
	default:
		validatePath(path)
	}

	r.registeredPaths[method] = append(r.registeredPaths[method], path)

	methodIndex := r.methodIndexOf(method)
	if methodIndex == -1 {
		tree := NewTree()
		tree.Mutable = r.treeMutable

		r.trees = append(r.trees, tree)
		methodIndex = len(r.trees) - 1
		r.customMethodsIndex[method] = methodIndex
	}

	tree := r.trees[methodIndex]
	if tree == nil {
		tree = NewTree()
		tree.Mutable = r.treeMutable

		r.trees[methodIndex] = tree
		r.globalAllowed = r.allowed("*", "")
	}

	if r.SaveMatchedRoutePath {
		handler = r.saveMatchedRoutePath(path, handler)
	}

	optionalPaths := getOptionalPaths(path)

	// if not has optional paths, adds the original
	if len(optionalPaths) == 0 {
		tree.Add(path, handler)
	} else {
		for _, p := range optionalPaths {
			tree.Add(p, handler)
		}
	}
}

// Lookup allows the manual lookup of a method + path combo.
// This is e.g. useful to build a framework around this router.
// If the path was found, it returns the handler function.
// Otherwise the second return value indicates whether a redirection to
// the same path with an extra / without the trailing slash should be performed.
func (r *Router) Lookup(method, path string, ctx *Ctx) (IHandler, bool) {
	methodIndex := r.methodIndexOf(method)
	if methodIndex == -1 {
		return nil, false
	}

	if tree := r.trees[methodIndex]; tree != nil {
		handler, tsr := tree.Get(path, ctx)
		if handler != nil || tsr {
			return handler, tsr
		}
	}

	if tree := r.trees[r.methodIndexOf(MethodWild)]; tree != nil {
		return tree.Get(path, ctx)
	}

	return nil, false
}

func (r *Router) recv(ctx *Ctx) {
	if rcv := recover(); rcv != nil {
		r.PanicHandler(ctx, rcv)
	}
}

func (r *Router) allowed(path, reqMethod string) (allow string) {
	allowed := make([]string, 0, 9)

	if path == "*" || path == "/*" { // server-wide{ // server-wide
		// empty method is used for internal calls to refresh the cache
		if reqMethod == "" {
			for method := range r.registeredPaths {
				if method == fasthttp.MethodOptions {
					continue
				}
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		} else {
			return r.globalAllowed
		}
	} else { // specific path
		for method := range r.registeredPaths {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == fasthttp.MethodOptions {
				continue
			}

			handle, _ := r.trees[r.methodIndexOf(method)].Get(path, nil)
			if handle != nil {
				// Add request method to list of allowed methods
				allowed = append(allowed, method)
			}
		}
	}

	if len(allowed) > 0 {
		// Add request method to list of allowed methods
		allowed = append(allowed, fasthttp.MethodOptions)

		// Sort allowed methods.
		// sort.Strings(allowed) unfortunately causes unnecessary allocations
		// due to allowed being moved to the heap and interface conversion
		for i, l := 1, len(allowed); i < l; i++ {
			for j := i; j > 0 && allowed[j] < allowed[j-1]; j-- {
				allowed[j], allowed[j-1] = allowed[j-1], allowed[j]
			}
		}

		// return as comma separated list
		return strings.Join(allowed, ", ")
	}
	return
}

func (r *Router) tryRedirect(ctx *Ctx, tree *Tree, tsr bool, method, path string) bool {
	// Moved Permanently, request with GET method
	code := fasthttp.StatusMovedPermanently
	if method != fasthttp.MethodGet {
		// Permanent Redirect, request with same method
		code = fasthttp.StatusPermanentRedirect
	}

	if tsr && r.RedirectTrailingSlash {
		uri := bytebufferpool.Get()

		if len(path) > 1 && path[len(path)-1] == '/' {
			uri.SetString(path[:len(path)-1])
		} else {
			uri.SetString(path)
			err := uri.WriteByte('/')
			if err != nil {
				return false
			}
		}

		if queryBuf := ctx.root.URI().QueryString(); len(queryBuf) > 0 {
			err := uri.WriteByte(questionMark)
			if err != nil {
				return false
			}
			_, err = uri.Write(queryBuf)
			if err != nil {
				return false
			}
		}

		ctx.root.Redirect(uri.String(), code)
		bytebufferpool.Put(uri)

		return true
	}

	// Try to fix the request path
	if r.RedirectFixedPath {
		path2 := strconv.B2S(ctx.root.Request.URI().Path())

		uri := bytebufferpool.Get()
		found := tree.FindCaseInsensitivePath(
			cleanPath(path2),
			r.RedirectTrailingSlash,
			uri,
		)

		if found {
			if queryBuf := ctx.root.URI().QueryString(); len(queryBuf) > 0 {
				err := uri.WriteByte(questionMark)
				if err != nil {
					return false
				}
				_, err = uri.Write(queryBuf)
				if err != nil {
					return false
				}
			}

			ctx.root.Redirect(uri.String(), code)
			bytebufferpool.Put(uri)

			return true
		}

		bytebufferpool.Put(uri)
	}

	return false
}

// Handler makes the router implement the http.Handler interface.
func (r *Router) Handler(ctx *Ctx) error {
	if r.PanicHandler != nil {
		defer r.recv(ctx)
	}

	path := strconv.B2S(ctx.root.Request.URI().PathOriginal())
	method := strconv.B2S(ctx.root.Request.Header.Method())
	methodIndex := r.methodIndexOf(method)

	if methodIndex > -1 {
		if tree := r.trees[methodIndex]; tree != nil {
			if handler, tsr := tree.Get(path, ctx); handler != nil {
				if err := handler.Validate(ctx); err != nil {
					return errorHandler(ctx, err, StatusBadRequest)
				}

				if err := handler.Handle(ctx); err != nil {
					return errorHandler(ctx, err, StatusInternalServerError)
				}

				return nil
			} else if method != fasthttp.MethodConnect && path != "/" {
				if ok := r.tryRedirect(ctx, tree, tsr, method, path); ok {
					return nil
				}
			}
		}
	}

	// Try to search in the wild method tree
	if tree := r.trees[r.methodIndexOf(MethodWild)]; tree != nil {
		if handler, tsr := tree.Get(path, ctx); handler != nil {
			if err := handler.Validate(ctx); err != nil {
				return errorHandler(ctx, err, StatusBadRequest)
			}

			if err := handler.Handle(ctx); err != nil {
				return errorHandler(ctx, err, StatusInternalServerError)
			}

			return nil
		} else if method != fasthttp.MethodConnect && path != "/" {
			if ok := r.tryRedirect(ctx, tree, tsr, method, path); ok {
				return nil
			}
		}
	}

	if r.HandleOPTIONS && method == fasthttp.MethodOptions {
		// Handle OPTIONS requests

		if allow := r.allowed(path, fasthttp.MethodOptions); allow != "" {
			ctx.root.Response.Header.Set("Allow", allow)
			if r.GlobalOPTIONS != nil {
				return r.GlobalOPTIONS(ctx)
			}
			return nil
		}
	} else if r.HandleMethodNotAllowed {
		// Handle 405

		if allow := r.allowed(path, method); allow != "" {
			ctx.root.Response.Header.Set("Allow", allow)
			if r.MethodNotAllowed != nil {
				return r.MethodNotAllowed(ctx)
			} else {
				ctx.root.SetStatusCode(fasthttp.StatusMethodNotAllowed)
				ctx.root.SetBodyString(fasthttp.StatusMessage(fasthttp.StatusMethodNotAllowed))
			}
			return nil
		}
	}

	// Handle 404
	if r.NotFound != nil {
		return r.NotFound(ctx)
	} else {
		ctx.root.Error(fasthttp.StatusMessage(fasthttp.StatusNotFound), fasthttp.StatusNotFound)
	}

	return nil
}

func errorHandler(ctx *Ctx, err error, code int) error {
	// Force status error
	if ctx.root.Response.StatusCode() == StatusOK {
		ctx.Status(code)
	}

	// Write log
	log.Errorf("Bad request: %v", err)

	// Set default body
	body := ctx.root.Response.Body()
	if len(body) == 0 {
		contentType := strings.ToLower(utils.UnsafeString(ctx.root.Request.Header.ContentType()))
		if strings.HasPrefix(contentType, MIMEApplicationJSON) {
			_ = ctx.JSON(map[string]string{
				"error": err.Error(),
			})
		} else {
			_ = ctx.HTML(err.Error())
		}
	}

	return err
}
