package main

import (
	"gin-play/cyn"
)

func main() {
	engine := cyn.New()
	engine.GET("/", func(c *cyn.Context) {
		c.HTML(200, "<h1>This is the index</h1>")
	})
	engine.GET("/hello", func(c *cyn.Context) {
		c.String(200, "hello")
	})
	engine.POST("/message", func(c *cyn.Context) {
		c.JSON(200, cyn.H{
			"name": c.PostForm("name"),
			"age":  c.PostForm("age"),
		})
	})

	engine.POST("/postSomething", func(context *cyn.Context) {
		context.String(200, context.PostBody())
	})
	engine.Run(":8848")
}
