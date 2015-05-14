package webapp

// Handler handler is an interface that objects can implement to be registered to serve as middleware
// in the Negroni middleware stack.
// ServeHTTP should yield to the next middleware in the chain by invoking the next http.HandlerFunc
// passed in.
//
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.
type MiddlewareHandler interface {
	//ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	ServeHTTP(c *Context, next HandlerFunc)
}

// HandlerFunc is an adapter to allow the use of ordinary functions as Negroni handlers.
// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler object that calls f.
//type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
type MiddlewareHandlerFunc func(c *Context, next HandlerFunc)

func (h MiddlewareHandlerFunc) ServeHTTP(c *Context, next HandlerFunc) {
	h(c, next)
}

type middleware struct {
	handler MiddlewareHandler
	next    *middleware
}

func (m middleware) ServeHTTP(c *Context) {
	m.handler.ServeHTTP(c, m.next.ServeHTTP)
}
