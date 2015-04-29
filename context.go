// Web App Helper functions.
package webapp

import (
	"fmt"

	"github.com/gorilla/schema"
	"github.com/zdebeer99/mux"
)

var decoder *schema.Decoder = schema.NewDecoder()

type Context struct {
	*mux.HandlerContext
	renderEngine Renderer
}

func (this Context) MuxContext() *mux.HandlerContext {
	return this.HandlerContext
}

func (this *Context) RenderString(txt string) {
	fmt.Fprint(this.Response, txt)
}

func (this *Context) Render(view string, model interface{}) {
	this.renderEngine.Render(this.Response, view, model)
}

func (this *Context) BindForm(model interface{}) {
	err := this.Request.ParseForm()
	if err != nil {
		panic(err)
	}
	err = decoder.Decode(model, this.Request.PostForm)
	if err != nil {
		panic(err)
	}
}
