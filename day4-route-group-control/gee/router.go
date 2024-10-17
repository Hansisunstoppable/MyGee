package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // 不同请求方式的根结点，如 roots['GET']、roots['GET']
	handlers map[string]HandlerFunc // 不同请求方式的HandleFunc，如 handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 将路由地址拆分
// 且只有第一个 * 会被保留，因为 * 后面的直接通配
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	// 第一个参数为 r.roots[method] 的值，第二个参数为 r.roots[method] 是否存在
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

// 返回匹配到的结点 与 :* 的解析结果
func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern) // 存放待匹配的路径，不是trie树上的路径
	params := make(map[string]string)    // 解析结果
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	// 在trie树上匹配
	n := root.search(searchParts, 0)

	// 匹配成功，进行参数匹配，得到解析结果
	// 如 /p/go/doc匹配到/p/:lang/doc，解析结果为：{lang: "go"}
	// 如 /static/css/geektutu.css匹配到/static/*filepath，解析结果为{filepath: "css/geektutu.css"}
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				// 若part为 :lang, part[1:]就是 lang
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				// searchParts[index:] 获取 searchParts切片从 index 到最后的所有元素
				// strings.Join(..., "/") 将这些元素用 / 连接起来，形成一个新的字符串
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

// 传入上下文到服务器，获取路由、解析参数，并向请求方返回响应
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		// 键值为context中的请求类型 与 node 中的完整路由路径
		// 注意：此时采用了trie树匹配路由，路径已经不存在context中，与day2有所不同
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
