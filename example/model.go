package main

import "math/rand"

// Model is a sample basic data model.
type Model struct {
	get    func(id string) interface{}
	getAll func() interface{}
	add    func(items ...interface{})
	put    func(id string, item interface{})
	delete func(id string)
}

func basicModel() Model {
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

func randString(l int) (str string) {
	const alphaNum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_1234567890"
	for i := 0; i < l; i++ {
		n := rand.Intn(len(alphaNum))
		str += alphaNum[n : n+1]
	}
	return
}
