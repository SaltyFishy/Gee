package main

import (
	"Gee/gee"
	"fmt"
	"net/http"
)

func main() {
	r := gee.New()
	r.GET("/", func (c *gee.Context) {
		fmt.Fprintf(c.Writer, "Request URL is: %s\n", c.Path)
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/kv", func (c *gee.Context) {
			for k, v := range c.Req.Header {
				fmt.Fprintf(c.Writer, "[key,value] is: %v %v\n", k, v)
			}
		})
		v1.GET("/html", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})

	}

	v2 := r.Group("/v2")
	{
		v2.GET("/hello", func(c *gee.Context) {
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":9999")
}
