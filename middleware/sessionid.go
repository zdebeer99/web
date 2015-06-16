package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/zdebeer99/webapp"
)

// NewSessionId Creates a SessionId and keeps track of the session using cookies
// This Middleware only manages a sessionid and does manage session data.
func NewSessionId() *sessionIdContext {
	durr, err := time.ParseDuration("300h")
	if err != nil {
		panic(err)
	}
	return &sessionIdContext{GenerateId, durr}
}

type sessionIdContext struct {
	GenerateSessionId func() string
	SessionTimeOut    time.Duration
}

func (this *sessionIdContext) ServeHTTP(c *webapp.Context, next webapp.HandlerFunc) {
	id := c.SessionId
	if id == "" {
		id = this.GetSessionId(c)
		if id == "" {
			id = GenerateId()
		}
	}
	c.SessionId = id

	cookie := &http.Cookie{
		Name:     webapp.KeySessionId,
		Value:    c.SessionId,
		HttpOnly: true,
		MaxAge:   0,
		Expires:  time.Now().Add(this.SessionTimeOut),
	}
	http.SetCookie(c.Response(), cookie)

	next(c)
}

func (this *sessionIdContext) GetSessionId(c *webapp.Context) string {
	cookie, err := c.Request().Cookie(webapp.KeySessionId)
	if err != nil {
		return ""
	}
	return cookie.Value
}

func GenerateId() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
