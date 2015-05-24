package webapp

// Renderer interface used to render templates
type Renderer interface {
	Render(c *Context, view string, model interface{})
}

// Session
type Session interface {
	SessionId() string
	Get(string) interface{}
	Set(string, interface{}) Session
	GetAll() map[string]interface{}
}

// User
type UserManager interface {
	Login(username string, password string) (bool, error)
	Logout()
	Authenticated() bool
	UserId() string
	UserName() string
	Info() *User
}
