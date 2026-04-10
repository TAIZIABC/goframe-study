// main.go
// -----------------------------------------------
// 应用入口文件
// 对标 GoFrame: main.go + internal/cmd/cmd.go
//
// GoFrame 入口：
//   import _ "github.com/gogf/gf/contrib/drivers/mysql/v2"  // 副作用导入 MySQL 驱动
//   func main() {
//       cmd.Main.Run(gctx.GetInitCtx())  // 启动命令（内部创建 Server + 注册路由 + Run）
//   }
//
// tRPC-Go 入口：
//   s := trpc.NewServer()                                    // 读取 trpc_go.yaml 创建 Server
//   thttp.RegisterNoProtocolServiceMux(s.Service(...), mux)  // 注册泛 HTTP 服务 + 路由
//   s.Serve()                                                // 启动服务
// -----------------------------------------------
package main

import (
	"fmt"

	trpc "trpc.group/trpc-go/trpc-go"
	thttp "trpc.group/trpc-go/trpc-go/http"

	"todo-api-trpc/internal/cmd"
	"todo-api-trpc/internal/database"
)

func main() {
	// ─── 初始化数据库 ───────────────────────────
	// 对标 GoFrame: 框架根据 config.yaml 中 database 配置自动初始化
	// tRPC-Go 没有内置数据库管理，需要手动初始化
	cfg := database.DefaultConfig()
	if err := database.Init(cfg); err != nil {
		panic(fmt.Sprintf("database init failed: %v", err))
	}
	fmt.Println("[DB] MySQL connected successfully")

	// ─── 创建 tRPC Server ────────────────────────
	// 对标 GoFrame: s := g.Server()
	// tRPC-Go: trpc.NewServer() 会自动读取 trpc_go.yaml 配置
	s := trpc.NewServer()

	// ─── 创建路由 ────────────────────────────────
	// 对标 GoFrame:
	//   s.Use(ghttp.MiddlewareHandlerResponse)
	//   s.Group("/api/v1", func(group *ghttp.RouterGroup) {
	//       group.Bind(controller.Todo)
	//   })
	router := cmd.NewRouter()

	// ─── 注册泛 HTTP 标准服务 ─────────────────────
	// 对标 GoFrame: s.Run() (内部完成路由绑定和服务启动)
	// tRPC-Go 使用 RegisterNoProtocolServiceMux 注册 Mux 路由
	// 参数中的 service name 必须与 trpc_go.yaml 中配置的 name 一致
	thttp.RegisterNoProtocolServiceMux(
		s.Service("trpc.todo.api.stdhttp"),
		router,
	)

	fmt.Println("[SERVER] Todo API (tRPC-Go) is running on :8000")
	fmt.Println("[SERVER] API 基础路径: http://localhost:8000/api/v1")

	// ─── 启动服务 ────────────────────────────────
	// 对标 GoFrame: s.Run()
	if err := s.Serve(); err != nil {
		panic(err)
	}
}
