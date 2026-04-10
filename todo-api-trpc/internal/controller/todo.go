// internal/controller/todo.go
// -----------------------------------------------
// Todo Controller 层 — 请求处理
// 对标 GoFrame: internal/controller/todo.go
//
// GoFrame 控制器通过 struct 方法 + 反射自动绑定：
//   func (c *cTodo) Create(ctx context.Context, req *v1.TodoCreateReq) (*v1.TodoCreateRes, error)
//   - 框架自动从 HTTP 请求中解析参数到 req struct
//   - 框架自动根据 v tag 校验参数
//   - 框架自动将返回值通过中间件包装为统一 JSON
//
// tRPC-Go 泛 HTTP 标准服务使用标准 http.Handler 风格：
//   func(w http.ResponseWriter, r *http.Request) error
//   - 需要手动解析 JSON body / URL 参数
//   - 需要手动调用校验函数
//   - 需要手动序列化 JSON 响应
// -----------------------------------------------
package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	v1 "todo-api-trpc/api/v1"
	"todo-api-trpc/internal/service"
)

// writeJSON 写入统一 JSON 响应
// 对标 GoFrame: ghttp.MiddlewareHandlerResponse 自动包装
// GoFrame 自动将 controller 返回值包装为: {"code":0,"message":"","data":{...}}
// tRPC-Go 版需要手动包装
func writeJSON(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	resp := v1.Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(resp)
}

// writeError 写入错误响应
// 对标 GoFrame: controller 返回 err 时，中间件自动设置非0 code
func writeError(w http.ResponseWriter, code int, message string) {
	writeJSON(w, code, message, nil)
}

// Create 创建 Todo
// -----------------------------------------------
// 对标 GoFrame:
//   func (c *cTodo) Create(ctx context.Context, req *v1.TodoCreateReq) (*v1.TodoCreateRes, error) {
//       id, err := service.Todo().Create(ctx, req.Title, req.Completed)
//       return &v1.TodoCreateRes{Id: id}, nil
//   }
//
// 路由: POST /api/v1/todos
// 对标 Go Meta: `path:"/todos" method:"post"`
// -----------------------------------------------
func Create(w http.ResponseWriter, r *http.Request) error {
	// 手动解析 JSON body — GoFrame 通过反射自动完成
	var req v1.TodoCreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 51, "请求参数解析失败: "+err.Error())
		return nil
	}

	// 手动校验 — GoFrame 通过 v tag 自动完成
	if err := v1.ValidateCreateReq(&req); err != nil {
		writeError(w, 51, err.Error())
		return nil
	}

	// 调用 service — 与 GoFrame 版完全一致
	id, err := service.Todo().Create(r.Context(), req.Title, req.Completed)
	if err != nil {
		writeError(w, 50, err.Error())
		return nil
	}

	// 返回 — 对标 &v1.TodoCreateRes{Id: id}
	writeJSON(w, 0, "", v1.TodoCreateRes{Id: id})
	return nil
}

// List 分页查询 Todo 列表
// -----------------------------------------------
// 对标 GoFrame:
//   func (c *cTodo) List(ctx context.Context, req *v1.TodoListReq) (*v1.TodoListRes, error) {
//       list, total, err := service.Todo().List(ctx, req.Page, req.Size)
//       return &v1.TodoListRes{List: list, Total: total}, nil
//   }
//
// 路由: GET /api/v1/todos
// 对标 Go Meta: `path:"/todos" method:"get"`
// -----------------------------------------------
func List(w http.ResponseWriter, r *http.Request) error {
	// 从 URL query 解析参数 — GoFrame 自动解析 + d tag 设默认值
	req := v1.TodoListReq{}
	if p := r.URL.Query().Get("page"); p != "" {
		req.Page, _ = strconv.Atoi(p)
	}
	if s := r.URL.Query().Get("size"); s != "" {
		req.Size, _ = strconv.Atoi(s)
	}

	// 校验并填充默认值
	if err := v1.ValidateListReq(&req); err != nil {
		writeError(w, 51, err.Error())
		return nil
	}

	// 调用 service — 与 GoFrame 版完全一致
	list, total, err := service.Todo().List(r.Context(), req.Page, req.Size)
	if err != nil {
		writeError(w, 50, err.Error())
		return nil
	}

	// 返回 — 对标 &v1.TodoListRes{List: list, Total: total}
	writeJSON(w, 0, "", v1.TodoListRes{List: list, Total: total})
	return nil
}

// Update 更新 Todo
// -----------------------------------------------
// 对标 GoFrame:
//   func (c *cTodo) Update(ctx context.Context, req *v1.TodoUpdateReq) (*v1.TodoUpdateRes, error) {
//       return nil, service.Todo().Update(ctx, req.Id, req.Title, req.Completed)
//   }
//
// 路由: PUT /api/v1/todos/{id}
// 对标 Go Meta: `path:"/todos/:id" method:"put"`
// -----------------------------------------------
func Update(w http.ResponseWriter, r *http.Request) error {
	// 从 URL path 获取 id — GoFrame 通过 in:"path" tag 自动获取
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		writeError(w, 51, "id 是必填项")
		return nil
	}

	// 解析 body
	var req v1.TodoUpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, 51, "请求参数解析失败: "+err.Error())
		return nil
	}
	req.Id = id

	// 校验
	if err := v1.ValidateUpdateReq(&req); err != nil {
		writeError(w, 51, err.Error())
		return nil
	}

	// 调用 service — 与 GoFrame 版完全一致
	if err := service.Todo().Update(r.Context(), req.Id, req.Title, req.Completed); err != nil {
		writeError(w, 50, err.Error())
		return nil
	}

	// 返回空对象 — 对标 TodoUpdateRes{}
	writeJSON(w, 0, "", nil)
	return nil
}

// Delete 删除 Todo
// -----------------------------------------------
// 对标 GoFrame:
//   func (c *cTodo) Delete(ctx context.Context, req *v1.TodoDeleteReq) (*v1.TodoDeleteRes, error) {
//       return nil, service.Todo().Delete(ctx, req.Id)
//   }
//
// 路由: DELETE /api/v1/todos/{id}
// 对标 Go Meta: `path:"/todos/:id" method:"delete"`
// -----------------------------------------------
func Delete(w http.ResponseWriter, r *http.Request) error {
	// 从 URL path 获取 id
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id <= 0 {
		writeError(w, 51, "id 是必填项")
		return nil
	}

	// 调用 service — 与 GoFrame 版完全一致
	if err := service.Todo().Delete(r.Context(), id); err != nil {
		writeError(w, 50, err.Error())
		return nil
	}

	// 返回空对象 — 对标 TodoDeleteRes{}
	writeJSON(w, 0, "", nil)
	return nil
}
