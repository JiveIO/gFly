package gfly

// MiddlewareHandler must return RequestHandler for continuing or error to stop on it.
type MiddlewareHandler func(ctx *Ctx) error

// IMiddleware Middleware interface
type IMiddleware interface {
	Group(middlewares ...MiddlewareHandler) func(IHandler) IHandler
}

// Middleware Middleware type
type Middleware struct{}

// Group Create a group Middleware functions. Implement for IMiddleware interface
func (m *Middleware) Group(middlewares ...MiddlewareHandler) func(IHandler) IHandler {
	return func(handler IHandler) IHandler {
		return &middlewareEndpoint{
			handler:     handler,
			middlewares: middlewares,
		}
	}
}

// middlewareEndpoint Default handler
// Need to wrap middleware as a implementation of IHandler interface
type middlewareEndpoint struct {
	Endpoint
	handler     IHandler
	middlewares []MiddlewareHandler
}

// Validate execute Handler's validation
// Middleware itself don't have any validation.
// Of course, we don't need this method because it was created in the embed struct Endpoint.
// But we override Endpoint because the handler needs to validate data as well.
func (m *middlewareEndpoint) Validate(c *Ctx) error {
	// But handler need to validate data. So, need to implement
	return m.handler.Validate(c)
}

// Handle check handler validation
func (m *middlewareEndpoint) Handle(c *Ctx) error {
	// Run middleware functions
	for _, m := range m.middlewares {
		err := m(c)
		if err != nil {
			return err
		}
	}

	// Execute Handler's handling
	return m.handler.Handle(c)
}

func NewMiddleware() IMiddleware {
	return &Middleware{}
}
