package river

import (
	"encoding/json"
	"fmt"
)

// Renderer renders data in a specified format.
// Render should set Content-Type accordingly.
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
