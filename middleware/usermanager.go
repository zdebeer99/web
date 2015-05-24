package middleware

import (
	"fmt"

	"github.com/zdebeer99/webapp"
	"gopkg.in/mgo.v2"
)

// NewUserCustomStoreMW Allows a custom store constructor.
// The storefactory method must return a struct of type UserStore Interface
func NewUserCustomStoreMW(storefactory func(*webapp.Context) UserStore) *userMW {
	return &userMW{storefactory}
}

// NewUserMW Create a User Authentication Middleware with mongodb as default user store.
func NewUserMW() *userMW {
	return &userMW{func(ctx *webapp.Context) UserStore {
		return &DefaultUserStore{ctx, ctx.DB()}
	}}
}

// UserStore allows for customizing user stores
type UserStore interface {
	VerifyUser(username string) bool
	Authenticate(username, password string) bool
	GetUser(username string) *webapp.User
}

type userMW struct {
	storefactory func(*webapp.Context) UserStore
}

func (this *userMW) ServeHTTP(c *webapp.Context, next webapp.HandlerFunc) {
	manager := &userManager{context: c, store: this.storefactory(c)}
	c.User = manager
	next(c)
}

type userOptions struct {
	LoginUrl string
}

type userManager struct {
	context *webapp.Context
	user    *webapp.User
	store   UserStore
}

func (this *userManager) Login(username string, password string) (bool, error) {
	var auth bool
	if this.store != nil {
		auth = this.store.Authenticate(username, password)
	} else {
		auth = false
	}
	fmt.Println("LOGIN", auth)
	if auth {
		user := &webapp.User{username, username, password}
		c := this.validateContext()
		c.Session.Set(webapp.KeyUser, user)
		return true, nil
	}
	return false, nil
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

type DefaultUserStore struct {
	ctx *webapp.Context
	db  *mgo.Database
}

func (this *DefaultUserStore) VerifyUser(username string) bool {
	user := this.GetUser(username)
	if user == nil {
		return false
	}
	return true
}
func (this *DefaultUserStore) Authenticate(username, password string) bool {
	user := this.GetUser(username)
	fmt.Println("Auth User", username, user)
	if user == nil {
		return false
	}
	return user.UserName == username && user.Password == password
}
func (this *DefaultUserStore) GetUser(username string) *webapp.User {
	c := this.db.C("users")
	var user webapp.User
	err := c.FindId(username).One(&user)
	if err != nil {
		if err.Error() == "not found" {
			return nil
		}
		panic(err)
	}
	return &user
}
