package main

import (
	"fmt"
	"ruoFramework/ruo"
)

func main()  {
	r := ruo.New()

	r.GET("/", func(c *ruo.Context) {
		fmt.Fprintf(c.Writer, "URL.Path = %q\n", c.Req.URL.Path)
	})

	r.POST("/hello", func(c *ruo.Context) {
		for k, v := range c.Req.Header {
			fmt.Fprintf(c.Writer, "Header[%q] = %q\n", k, v)
		}
	})
	r.Run(":9999")
}


