## day6 
Web 框架如何支持服务端渲染
#### 服务器处理静态文件的请求
- 通过定义类似 `/assets/*filepath` 的路由规则，服务器可以确定哪些 URL 是对静态资源的请求 （在这里即以 `/assets` 开头的）
  - 如用户请求的文件路径为 `/assets/js/geektutu.js` , 此时解析到的 `filepath` 为 `js/geektutu.js` , 这就是一个文件相对路径
- 通过以上操作得到相对路径，再与文件的默认路径前缀 `/usr/web` (仅在这个例子中)拼接，得到文件的绝对路径，再直接将这个绝对路径传入 `http.FileServer` 去访问这个静态文件
```go
// 添加处理静态文件的路由
func (group *RouterGroup) Static(relativePath string, root string) {
  // 传入相对路径 与 文件绝对路径前缀，构建静态文件的处理方法
	handler := group.createStaticHandler(relativePath, http.Dir(root))
  // 将自定义的静态文件前缀 与 /*filepath 拼接，用于后续捕获用户请求的文件路径（*代表通配符）
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

```
- 如 `r.Static("/assets", "./static")` 
- 此时访问一个静态`css` 文件,给出的路径为 `/assets/css/geektutu.css` 
- 解析得到 `filepath` 为 `css/geektutu.css`
- 与 `root` 拼接得到绝对地址 `/static/css/geektutu.css`
#### HTML 模板渲染 
- 存储在 `./templates/example.tmpl` 中，用于支持根据不同模板进行渲染
