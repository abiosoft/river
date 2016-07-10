package main

import "github.com/abiosoft/river"

func main() {
	rv := river.New(river.Logger(), river.Recovery(nil))
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
				panic("whatever")
			}),
		)
	}
	{
		rv.Handle("/school", river.NewEndpoint().
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
