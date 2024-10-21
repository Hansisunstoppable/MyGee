package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s/n/n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}

// 打印堆栈追踪（stack trace）
func trace(message string) string {
	var pcs [32]uintptr // 存储程序计数器 pc 的指针
	// 获取调用栈信息，写入到pcs切片中
	// 3 表示跳过的调用帧数量。通常，前面的调用栈帧会包含 runtime.Callers 的调用本身，因此要跳过这几层
	// 返回值 n：表示成功写入 pcs 数组的程序计数器的数量
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder                  // 声明一个高效字符串构建器
	str.WriteString(message + "\nTraceback") // 写入初始信息
	for _, pc := range pcs[:n] {             // 遍历程序计数器
		fn := runtime.FuncForPC(pc)                           // 根据当前的程序计数器 pc，获取与之对应的函数信息
		file, line := fn.FileLine(pc)                         // 获取文件名与行号
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line)) // 写入栈信息
	}
	return str.String()
}
