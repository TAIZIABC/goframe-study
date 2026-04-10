# Todo API — JavaScript 版本

> 本项目是 `todo-api`（GoFrame 版）的 **JavaScript/Express 对照实现**，方便对比学习两种语言和框架的差异。

## 快速启动

```bash
cd todo-api-js
npm install
node main.js    # 启动在 :8000
```

## API 接口（与 Go 版完全一致）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/todos` | 创建 Todo |
| `GET` | `/api/v1/todos?page=1&size=10` | 分页查询列表 |
| `PUT` | `/api/v1/todos/:id` | 更新 Todo |
| `DELETE` | `/api/v1/todos/:id` | 删除 Todo |

统一响应格式：`{ "code": 0, "message": "", "data": { ... } }`

---

## 目录结构对比

```
GoFrame (todo-api/)                    JavaScript (todo-api-js/)
─────────────────────────────          ─────────────────────────────
main.go                         ←→    main.js                         # 入口文件
go.mod                          ←→    package.json                    # 依赖管理
│                                      │
├── api/v1/                            ├── api/v1/
│   ├── model.go                ←→    │   └── todo.js                 # 数据模型 + 校验规则
│   └── todo.go                        │     (Go 用 struct tag 校验)
│       (g.Meta 定义路由+校验)          │     (JS 用函数手动校验)
│                                      │
├── internal/                          ├── internal/
│   ├── cmd/cmd.go              ←→    │   ├── cmd/cmd.js              # 路由注册
│   │   (s.Group + group.Bind)         │   │   (express.Router)
│   │                                  │   │
│   ├── controller/todo.go      ←→    │   ├── controller/todo.js      # 控制器层
│   │   (参数自动绑定到 struct)         │   │   (手动从 req 提取参数)
│   │                                  │   │
│   ├── service/todo.go         ←→    │   ├── service/todo.js         # 业务逻辑层
│   │   (GoFrame ORM)                  │   │   (mysql2 原生 SQL)
│   │                                  │   │
│   └── consts/                        │   ├── config/config.js        # 配置加载
│       packed/                        │   ├── database/db.js          # 数据库连接
│                                      │   └── middleware/response.js  # 统一响应中间件
│                                      │
├── manifest/                          ├── manifest/
│   ├── config/config.yaml      ←→    │   ├── config/config.yaml      # 配置文件
│   └── sql/init.sql            ←→    │   └── sql/init.sql            # 数据库 SQL
```

---

## 核心概念对比

### 1. 入口文件

```go
// ─── Go (main.go) ─────────────────────
func main() {
    cmd.Main.Run(gctx.GetInitCtx())
}
// GoFrame 通过 Command 模式启动，自动加载配置
```

```javascript
// ─── JS (main.js) ─────────────────────
const app = express();
app.use(express.json());
app.use('/api/v1', apiRouter);
app.listen(8000);
// Express 手动组装中间件和路由
```

### 2. 路由定义

```go
// ─── Go: 通过 struct Meta tag 声明式定义 ───
type TodoCreateReq struct {
    g.Meta `path:"/todos" method:"post" tags:"Todo"`
    Title  string `v:"required|length:1,200" json:"title"`
}
// 路由、方法、校验规则全部在 struct tag 中，框架自动解析
```

```javascript
// ─── JS: 命令式手动注册路由 ───────────────
router.post('/todos', todoController.create);
router.get('/todos', todoController.list);
router.put('/todos/:id', todoController.update);
router.delete('/todos/:id', todoController.delete);
// 路由和处理函数需要手动绑定
```

### 3. 参数校验

```go
// ─── Go: struct tag 自动校验 ──────────────
Title string `v:"required|length:1,200" json:"title"`
Page  int    `d:"1" v:"min:1" json:"page"`
// GoFrame 框架自动完成校验，开发者只需写 tag
```

```javascript
// ─── JS: 手动编写校验函数 ─────────────────
if (!body.title || body.title.length < 1 || body.title.length > 200) {
  errors.push('title 长度需在 1-200 之间');
}
// 需要自己写校验逻辑（或用 Joi/Zod 等库）
```

### 4. 数据库操作 (ORM vs 原生 SQL)

```go
// ─── Go: GoFrame ORM 链式调用 ─────────────
// 创建
result, _ := g.Model("todos").Data(g.Map{
    "title": title, "completed": completed,
}).Insert()
id, _ := result.LastInsertId()

// 分页查询
total, _ := model.Count()
model.Page(page, size).OrderDesc("id").Scan(&list)

// 更新
g.Model("todos").Where("id", id).Data(data).Update()
```

```javascript
// ─── JS: mysql2 原生 SQL ──────────────────
// 创建
const result = await execute(
  'INSERT INTO `todos` (`title`, `completed`) VALUES (?, ?)',
  [title, completed ? 1 : 0]
);
const id = result.insertId;

// 分页查询
const countResult = await query('SELECT COUNT(*) AS total FROM `todos`');
const list = await query(
  'SELECT * FROM `todos` ORDER BY `id` DESC LIMIT ? OFFSET ?',
  [size, offset]
);

// 更新
await execute('UPDATE `todos` SET `title` = ? WHERE `id` = ?', [title, id]);
```

### 5. 统一响应中间件

```go
// ─── Go: 框架内置中间件自动包装 ──────────
s.Use(ghttp.MiddlewareHandlerResponse)
// 控制器只需返回 (res, err)，中间件自动包装为：
// { "code": 0, "message": "", "data": { ... } }
```

```javascript
// ─── JS: 在 controller 中手动包装 ────────
res.json({ code: 0, message: '', data: { id } });
// 需要在每个 handler 中手动构造统一格式
```

### 6. 部分更新 (*bool 指针 vs undefined)

```go
// ─── Go: 用指针类型区分 "未传值" 和 "传了false" ─
Completed *bool `json:"completed"`
// nil  → 未传值，不更新
// false → 显式传了 false，更新为 false
```

```javascript
// ─── JS: 用 undefined 和 'in' 操作符 ──────────
const completed = 'completed' in req.body ? req.body.completed : undefined;
// undefined → 未传值，不更新
// false     → 显式传了 false，更新为 false
```

---

## 关键差异总结

| 方面 | GoFrame (Go) | Express (JavaScript) |
|------|-------------|---------------------|
| **类型系统** | 静态类型，编译期检查 | 动态类型，运行时检查 |
| **路由定义** | 声明式 (struct Meta tag) | 命令式 (router.method) |
| **参数校验** | 框架自动 (v tag) | 手动编写或第三方库 |
| **数据库** | 内置 ORM，链式调用 | 原生 SQL 或第三方 ORM |
| **响应包装** | 中间件自动包装 | 手动在 controller 中包装 |
| **依赖注入** | 框架全局对象 g.Server()、g.DB() | require/import 模块 |
| **错误处理** | 多返回值 (result, err) | try/catch + async/await |
| **单例模式** | var instance = &struct{} | module.exports（天然单例） |
| **部分更新** | *bool 指针类型 | undefined + 'in' 操作符 |
| **配置管理** | 框架自动读取 YAML | 手动读取或环境变量 |

---

## 数据库

两个版本共用同一个 MySQL 数据库 `todo_db`，表结构完全一致。
