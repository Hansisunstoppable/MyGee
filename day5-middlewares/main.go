package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	// 创建v1路由分组
	v1 := r.NewGroup("/v1")
	v1.Use(gee.Logger())
	v1.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK, "<h1>hello logger</h1>\n")
	})
	r.Run(":9999")

}
