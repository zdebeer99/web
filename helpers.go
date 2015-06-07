package webapp

type MuxHandlerAdapter struct {
	handler Handler
}

func NewMuxHandlerAdapter(handler Handler) *MuxHandlerAdapter {
	return &MuxHandlerAdapter{handler}
}

func (this *MuxHandlerAdapter) ServeHTTP(c interface{}) {
	this.handler.ServeHTTP(c.(*Context))
}
