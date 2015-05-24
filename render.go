package webapp

import "github.com/zdebeer99/gojade"

type JadeRenderer struct {
	jadeEngine *gojade.Engine
}

func (this *JadeRenderer) Render(c *Context, view string, model interface{}) {
	m := ViewModel{model, c.User}
	this.jadeEngine.RenderFileW(c.Response(), view, m)
}

func NewJadeRender(viewpath string) *JadeRenderer {
	e := &JadeRenderer{gojade.New()}
	e.jadeEngine.ViewPath = viewpath
	return e
}
