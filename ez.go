package main

import (
	"fmt"
	"net/http"
)

type handleFunc func(writer http.ResponseWriter, request *http.Request)

type Engine struct {
	trees map[string]handleFunc
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path
	handleFunction, ok := e.trees[key]
	if ok {
		handleFunction(w, r)
	} else {
		err := fmt.Errorf("could not find the route: %s", r.URL)
		fmt.Println(err)
	}
}

func New() *Engine {
	return &Engine{trees: make(map[string]handleFunc)}
}

func (e *Engine) addRoute(method string, path string, handleFunction handleFunc) {
	key := method + "-" + path
	e.trees[key] = handleFunction
}

func (e *Engine) POST(path string, handleFunction handleFunc) {
	e.addRoute("POST", path, handleFunction)
}

func (e *Engine) GET(path string, handleFunction handleFunc) {
	e.addRoute("GET", path, handleFunction)
}

func (e *Engine) Run(addr string) {
	err := http.ListenAndServe(addr, e)
	if err != nil {
		panic("error to Listen")
	}
}
