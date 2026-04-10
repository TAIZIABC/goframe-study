/**
 * 应用入口文件
 * -----------------------------------------------
 * 对标 GoFrame: main.go + internal/cmd/cmd.go
 *
 * GoFrame 入口：
 *   func main() {
 *       cmd.Main.Run(gctx.GetInitCtx())
 *   }
 *
 *   var Main = gcmd.Command{
 *       Func: func(ctx, parser) error {
 *           s := g.Server()
 *           s.Use(ghttp.MiddlewareHandlerResponse)  // 注册统一响应中间件
 *           s.Group("/api/v1", func(group) {
 *               group.Bind(controller.Todo)          // 绑定控制器路由
 *           })
 *           s.Run()                                  // 启动服务
 *       },
 *   }
 *
 * JS 版本用 Express 实现完全相同的启动流程
 * -----------------------------------------------
 */

const express = require('express');
const config = require('./internal/config/config');
const apiRouter = require('./internal/cmd/cmd');
const { errorHandler, notFoundHandler } = require('./internal/middleware/response');

const app = express();

// ─── 中间件注册 ──────────────────────────────────
// 对标 GoFrame: s.Use(ghttp.MiddlewareHandlerResponse)
app.use(express.json());

// ─── 路由组注册 ──────────────────────────────────
// 对标 GoFrame: s.Group("/api/v1", func(group) { group.Bind(controller.Todo) })
app.use('/api/v1', apiRouter);

// ─── 404 & 错误处理 ──────────────────────────────
app.use(notFoundHandler);
app.use(errorHandler);

// ─── 启动服务 ────────────────────────────────────
// 对标 GoFrame: s.Run()
const PORT = config.server.port;
app.listen(PORT, () => {
  console.log(`[SERVER] Todo API (JS) is running on :${PORT}`);
  console.log(`[SERVER] API 基础路径: http://localhost:${PORT}/api/v1`);
});
