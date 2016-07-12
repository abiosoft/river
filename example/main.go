package main

import "github.com/abiosoft/river"

func main() {
	rv := river.New(river.Logger(), river.Recovery(nil))
	{
		rv.Handle("/user", river.NewEndpoint().
			Get("/:id", func(c *river.Context) {
				c.Render(200, c.Param("id"))
			}).
			Post("/:id", func(c *river.Context) {
				c.Render(200, c.Param("id"))
			}).
			Get("/", func(c *river.Context) {
				c.Render(200, "It works")
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
	c.Render(200, river.M{"name": c.Param("name")})
}

func nameID(c *river.Context) {
	c.Render(200, river.M{
		"name": c.Param("name"),
		"id":   c.Param("id"),
	})
}
