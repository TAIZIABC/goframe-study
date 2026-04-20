package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

func main() {
	resource := flag.String("name", "", "资源名称 (如: user, product, order)")
	output := flag.String("out", "", "输出目录 (默认: ./<资源名>-api)")
	module := flag.String("module", "", "Go module 名 (默认: <资源名>-api)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: crudgen -name <资源名> [选项]\n\n")
		fmt.Fprintf(os.Stderr, "自动生成 Gin + GORM + SQLite CRUD 项目脚手架。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  crudgen -name user\n")
		fmt.Fprintf(os.Stderr, "  crudgen -name product -out ./my-api -module my-api\n")
	}
	flag.Parse()

	// 交互模式
	if *resource == "" {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("请输入资源名称 (如 user, product): ")
		line, _ := reader.ReadString('\n')
		*resource = strings.TrimSpace(line)
	}

	if *resource == "" {
		fmt.Fprintln(os.Stderr, "✗ 资源名称不能为空")
		os.Exit(1)
	}

	name := strings.ToLower(*resource)
	if *output == "" {
		*output = "./" + name + "-api"
	}
	if *module == "" {
		*module = name + "-api"
	}

	data := TplData{
		Module:    *module,
		Name:      name,
		NameTitle: toTitle(name),
		NamePlural: name + "s",
	}

	fmt.Printf("\n  生成 CRUD 项目: %s\n", data.NameTitle)
	fmt.Printf("  输出目录: %s\n", *output)
	fmt.Printf("  模块名: %s\n\n", *module)

	files := []struct {
		path    string
		content string
	}{
		{"go.mod", tplGoMod},
		{"main.go", tplMain},
		{"config/config.go", tplConfig},
		{"database/database.go", tplDatabase},
		{"model/" + name + ".go", tplModel},
		{"handler/" + name + ".go", tplHandler},
		{"router/router.go", tplRouter},
		{"README.md", tplReadme},
	}

	for _, f := range files {
		fullPath := filepath.Join(*output, f.path)

		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ 创建目录失败: %v\n", err)
			os.Exit(1)
		}

		out, err := os.Create(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ 创建文件失败 %s: %v\n", f.path, err)
			os.Exit(1)
		}

		tmpl, err := template.New(f.path).Parse(f.content)
		if err != nil {
			out.Close()
			fmt.Fprintf(os.Stderr, "  ✗ 模板解析失败 %s: %v\n", f.path, err)
			os.Exit(1)
		}

		if err := tmpl.Execute(out, data); err != nil {
			out.Close()
			fmt.Fprintf(os.Stderr, "  ✗ 模板渲染失败 %s: %v\n", f.path, err)
			os.Exit(1)
		}
		out.Close()

		fmt.Printf("  ✓ %s\n", f.path)
	}

	fmt.Println()
	fmt.Println("  ────────────────────────────────────────")
	fmt.Println("  项目生成完成！开始使用:")
	fmt.Println()
	fmt.Printf("    cd %s\n", *output)
	fmt.Println("    go mod tidy")
	fmt.Println("    go run main.go")
	fmt.Println()
	fmt.Printf("  API 地址: http://localhost:8080/api/v1/%s\n", data.NamePlural)
	fmt.Println("  ────────────────────────────────────────")
}

type TplData struct {
	Module     string
	Name       string // user
	NameTitle  string // User
	NamePlural string // users
}

