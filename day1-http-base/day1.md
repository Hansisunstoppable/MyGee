## day1
#### go接口的隐式实现
1. 结构体 `Engine` 定义如下
   ```go
   type Engine struct {
        router map[string]func(http.ResponseWriter, *http.Request)
    }
2. 接口 `Handler` 的定义
   ```go
   type Handler interface {
        ServeHTTP(w http.ResponseWriter, r *http.Request)
    }
3. 结构体 `Engine` 是如何隐式的实现接口 `Handler` 的？
   - 直接为 `Engine` 添加一个 `ServeHTTP` 方法，便实现了该接口
   - go 中是定义一个方法，并在前面加上 `(engine *Engine)` , 代表该方法属于 `Engine` 
    ```go
    func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        // 实现处理 HTTP 请求的逻辑
        // 例如，使用 engine.router 来路由请求
        handler, ok := engine.router[r.URL.Path]
        if ok {
            handler(w, r) // 调用路由对应的处理函数
        } else {
            http.NotFound(w, r) // 如果没有匹配的路由，则返回 404
        }
    }
    ```
4. Go 中的接口是隐式实现的。`Engine` 类型实现了`Handler` 接口的所有方法，便自动实现了该接口。
   - 需要传入该接口实例的地方，直接传入该结构体即可，go 会实现自动转换
    ```go
    // 根据源码，该方法第二个参数需要传入type Handler interface
    func ListenAndServe(addr string, handler Handler) error {
        server := &Server{Addr: addr, Handler: handler}
        return server.ListenAndServe()
    }
    // 直接传入结构体Engine实例engine，go会自动将其转换为接口类型
    http.ListenAndServe(addr, engine)

 