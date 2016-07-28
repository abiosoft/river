package river

import "net/http"

// Middleware is River middleware.
// A middleware needs to call c.Next()
// for request to continue to other middlewares.
//  func (c *river.Context){
//    // do something before
//    c.Next()
//    // do something after
//  }
type Middleware func(*Context)

// middlewareChain is chain of middlewares.
type middlewareChain []Middleware

// Use adds middlewares to the middleware chain.
func (c *middlewareChain) Use(middlewares ...Middleware) {
	*c = append(*c, middlewares...)
}

// UseHandler adds any http.Handler as middleware to the middleware chain.
func (c *middlewareChain) UseHandler(middlewares ...http.Handler) {
	for i := range middlewares {
		c.Use(toMiddleware(middlewares[i]))
	}
}

func toMiddleware(h http.Handler) Middleware {
	return func(c *Context) {
		h.ServeHTTP(c, c.Request)
		c.Next()
	}
}

// Recovery creates a panic recovery middleware.
// handlers are called after recovery.
func Recovery(handlers ...func(c *Context, err interface{})) Middleware {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if handlers != nil {
					for i := range handlers {
						handlers[i](c, err)
					}
				} else {
					c.Render(http.StatusInternalServerError, err)
				}
			}
		}()
		c.Next()
	}
}
