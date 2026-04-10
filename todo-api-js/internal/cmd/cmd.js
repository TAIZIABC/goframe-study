/**
 * 路由注册模块
 * -----------------------------------------------
 * 对标 GoFrame: internal/cmd/cmd.go
 *
 * GoFrame 通过 Meta tag 自动路由绑定：
 *   s.Group("/api/v1", func(group *ghttp.RouterGroup) {
 *       group.Bind(controller.Todo)  // 自动解析 g.Meta 中的 path 和 method
 *   })
 *
 * JS/Express 需要手动注册每个路由
 * -----------------------------------------------
 */

const express = require('express');
const todoController = require('../controller/todo');

const router = express.Router();

/**
 * 路由定义 — 对标 GoFrame api/v1/todo.go 中的 g.Meta 标签
 *
 * Go:  g.Meta `path:"/todos"    method:"post"`    → POST   /api/v1/todos
 * Go:  g.Meta `path:"/todos"    method:"get"`     → GET    /api/v1/todos
 * Go:  g.Meta `path:"/todos/:id" method:"put"`    → PUT    /api/v1/todos/:id
 * Go:  g.Meta `path:"/todos/:id" method:"delete"` → DELETE /api/v1/todos/:id
 */
router.post('/todos', todoController.create);
router.get('/todos', todoController.list);
router.put('/todos/:id', todoController.update);
router.delete('/todos/:id', todoController.delete);

module.exports = router;
