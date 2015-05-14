// Web App Helper functions.
package webapp

import (
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/schema"
	"github.com/zdebeer99/mux"
)

var decoder *schema.Decoder = schema.NewDecoder()

type Context struct {
	*mux.HandlerContext
	app      *Webapp
	register map[string]interface{}
}

func NewContext(app *Webapp, w http.ResponseWriter, req *http.Request) *Context {
	c := &Context{}
	c.HandlerContext = mux.NewContext(NewResponseWriter(w), req)
	c.app = app
	return c
}

func (this *Context) Http() (ResponseWriter, *http.Request) {
	return this.ResponseWriter(), this.Request()
}

func (this *Context) ResponseWriter() ResponseWriter {
	return this.Response().(ResponseWriter)
}

// Get a value that was set on this request context.
func (this *Context) Get(name string) interface{} {
	if len(this.register) == 0 {
		return nil
	}
	return this.register[name]
}

// Get a value that was set on this request context.
func (this *Context) GetAll() map[string]interface{} {
	return this.register
}

// Set a value on this request context.
func (this *Context) Set(name string, value interface{}) {
	if len(this.register) == 0 {
		this.register = make(map[string]interface{})
	}
	this.register[name] = value
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
	db := this.Get("MongoDB")
	if db == nil {
		panic("Database connection was not establish. Use MongoDB Middleware to connect the initial connection.")
	}
	return db.(*mgo.Database)
}

func (this *Context) Error(errormessage string, code int) {

}
