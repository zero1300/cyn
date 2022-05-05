package main

import (
	"net/http"
)

func main() {
	engine := New()
	engine.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte("hello"))
		if err != nil {
			return
		}
	})
	engine.Run(":8848")
}
