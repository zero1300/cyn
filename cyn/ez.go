package cyn

import (
	"log"
	"net/http"
	"strings"
)

type HandleFunc func(*Context)

type RouteGroup struct {
	prefix      string
	middlewares []HandleFunc
	parent      *RouteGroup // support nesting
	engine      *Engine     // All group share an  Engine instance
}

func (g *RouteGroup) Group(prefix string) *RouteGroup {
	engine := g.engine
	newEngine := &RouteGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
	}
	return newEngine
}

type Engine struct {
	*RouteGroup // Engine possess all features of RouteGroup
	router      *router
	groups      []*RouteGroup // store all groups
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandleFunc
	for _, group := range e.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := NewContext(w, r)
	c.handlers = middlewares
	e.router.handle(c)
}

// New constructor of ez.Engine
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouteGroup = &RouteGroup{engine: engine}
	engine.groups = []*RouteGroup{engine.RouteGroup}
	return engine
}

func (g *RouteGroup) addRoute(method string, comp string, handleFunction HandleFunc) {
	path := g.prefix + comp
	log.Printf("Route %4s - %s", method, path)
	g.engine.router.addRoute(method, path, handleFunction)
}

func (g *RouteGroup) POST(path string, handleFunction HandleFunc) {
	g.addRoute("POST", path, handleFunction)
}

func (g *RouteGroup) GET(path string, handleFunction HandleFunc) {
	g.addRoute("GET", path, handleFunction)
}

// Use is defined to add middleware to the group
func (g *RouteGroup) Use(middlewares ...HandleFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (e *Engine) Run(addr string) {
	err := http.ListenAndServe(addr, e)
	if err != nil {
		panic("error to Listen")
	}
}
