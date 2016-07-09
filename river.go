package river

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/julienschmidt/httprouter"
)

var (
	logger = log.New(os.Stdout, "[River] ", 0)
)

// LogOutput sets the output to use for logger.
func LogOutput(w io.Writer) {
	logger.SetOutput(w)
}

// River is a REST handler.
type River struct {
	r           *httprouter.Router
	endpoints   map[string]*Endpoint
	middlewares []Handler
	after       []Handler
	notAllowed  Handler
}

// New creates a new River.
func New() *River {
	return &River{
		r:          httprouter.New(),
		notAllowed: notAllowed,
		endpoints:  make(map[string]*Endpoint),
	}
}

func (rv *River) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rv.r.ServeHTTP(w, r)
}

// Handle handles endpoint at path p.
func (rv *River) Handle(p string, e *Endpoint) *River {
	rv.handle(p, e)
	return rv
}

func (rv *River) routerHandle(handler Handler, renderer Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// create request context
		c := &Context{
			rw:      w,
			Request: r,
			params:  p,
		}
		// set renderer
		c.renderer = renderer

		// run middlewares
		rv.runMiddlewares(c)

		// return if response has been written
		if c.responseWritten {
			return
		}

		// handle request
		handler(c)

		// run after functions
		rv.afterFuncs(c)
	}
}

func (rv *River) handle(p string, e *Endpoint) {
	for subpath := range e.handlers {
		fullPath := path.Join(p, subpath)
		rv.endpoints[fullPath] = e
		for method, handler := range e.handlers[subpath] {
			rv.r.Handle(method, fullPath, rv.routerHandle(handler, e.renderer))
		}
	}
}

// Use adds middlewares to the middleware chain.
func (rv *River) Use(middlewares ...Handler) *River {
	rv.middlewares = append(rv.middlewares, middlewares...)
	return rv
}

// UseHandler adds any http.Handler as middleware to the middleware chain.
func (rv *River) UseHandler(middlewares ...http.Handler) *River {
	for i := range middlewares {
		rv.Use(toHandler(middlewares[i]))
	}
	return rv
}

// After executes after the request has been responded to.
// Useful for logging e.t.c.
func (rv *River) After(h Handler) *River {
	rv.after = append(rv.after, h)
	return rv
}

// Run starts River as an http server.
func (rv *River) Run(addr string) error {
	logger.Printf("Server started on %s", addr)
	return http.ListenAndServe(addr, rv)
}

// NotAllowed replaces the default handler for methods not handled by
// any endpoint with f.
func (rv *River) NotAllowed(f Handler) *River {
	rv.notAllowed = f
	return rv
}

func (rv *River) runMiddlewares(c *Context) {
	for i := range rv.middlewares {
		rv.middlewares[i](c)
		// stop middleware chain if
		// response is written.
		if c.responseWritten {
			return
		}
	}
}

func (rv *River) afterFuncs(c *Context) {
	for i := range rv.after {
		rv.after[i](c)
	}
}

func toHandler(h http.Handler) Handler {
	return func(c *Context) {
		h.ServeHTTP(c, c.Request)
	}
}

func notAllowed(c *Context) {
	http.Error(
		c,
		http.StatusText(http.StatusMethodNotAllowed),
		http.StatusMethodNotAllowed,
	)
}

// type handledPaths map[string][]string

// func (h handledPaths) Dump() {
// 	logger.Println()
// 	logger.Println("Routes")
// 	logger.Println("-------")
// 	for path, methods := range h {
// 		logger.Printf("%s -> %s \n", path, strings.Join(methods, ", "))
// 	}
// 	logger.Println("-------")
// }
