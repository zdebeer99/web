package webapp

import (
	"fmt"
	"net/http"
)

func WebappHandler(f func(*Context)) func(interface{}) {
	return func(mx interface{}) {
		c := mx.(*Context)
		f(c)
	}
}

func LoginHandler() func(*Context) {
	return func(c *Context) {
		req := c.Request()
		req.ParseForm()
		fmt.Println(req.Form)
		username := req.PostFormValue("UserName")
		password := req.PostFormValue("Password")
		cookie := &http.Cookie{
			Name:  "login",
			Value: username + ":" + password,
			Path:  "/",
		}
		http.SetCookie(c.Response(), cookie)
	}
}
