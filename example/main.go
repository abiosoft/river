package main

import (
	"net/http"

	"github.com/abiosoft/river"
)

func main() {
	rv := river.New()
	{
		var user river.EmptyModel
		user.GetFunc(func(r *http.Request, v river.Vars) (interface{}, int) {
			return map[string]interface{}{
				"name": "Hello",
				"id":   v["id"],
			}, 200
		})
		rv.Handle("/user", river.NewEndpoint().
			Use("/{id}", user).
			Use("/json", user).
			Renderer(river.PlainRenderer),
		)
	}
	rv.Run(":8080")
}
