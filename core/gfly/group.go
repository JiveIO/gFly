package gfly

// Group is a sub-router to group paths
type Group struct {
	router *Router
	prefix string
	// Global middleware handler for router group.
	middlewares []MiddlewareHandler
}

// ===========================================================================================================
// 										Group Middleware
// ===========================================================================================================

// IGroupMiddleware Interface to declare all Middleware methods for gFly struct.
type IGroupMiddleware interface {
	// Use apply middleware for all router group.
	// Important: Should put the code at the top of group router.
	Use(middlewares ...MiddlewareHandler)
}

// Use apply middleware for all router group.
func (g *Group) Use(middlewares ...MiddlewareHandler) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// ===========================================================================================================
// 										Group Router
// ===========================================================================================================

// IGroupRouter Interface to declare all HTTP methods.
type IGroupRouter interface {
	// GET Http GET method
	GET(path string, handler IHandler)
	// HEAD Http HEAD method
	HEAD(path string, handler IHandler)
	// POST Http POST method
	POST(path string, handler IHandler)
	// PUT Http PUT method
	PUT(path string, handler IHandler)
	// PATCH Http PATCH method
	PATCH(path string, handler IHandler)
	// DELETE Http DELETE method
	DELETE(path string, handler IHandler)
	// CONNECT Http CONNECT method
	CONNECT(path string, handler IHandler)
	// OPTIONS Http OPTIONS method
	OPTIONS(path string, handler IHandler)
	// TRACE Http TRACE method
	TRACE(path string, handler IHandler)
	// Group multi routers
	Group(path string, groupFunc func(*Group))
}

// GET is a shortcut for Router.GET(path, handler)
func (g *Group) GET(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.GET(g.prefix+path, g.wrapMiddlewares(handler))
}

// HEAD is a shortcut for Router.HEAD(path, handler)
func (g *Group) HEAD(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.HEAD(g.prefix+path, g.wrapMiddlewares(handler))
}

// POST is a shortcut for Router.POST(path, handler)
func (g *Group) POST(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.POST(g.prefix+path, g.wrapMiddlewares(handler))
}

// PUT is a shortcut for Router.PUT(path, handler)
func (g *Group) PUT(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.PUT(g.prefix+path, g.wrapMiddlewares(handler))
}

// PATCH is a shortcut for Router.PATCH(path, handler)
func (g *Group) PATCH(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.PATCH(g.prefix+path, g.wrapMiddlewares(handler))
}

// DELETE is a shortcut for Router.DELETE(path, handler)
func (g *Group) DELETE(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.DELETE(g.prefix+path, g.wrapMiddlewares(handler))
}

// CONNECT is a shortcut for Router.CONNECT(path, handler)
func (g *Group) CONNECT(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.CONNECT(g.prefix+path, g.wrapMiddlewares(handler))
}

// OPTIONS is a shortcut for Router.OPTIONS(path, handler)
func (g *Group) OPTIONS(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.OPTIONS(g.prefix+path, g.wrapMiddlewares(handler))
}

// TRACE is a shortcut for Router.TRACE(path, handler)
func (g *Group) TRACE(path string, handler IHandler) {
	if g.prefix == "" {
		validatePath(path)
	}

	g.router.TRACE(g.prefix+path, g.wrapMiddlewares(handler))
}

// Group Create a group Handler functions.
func (g *Group) Group(path string, groupFunc func(*Group)) {
	group := g.router.Group(g.prefix + path)
	// Auto append middleware from parent to children
	// Let example have a parent group have prefix url "/user" have 2 middleware functions: A, B.
	// Now So, now create a new subgroup have prefix url "/info". So, all handlers of subgroup "/info"
	// must be affected by 2 middlewares A, B from a parent group "/user"
	group.middlewares = g.middlewares

	groupFunc(group)
}

func (g *Group) wrapMiddlewares(handler IHandler) IHandler {
	if len(g.middlewares) > 0 {
		middlewareGroup := NewMiddleware()

		return middlewareGroup.Group(g.middlewares...)(handler)
	}

	return handler
}
