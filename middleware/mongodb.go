package middleware

import (
	"github.com/zdebeer99/webapp"
	"gopkg.in/mgo.v2"
)

type M map[string]interface{}

func NewMongoDb(path string) *mongoDBContext {
	mongodb := new(mongoDBContext)
	mongodb.ConnectDatabase(path)
	return mongodb
}

type mongoDBContext struct {
	MainSession *mgo.Session
}

func (this *mongoDBContext) ConnectDatabase(path string) {
	if path == "" {
		path = "localhost"
	}

	var err error
	this.MainSession, err = mgo.Dial(path)
	if err != nil {
		panic(err)
	}
}

func (this *mongoDBContext) Session() *mgo.Session {
	return this.MainSession.Copy()
}

func (this *mongoDBContext) ServeHTTP(c *webapp.Context, next webapp.HandlerFunc) {
	session := this.MainSession.Copy()
	defer session.Close()
	c.Set("MongoSession", session)
	c.Set(webapp.KeyDatabaseObject, session.DB(""))
	next(c)
}
