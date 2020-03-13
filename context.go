package mio

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// H don't know
type H map[string]interface{}

// Context is the mio-qq's context
type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int

	// engine
	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

// Next defines the middleware chains next method
func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// PostForm return the key value of post form data
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// Query return the url's key value
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// Param return the get param
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// Status set the response's status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// SetHeader set the response's header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// JSON reply the client using json data
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// String reply String data
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// Data reply []byte data
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// Fail is defined to send a fail json data
func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// HTML reply html string data
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		c.Fail(500, err.Error())
	}
}
