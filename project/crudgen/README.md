# crudgen — Gin CRUD 代码生成器

输入资源名称，自动生成包含 CRUD 接口的 Go 项目代码（Gin + GORM + SQLite）。

## 使用方式

```bash
# 指定资源名生成项目
go run main.go -name user
go run main.go -name product -out ./my-api -module my-api

# 交互模式
go run main.go
```

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-name` | (必填) | 资源名称，如 user、product、order |
| `-out` | `./<name>-api` | 输出目录 |
| `-module` | `<name>-api` | Go module 名 |

## 生成的项目结构

```
<name>-api/
├── main.go              # 入口文件
├── go.mod               # Gin + GORM + SQLite 依赖
├── config/config.go     # 配置加载 (PORT, DB_PATH 环境变量)
├── database/database.go # GORM 连接 + AutoMigrate 自动建表
├── model/<name>.go      # 数据模型 + 请求结构体
├── handler/<name>.go    # 5 个 CRUD Handler
├── router/router.go     # RESTful 路由注册
└── README.md            # 使用文档
```

## 生成的 API

| 方法 | 路径 | 说明 |
|------|------|------|
| `POST` | `/api/v1/<name>s` | 创建 |
| `GET` | `/api/v1/<name>s` | 分页列表 |
| `GET` | `/api/v1/<name>s/:id` | 查询单个 |
| `PUT` | `/api/v1/<name>s/:id` | 更新 |
| `DELETE` | `/api/v1/<name>s/:id` | 删除 |

生成后 `go mod tidy && go run main.go` 即可运行，SQLite 自动建表。
