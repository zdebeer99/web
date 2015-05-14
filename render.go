package webapp

import (
	"net/http"

	"github.com/zdebeer99/gojade"
)

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
