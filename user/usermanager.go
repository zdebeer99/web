package app

import "github.com/zdebeer99/webapp"

func NewUserMW() *userMW {
	return new(userMW)
}

type userMW struct {
}

func (this *userMW) ServeHTTP(c *webapp.Context, next webapp.HandlerFunc) {
	manager := &userManager{context: c}
	c.User = manager
	next(c)
}

type userOptions struct {
	LoginUrl string
}

type userManager struct {
	context *webapp.Context
	user    *webapp.User
}

func (this *userManager) Login(username string, password string) (bool, error) {
	user := &webapp.User{username, username, password}
	c := this.validateContext()
	c.Session.Set(webapp.KeyUser, user)
	return true, nil
}

func (this *userManager) Logout() {
	c := this.validateContext()
	this.user = nil
	c.Session.Set(webapp.KeyUser, nil)
}

func (this *userManager) Authenticated() bool {
	user := this.Info()
	if user == nil {
		return false
	}
	return true
}

func (this *userManager) UserId() string {
	user := this.Info()
	if user == nil {
		return ""
	}
	return user.UserId
}

func (this *userManager) UserName() string {
	user := this.Info()
	if user == nil {
		return ""
	}
	return user.UserName
}

func (this *userManager) Info() *webapp.User {
	if this.user != nil {
		return this.user
	}
	c := this.validateContext()
	userobj := c.Session.Get(webapp.KeyUser)
	if userobj == nil {
		return nil
	}
	user, ok := userobj.(*webapp.User)
	if !ok {
		panic("The object in session variable 'User' is not of type webapp.User")
	}
	this.user = user
	return user
}

func (this *userManager) Detail() interface{} {
	return nil
}

func (this *userManager) validateContext() *webapp.Context {
	c := this.context
	if c.Session == nil {
		panic("User Middleware Requires Session Middleware.")
	}
	return c
}

func validateContextForUser(c *webapp.Context) bool {
	if c.Session == nil {
		panic("User Middleware Requires Session Middleware.")
	}
	return true
}
