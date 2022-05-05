package cyn

import (
	"net/http"
)

type HandleFunc func(*Context)

type Engine struct {
	router *router
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)
	e.router.handle(c)
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) addRoute(method string, path string, handleFunction HandleFunc) {
	e.router.addRoute(method, path, handleFunction)
}

func (e *Engine) POST(path string, handleFunction HandleFunc) {
	e.addRoute("POST", path, handleFunction)
}

func (e *Engine) GET(path string, handleFunction HandleFunc) {
	e.addRoute("GET", path, handleFunction)
}

func (e *Engine) Run(addr string) {
	err := http.ListenAndServe(addr, e)
	if err != nil {
		panic("error to Listen")
	}
}
