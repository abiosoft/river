package river

// Endpoint is a REST endpoint.
type Endpoint struct {
	models   map[string]Model
	renderer Renderer
}

// NewEndpoint creates a new Endpoint with p as path prefix.
// Renderer defaults to JSONRenderer.
func NewEndpoint() *Endpoint {
	return &Endpoint{
		models:   make(map[string]Model),
		renderer: JSONRenderer,
	}
}

// Use uses Model m at subpath.
// Supports variable templates. e.g. /{category}/{id:[0-9]+}
// and accessible via river.Vars(r).
func (e *Endpoint) Use(subpath string, m Model) *Endpoint {
	e.models[subpath] = m
	return e
}

// Renderer sets the output render for Endpoint.
func (e *Endpoint) Renderer(r Renderer) *Endpoint {
	e.renderer = r
	return e
}
