# ginspec — Gin 路由扫描 & OpenAPI 文档生成器

扫描 Gin 项目中注册的路由，解析 Handler 注释，自动生成 OpenAPI 3.0 文档，可选启动 Swagger UI 查看。

## 使用方式

```bash
# 扫描项目，输出 JSON
go run main.go -dir ./my-gin-api -out api.json

# 输出 YAML
go run main.go -dir ./my-gin-api -out api.yaml -format yaml

# 启动 Swagger UI 在线查看
go run main.go -dir ./my-gin-api -serve -port 8088

# 自定义标题和版本
go run main.go -dir ./my-gin-api -title "Product API" -version 2.0.0
```

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-dir` | `.` | Go 项目目录路径 |
| `-out` | - | 输出文件路径 |
| `-format` | `json` | 输出格式：json 或 yaml |
| `-title` | `API Documentation` | API 文档标题 |
| `-version` | `1.0.0` | API 版本号 |
| `-serve` | `false` | 启动 Swagger UI 服务 |
| `-port` | `8088` | Swagger UI 端口 |

## 项目结构

```
ginspec/
├── main.go                # CLI 入口
├── parser/parser.go       # Go AST 解析器（提取路由 + 注释）
├── generator/openapi.go   # OpenAPI 3.0 文档生成器
└── swagger/swagger.go     # 内嵌 Swagger UI HTTP 服务
```

## 功能特点

- **AST 解析**：用 `go/ast` 解析 Handler 函数的 Go Doc 注释
- **正则匹配**：识别 `r.GET("/path", handler)` 等 Gin 路由注册模式
- **路由组解析**：自动识别 `r.Group("/api/v1")` 前缀拼接
- **路径参数转换**：Gin 的 `:id` 自动转为 OpenAPI 的 `{id}`
- **Tag 推断**：从路径自动推断分组标签
- **Swagger UI**：内嵌页面，通过 CDN 加载无需额外依赖
