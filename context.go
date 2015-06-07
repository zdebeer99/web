// Web App Helper functions.
package webapp

import (
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/mgo.v2"

	"github.com/gorilla/schema"
	"github.com/zdebeer99/mux"
)

var decoder *schema.Decoder = schema.NewDecoder()

type Context struct {
	*mux.HandlerContext
	app       *Webapp
	register  map[string]interface{}
	Session   Session
	SessionId string
	User      UserManager
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

// Get all values that was set on this request context.
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

//Return a String to the client.
func (this *Context) ViewString(txt string) {
	fmt.Fprint(this.Response(), txt)
}

// View Render a template to html.
// By default gojade rendering engine is used, this can be customized.
func (this *Context) View(view string, model interface{}) {
	this.app.RenderEngine.Render(this, view, model)
}

func (this *Context) Redirect(path string) {
	fmt.Println("Redirected To", path)
	http.Redirect(this.Response(), this.Request(), path, http.StatusMovedPermanently)
}

// BindForms binds a go structure to a html form
// Uses gorilla.schema
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

func (this *Context) Form() url.Values {
	this.Request().ParseForm()
	return this.Request().Form
}

func (this *Context) PostForm() url.Values {
	this.Request().ParseForm()
	return this.Request().PostForm
}

// DB get a mgo.Database instance for a mongo database.
// This function can be modified to return your database instance.
// The MongoDB Middleware must be used for this function to work.
func (this *Context) DB() *mgo.Database {
	db := this.Get(KeyDatabaseObject)
	if db == nil {
		panic("Database connection was not establish. Use MongoDB Middleware to connect the initial connection.")
	}
	return db.(*mgo.Database)
}

func (this *Context) Auhtenticate() bool {
	return this.User.Authenticated()
}

func (this *Context) Error(errormessage string, code int) {

}
