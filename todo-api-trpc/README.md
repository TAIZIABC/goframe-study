# Todo API — tRPC-Go 版本

> 本项目是 `todo-api`（GoFrame 版）的 **tRPC-Go 对照实现**，使用泛 HTTP 标准服务（无需 protobuf），方便对比学习两种 Go 框架的差异。

## 快速启动

```bash
cd todo-api-trpc
go run main.go        # 启动在 :8000，自动读取 trpc_go.yaml
```

## API 接口（与 GoFrame 版完全一致）

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/todos` | 创建 Todo |
| `GET` | `/api/v1/todos?page=1&size=10` | 分页查询列表 |
| `PUT` | `/api/v1/todos/{id}` | 更新 Todo |
| `DELETE` | `/api/v1/todos/{id}` | 删除 Todo |

统一响应格式：`{ "code": 0, "message": "", "data": { ... } }`

---

## 目录结构对比

```
GoFrame (todo-api/)                    tRPC-Go (todo-api-trpc/)
─────────────────────────────          ─────────────────────────────
main.go                         ←→    main.go                         # 入口文件
go.mod                          ←→    go.mod                          # 依赖管理
manifest/config/config.yaml     ←→    trpc_go.yaml                    # 框架配置文件
│                                      │
├── api/v1/                            ├── api/v1/
│   ├── model.go                ←→    │   ├── model.go                # 数据模型 + 请求/响应结构
│   └── todo.go                        │   └── validate.go             # 校验函数
│       (g.Meta 定义路由+校验)          │       (GoFrame 用 tag, tRPC 用函数)
│                                      │
├── internal/                          ├── internal/
│   ├── cmd/cmd.go              ←→    │   ├── cmd/cmd.go              # 路由注册
│   │   (s.Group + group.Bind)         │   │   (mux.Router + HandleFunc)
│   │                                  │   │
│   ├── controller/todo.go      ←→    │   ├── controller/todo.go      # 控制器层
│   │   (参数自动绑定到 struct)         │   │   (手动解析 JSON + URL params)
│   │                                  │   │
│   ├── service/todo.go         ←→    │   ├── service/todo.go         # 业务逻辑层
│   │   (GoFrame ORM)                  │   │   (database/sql 原生 SQL)
│   │                                  │   │
│   └── (consts/packed/...)            │   └── database/db.go          # 数据库连接管理
│                                      │       (GoFrame 自动管理, tRPC 需手动)
│                                      │
├── manifest/                          ├── manifest/
│   ├── config/config.yaml             │   └── sql/init.sql            # 数据库 SQL
│   └── sql/init.sql
```

---

## 核心概念对比（GoFrame vs tRPC-Go）

### 1. 配置文件

```yaml
# ─── GoFrame: manifest/config/config.yaml ──────
server:
  address: ":8000"
database:
  default:
    link: "mysql:root:xxx@tcp(...)/todo_db"
# GoFrame 读取 config.yaml，自动初始化服务器和数据库
```

```yaml
# ─── tRPC-Go: trpc_go.yaml ────────────────────
server:
  service:
    - name: trpc.todo.api.stdhttp
      protocol: http_no_protocol    # 泛 HTTP 标准服务
      port: 8000
# tRPC-Go 读取 trpc_go.yaml，只管理服务，数据库需手动
```

### 2. 入口文件

```go
// ─── GoFrame (main.go) ───────────────────────
import _ "github.com/gogf/gf/contrib/drivers/mysql/v2"  // 副作用导入
func main() {
    cmd.Main.Run(gctx.GetInitCtx())  // 一行启动，框架自动做一切
}

