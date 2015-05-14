// Web App Helper functions.
package webapp

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/schema"
	"github.com/zdebeer99/mux"
)

func (this *Webapp) ConnectDatabase(path string) {
	this.dataContext = &DataContext{}
	this.dataContext.ConnectDatabase(path)
}

var decoder *schema.Decoder = schema.NewDecoder()

type Context struct {
	*mux.HandlerContext
	app *Webapp
}

func NewContext(app *Webapp, w http.ResponseWriter, req *http.Request) *Context {
	c := &Context{}
	c.HandlerContext = mux.NewContext(NewResponseWriter(w), req)
	c.app = app
	return c
}

func (this *Context) ResponseWriter() ResponseWriter {
	return this.Response().(ResponseWriter)
}

func (this *Context) RenderString(txt string) {
	fmt.Fprint(this.Response(), txt)
}

func (this *Context) Render(view string, model interface{}) {
	this.app.RenderEngine.Render(this.Response(), view, model)
}

func (this *Context) BindForm(model interface{}) {
	err := this.Request().ParseForm()
	if err != nil {
		panic(err)
	}
	err = decoder.Decode(model, this.Request().PostForm)
	if err != nil {
		panic(err)
	}
}

func (this *Context) DB() *mgo.Database {
	return this.app.dataContext.Session().DB("")
}

func (this *Context) Error(errormessage string, code int) {

}
