package river

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Renderer is output renderer.
type Renderer func(w http.ResponseWriter, r *http.Request, data interface{}) error

// JSONRenderer is a json renderer.
func JSONRenderer(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(data)
}

// PlainRenderer is plain text renderer.
func PlainRenderer(w http.ResponseWriter, r *http.Request, data interface{}) error {
	w.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprint(w, data)
	return err
}
