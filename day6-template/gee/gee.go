package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
	router        *router
	groups        []*RouterGroup     // 存储所有分组
	htmlTemplates *template.Template // 提供 html 支持
	funcMap       template.FuncMap   // 提供 html 支持
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
	c.engine = engine
	engine.router.handle(c)

}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	// Strip:去除，去掉absolutePath前缀来访问，以免对用户暴露完整文件路径
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		// 得到解析出来的路径
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 添加处理静态文件的路由
func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 用于批量加载 HTML 模板文件
// pattern：文件路径模式。可以使用通配符（如 * 或 **）来匹配文件。
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
