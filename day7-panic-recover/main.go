package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	// 使用 Logger 与 Recovery 中间件, 注意：添加中间件的顺序会影响执行顺序
	r.Use(gee.Logger(), gee.Recovery())
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "hello geektutu\n")
	})
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"geektutu"}
		// bash
		// $ curl "http://localhost:9999/panic"
		// {"message":"Internal Server Error"}
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")

}