func toTitle(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

// ─── 模板定义 ─────────────────────────────────────

var tplGoMod = `module {{.Module}}

go 1.23.0

require (
	github.com/gin-gonic/gin v1.10.0
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)
`

var tplMain = `package main

import (
	"fmt"
	"{{.Module}}/config"
	"{{.Module}}/database"
	"{{.Module}}/router"
)

func main() {
	cfg := config.Load()

	// 初始化数据库
	database.Init(cfg.DBPath)

	// 启动路由
	r := router.Setup()

	fmt.Printf("🚀 {{.NameTitle}} API 启动成功\n")
	fmt.Printf("   地址: http://localhost:%s\n", cfg.Port)
	fmt.Printf("   数据库: %s\n\n", cfg.DBPath)

	r.Run(":" + cfg.Port)
}
`

var tplConfig = `package config

import "os"

type Config struct {
	Port   string
	DBPath string
}

func Load() *Config {
	return &Config{
		Port:   getEnv("PORT", "8080"),
		DBPath: getEnv("DB_PATH", "data.db"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
`

var tplDatabase = `package database

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"{{.Module}}/model"
)

var DB *gorm.DB

func Init(dbPath string) {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "数据库连接失败: %v\n", err)
		os.Exit(1)
	}

	// 自动迁移
	if err := DB.AutoMigrate(&model.{{.NameTitle}}{}); err != nil {
		fmt.Fprintf(os.Stderr, "数据库迁移失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ 数据库初始化完成")
}
`

var tplModel = `package model

import "time"

type {{.NameTitle}} struct {
	ID        uint       ` + "`json:\"id\" gorm:\"primaryKey\"`" + `
	Name      string     ` + "`json:\"name\" gorm:\"not null\" binding:\"required\"`" + `
	CreatedAt time.Time  ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time  ` + "`json:\"updated_at\"`" + `
}

// Create{{.NameTitle}}Req 创建请求
type Create{{.NameTitle}}Req struct {
	Name string ` + "`json:\"name\" binding:\"required\"`" + `
}

// Update{{.NameTitle}}Req 更新请求
type Update{{.NameTitle}}Req struct {
	Name string ` + "`json:\"name\"`" + `
}
`

var tplHandler = `package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"{{.Module}}/database"
	"{{.Module}}/model"
)

// Create{{.NameTitle}} 创建
// POST /api/v1/{{.NamePlural}}
func Create{{.NameTitle}}(c *gin.Context) {
	var req model.Create{{.NameTitle}}Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := model.{{.NameTitle}}{Name: req.Name}
	if err := database.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": item})
}

// List{{.NameTitle}}s 列表查询（分页）
// GET /api/v1/{{.NamePlural}}?page=1&size=10
func List{{.NameTitle}}s(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	var items []model.{{.NameTitle}}
	var total int64

	database.DB.Model(&model.{{.NameTitle}}{}).Count(&total)
	database.DB.Order("id DESC").Offset((page - 1) * size).Limit(size).Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"total": total,
		"page":  page,
		"size":  size,
	})
}

// Get{{.NameTitle}} 查询单个
// GET /api/v1/{{.NamePlural}}/:id
func Get{{.NameTitle}}(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var item model.{{.NameTitle}}
	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// Update{{.NameTitle}} 更新
// PUT /api/v1/{{.NamePlural}}/:id
func Update{{.NameTitle}}(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var item model.{{.NameTitle}}
	if err := database.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到"})
		return
	}

	var req model.Update{{.NameTitle}}Req
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		item.Name = req.Name
	}

	if err := database.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": item})
}

// Delete{{.NameTitle}} 删除
// DELETE /api/v1/{{.NamePlural}}/:id
func Delete{{.NameTitle}}(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	if err := database.DB.Delete(&model.{{.NameTitle}}{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
`

var tplRouter = `package router

import (
	"github.com/gin-gonic/gin"

	"{{.Module}}/handler"
)

func Setup() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/{{.NamePlural}}", handler.Create{{.NameTitle}})
		v1.GET("/{{.NamePlural}}", handler.List{{.NameTitle}}s)
		v1.GET("/{{.NamePlural}}/:id", handler.Get{{.NameTitle}})
		v1.PUT("/{{.NamePlural}}/:id", handler.Update{{.NameTitle}})
		v1.DELETE("/{{.NamePlural}}/:id", handler.Delete{{.NameTitle}})
	}

	return r
}
`

var tplReadme = `# {{.NameTitle}} API

> 由 crudgen 自动生成的 Gin + GORM + SQLite CRUD 项目

## 快速启动

` + "```bash" + `
go mod tidy
go run main.go
` + "```" + `

服务启动在 http://localhost:8080

## API 接口

| 方法     | 路径                      | 说明         |
|----------|---------------------------|-------------|
| POST     | /api/v1/{{.NamePlural}}        | 创建 {{.NameTitle}} |
| GET      | /api/v1/{{.NamePlural}}        | 列表查询 (分页) |
| GET      | /api/v1/{{.NamePlural}}/:id    | 查询单个     |
| PUT      | /api/v1/{{.NamePlural}}/:id    | 更新         |
| DELETE   | /api/v1/{{.NamePlural}}/:id    | 删除         |

## 请求示例

` + "```bash" + `
# 创建
curl -X POST http://localhost:8080/api/v1/{{.NamePlural}} \
  -H "Content-Type: application/json" \
  -d '{"name": "示例{{.NameTitle}}"}'

# 列表
curl http://localhost:8080/api/v1/{{.NamePlural}}?page=1&size=10

# 查询
curl http://localhost:8080/api/v1/{{.NamePlural}}/1

# 更新
curl -X PUT http://localhost:8080/api/v1/{{.NamePlural}}/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "新名称"}'

# 删除
curl -X DELETE http://localhost:8080/api/v1/{{.NamePlural}}/1
` + "```" + `

## 项目结构

` + "```" + `
{{.Name}}-api/
├── main.go              # 入口文件
├── go.mod               # 依赖管理
├── config/
│   └── config.go        # 配置加载
├── database/
│   └── database.go      # 数据库初始化 + 自动迁移
├── model/
│   └── {{.Name}}.go          # 数据模型
├── handler/
│   └── {{.Name}}.go          # CRUD Handler
└── router/
    └── router.go        # 路由注册
` + "```" + `

## 环境变量

| 变量     | 默认值   | 说明       |
|----------|---------|-----------|
| PORT     | 8080    | 监听端口   |
| DB_PATH  | data.db | SQLite 路径 |
`
