package user

type UserStore interface {
	SaveUser()
	LoadUser(username string)
}

// User
type UserManager interface {
	Login(username string, password string) (bool, error)
	Authenticated() bool
	UserId() string
	User() *User
}
