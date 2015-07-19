package webapp

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zdebeer99/gojade"
)

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	//	if val := header["Content-Type"]; len(val) == 0 {
	//		header["Content-Type"] = value
	//	}
	header["Content-Type"] = value
}

var htmlContentType = []string{"text/html; charset=utf-8"}

type JadeRenderer struct {
	jadeEngine *gojade.Engine
}

func (this *JadeRenderer) Render(c *Context, view string, model interface{}) {
	m := ViewModel{Model: model, User: c.User}
	m.Html = &Html{&m, c}
	writeContentType(c.ResponseWriter(), htmlContentType)
	this.jadeEngine.RenderFileW(c.Response(), view, m)
}

func NewJadeRender(viewpath string) *JadeRenderer {
	e := &JadeRenderer{gojade.New()}
	e.jadeEngine.ViewPath = viewpath
	return e
}

var jsonContentType = []string{"application/json; charset=utf-8"}

func WriteJson(w http.ResponseWriter, model interface{}) error {
	writeContentType(w, jsonContentType)
	return json.NewEncoder(w).Encode(model)
}

var plainContentType = []string{"text/plain; charset=utf-8"}

func WriteString(w http.ResponseWriter, format string, data ...interface{}) {
	writeContentType(w, plainContentType)

	if len(data) > 0 {
		fmt.Fprintf(w, format, data...)
	} else {
		io.WriteString(w, format)
	}
}
