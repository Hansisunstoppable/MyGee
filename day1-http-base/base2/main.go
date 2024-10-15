package main

import (
	"fmt"
	"log"
	"net/http"
)

// 定义一个结构体并为其实现方法是go中的常用行为
type Engine struct{}

// 实现http包中的ServeHTTP接口
// (engine *Engine) 表示方法接收者，它告诉 Go 编译器该方法属于 *Engine 类型
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/hello":
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND: %S\n", req.URL)
	}
}

func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9999", engine))
}
