package cyn

import (
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandleFunc
}

func newRouter() *router {
	return &router{make(map[string]HandleFunc)}
}

func (r *router) addRoute(method string, path string, handler HandleFunc) {
	log.Printf("Route: %4s  -  %s", method, path)
	key := method + "-" + path
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND")
	}
}
