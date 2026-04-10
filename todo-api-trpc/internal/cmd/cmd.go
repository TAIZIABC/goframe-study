// internal/cmd/cmd.go
// -----------------------------------------------
// 路由注册 & 服务启动命令
// 对标 GoFrame: internal/cmd/cmd.go
//
// GoFrame 通过 g.Meta tag 自动路由绑定：
//   var Main = gcmd.Command{
//       Func: func(ctx, parser) error {
//           s := g.Server()
//           s.Use(ghttp.MiddlewareHandlerResponse)
//           s.Group("/api/v1", func(group *ghttp.RouterGroup) {
//               group.Bind(controller.Todo)  // 反射解析 g.Meta 自动绑定
//           })
//           s.Run()
//       },
//   }
//
// tRPC-Go 使用 gorilla/mux 手动注册路由：
//   router := mux.NewRouter()
//   api := router.PathPrefix("/api/v1").Subrouter()
//   api.HandleFunc("/todos", controller.Create).Methods("POST")
//   thttp.RegisterNoProtocolServiceMux(s.Service(...), router)
// -----------------------------------------------
package cmd

import (
	"net/http"

	"github.com/gorilla/mux"

	"todo-api-trpc/internal/controller"
)

// NewRouter 创建并注册所有路由
// 对标 GoFrame:
//
//	s.Group("/api/v1", func(group *ghttp.RouterGroup) {
//	    group.Bind(controller.Todo)
//	})
//
// GoFrame 通过 group.Bind() + g.Meta tag 自动解析路由和 HTTP 方法
// tRPC-Go 需要手动注册每个路由
func NewRouter() *mux.Router {
	router := mux.NewRouter()

	// 路由组 /api/v1 — 对标 GoFrame: s.Group("/api/v1", ...)
	api := router.PathPrefix("/api/v1").Subrouter()

	// 路由定义 — 对标 GoFrame api/v1/todo.go 中的 g.Meta 标签
	//
	// GoFrame:  g.Meta `path:"/todos"    method:"post"`    → POST   /api/v1/todos
	// tRPC-Go:  api.HandleFunc("/todos", Create).Methods("POST")
	api.HandleFunc("/todos", wrapHandler(controller.Create)).Methods(http.MethodPost)

	// GoFrame:  g.Meta `path:"/todos"    method:"get"`     → GET    /api/v1/todos
	// tRPC-Go:  api.HandleFunc("/todos", List).Methods("GET")
	api.HandleFunc("/todos", wrapHandler(controller.List)).Methods(http.MethodGet)

	// GoFrame:  g.Meta `path:"/todos/:id" method:"put"`    → PUT    /api/v1/todos/{id}
	// tRPC-Go:  api.HandleFunc("/todos/{id}", Update).Methods("PUT")
	// 注意：GoFrame 用 :id，gorilla/mux 用 {id}
	api.HandleFunc("/todos/{id:[0-9]+}", wrapHandler(controller.Update)).Methods(http.MethodPut)

	// GoFrame:  g.Meta `path:"/todos/:id" method:"delete"` → DELETE /api/v1/todos/{id}
	// tRPC-Go:  api.HandleFunc("/todos/{id}", Delete).Methods("DELETE")
	api.HandleFunc("/todos/{id:[0-9]+}", wrapHandler(controller.Delete)).Methods(http.MethodDelete)

	return router
}

// wrapHandler 将 tRPC-Go 风格的 handler (返回 error) 包装为标准 http.HandlerFunc
// tRPC-Go 泛 HTTP 的 handler 签名是 func(w, r) error
// gorilla/mux 需要标准的 http.HandlerFunc
func wrapHandler(fn func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
