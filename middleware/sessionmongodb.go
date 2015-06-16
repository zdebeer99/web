package middleware

import (
	"sync"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/zdebeer99/webapp"
)

// SessionMongoDB store session data in mangodb
type SessionMongoDB struct {
	id   string
	db   *mgo.Database
	data *BasicSession
}

// SessionId return the sessionid associated with this session.
func (this *SessionMongoDB) SessionId() string {
	return this.id
}

func (this *SessionMongoDB) SetSessionId(id string) {
	if this.data != nil {
		this.data.SetSessionId(id)
	}
	this.id = id
}

func (this *SessionMongoDB) setData(data *BasicSession) {
	data.SetSessionId(this.SessionId())
	this.data = data
}

// Get get a session variable
func (this *SessionMongoDB) Get(name string) interface{} {
	return this.data.Get(name)
}

// GetAll Get all session variables.
func (this *SessionMongoDB) GetAll() map[string]interface{} {
	return this.data.GetAll()
}

// Set a session variable
// if value is nil, the key will be removed from the Session.
func (this *SessionMongoDB) Set(name string, value interface{}) webapp.Session {
	this.data.Set(name, value)
	saveMongoSession(this.db, this)
	return this
}

func NewSessionMongoDB() *sessionMongoDBMW {
	return new(sessionMongoDBMW)
}

type sessionMongoDBMW struct {
	mutex    sync.RWMutex
	sessions map[string]webapp.Session
}

func (this *sessionMongoDBMW) ServeHTTP(c *webapp.Context, next webapp.HandlerFunc) {
	id := c.SessionId
	if len(id) == 0 {
		panic("SessionId Empty, make sure the SessionId Middleware is set before this middleware.")
	}
	msession := &SessionMongoDB{id: id, db: c.DB()}
	found := this.get(msession, id)
	if found {
		c.Session = msession
	} else {
		msession.setData(new(BasicSession))
		c.Session = msession
		this.set(id, msession)
	}
	next(c)
	msession.db = nil
}

func (this *sessionMongoDBMW) get(msession *SessionMongoDB, id string) bool {
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	session := this.sessions[id]
	if session == nil {
		err := loadMongoSession(msession.db, msession)
		if err != nil {
			return false
		}
		if len(msession.SessionId()) == 0 {
			panic("SessionId not set")
		}
		return true
	} else {
		msession.setData(session.(*BasicSession))
		return true
	}
}

func (this *sessionMongoDBMW) set(id string, session *SessionMongoDB) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	if len(this.sessions) == 0 {
		this.sessions = make(map[string]webapp.Session)
	}
	this.sessions[id] = session.data
}

func saveMongoSession(db *mgo.Database, session *SessionMongoDB) {
	if len(session.SessionId()) == 0 {
		panic("SessionId Not Set")
	}
	data := make(map[string]interface{})
	for k, v := range session.GetAll() {
		data[k] = v
	}
	data["_id"] = session.SessionId()
	db.C("sessions").UpsertId(session.SessionId(), data)
}

func loadMongoSession(db *mgo.Database, session *SessionMongoDB) error {
	if len(session.SessionId()) == 0 {
		panic("SessionId Not Set")
	}
	var data bson.M
	err := db.C("sessions").FindId(session.SessionId()).One(&data)
	if err != nil {
		return err
	}
	bsession := new(BasicSession)
	session.setData(bsession)
	for k, v := range data {
		bsession.Set(k, v)
	}
	return err
}
