package gee

import (
	"log"
	"net/http"
	"strings"
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

// 添加GET的路由
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// 添加POST的路由
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 添加中间件,中间件是对于路由分组来来执行的，因此是 路由分组 的方法
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 加入中间件后，收到一个请求后，需要判断该请求适用哪些中间件
// 通过URL前缀(即路由分组)进行判断，得到中间件列表后，赋值给 c.mid_handlers
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var mid_handlers []HandlerFunc
	// 遍历所有分组
	for _, group := range engine.groups {
		// func HasPrefix(s, prefix string) bool:返回s是否以prefix为前缀
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			mid_handlers = append(mid_handlers, group.middlewares...)
		}
	}
	c := newContext(w, req)
	// 将需要执行的中间件加入context
	c.handlers = mid_handlers
	engine.router.handle(c)
}
