// Web App Helper functions.
package webapp

import (
	"net/http"

	"github.com/zdebeer99/gojade"
	"github.com/zdebeer99/mux"
)

type Engine struct {
	*mux.Router
	RenderEngine Renderer
}

func New() *Engine {
	e := &Engine{}
	e.Router = mux.NewRouter()
	e.SetContextFactory(e.contextFactory)
	e.RenderEngine = NewJadeRender("./views")
	return e
}

func (this *Engine) Run(addr string) {
	http.ListenAndServe(addr, this)
}

func (this *Engine) HandleFunc(path string, f func(*Context)) *mux.Route {
	return this.Router.HandleFunc(path, Handler(f))
}

func (this *Engine) contextFactory(w http.ResponseWriter, req *http.Request) interface{} {
	c := &Context{}
	c.HandlerContext = mux.NewContext(w, req)
	c.renderEngine = this.RenderEngine
	return c
}

func (this *Engine) handleAdapter(f func(interface{})) {

}

type Renderer interface {
	Render(w http.ResponseWriter, view string, model interface{})
}

type JadeRenderer struct {
	jadeEngine *gojade.Engine
}

func (this *JadeRenderer) Render(w http.ResponseWriter, view string, model interface{}) {
	this.jadeEngine.RenderFileW(w, view, model)
}

func NewJadeRender(viewpath string) *JadeRenderer {
	e := &JadeRenderer{gojade.New()}
	e.jadeEngine.ViewPath = viewpath
	return e
}
