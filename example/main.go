package main

import (
	"fmt"
	"time"

	"github.com/abiosoft/river"
)

func main() {
	rv := river.New(river.Logger())
	{
		rv.Handle("/user", river.NewEndpoint().
			Get("/:id", func(c *river.Context) {
				c.Render(c.Param("id"), 200)
			}).
			Post("/:id", func(c *river.Context) {
				c.Render(c.Param("id"), 201)
			}).
			Get("/", func(c *river.Context) {
				c.Render("It works", 200)
			}).
			Get("/:id/:name", func(c *river.Context) {
				time.Sleep(time.Second * 3)
				panic("whatever")
			}).
			Renderer(river.JSONRenderer),
		)
	}
	{
		rv.Use(river.Recovery(nil))
	}
	{
		e := river.NewEndpoint()
		e.Use(func(c *river.Context) {
			fmt.Println("Request endpoint for", c.URL.Path)
			c.Next()
		})
		rv.Handle("/school", e.
			Get("/:name", name).
			Put("/:name/:id", nameID),
		)
	}

	rv.Run(":8080")

}

func name(c *river.Context) {
	c.Render(river.M{"name": c.Param("name")}, 200)
}

func nameID(c *river.Context) {
	c.Render(river.M{
		"name": c.Param("name"),
		"id":   c.Param("id"),
	}, 200)
}
