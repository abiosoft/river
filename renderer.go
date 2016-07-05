package river

import (
	"encoding/json"
	"net/http"
)

// RendererFunc is returned by Renderers
type RendererFunc func(w http.ResponseWriter, r *http.Request) error

func (rf RendererFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rf(w, r)
}

// Renderer is output Renderer.
type Renderer func(f InputFunc) RendererFunc

// InputFunc is passed to a Renderer.
type InputFunc func(*http.Request) (interface{}, error)

// JSON is a json renderer.
var JSON Renderer = func(f InputFunc) RendererFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		data, _ := f(r)
		return json.NewEncoder(w).Encode(data)
	}
}
