package main

import (
	"fmt"

	"github.com/abiosoft/river"
)

func main() {
	rv := river.New()
	{
		rv.Handle("/user", river.NewEndpoint().
			Get("/:id", func(c *river.Context) {
				c.Render(c.Param("id"), 201)
			}).
			Get("/", func(c *river.Context) {
				c.Render("It works", 200)
			}).
			Renderer(river.PlainRenderer),
		)
	}
	{
		rv.Use(func(c *river.Context) {
			defer func() {
				if err := recover(); err != nil {
					c.Render(
						river.M{
							"status":  "error",
							"message": err,
						}, 500,
					)
				}
			}()
			c.Next()
		})
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
