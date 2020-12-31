package main

import (
	"net/http"
	"ruo"
)

func main()  {
	r := ruo.New()

	r.GET("/", func(c *ruo.Context) {
		c.HTML(http.StatusOK, "<h1>Hello ruo</h1>")
	})
	r.GET("/hello", func(c *ruo.Context) {
		// expect /hello?name=eddie
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	r.POST("/login", func(c *ruo.Context) {
		c.Json(http.StatusOK, ruo.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.GET("/hello/:name", func(c *ruo.Context) {
		// expect /hello/eddie
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	r.GET("/assets/*filepath", func(c *ruo.Context) {
		c.Json(http.StatusOK, ruo.H{"filepath": c.Param("filepath")})
	})

	r.GET("/header", func(c *ruo.Context) {
		h := c.GetHeader("token")
		c.String(http.StatusOK, "token is %s", h)
	})

	r.Run(":9999")
}


