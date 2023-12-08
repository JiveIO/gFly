package gfly

import (
	"app/core/log"
	"app/core/utils"
	"fmt"
	_ "github.com/joho/godotenv/autoload" // load .env file automatically
	"github.com/valyala/fasthttp"
	"os"
)

// ===========================================================================================================
// 												gFly
// ===========================================================================================================

var (
	ViewPath     = os.Getenv("VIEW_PATH")
	ViewExt      = os.Getenv("VIEW_EXT")
	TemporaryDir = os.Getenv("TEMP_DIR")
)

// IFly Interface to declare all methods for gFly struct.
type IFly interface {
	Start()
	Router() *Router
	IFlyRouter
	IFlyMiddleware
}

// gFly Struct define main elements in app.
type gFly struct {
	router      *Router             // Keep reference router
	server      *fasthttp.Server    // Keep reference web server
	config      Config              // App configuration
	middleware  IMiddleware         // Keep referenceMiddleware
	middlewares []MiddlewareHandler // Global middleware handler
}

// Router Get root router in gFly app.
func (fly *gFly) Router() *Router { // Get root Router
	return fly.router
}

// Start Run gFly app.
func (fly *gFly) Start() {
	// --------------- Setup Logs ---------------
	setupLog()

	// --------------- Setup Session ---------------
	setupSession()

	// --------------- Setup OAuth ---------------
	setupOAuth()

	// --------------- Setup DB ---------------
	setupDB()

	// --------------- Setup Server ---------------
	fly.server = &fasthttp.Server{
		Handler:                       fasthttp.CompressHandler(fly.serveFastHTTP),
		ErrorHandler:                  fly.errorHandler,
		Name:                          fly.config.Name,
		Concurrency:                   fly.config.Concurrency,
		ReadTimeout:                   fly.config.ReadTimeout,
		WriteTimeout:                  fly.config.WriteTimeout,
		IdleTimeout:                   fly.config.IdleTimeout,
		ReadBufferSize:                fly.config.ReadBufferSize,
		WriteBufferSize:               fly.config.WriteBufferSize,
		NoDefaultDate:                 fly.config.NoDefaultDate,
		NoDefaultContentType:          fly.config.NoDefaultContentType,
		DisableHeaderNamesNormalizing: fly.config.DisableHeaderNamesNormalizing,
		DisableKeepalive:              fly.config.DisableKeepalive,
		MaxRequestBodySize:            fly.config.MaxRequestBodySize,
		NoDefaultServerHeader:         fly.config.NoDefaultServerHeader, // True when `Name` Empty
		GetOnly:                       fly.config.GetOnly,
		ReduceMemoryUsage:             fly.config.ReduceMemoryUsage,
		StreamRequestBody:             fly.config.StreamRequestBody,
		DisablePreParseMultipartForm:  fly.config.DisablePreParseMultipartForm,
	}

	url := fmt.Sprintf(
		"%s:%s",
		os.Getenv("SERVER_HOST"),
		os.Getenv("SERVER_PORT"),
	)

	// Print startup message
	if !fly.config.DisableStartupMessage {
		startupMessage(url)
	}

	certFile := utils.Getenv("SERVER_TLS_CERT", "")
	keyFile := utils.Getenv("SERVER_TLS_KEY", "")

	switch {
	case certFile != "" && keyFile != "":
		if err := fly.server.ListenAndServeTLS(url, certFile, keyFile); err != nil {
			log.Fatalf("Error start server %v", err)
		}
	default:
		log.Fatal(fly.server.ListenAndServe(url))
	}
}

// serveFastHTTP Serve FastHTTP via HTTP function
// The linking between fasthttp.RequestHandler to gFly's Ctx
func (fly *gFly) serveFastHTTP(ctx *fasthttp.RequestCtx) {
	handlerCtx := &Ctx{
		app:  fly,
		root: ctx,
		data: map[string]any{},
	}

	// Handle global middlewares
	var err error = nil
	for _, m := range fly.middlewares {
		err = m(handlerCtx)
		if err != nil {
			break
		}
	}

	if err == nil {
		_ = fly.router.Handler(handlerCtx)
	}
}