// ─── GoFrame (cmd.go) ────────────────────────
s := g.Server()
s.Use(ghttp.MiddlewareHandlerResponse)  // 注册统一响应中间件
s.Group("/api/v1", func(group *ghttp.RouterGroup) {
    group.Bind(controller.Todo)          // 反射自动绑定路由
})
s.Run()
```

```go
// ─── tRPC-Go (main.go) ──────────────────────
func main() {
    database.Init(cfg)                   // 手动初始化数据库

    s := trpc.NewServer()                // 读取 trpc_go.yaml
    router := cmd.NewRouter()            // 手动创建路由
    thttp.RegisterNoProtocolServiceMux(  // 注册泛 HTTP 服务
        s.Service("trpc.todo.api.stdhttp"),
        router,
    )
    s.Serve()                            // 启动服务
}
```

### 3. 路由定义

```go
// ─── GoFrame: 声明式路由（struct Meta tag）─────
type TodoCreateReq struct {
    g.Meta `path:"/todos" method:"post" tags:"Todo"`
    Title  string `v:"required|length:1,200" json:"title"`
}
// group.Bind(controller.Todo) 自动解析所有 Meta tag
```

```go
// ─── tRPC-Go: 命令式路由（gorilla/mux）────────
api := router.PathPrefix("/api/v1").Subrouter()
api.HandleFunc("/todos", Create).Methods("POST")
api.HandleFunc("/todos", List).Methods("GET")
api.HandleFunc("/todos/{id:[0-9]+}", Update).Methods("PUT")
api.HandleFunc("/todos/{id:[0-9]+}", Delete).Methods("DELETE")
// 手动注册每个路由和 HTTP 方法
```

### 4. Controller 层

```go
// ─── GoFrame: 框架自动完成参数解析+校验+响应包装 ──
func (c *cTodo) Create(ctx context.Context, req *v1.TodoCreateReq) (*v1.TodoCreateRes, error) {
    id, err := service.Todo().Create(ctx, req.Title, req.Completed)
    return &v1.TodoCreateRes{Id: id}, nil
    // 框架自动: 解析HTTP请求→填充req→校验→调用方法→包装响应JSON
}
```

```go
// ─── tRPC-Go: 手动完成参数解析+校验+响应包装 ────
func Create(w http.ResponseWriter, r *http.Request) error {
    var req v1.TodoCreateReq
    json.NewDecoder(r.Body).Decode(&req)   // 手动解析 JSON
    v1.ValidateCreateReq(&req)             // 手动校验
    id, _ := service.Todo().Create(...)    // 调用 service
    writeJSON(w, 0, "", TodoCreateRes{Id: id})  // 手动包装响应
    return nil
}
```

### 5. 数据库操作

```go
// ─── GoFrame: 内置 ORM，链式调用 ──────────────
// 框架根据 config.yaml 自动创建连接池
result, _ := g.Model("todos").Data(g.Map{"title": title}).Insert()
model.Page(page, size).OrderDesc("id").Scan(&list)
g.Model("todos").Where("id", id).Data(data).Update()
```

```go
// ─── tRPC-Go: 手动管理连接 + 原生 SQL ────────
// 需要手动创建连接池
db, _ := sql.Open("mysql", dsn)

db.ExecContext(ctx, "INSERT INTO `todos` (...) VALUES (?, ?)", ...)
db.QueryContext(ctx, "SELECT ... LIMIT ? OFFSET ?", size, offset)
db.ExecContext(ctx, "UPDATE `todos` SET ... WHERE `id` = ?", ...)
```

### 6. 响应包装

```go
// ─── GoFrame: 中间件自动包装 ─────────────────
s.Use(ghttp.MiddlewareHandlerResponse)
// controller 只需 return (res, nil)
// 中间件自动输出: {"code":0,"message":"","data":{...}}
```

```go
// ─── tRPC-Go: 手动包装 ─────────────────────
func writeJSON(w http.ResponseWriter, code int, msg string, data interface{}) {
    json.NewEncoder(w).Encode(Response{Code: code, Message: msg, Data: data})
}
// 每个 handler 都需要调用 writeJSON
```

---

## 三框架对比总结

| 方面 | GoFrame | tRPC-Go (泛HTTP) | Express (JS) |
|------|---------|-------------------|-------------|
| **定位** | 全功能 Web 框架 | 高性能 RPC 框架 + HTTP 扩展 | 轻量 Web 框架 |
| **配置文件** | config.yaml（自动加载） | trpc_go.yaml（自动加载） | 手动读取 |
| **路由定义** | 声明式 (struct Meta tag) | gorilla/mux 手动注册 | router.method 手动注册 |
| **参数绑定** | 自动反射绑定 | 手动解析 JSON/Query | 手动从 req 提取 |
| **参数校验** | v tag 自动校验 | 手动编写校验函数 | 手动编写校验函数 |
| **数据库** | 内置 ORM，链式调用 | 无内置，用 database/sql | 无内置，用 mysql2 |
| **响应包装** | 中间件自动包装 | 手动包装 | 手动包装 |
| **协议支持** | HTTP | tRPC/HTTP/HTTP2/gRPC | HTTP |
| **服务治理** | 基础支持 | 完整生态（注册/发现/熔断/限流） | 需第三方库 |
| **学习曲线** | 中等（约定多） | 较高（概念多） | 低（灵活自由） |
| **适用场景** | 中小型 Web 应用 | 微服务/大规模分布式 | 全栈/快速原型 |

---

## 数据库

三个版本共用同一个 MySQL 数据库 `todo_db`，表结构完全一致。
