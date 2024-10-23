# MyGee
记录对 Gin 源码的学习
### 实现功能
- 将 `Request` 与 `Response` 封装为 `context` 上下文
- 通过 trie 树，实现匹配某一类型而不是固定路由
  - 如 包含 : 或者 * 的路由规则
- 支持分组控制
  - 增加了路由分组，一切对路由的操作变成 `RouterGroup` 的方法
- 支持中间件
  - 中间件的本质是一组 HandlerFunc
- 支持 html 渲染模板
- 支持错误恢复