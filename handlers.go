package webapp

import (
	"fmt"

	"github.com/zdebeer99/mux"
)

func LoginHandler() mux.HandlerFunc {
	return func(mc mux.Context) {
		c := mc.MuxContext()
		c.Request.ParseForm()
		fmt.Println(c.Request.Form)
		username := c.Request.PostFormValue("UserName")
		password := c.Request.FormValue("Password")
	}
}
