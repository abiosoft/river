package river

import (
	"errors"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

var errNoRenderer = errors.New("Renderer is nil")

// Context is a request scope context.
// Context implements http.ResponseWriter and embeds http.Request.
// It can be adapted for use in an http.Handler e.g.
//  handler.ServeHTTP(c, c.Request)
type Context struct {
	*http.Request
	rw          http.ResponseWriter
	params      httprouter.Params
	values      map[string]interface{}
	renderer    Renderer
	middlewares []Handler
}

// Param returns URI parameters. If key is not found,
// empty string is returned.
// Params are set with :key in the handle path.
// e.g. /:category/:id
func (c *Context) Param(key string) string {
	return c.params.ByName(key)
}

// Next calls the next handler in the middleware chain.
// A middleware must call Next, otherwise the request stops
// at the middleware.
// Next has no effect if called in an endpoint handler.
func (c *Context) Next() {
	if len(c.middlewares) < 1 {
		return
	}
	current := c.middlewares[0]
	c.middlewares = c.middlewares[1:]
	current(c)
}

// Header returns the header map that will be sent by
// WriteHeader. Changing the header after a call to
// WriteHeader (or Write) has no effect unless the modified
// headers were declared as trailers by setting the
// "Trailer" header before the call to WriteHeader (see example).
// To suppress implicit response headers, set their value to nil.
func (c *Context) Header() http.Header {
	return c.rw.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (c *Context) Write(b []byte) (int, error) {
	return c.rw.Write(b)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (c *Context) WriteHeader(status int) {
	c.rw.WriteHeader(status)
}

// Get gets the value for key in the context. Key must have been
// previously set using c.Set.
func (c *Context) Get(key string) interface{} {
	return c.values[key]
}

// Set sets key in context to value.
func (c *Context) Set(key string, value interface{}) {
	if c.values == nil {
		c.values = make(map[string]interface{})
	}
	c.values[key] = value
}

// Render renders data using the current endpoint's renderer.
// status is HTTP status code to respond with.
func (c *Context) Render(data interface{}, status int) error {
	if c.renderer == nil {
		return errNoRenderer
	}
	c.WriteHeader(status)
	return c.renderer(c, data)
}

// net/context / Go 1.7 Request.Context

// Deadline satisfies Context.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done satisfies Context.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err satisfies Context.
func (c *Context) Err() error {
	return nil
}

// Value satisfies Context.
func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if k, ok := key.(string); ok {
		return c.Get(k)
	}
	return nil
}
