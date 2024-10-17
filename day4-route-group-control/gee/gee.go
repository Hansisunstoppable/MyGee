package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	// 指向对应的Engine实例
	engine *Engine
}

// Engine 相当于是最顶层的分组
type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 存储所有分组
}

func New() *Engine {
	// engine 是Engine指针类型，加&的原因就是取Engine的地址
	engine := &Engine{router: newRouter()}
	// 指向最顶层的分组
	engine.RouterGroup = &RouterGroup{engine: engine}
	// 将顶层分组加入切片
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 创建一个新分组
// 所有分组共享一个Engine实例
func (group *RouterGroup) NewGroup(prefix string) *RouterGroup {
	engine := group.engine
	// 创建一个子分组
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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
