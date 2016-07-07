package main

import (
	"net/http"

	"github.com/abiosoft/river"
)

func main() {
	rv := river.New().
		AfterHandle(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-River-Afer", "Hello World After")
		}).
		BeforeHandle(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-River-Before", "Hello World Before")
		}).
		BeforeWrite(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-River-BeforeWrite", "Hello World Before Write")
		})
	{
		var user river.EmptyModel
		user.GetFunc(func(r *http.Request) (interface{}, int) {
			v := river.Vars(r)
			return map[string]interface{}{
				"name": "Hello",
				"id":   v["id"],
			}, 201
		})
		rv.Handle("/user", river.NewEndpoint().
			Use("/{id}", user).
			Use("/json", user).
			Renderer(river.JSONRenderer),
		)
	}

	rv.Run(":8080")

}
