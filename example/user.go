package main

import (
	"net/http"

	"github.com/abiosoft/river"
)

// User is user data.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// getUser handles GET /user/:id.
func getUser(c *river.Context, model Model) {
	user := model.get(c.Param("id"))
	if user == nil {
		c.RenderEmpty(http.StatusNotFound)
		return
	}
	c.Render(http.StatusOK, user)
}

// getAllUser handles GET /user.
func getAllUser(c *river.Context, model Model) {
	c.Render(http.StatusOK, model.getAll())
}

// addUser handles POST /user.
func addUser(c *river.Context, model Model) {
	var users []User
	if err := c.DecodeJSONBody(&users); err != nil {
		c.Render(http.StatusBadRequest, err)
		return
	}
	for i := range users {
		model.add(users[i])
	}
	c.Render(http.StatusCreated, users)
}

// updateUser handles PUT /user/:id.
func updateUser(c *river.Context, model Model) {
	id := c.Param("id")
	var user User
	if err := c.DecodeJSONBody(&user); err != nil {
		c.Render(http.StatusBadRequest, err)
		return
	}
	model.put(id, user)
	c.Render(http.StatusOK, user)
}

// deleteUser handles DELETE /user/:id.
func deleteUser(c *river.Context, model Model) {
	model.delete(c.Param("id"))
	c.RenderEmpty(http.StatusNoContent)
}
