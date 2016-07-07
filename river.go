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

// LogOutput sets the output to use for logger.
func LogOutput(w io.Writer) {
	logger.SetOutput(w)
}

// River is a REST handler.
type River struct {
	r            *mux.Router
	beforeHandle []http.HandlerFunc
	beforeWrite  []http.HandlerFunc
	afterHandle  []http.HandlerFunc
	err          ErrorFunc
	notAllowed   http.HandlerFunc
}

// New creates a new River.
func New() *River {
	return &River{
		r:          mux.NewRouter(),
		notAllowed: notAllowed,
	}
}

func (rv *River) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rv.beforeHandleFuncs(noWriteRW(w), r)
	rv.r.ServeHTTP(beforeWriteRW(w, func() { rv.beforeWriteFuncs(w, r) }), r)
	rv.afterHandleFuncs(w, r)
}

// Handle handles endpoint at path p.
func (rv *River) Handle(p string, e *Endpoint) *River {
	for subpath, model := range e.models {
		fullPath := path.Join(p, subpath)
		rv.r.HandleFunc(fullPath, func(w http.ResponseWriter, r *http.Request) {
			mf := modelFunc(r.Method, model)
			if mf == nil {
				rv.notAllowed(w, r)
				return
			}

			// render
			data, status := mf(r)
			if status == 0 {
				rv.notAllowed(w, r)
				return
			}

			err := e.renderer(staticStatusRW(w, status), r, data)
			if err != nil && rv.err != nil {
				rv.err(w, r, err)
			}
		})
	}
	return rv
}

// BeforeHandle executes before handler handles the request.
// The passed ResponseWriter to the HandlerFunc can only modify the headers
// and has Write() and WriteHeader() invalidated.
func (rv *River) BeforeHandle(f http.HandlerFunc) *River {
	rv.beforeHandle = append(rv.beforeHandle, f)
	return rv
}

// BeforeWrite executes before handler writes to the ResponseWriter.
func (rv *River) BeforeWrite(f http.HandlerFunc) *River {
	rv.beforeWrite = append(rv.beforeWrite, f)
	return rv
}

// AfterHandle executes after handler has handled the request.
func (rv *River) AfterHandle(f http.HandlerFunc) *River {
	rv.afterHandle = append(rv.afterHandle, f)
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

// NotAllowed replaces the default handler for methods not handled by
// any endpoint with f.
func (rv *River) NotAllowed(f http.HandlerFunc) *River {
	rv.notAllowed = f
	return rv
}

func (rv *River) beforeHandleFuncs(w http.ResponseWriter, r *http.Request) {
	for i := range rv.beforeHandle {
		rv.beforeHandle[i](w, r)
	}
}

func (rv *River) beforeWriteFuncs(w http.ResponseWriter, r *http.Request) {
	for i := range rv.beforeWrite {
		rv.beforeWrite[i](w, r)
	}
}

func (rv *River) afterHandleFuncs(w http.ResponseWriter, r *http.Request) {
	for i := range rv.afterHandle {
		rv.afterHandle[i](w, r)
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
		http.StatusMethodNotAllowed,
	)
}
