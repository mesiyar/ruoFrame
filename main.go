package main

import (
	"log"
	"net/http"
	"ruo"
	"time"
)


func onlyForV2() ruo.HandlerFunc {
	return func(c *ruo.Context) {
		// Start timer
		t := time.Now()
		// Calculate resolution time
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := ruo.New()
	r.Use(ruo.Logger()) // global midlleware
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

	v1 := r.Group("/v1")
	{
		v1.GET("/abc", func(c *ruo.Context) {
			c.HTML(http.StatusOK, "<h1> adbcdef </h1>")
		})
	}


	r.GET("/", func(c *ruo.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	v2.Use(onlyForV2()) // v2 group middleware
	{
		v2.GET("/hello/:name", func(c *ruo.Context) {
			// expect /hello/eddie
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	r.Run(":9999")
}
