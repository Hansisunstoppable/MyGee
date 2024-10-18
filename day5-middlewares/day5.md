## day5
实现 `Logger` 中间件，能够记录请求到响应花费的时间
#### c.Next() 方法
```go
func A(c *Context) {
    part1
    c.Next()
    part2
}
func B(c *Context) {
    part3
    c.Next()
    part4
}
```
- 假设我们应用了中间件 A 和 B，和路由映射的 `Handler` 。`c.handlers` 是这样的 `[A, B, Handler]` ，`c.index` 初始化为 `-1` 。调用 `c.Next()` ，接下来的流程是这样的：
  - c.index++，c.index 变为 0
  - 0 < 3，调用 c.handlers[0]，即 A
  - 执行 part1，调用 c.Next()
  - c.index++，c.index 变为 1
  - 1 < 3，调用 c.handlers[1]，即 B
  - 执行 part3，调用 c.Next()
  - c.index++，c.index 变为 2
  - 2 < 3，调用 c.handlers[2]，即Handler
  - Handler 调用完毕，返回到 B 中的 part4，执行 part4
  - part4 执行完毕，返回到 A 中的 part2，执行 part2
  - part2 执行完毕，结束。
- 最终的顺序是 `part1 -> part3 -> Handler -> part 4 -> part2` 

