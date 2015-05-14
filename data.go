package webapp

import "gopkg.in/mgo.v2"

type DataContext struct {
	MainSession *mgo.Session
}

func (this *DataContext) ConnectDatabase(path string) {
	if path == "" {
		path = "localhost"
	}

	var err error
	this.MainSession, err = mgo.Dial(path)
	if err != nil {
		panic(err)
	}
}

func (this *DataContext) Session() *mgo.Session {
	return this.MainSession.Copy()
}
