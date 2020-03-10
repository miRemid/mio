package mioqq

import (
	"net/http"
	"strings"
	"text/template"
)

// HandlerFunc defines the request handler used by mioqq
type HandlerFunc func(c *Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	*RouterGroup
	router        *router
	groups        []*RouterGroup
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

// New is the constructor of mioqq.Engine
// It will return a ptr struct of mioqq.Engine
func New() *Engine {
	engine := &Engine{
		router: newRouter(),
	}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.groups = []*RouterGroup{
		engine.RouterGroup,
	}
	return engine
}

// Default will return a engine set logger and recovery middleware
func Default() *Engine {
	engine := New()
	engine.Use(Logger())
	engine.Use(Recovery())
	return engine
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST defines the method to add GET request
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// Server defines the method to start a http server
func (engine *Engine) Server(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServerTLS defines the method to start a https server
func (engine *Engine) ServerTLS(addr, certFile, keyFile string) (err error) {
	return http.ListenAndServeTLS(addr, certFile, keyFile, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
}

// SetFuncMap is defined to set funcmap
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob is defined to set htmlTemplates pattern
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
