package main

import (
	"net/http"
	"time"

	"github.com/abiosoft/river"
)

// Session is basic token based session.
type Session struct {
	Token   string
	Created time.Time
	Expires time.Time
}

// authMid is sample authentication middleware.
func authMid(c *river.Context) {
	token := c.Query("token")
	session := getSession(token)
	if !session.Valid() {
		c.Render(http.StatusUnauthorized, river.M{"error": "Unauthorized"})
		return
	}
	c.Register(session)
	c.Next()
}

// newAuthToken handles GET /auth.
func newAuthToken(c *river.Context) {
	session := newSession()
	c.Render(200, river.M{"token": session.Token, "expires": session.Expires})
}

// sessionInfo handles GET /session.
func sessionInfo(c *river.Context, session Session) {
	c.Render(http.StatusOK, session)
}

var sessions = map[string]Session{}

func getSession(token string) (session Session) {
	if s, ok := sessions[token]; ok {
		session = s
	}
	return
}

// Valid checks if the current session is valid.
func (s Session) Valid() bool {
	return time.Now().Before(s.Expires)
}

func newSession() Session {
	session := Session{
		Token:   randString(10),
		Created: time.Now(),
		Expires: time.Now().Add(time.Minute * 5),
	}
	sessions[session.Token] = session
	return session
}
