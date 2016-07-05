package river

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	logger = log.New(os.Stdout, "[River] ", 0)
)

// River is a REST handler.
type River struct {
	r      *mux.Router
	before []func(*http.Request)
	after  []http.HandlerFunc
	err    http.HandlerFunc
}

// New creates a new Bracket.
func New() *River {
	return &River{r: mux.NewRouter()}
}

func (rv *River) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rv.beforeFuncs(r)
	rv.r.ServeHTTP(w, r)
	rv.afterFuncs(w, r)
}

// Handle hadnles path with h.
func (rv *River) Handle(path string, h http.Handler) {
	if _, ok := h.(Endpoint); ok {
		rv.r.Handle(path, h)
		rv.r.Handle(path+"/{id}", h)
	} else {
		rv.r.Handle(path, h)
	}
}

// BeforeHandle executes before handler handles the request.
func (rv *River) BeforeHandle(f func(*http.Request)) {
	rv.before = append(rv.before, f)
}

// AfterHandle executes after handler has handled the request.
func (rv *River) AfterHandle(f http.HandlerFunc) {
	rv.after = append(rv.after, f)
}

// Run starts Bracket as an http server.
func (rv *River) Run(addr string) error {
	logger.Printf("Server started on %s", addr)
	return http.ListenAndServe(addr, rv)
}

func (rv *River) beforeFuncs(r *http.Request) {
	for i := range rv.before {
		rv.before[i](r)
	}
}

func (rv *River) afterFuncs(w http.ResponseWriter, r *http.Request) {
	for i := range rv.after {
		rv.after[i](w, r)
	}
}
