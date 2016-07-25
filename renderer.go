package river

import (
	"encoding/json"
	"fmt"
)

// M is a convenience wrapper for map[string]interface{}.
//  M{"status": "success, "data": M{"id": 1, "type": "complex"}}
type M map[string]interface{}

// Renderer renders data in a specified format.
// Renderer should set Content-Type accordingly.
type Renderer func(c *Context, data interface{}) error

// JSONRenderer is json renderer.
func JSONRenderer(c *Context, data interface{}) error {
	c.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(c).Encode(data)
}

// PlainRenderer is plain text renderer.
func PlainRenderer(c *Context, data interface{}) error {
	c.Header().Set("Content-Type", "text/plain")
	_, err := fmt.Fprint(c, data)
	return err
}

// ErrHandler handles error returned by Renderer.
type ErrHandler func(c *Context, err error)
