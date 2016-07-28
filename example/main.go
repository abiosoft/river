package main

import "github.com/abiosoft/river"

func main() {
	rv := river.New()

	userEndpoint := river.NewEndpoint().
		Post("/", addUser).
		Get("/", getAllUser).
		Get("/:id", getUser).
		Put("/:id", updateUser).
		Delete("/:id", deleteUser)
	userEndpoint.Use(authMid)
	userEndpoint.Register(basicModel())
	rv.Handle("/user", userEndpoint)

	infoEndpoint := river.NewEndpoint().
		Get("/", sessionInfo)
	infoEndpoint.Use(authMid)
	rv.Handle("/session", infoEndpoint)

	authEndpoint := river.NewEndpoint().
		Get("/", newAuthToken)
	rv.Handle("/auth", authEndpoint)

	panicEndpoint := river.NewEndpoint().
		Get("/", func() { panic("Recovery middleware has handled this") })
	rv.Handle("/panic", panicEndpoint)

	rv.Use(river.Recovery())
	rv.Run(":8080")
}
