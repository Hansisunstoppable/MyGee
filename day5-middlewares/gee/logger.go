package gee

import (
	"log"
	"time"
)

// Logger 中间件，用于记录请求到响应的时间
// 中间件本质都是传入Context的方法

func Logger() HandlerFunc {
	return func(ctx *Context) {
		// 起始时间
		t := time.Now()
		// 处理所有用户请求与中间件
		// 等待用户自己定义的 Handler 处理结束后，再做一些额外的操作
		ctx.Next()
		// 记录响应时间
		log.Printf("[%d] %s in %v\n", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}
