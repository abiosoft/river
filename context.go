package river

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Context is a request scope context.
// Context implements http.ResponseWriter and embeds http.Request.
//
// It can be adapted for use in an http.Handler e.g.
//  handler.ServeHTTP(c, c.Request)
type Context struct {
	*http.Request
	rw            http.ResponseWriter
	params        httprouter.Params
	values        map[string]interface{}
	renderer      Renderer
	middlewares   []Handler
	jsonDecoder   jsonDecoder
	headerWritten bool
	status        int
	serviceInjector
}

// Param returns URL parameters. If key is not found,
// empty string is returned.
//
// Params are set with :key in the handle path.
// e.g. /:category/:id
func (c *Context) Param(key string) string {
	return c.params.ByName(key)
}

// Query returns URL query parameters. If key not found,
// empty string is returned.
func (c *Context) Query(key string) string {
	return c.URL.Query().Get(key)
}

// Redirect performs HTTP redirect to url with code as redirect code.
// code must be 3xx, otherwise http.StatusFound (302) will be used.
func (c *Context) Redirect(url string, code int) {
	if code < 300 || code > 399 {
		code = http.StatusFound
	}
	http.Redirect(c, c.Request, url, code)
}

// Next calls the next handler in the middleware chain.
// A middleware must call Next, otherwise the request stops
// at the middleware.
//
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
	if !c.headerWritten {
		c.WriteHeader(http.StatusOK)
	}
	return c.rw.Write(b)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (c *Context) WriteHeader(status int) {
	c.status = status
	c.headerWritten = true
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

// Render renders data using the current endpoint's renderer (if any)
// or global renderer (if any) or PlainRenderer; in that preference order.
// status is HTTP status code to respond with.
func (c *Context) Render(status int, data interface{}) error {
	c.WriteHeader(status)
	return c.renderer(c, data)
}

// RenderEmpty renders status text for status as body.
// status is HTTP status code to respond with.
func (c *Context) RenderEmpty(status int) error {
	c.WriteHeader(status)
	return PlainRenderer(c, http.StatusText(status))
}

// Status returns the response status code. This returns 0 unless response
// has been written.
func (c *Context) Status() int {
	return c.status
}

// DecodeJSONBody decodes the request body as JSON into v.
//
// The request body must be JSON and v must be a pointer to
// a compatible type for the JSON body.
//
// Type conversion is done if required; based on v's underlying type.
// If v points to a struct and request body is a json array,
// an attempt is made to decode to a slice of the struct and the
// first element of the slice will be stored in v.
//
// Likewise if v points to a slice and request body is a json object,
// an attempt is made to decode to the item type of the slice and a slice
// containing the item will be stored in v.
//
//  var v []Type // c.DecodeJSONBody(&v) works even if body is a json object.
//  var v Type // c.DecodeJSONBody(&v) works even if body is a json array.
func (c *Context) DecodeJSONBody(v interface{}) error {
	if c.jsonDecoder == nil {
		if err := c.jsonDecoder.init(c.Request.Body); err != nil {
			return err
		}
	}
	return c.jsonDecoder.decode(v)
}

/* net/context / Go 1.7 Request.Context */

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
func (c *Context) Done() <-chan struct{} {
	return nil
}

// Err returns a non-nil error value after Done is closed. Err returns
// Canceled if the context was canceled or DeadlineExceeded if the
// context's deadline passed. No other values for Err are defined.
// After Done is closed, successive calls to Err return the same value.
func (c *Context) Err() error {
	return nil
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
func (c *Context) Value(key interface{}) interface{} {
	if k, ok := key.(string); ok {
		return c.Get(k)
	}
	return nil
}
