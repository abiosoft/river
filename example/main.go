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
		var user river.Model
		user.Get(func(r *http.Request) (interface{}, int) {
			v := river.Vars(r)
			return map[string]interface{}{
				"name": "Hello",
				"id":   v["id"],
			}, 201
		})
		user.Delete(func(r *http.Request) (interface{}, int) {
			return nil, 201
		})
		var people river.Model
		people.Get(func(r *http.Request) (interface{}, int) {
			v := river.Vars(r)
			return map[string]interface{}{
				"type": "People",
				"name": v["name"],
			}, 201
		})
		rv.Handle("/user", river.NewEndpoint().
			Use("/{id}", user),
		)
		rv.Handle("/people", river.NewEndpoint().
			Use("/{name}", people),
		)
		rv.Handle("/people", river.NewEndpoint().
			Use("/", people),
		)
	}

	rv.Run(":8080")

}
