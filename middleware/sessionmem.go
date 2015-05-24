package middleware

import (
	"sync"

	"github.com/zdebeer99/webapp"
)

// BasicSession is a basic session store that stores session data in local memory
type BasicSession struct {
	sessionid string
	data      map[string]interface{}
}

// SessionId return the sessionid associated with this session.
func (this *BasicSession) SessionId() string {
	return this.sessionid
}

func (this *BasicSession) SetSessionId(id string) {
	this.sessionid = id
}

// Get get a session variable
func (this *BasicSession) Get(name string) interface{} {
	if len(this.data) == 0 {
		return nil
	}
	return this.data[name]
}

// GetAll Get all session variables.
func (this *BasicSession) GetAll() map[string]interface{} {
	return this.data
}

// Set a session variable
// if value is nil, the key will be removed from the Session.
func (this *BasicSession) Set(name string, value interface{}) webapp.Session {
	if len(this.data) == 0 {
		this.data = make(map[string]interface{})
	}
	if value == nil {
		delete(this.data, name)
	} else {
		this.data[name] = value
	}
	return this
}

func NewSession() *sessionMemoryMW {
	return new(sessionMemoryMW)
}

type sessionMemoryMW struct {
	mutex    sync.RWMutex
	sessions map[string]webapp.Session
}

func (this *sessionMemoryMW) ServeHTTP(c *webapp.Context, next webapp.HandlerFunc) {
	id := c.SessionId
	if id == "" {
		panic("SessionId Empty, make sure the SessionId Middleware is set before this middleware.")
	}
	s := this.get(id)
	if s != nil {
		c.Session = s
	} else {
		session := new(BasicSession)
		session.SetSessionId(id)
		c.Session = session
		this.set(id, session)
	}
	next(c)
}

func (this *sessionMemoryMW) get(id string) webapp.Session {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	return this.sessions[id]
}

func (this *sessionMemoryMW) set(id string, session webapp.Session) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if len(this.sessions) == 0 {
		this.sessions = make(map[string]webapp.Session)
	}
	this.sessions[id] = session
}
