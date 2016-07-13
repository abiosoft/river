package main

import (
	"net/http"

	"github.com/abiosoft/river"
)

func main() {
	rv := river.New(river.Logger(), river.Recovery(nil))

	userEndpoint := river.NewEndpoint().
		Get("/:id", getUser).
		Post("/:id", addUser).
		Get("/", getAllUser).
		Put("/:id", updateUser).
		Delete("/:id", deleteUser)

	rv.Handle("/user", userEndpoint)

	rv.Run(":8080")
}

func getUser(c *river.Context) {
	user := userModel.get(c.Param("id"))
	if user == nil {
		c.RenderEmpty(http.StatusNotFound)
		return
	}
	c.Render(http.StatusOK, user)
}

func getAllUser(c *river.Context) {
	c.Render(http.StatusOK, userModel.getAll())
}

func addUser(c *river.Context) {
	id := c.FormValue("id")
	name := c.FormValue("name")
	if id == "" || name == "" {
		c.RenderEmpty(http.StatusBadRequest)
		return
	}
	user := User{ID: id, Name: name}
	userModel.add(user)
	c.Render(http.StatusCreated, user)
}

func updateUser(c *river.Context) {
	id := c.FormValue("id")
	name := c.FormValue("name")
	if id == "" || name == "" {
		c.RenderEmpty(http.StatusBadRequest)
		return
	}
	user := User{ID: id, Name: name}
	userModel.put(id, user)
	c.Render(http.StatusOK, user)
}

func deleteUser(c *river.Context) {
	userModel.delete(c.Param("id"))
	c.RenderEmpty(http.StatusNoContent)
}

/*
Sample basic data model
*/
type model struct {
	get    func(id string) interface{}
	getAll func() interface{}
	add    func(item interface{})
	put    func(id string, item interface{})
	delete func(id string)
}

var userModel model

// User is user data.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var users = []User{}

func init() {
	search := func(id string) int {
		for i := range users {
			if users[i].ID == id {
				return i
			}
		}
		return -1
	}
	userModel.get = func(id string) interface{} {
		if i := search(id); i > -1 {
			return users[i]
		}
		return nil
	}
	userModel.getAll = func() interface{} {
		return users
	}
	userModel.add = func(item interface{}) {
		users = append(users, item.(User))
	}
	userModel.put = func(id string, item interface{}) {
		if i := search(id); i > -1 {
			users[i] = item.(User)
		}
	}
	userModel.delete = func(id string) {
		if i := search(id); i > -1 {
			part := append(users[:i])
			if i < len(users)-1 {
				part = append(part, users[i+1:]...)
			}
			users = part
		}
	}
}
