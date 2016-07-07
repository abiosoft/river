package river

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Vars returns URI variables.
func Vars(r *http.Request) map[string]string {
	return mux.Vars(r)
}

// ModelFunc is a function executed by models to handle requests.
// data is the data return by the model and status is HTTP status code.
// 0 status code should be returned if request is ignored or not handled.
type ModelFunc func(*http.Request) (data interface{}, status int)

// Model is a REST endpoint data model.
// It maps request types (e.g. "GET", POST") as key to ModelFunc.
// Empty value of Model is usable.
type Model map[string]ModelFunc

// Get sets the model function for Get requests.
func (e *Model) Get(f ModelFunc) {
	e.set("GET", f)
}

// Post sets the model function for Get requests.
func (e *Model) Post(f ModelFunc) {
	e.set("POST", f)
}

// Put sets the model function for Get requests.
func (e *Model) Put(f ModelFunc) {
	e.set("PUT", f)
}

// Patch sets the model function for Get requests.
func (e *Model) Patch(f ModelFunc) {
	e.set("PATCH", f)
}

// Delete sets the model function for Get requests.
func (e *Model) Delete(f ModelFunc) {
	e.set("DELETE", f)
}

func (e *Model) set(method string, f ModelFunc) {
	if *e == nil {
		*e = make(Model)
	}
	(*e)[method] = f
}

func (e Model) handle(r *http.Request) (interface{}, int) {
	if modelFunc, ok := e[r.Method]; ok {
		return modelFunc(r)
	}
	return nil, 0
}
