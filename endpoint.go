package river

import (
	"net/http"

	"github.com/gorilla/mux"
)

const wildcardMethod = "*"

// Endpoint is a REST endpoint.
type Endpoint map[string]http.HandlerFunc

// NewEndpoint creates a new Endpoint.
func NewEndpoint() Endpoint {
	return make(Endpoint)
}

func (e Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	method := wildcardMethod

	_, ok := mux.Vars(r)["id"]

	// Attempt to determine if valid request type.
	switch r.Method {
	case "PUT", "PATCH", "DELETE":
		if ok {
			method = r.Method
		}
	case "POST":
		if !ok {
			method = r.Method
		}
	case "GET":
		method = r.Method
	default:
		method = r.Method
	}

	if h, ok := e[method]; ok {
		h(w, r)
	} else if h, ok := e[wildcardMethod]; ok {
		h(w, r)
	} else {
		notImplemented(w, r)
	}
}

// Get sets the handler for a GET request to the endpoint.
func (e *Endpoint) Get(h http.Handler) {
	e.handleMethod("GET", h)
}

// GetFunc sets the handler function for a GET request to the endpoint.
func (e *Endpoint) GetFunc(h http.HandlerFunc) {
	e.handleMethod("GET", h)
}

// Post sets the handler for a POST request to the endpoint.
func (e *Endpoint) Post(h http.Handler) {
	e.handleMethod("POST", h)
}

// PostFunc sets the handler function for a POST request to the endpoint.
func (e *Endpoint) PostFunc(h http.HandlerFunc) {
	e.handleMethod("POST", h)
}

// Put sets the handler for a PUT request to the endpoint.
func (e *Endpoint) Put(h http.Handler) {
	e.handleMethod("PUT", h)
}

// PutFunc sets the handler function for a PUT request to the endpoint.
func (e *Endpoint) PutFunc(h http.HandlerFunc) {
	e.handleMethod("PUT", h)
}

// Delete sets the handler for a DELETE request to the endpoint.
func (e *Endpoint) Delete(h http.Handler) {
	e.handleMethod("DELETE", h)
}

// DeleteFunc sets the handler function for a DELETE request to the endpoint.
func (e *Endpoint) DeleteFunc(h http.HandlerFunc) {
	e.handleMethod("DELETE", h)
}

// Patch sets the handler for a PATCH request to the endpoint.
func (e *Endpoint) Patch(h http.Handler) {
	e.handleMethod("PATCH", h)
}

// PatchFunc sets the handler function for a PATCH request to the endpoint.
func (e *Endpoint) PatchFunc(h http.HandlerFunc) {
	e.handleMethod("PATCH", h)
}

// All sets the handler for all other requests to the endpoint
// that has not been handled.
func (e *Endpoint) All(h http.Handler) {
	e.handleMethod(wildcardMethod, h)
}

// AllFunc sets the handler function for all other requests to the endpoint
// that has not been handled.
func (e *Endpoint) AllFunc(h http.HandlerFunc) {
	e.handleMethod(wildcardMethod, h)
}

// HandleMethod handles the request method using h.
func (e *Endpoint) HandleMethod(method string, h http.Handler) {
	e.handleMethod(method, h)
}

func (e *Endpoint) handleMethod(method string, h http.Handler) {
	(*e)[method] = h.ServeHTTP
}

func notImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(
		w,
		http.StatusText(http.StatusNotImplemented),
		http.StatusNotImplemented,
	)
}
