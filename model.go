package river

import "net/http"

// Vars is URI variables.
type Vars map[string]string

// ModelFunc is a function returned by a model.
// data is the data return by the model and status is HTTP status code.
type ModelFunc func(*http.Request, Vars) (data interface{}, status int)

// Model is a REST endpoint data model.
type Model interface {
	Get(*http.Request, Vars) (interface{}, int)
	Post(*http.Request, Vars) (interface{}, int)
	Put(*http.Request, Vars) (interface{}, int)
	Patch(*http.Request, Vars) (interface{}, int)
	Delete(*http.Request, Vars) (interface{}, int)
}

func modelFunc(method string, model Model) ModelFunc {
	var m ModelFunc
	switch method {
	case "GET":
		m = model.Get
	case "POST":
		m = model.Post
	case "PUT":
		m = model.Put
	case "PATCH":
		m = model.Patch
	case "DELETE":
		m = model.Delete
	}
	return m
}

// EmptyModel is an empty data model implementation.
type EmptyModel map[string]ModelFunc

// Get satisfies Model.
func (e EmptyModel) Get(r *http.Request, v Vars) (interface{}, int) {
	return e.handle("GET", r, v)
}

// Post satisfies Model.
func (e EmptyModel) Post(r *http.Request, v Vars) (interface{}, int) {
	return e.handle("POST", r, v)
}

// Put satisfies Model.
func (e EmptyModel) Put(r *http.Request, v Vars) (interface{}, int) {
	return e.handle("PUT", r, v)
}

// Patch satisfies Model.
func (e EmptyModel) Patch(r *http.Request, v Vars) (interface{}, int) {
	return e.handle("PATCH", r, v)
}

// Delete satisfies Model.
func (e EmptyModel) Delete(r *http.Request, v Vars) (interface{}, int) {
	return e.handle("DELETE", r, v)
}

// GetFunc sets the model function for Get requests.
func (e *EmptyModel) GetFunc(f ModelFunc) {
	e.setFunc("GET", f)
}

// PostFunc sets the model function for Get requests.
func (e *EmptyModel) PostFunc(f ModelFunc) {
	e.setFunc("POST", f)
}

// PutFunc sets the model function for Get requests.
func (e *EmptyModel) PutFunc(f ModelFunc) {
	e.setFunc("PUT", f)
}

// PatchFunc sets the model function for Get requests.
func (e *EmptyModel) PatchFunc(f ModelFunc) {
	e.setFunc("PATCH", f)
}

// DeleteFunc sets the model function for Get requests.
func (e *EmptyModel) DeleteFunc(f ModelFunc) {
	e.setFunc("DELETE", f)
}

func (e *EmptyModel) setFunc(method string, f ModelFunc) {
	if *e == nil {
		*e = make(EmptyModel)
	}
	(*e)[method] = f
}

func (e EmptyModel) handle(method string, r *http.Request, v Vars) (interface{}, int) {
	if modelFunc, ok := e[method]; ok {
		return modelFunc(r, v)
	}
	return nil, 0
}
