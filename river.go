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

// River is a REST server handler and toolkit.
type River struct {
	r *httprouter.Router
	HandlerChain
	renderer Renderer
	verbose
}

// New creates a new River.
// Optional params middlewares are the middlewares to initiate with.
// Middlewares can also be added with river.Use* methods.
func New(middlewares ...Handler) *River {
	r := httprouter.New()
	r.HandleMethodNotAllowed = true
	r.HandleOPTIONS = true
	r.RedirectTrailingSlash = true

	return (&River{r: r, HandlerChain: middlewares}).
		NotFound(notFound).
		NotAllowed(notAllowed)
}

func (rv *River) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rv.r.ServeHTTP(w, r)
}

// Handle handles endpoint at path p.
// This should only be called after Endpoint requests have been handled.
func (rv *River) Handle(p string, e *Endpoint) *River {
	rv.handle(p, e)
	return rv
}

func (rv *River) routerHandle(handler Handler, e *Endpoint) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c := &Context{
			rw:          w,
			Request:     r,
			params:      p,
			eRenderer:   e.renderer,
			gRenderer:   rv.renderer,
			middlewares: append(rv.HandlerChain, append(e.HandlerChain, handler)...),
		}
		c.Next()
	}
}

func (rv *River) routerHandleNoEndpoint(handler Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &Context{
			rw:          w,
			Request:     r,
			middlewares: append(rv.HandlerChain, handler),
			gRenderer:   rv.renderer,
		}
		c.Next()
	}
}

func (rv *River) handle(p string, e *Endpoint) {
	for subpath := range e.handlers {
		fullPath := path.Join(p, subpath)
		for method, handler := range e.handlers[subpath] {
			rv.r.Handle(method, fullPath, rv.routerHandle(handler, e))
			rv.handledPaths.add(method, fullPath, handler)
		}
	}
}

// Run starts River as an http server.
func (rv *River) Run(addr string) error {
	logger.Printf("Server started on %s", addr)
	rv.dump()
	return http.ListenAndServe(addr, rv)
}

// Renderer sets output renderer.
// Endpoint specific Renderer overrules this.
func (rv *River) Renderer(r Renderer) *River {
	rv.renderer = r
	return rv
}

// NotAllowed replaces the default handler for methods not handled by
// any endpoint with h.
func (rv *River) NotAllowed(h Handler) *River {
	rv.r.MethodNotAllowed = rv.routerHandleNoEndpoint(h)
	return rv
}

// NotFound replaces the default handler for request paths without
// any endpoint.
func (rv *River) NotFound(h Handler) *River {
	rv.r.NotFound = rv.routerHandleNoEndpoint(h)
	return rv
}

func notFound(c *Context) {
	c.RenderEmpty(http.StatusNotFound)
}

func notAllowed(c *Context) {
	c.RenderEmpty(http.StatusMethodNotAllowed)
}
