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
func (m *Model) Get(f ModelFunc) {
	m.set("GET", f)
}

// Post sets the model function for Get requests.
func (m *Model) Post(f ModelFunc) {
	m.set("POST", f)
}

// Put sets the model function for Get requests.
func (m *Model) Put(f ModelFunc) {
	m.set("PUT", f)
}

// Patch sets the model function for Get requests.
func (m *Model) Patch(f ModelFunc) {
	m.set("PATCH", f)
}

// Delete sets the model function for Get requests.
func (m *Model) Delete(f ModelFunc) {
	m.set("DELETE", f)
}

func (m *Model) set(method string, f ModelFunc) {
	if *m == nil {
		*m = make(Model)
	}
	(*m)[method] = f
}

func (m Model) handle(r *http.Request) (interface{}, int) {
	if modelFunc, ok := m[r.Method]; ok {
		return modelFunc(r)
	}
	return nil, 0
}

func (m Model) methods() []string {
	var methods []string
	for method := range m {
		methods = append(methods, method)
	}
	return methods
}
