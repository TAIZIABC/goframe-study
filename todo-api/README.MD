# 📝 Todo List API

基于 **GoFrame** 框架的 Todo 待办事项 RESTful API，是 [GoFrame 学习之路](https://github.com/TAIZIABC/goframe-study) 的第一个实战项目。

## 🚀 快速开始

### 1. 环境准备

- Go 1.18+
- MySQL 5.7+
- GoFrame CLI（可选）：`go install github.com/gogf/gf/cmd/gf/v2@latest`

### 2. 初始化数据库

```bash
# 登录 MySQL 并执行初始化脚本
mysql -u root -p < manifest/sql/init.sql
```

或者手动登录 MySQL 后执行：

```sql
source /path/to/todo-api/manifest/sql/init.sql;
```

该脚本会自动完成：
- 创建 `todo_db` 数据库
- 创建 `todos` 表（含 id、title、completed、created_at、updated_at 字段）
- 插入 8 条示例数据

### 3. 修改数据库配置

编辑 `manifest/config/config.yaml`，修改数据库连接信息：

```yaml
database:
  default:
    link: "mysql:root:你的密码@tcp(127.0.0.1:3306)/todo_db"
```

### 4. 启动服务

```bash
# 安装依赖
go mod tidy

# 启动
go run main.go

# 或使用 GoFrame CLI（支持热重载）
gf run main.go
```

服务启动后访问：
- API 地址：http://localhost:8000/api/v1/todos
- Swagger 文档：http://localhost:8000/swagger

## 📖 API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST`   | `/api/v1/todos`     | 创建 Todo |
| `GET`    | `/api/v1/todos`     | 查询列表（支持分页） |
| `PUT`    | `/api/v1/todos/:id` | 更新 Todo |
| `DELETE` | `/api/v1/todos/:id` | 删除 Todo |

### 请求示例

```bash
# 创建 Todo
curl -X POST http://localhost:8000/api/v1/todos \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 GoFrame","completed":false}'

# 查询列表（分页）
curl "http://localhost:8000/api/v1/todos?page=1&size=10"

# 更新 Todo
curl -X PUT http://localhost:8000/api/v1/todos/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"学习 GoFrame（已完成）","completed":true}'

# 删除 Todo
curl -X DELETE http://localhost:8000/api/v1/todos/1
```

## 📁 项目结构

```
todo-api/
├── api/v1/                  # API 接口定义（请求/响应结构体）
│   ├── todo.go
│   └── model.go
├── internal/
│   ├── cmd/cmd.go           # 路由注册与服务启动
│   ├── controller/todo.go   # 控制器（处理请求）
│   └── service/todo.go      # 业务逻辑层（数据库操作）
├── manifest/
│   ├── config/config.yaml   # 应用配置（数据库、端口等）
│   └── sql/init.sql         # 数据库初始化脚本
├── main.go                  # 入口文件
├── go.mod
└── README.MD
```

## 🔗 相关链接

- [GoFrame 官方文档](https://goframe.org/)
- [GoFrame 学习之路](https://github.com/TAIZIABC/goframe-study)
