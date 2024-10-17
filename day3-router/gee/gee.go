package gee

import "net/http"

type HandlerFunc func(*Context)

// 将router代码独立封装
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

// 与day1相比，router相关代码实现了封装
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 该方法由 Go 标准库中的 http.ListenAndServe 自动调用
// 当你将 Engine 作为处理器传递给服务器时, 实现了这个接口 可以将Engine直接作为一组接口传入
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
