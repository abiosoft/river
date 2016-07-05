package main

import (
	"log"
	"net/http"

	"github.com/abiosoft/river"
)

func main() {
	rv := river.New()
	rv.BeforeHandle(func(r *http.Request) {
		log.Println(r.URL.Path)
	})
	rv.AfterHandle(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Length", w.Header())
	})
	{
		e := river.NewEndpoint()
		e.GetFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello world"))
		})

		rv.Handle("/user", e)
	}
	{
		e := river.NewEndpoint()
		e.Get(river.JSON(func(r *http.Request) (interface{}, error) {
			return map[string]string{
				"name": "Hello World",
			}, nil
		}))

		e.PutFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello from Put"))
		})

		rv.Handle("/json", e)
	}

	rv.Run(":8080")
}