// errorHandler Server error handler.
func (fly *gFly) errorHandler(ctx *fasthttp.RequestCtx, err error) {
	log.Debug(ctx.String())
	log.Error(err)
}

// New Create new gFly app.
func New(config ...Config) IFly {
	app := &gFly{
		router:     NewRouter(),
		middleware: NewMiddleware(),
	}

	// Override config if provided
	if len(config) > 0 {
		app.config = config[0]
	} else {
		app.config = DefaultConfig
	}

	return app
}

// ===========================================================================================================
// 										gFly - Middleware methods
// ===========================================================================================================

// IFlyMiddleware Interface to declare all Middleware methods for gFly struct.
type IFlyMiddleware interface {
	// Use middleware to global (All requests)
	Use(middlewares ...MiddlewareHandler)
	// Middleware is a shortcut for Middleware.Group(middlewares ...MiddlewareHandler)
	Middleware(middleware ...MiddlewareHandler) func(IHandler) IHandler
}

// Use Middleware for global (All requests)
func (fly *gFly) Use(middlewares ...MiddlewareHandler) {
	fly.middlewares = append(fly.middlewares, middlewares...)
}

// Middleware is a shortcut for Middleware.Group(middlewares ...MiddlewareHandler)
func (fly *gFly) Middleware(middlewares ...MiddlewareHandler) func(IHandler) IHandler {
	return fly.middleware.Group(middlewares...)
}

// ===========================================================================================================
// 										gFly - HTTP methods
// ===========================================================================================================

// IFlyRouter Interface to declare all HTTP methods for gFly struct.
type IFlyRouter interface {
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
	// ServeFiles Serve static files
	ServeFiles()
}

// GET is a shortcut for Router.GET(path, handler)
func (fly *gFly) GET(path string, handler IHandler) {
	fly.router.GET(path, handler)
}

// HEAD is a shortcut for Router.HEAD(path, handler)
func (fly *gFly) HEAD(path string, handler IHandler) {
	fly.router.HEAD(path, handler)
}

// POST is a shortcut for Router.POST(path, handler)
func (fly *gFly) POST(path string, handler IHandler) {
	fly.router.POST(path, handler)
}

// PUT is a shortcut for Router.PUT(path, handler)
func (fly *gFly) PUT(path string, handler IHandler) {
	fly.router.PUT(path, handler)
}

// PATCH is a shortcut for Router.PATCH(path, handler)
func (fly *gFly) PATCH(path string, handler IHandler) {
	fly.router.PATCH(path, handler)
}

// DELETE is a shortcut for Router.DELETE(path, handler)
func (fly *gFly) DELETE(path string, handler IHandler) {
	fly.router.DELETE(path, handler)
}

// CONNECT is a shortcut for Router.CONNECT(path, handler)
func (fly *gFly) CONNECT(path string, handler IHandler) {
	fly.router.CONNECT(path, handler)
}

// OPTIONS is a shortcut for Router.OPTIONS(path, handler)
func (fly *gFly) OPTIONS(path string, handler IHandler) {
	fly.router.OPTIONS(path, handler)
}

// TRACE is a shortcut for Router.TRACE(path, handler)
func (fly *gFly) TRACE(path string, handler IHandler) {
	fly.router.TRACE(path, handler)
}

// Group Create a group Handler functions.
func (fly *gFly) Group(path string, groupFunc func(*Group)) {
	group := fly.router.Group(path)

	groupFunc(group)
}

// ServeFiles Serve static files from the given file system root is `./public`
// You can change parameter name STATIC_PATH.
// Use:
//
//	app.ServeFiles()
func (fly *gFly) ServeFiles() {
	// Default static file path
	rootPath := os.Getenv("STATIC_PATH")

	fly.router.ServeFiles("/{filepath:*}", rootPath)
}
