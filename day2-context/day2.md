## day2
#### context实现对动态填入返回数据的支持
-  `context.go` 中的 `String` 方法
  ```go
  func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	// Sprintf: 用于格式化的字符串，可以理解为将一组参数填入一个字符串模板
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
```
#### router 与 context
- `router` 是对day1中的路由映射表的封装
- `context` 是对 `Request` 与 `Response` 的封装，提供各类返回类型的支持（如 `JSON` 、 `HTML` ） 
