package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 通过设计 上下文(Context)，对Request 与 Response 进行封装
// 提供各类返回类型的支持（如JSON、HTML）
// 相当于需要什么返回类型，直接调用对应方法即可，不需要重写一遍

// day3 添加了 Params，用于存储解析参数

type H map[string]interface{}

type Context struct {
	// 原始的对象
	Writer http.ResponseWriter
	Req    *http.Request
	// 请求信息
	Path   string
	Method string
	Params map[string]string
	// 响应信息
	StatusCode int
	// 中间件 与 路由匹配得到的Handler
	handlers []HandlerFunc
	index    int
	// engine 指针
	engine *Engine
}

// 返回解析到的参数
func (c *Context) Param(key string) string {
	value := c.Params[key]
	return value
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	// c.index 初始化为 -1
	c.index++
	// 从第index开始执行
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) PostForm(key string) string {
	return c.Req.PostFormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	// 设置 http.ResponseWriter 的header
	c.SetHeader("Content-Type", "text/plain")
	// 设置 Context 与 http.ResponseWriter 的Status
	c.Status(code)
	// Sprintf: 用于格式化的字符串，可以理解为将一组参数填入一个字符串模板
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	// []byte(html), 将 string 转换为 []byte，字节切片
	c.Writer.Write([]byte(html))
}

// 在某种错误条件下立即停止执行，并返回一个JSON响应给客户端
func (c *Context) Fail(code int, err string) {
	// 终止当前请求的处理流程, 将index直接设置为最大值，跳过了所有中间件
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

// 支持根据模板文件名选择模板进行渲染
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	if err := c.engine.htmlTemplates.ExecuteTemplate(c.Writer, name, data); err != nil {
		// 500 是 服务器内部错误（Internal Server Error）的状态码
		c.Fail(500, err.Error())
	}

}
