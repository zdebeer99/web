package middleware

import (
	"github.com/zdebeer99/webapp"
	"gopkg.in/mgo.v2"
)

func NewMongoDb(path string) *MongoDBContext {
	mongodb := new(MongoDBContext)
	mongodb.ConnectDatabase(path)
	return mongodb
}

type MongoDBContext struct {
	MainSession *mgo.Session
}

func (this *MongoDBContext) ConnectDatabase(path string) {
	if path == "" {
		path = "localhost"
	}

	var err error
	this.MainSession, err = mgo.Dial(path)
	if err != nil {
		panic(err)
	}
}

func (this *MongoDBContext) Session() *mgo.Session {
	return this.MainSession.Copy()
}

func (this *MongoDBContext) ServeHTTP(c *webapp.Context, next webapp.HandlerFunc) {
	session := this.MainSession.Copy()
	defer session.Close()
	c.Set("MongoSession", session)
	c.Set("MongoDB", session.DB(""))
	next(c)
}
