package main

import "github.com/abiosoft/river"

func main() {
	rv := river.New()
	{
		rv.Handle("/user", river.NewEndpoint().
			Get("/:id", func(c *river.Context) {
				c.Render(c.Param("id"), 201)
			}).
			Get("/", func(c *river.Context) {
				c.Render("You are not serious", 201)
			}).
			Renderer(river.PlainRenderer),
		)
	}

	rv.Run(":8080")

}
