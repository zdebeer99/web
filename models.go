package webapp

type User struct {
	UserId   string `mgo:"_id"`
	UserName string
	Password string
}

type ViewModel struct {
	Model interface{}
	User  UserManager
}
