package river

import (
	"net/http"
	"reflect"
)

// EndpointHandler is an endpoint handler with support for dependency injection.
// Any function type is a valid EndpointHandler; including river.Handler and
// http.Handler. Function parameters will be injected accordingly.
//
// If a service is not previously registered and it is not one of
// *river.Context, http.ResponseWriter and *http.Request, zero value of the type
// will be passed as the parameter.
//
// If a non function type is passed as EndpointHandler to Endpoint
// request functions (Get, Post e.t.c.), a panic occurs.
//
// The return values of the function (if any) are discarded.
type EndpointHandler interface{}

// Endpoint is a REST endpoint.
type Endpoint struct {
	handlers map[string]endpointFuncs
	renderer Renderer
	handlerChain
	serviceInjector
}

// NewEndpoint creates a new Endpoint.
func NewEndpoint() *Endpoint {
	return &Endpoint{
		handlers: make(map[string]endpointFuncs),
	}
}

// Get sets the function for Get requests.
func (e *Endpoint) Get(p string, h EndpointHandler) *Endpoint {
	e.set(p, "GET", h)
	return e
}

// Post sets the function for Post requests.
func (e *Endpoint) Post(p string, h EndpointHandler) *Endpoint {
	e.set(p, "POST", h)
	return e
}

// Put sets the function for Put requests.
func (e *Endpoint) Put(p string, h EndpointHandler) *Endpoint {
	e.set(p, "PUT", h)
	return e
}

// Patch sets the function for Patch requests.
func (e *Endpoint) Patch(p string, h EndpointHandler) *Endpoint {
	e.set(p, "PATCH", h)
	return e
}

// Delete sets the function for Delete requests.
func (e *Endpoint) Delete(p string, h EndpointHandler) *Endpoint {
	e.set(p, "DELETE", h)
	return e
}

// Options sets the function for Options requests.
func (e *Endpoint) Options(p string, h EndpointHandler) *Endpoint {
	e.set(p, "OPTIONS", h)
	return e
}

// Renderer sets the output renderer for endpoint.
func (e *Endpoint) Renderer(r Renderer) *Endpoint {
	e.renderer = r
	return e
}

// Handle sets the function for a custom requests.
func (e *Endpoint) Handle(requestMethod, p string, h EndpointHandler) *Endpoint {
	e.set(p, requestMethod, h)
	return e
}

func (e *Endpoint) set(subpath string, method string, h EndpointHandler) {
	if e.handlers[subpath] == nil {
		e.handlers[subpath] = make(endpointFuncs)
	}
	if reflect.TypeOf(h).Kind() != reflect.Func {
		// this is the beginning of the app, safer to panic here
		// and prevent possible request time panic.
		panic("Cannot use non function type as EndpointHandler")
	}
	e.handlers[subpath][method] = h
}

// endpointFuncs maps request method to EndpointFunc.
type endpointFuncs map[string]EndpointHandler

func endpointToHandler(h EndpointHandler) Handler {
	if handler, ok := h.(Handler); ok {
		return handler
	}

	return func(c *Context) {
		/* default injections */
		// context
		c.register(c)

		// responseWriter
		var rw http.ResponseWriter = c
		c.register(rw)

		// request
		c.register(c.Request)

		/* handle request */
		c.invoke(h)
	}

}
