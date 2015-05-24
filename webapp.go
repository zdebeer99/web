package webapp

import (
	"log"
	"net/http"
	"os"

	"github.com/zdebeer99/mux"
)

const (
	KeySessionId      = "SessionId"
	KeyDatabaseObject = "DatabaseObject"
	KeyUser           = "UserName"
)

type Handler interface {
	ServeHTTP(c *Context)
}

type HandlerFunc func(c *Context)

func (h HandlerFunc) ServeHTTP(c *Context) {
	h(c)
}

// Wrap converts a http.Handler into a negroni.Handler so it can be used as a Negroni
// middleware. The next http.HandlerFunc is automatically called after the Handler
// is executed.
func Wrap(handler Handler) MiddlewareHandler {
	return MiddlewareHandlerFunc(func(c *Context, next HandlerFunc) {
		handler.ServeHTTP(c)
		next(c)
	})
}

//WebappHandler Wrap a mux handler and calls a webapp handler
func WebappHandler(f func(*Context)) func(interface{}) {
	return func(mx interface{}) {
		c := mx.(*Context)
		f(c)
	}
}

// Negroni is a stack of Middleware Handlers that can be invoked as an http.Handler.
// Negroni middleware is evaluated in the order that they are added to the stack using
// the Use and UseHandler methods.
type Webapp struct {
	middleware   middleware
	handlers     []MiddlewareHandler
	router       *RouterContext
	RenderEngine Renderer
}

// New returns a new Negroni instance with no middleware preconfigured.
func New(handlers ...MiddlewareHandler) *Webapp {
	web := &Webapp{
		handlers:   handlers,
		middleware: build(handlers),
	}
	web.router = NewRouter()
	web.router.SetContextFactory(web.contextFactory)
	web.RenderEngine = NewJadeRender("./views")
	return web
}

func (this *Webapp) HandleFunc(path string, f func(*Context)) *mux.Route {
	return this.router.HandleFunc(path, WebappHandler(f))
}

func (this *Webapp) Get(path string, f func(*Context)) *mux.Route {
	return this.router.HandleFunc(path, WebappHandler(f)).Methods("GET")
}

func (this *Webapp) Post(path string, f func(*Context)) *mux.Route {
	return this.router.HandleFunc(path, WebappHandler(f)).Methods("POST")
}

func (this *Webapp) FileServer(path_prefix string, file_path string) {
	this.router.FileServer(path_prefix, file_path)
}

// Classic returns a new Negroni instance with the default middleware already
// in the stack.
//
// Recovery - Panic Recovery Middleware
// Logger - Request/Response Logging
func Classic() *Webapp {
	return New(NewRecovery(), NewLogger())
}

func (n *Webapp) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := NewContext(n, rw, r)
	n.middleware.ServeHTTP(c)
}

// Use adds a Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Webapp) Use(handler MiddlewareHandler) {
	n.handlers = append(n.handlers, handler)
	n.middleware = build(n.handlers)
}

// UseFunc adds a Negroni-style handler function onto the middleware stack.
func (n *Webapp) UseFunc(handlerFunc func(c *Context, next HandlerFunc)) {
	n.Use(MiddlewareHandlerFunc(handlerFunc))
}

// UseHandler adds a http.Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Webapp) UseHandler(handler Handler) {
	n.Use(Wrap(handler))
}

// UseHandler adds a http.HandlerFunc-style handler function onto the middleware stack.
func (n *Webapp) UseHandlerFunc(handlerFunc func(c *Context)) {
	n.UseHandler(HandlerFunc(handlerFunc))
}

// Run is a convenience function that runs the negroni stack as an HTTP
// server. The addr string takes the same format as http.ListenAndServe.
func (n *Webapp) Run(addr string) {
	n.UseHandler(n.router)
	l := log.New(os.Stdout, "[webapp] ", 0)
	l.Printf("listening on %s", addr)
	l.Fatal(http.ListenAndServe(addr, n))
}

// Returns a list of all the handlers in the current Negroni middleware chain.
func (n *Webapp) Handlers() []MiddlewareHandler {
	return n.handlers
}

func build(handlers []MiddlewareHandler) middleware {
	var next middleware

	if len(handlers) == 0 {
		return voidMiddleware()
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = voidMiddleware()
	}

	return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
	return middleware{
		MiddlewareHandlerFunc(func(c *Context, next HandlerFunc) {}),
		&middleware{},
	}
}

func (this *Webapp) contextFactory(w http.ResponseWriter, req *http.Request) interface{} {
	c := NewContext(this, w, req)
	c.app = this
	return c
}

//RouterContext Wrap mux.Router to support ServeHTTP(*Context)
type RouterContext struct {
	*mux.Router
}

//NewRouter Create a new mux router adapter
func NewRouter() *RouterContext {
	return &RouterContext{mux.NewRouter()}
}

//Wrapped ServeHttp
func (this *RouterContext) ServeHTTP(c *Context) {
	this.Router.ServeHTTPContext(c)
}
