package main

import (
	"fmt"
	"gin-play/cyn"
	"net/http"
)

func main() {
	engine := cyn.New()
	engine.GET("/index", func(context *cyn.Context) {
		context.HTML(200, "<h1>Index Page</h1>")
	})
	v1 := engine.Group("/v1")
	{
		v1.GET("/", func(context *cyn.Context) {
			context.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})
		v1.GET("/hello", func(c *cyn.Context) {
			// expect /hello?name=geektutu
			c.String(http.StatusOK, fmt.Sprintf("hello %s, you're at %s\n", c.Query("name"), c.Path))
		})
	}

	v2 := engine.Group("/v2")
	{
		v2.GET("/hello/:name", func(c *cyn.Context) {
			// expect /hello/geektutu
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *cyn.Context) {
			c.JSON(http.StatusOK, cyn.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})

	}

	engine.Run(":8848")
}
