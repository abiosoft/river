package river

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

var (
	logger = log.New(os.Stdout, "[River] ", 0)
)

// LogOutput sets the writer to use for logger.
func LogOutput(w io.Writer) {
	logger.SetOutput(w)
}

// River is a REST handler.
type River struct {
	r      *mux.Router
	before []func(*http.Request)
	after  []http.HandlerFunc
	err    ErrorFunc
}

// New creates a new River.
func New() *River {
	return &River{r: mux.NewRouter()}
}

func (rv *River) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rv.beforeFuncs(r)
	rv.r.ServeHTTP(w, r)
	rv.afterFuncs(w, r)
}

// Handle handles endpoint at path p.
func (rv *River) Handle(p string, e *Endpoint) *River {
	for subpath, model := range e.models {
		fullPath := path.Join(p, subpath)
		rv.r.HandleFunc(fullPath, func(w http.ResponseWriter, r *http.Request) {
			vars := map[string]string{}
			if m := new(mux.RouteMatch); rv.r.Match(r, m) {
				vars = m.Vars
			}

			mFunc := modelFunc(r.Method, model)
			if mFunc == nil {
				notAllowed(w, r)
				return
			}

			// render
			err := e.renderer(mFunc(r, vars))(w, r)
			if err != nil && rv.err != nil {
				rv.err(w, r, err)
			}
		})
	}
	return rv
}

// BeforeHandle executes before handler handles the request.
func (rv *River) BeforeHandle(f func(*http.Request)) *River {
	rv.before = append(rv.before, f)
	return rv
}

// AfterHandle executes after handler has handled the request.
func (rv *River) AfterHandle(f http.HandlerFunc) *River {
	rv.after = append(rv.after, f)
	return rv
}

// Run starts Bracket as an http server.
func (rv *River) Run(addr string) error {
	logger.Printf("Server started on %s", addr)
	return http.ListenAndServe(addr, rv)
}

// Err registers f as error handler.
func (rv *River) Err(f ErrorFunc) *River {
	rv.err = f
	return rv
}

func (rv *River) beforeFuncs(r *http.Request) *River {
	for i := range rv.before {
		rv.before[i](r)
	}
	return rv
}

func (rv *River) afterFuncs(w http.ResponseWriter, r *http.Request) {
	for i := range rv.after {
		rv.after[i](w, r)
	}
}

// HandlerFunc is a function definition for Endpoint handlers.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

func (h HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(w, r)
}

// ErrorFunc handles error that occurs during request handling.
type ErrorFunc func(http.ResponseWriter, *http.Request, error)

func notAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(
		w,
		http.StatusText(http.StatusMethodNotAllowed),
		http.StatusNotImplemented,
	)
}
