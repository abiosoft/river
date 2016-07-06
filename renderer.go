package river

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Renderer is output renderer.
type Renderer func(data interface{}, status int) HandlerFunc

// JSONRenderer is a json renderer.
var JSONRenderer Renderer = func(data interface{}, status int) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		return json.NewEncoder(w).Encode(data)
	}
}

// PlainRenderer is plain text renderer.
var PlainRenderer Renderer = func(data interface{}, status int) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(status)
		_, err := fmt.Fprint(w, data)
		return err
	}
}
