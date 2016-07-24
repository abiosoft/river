package main

import (
	"net/http"

	"github.com/abiosoft/river"
)

func main() {
	rv := river.New(river.Logger()) //, river.Recovery())

	userEndpoint := river.NewEndpoint().
		Get("/:id", getUser).
		Post("/", addUser).
		Get("/", getAllUser).
		Put("/:id", updateUser).
		Delete("/:id", deleteUser)

	rv.Handle("/user", userEndpoint)
	rv.Register(newBasicModel())
	rv.Run(":8080")
}

func getUser(c *river.Context, model Model) {
	user := model.get(c.Param("id"))
	if user == nil {
		c.RenderEmpty(http.StatusNotFound)
		return
	}
	c.Render(http.StatusOK, user)
}

func getAllUser(c *river.Context, model Model) {
	c.Render(http.StatusOK, model.getAll())
}

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

func deleteUser(c *river.Context, model Model) {
	model.delete(c.Param("id"))
	c.RenderEmpty(http.StatusNoContent)
}

// Model is a sample basic data model.
type Model struct {
	get    func(id string) interface{}
	getAll func() interface{}
	add    func(items ...interface{})
	put    func(id string, item interface{})
	delete func(id string)
}

// User is user data.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func newBasicModel() Model {
	var model Model
	var users = []User{}
	search := func(id string) int {
		for i := range users {
			if users[i].ID == id {
				return i
			}
		}
		return -1
	}
	model.get = func(id string) interface{} {
		if i := search(id); i > -1 {
			return users[i]
		}
		return nil
	}
	model.getAll = func() interface{} {
		return users
	}
	model.add = func(items ...interface{}) {
		for i := range items {
			users = append(users, items[i].(User))
		}
	}
	model.put = func(id string, item interface{}) {
		if i := search(id); i > -1 {
			users[i] = item.(User)
		}
	}
	model.delete = func(id string) {
		if i := search(id); i > -1 {
			part := append(users[:i])
			if i < len(users)-1 {
				part = append(part, users[i+1:]...)
			}
			users = part
		}
	}
	return model
}
