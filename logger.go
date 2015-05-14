package webapp

import (
	"log"
	"net/http"
	"os"
	"time"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// NewLogger returns a new Logger instance
func NewLogger() *Logger {
	return &Logger{log.New(os.Stdout, "[webapp] ", 0)}
}

func (l *Logger) ServeHTTP(c *Context, next HandlerFunc) {
	start := time.Now()
	r := c.Request()
	l.Printf("Started %s %s", r.Method, r.URL.Path)

	next(c)

	res := c.ResponseWriter()
	l.Printf("Completed %v %s in %v", res.Status(), http.StatusText(res.Status()), time.Since(start))
}
