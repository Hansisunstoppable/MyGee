package gee

import (
	"fmt"
	"net/http"
)

// 处理方法
// 第一个参数 http.ResponseWriter 是一个接口，负责向客户端发送 HTTP 响应。
// 第二个参数 *http.Request 是指向 http.Request 结构体的指针，它代表从客户端发来的 HTTP 请求。
type HandlerFunc func(http.ResponseWriter, *http.Request)

// router 是一张路由映射表，路由地址string 与 处理方法HandlerFunc对应
type Engine struct {
	router map[string]HandlerFunc
}

// 在go中，通常以New命名的函数用于创建和初始化某个类型的实例
// *Engine：函数类型，返回一个指向Engine类型的指针
func New() *Engine {
	// 返回一个指针（& 代表取指针）
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// 添加路由
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 启动一个http服务器，在addr端口上的监听
// 返回值是error类型，出错返回err
func (engine *Engine) Run(addr string) (err error) {
	// 在addr端口监听
	// 使用engine作为HTTP全球处理器
	// 在http.ListenAndServe 函数中，端口号需要以":<port>", 如 ":9999"
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	// 这是一种go典型的map查找方式
	// handler 为 engine.router[key] 对应的函数
	// ok 代表是否查找成功
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
