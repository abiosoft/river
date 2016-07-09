package river

// Endpoint is a REST endpoint.
type Endpoint struct {
	handlers map[string]endpointFuncs
	renderer Renderer
}

// NewEndpoint creates a new Endpoint.
// Renderer defaults to JSONRenderer.
func NewEndpoint() *Endpoint {
	return &Endpoint{
		handlers: make(map[string]endpointFuncs),
		renderer: JSONRenderer,
	}
}

// Renderer sets the output render for Endpoint.
func (e *Endpoint) Renderer(r Renderer) *Endpoint {
	e.renderer = r
	return e
}

// Handler is request handler for endpoints and middlewares.
// The first middleware to write a response stops both
// middleware chain and the request from reaching the endpoint.
// This can be done either by context.Write or context.Render.
// Otherwise, every middleware will be called in order in which they are
// added.
// Endpoints handlers are called after all middleware.
type Handler func(*Context)

// Get sets the model function for Get requests.
func (e *Endpoint) Get(p string, f Handler) *Endpoint {
	e.set(p, "GET", f)
	return e
}

// Post sets the model function for Post requests.
func (e *Endpoint) Post(p string, f Handler) *Endpoint {
	e.set(p, "POST", f)
	return e
}

// Put sets the model function for Put requests.
func (e *Endpoint) Put(p string, f Handler) *Endpoint {
	e.set(p, "PUT", f)
	return e
}

// Patch sets the model function for Patch requests.
func (e *Endpoint) Patch(p string, f Handler) *Endpoint {
	e.set(p, "PATCH", f)
	return e
}

// Delete sets the model function for Delete requests.
func (e *Endpoint) Delete(p string, f Handler) *Endpoint {
	e.set(p, "DELETE", f)
	return e
}

// Options sets the model function for Options requests.
func (e *Endpoint) Options(p string, f Handler) *Endpoint {
	e.set(p, "OPTIONS", f)
	return e
}

func (e *Endpoint) set(subpath string, method string, f Handler) {
	if e.handlers[subpath] == nil {
		e.handlers[subpath] = make(endpointFuncs)
	}
	e.handlers[subpath][method] = f
}

// endpointFuncs maps request method to EndpointFunc.
type endpointFuncs map[string]Handler

// runtime.FuncForPC(reflect.ValueOf(a).Pointer()).Name()
