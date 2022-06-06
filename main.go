package main

import (
	"fmt"
	"gin-play/cyn"
	"log"
	"net/http"
	"time"
)

func onlyForV2() cyn.HandleFunc {
	return func(c *cyn.Context) {
		// Start timer
		t := time.Now()
		// if a server error occurred
		//c.Fail(500, "Internal Server Error")
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := cyn.New()
	r.Use(cyn.Logger()) // global midlleware
	r.GET("/", func(c *cyn.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *cyn.Context) {
			c.String(http.StatusOK, fmt.Sprintf("hello %s, you're at %s\n", c.Param("name"), c.Path))
		})
	}

	r.Run(":8848")
}
