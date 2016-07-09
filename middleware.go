package river

import "net/http"

// Handler is request handler for endpoints and middlewares.
type Handler func(*Context)

// Chain is middleware chain.
type Chain []Handler

// Use adds middlewares to the middleware chain.
func (c *Chain) Use(middlewares ...Handler) {
	*c = append(*c, middlewares...)
}

// UseHandler adds any http.Handler as middleware to the middleware chain.
func (c *Chain) UseHandler(middlewares ...http.Handler) {
	for i := range middlewares {
		c.Use(httpHandlerToHandler(middlewares[i]))
	}
}

// UseHandlerNext adds any func(rw, r, next) as a middleware to the middleware
// chain. This adds compatibility to generic middlewares.
func (c *Chain) UseHandlerNext(middlewares ...func(w http.ResponseWriter, r *http.Request, next http.Handler)) {
	for i := range middlewares {
		c.Use(handlerNextToHandler(middlewares[i]))
	}
}

func httpHandlerToHandler(h http.Handler) Handler {
	return func(c *Context) {
		h.ServeHTTP(c, c.Request)
		c.Next()
	}
}

func handlerNextToHandler(handlerFunc func(w http.ResponseWriter, r *http.Request, next http.Handler)) Handler {
	return func(c *Context) {
		handlerFunc(
			c,
			c.Request,
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.Next()
			}),
		)
	}
}

func handlerToHTTPHandler(h Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(&Context{rw: w, Request: r, renderer: JSONRenderer})
	}
}
