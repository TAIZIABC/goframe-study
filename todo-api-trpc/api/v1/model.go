// api/v1/model.go
// -----------------------------------------------
// 对标 GoFrame: api/v1/model.go
//
// GoFrame 版:
//   type TodoItem struct {
//       Id        int    `json:"id"`
//       Title     string `json:"title"`
//       Completed bool   `json:"completed"`
//   }
//
// tRPC-Go 版完全一致，这是纯数据结构，与框架无关
// -----------------------------------------------
package v1

// TodoItem 是列表返回的单个条目
type TodoItem struct {
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// ─── 统一响应结构 ───────────────────────────────
// 对标 GoFrame: ghttp.MiddlewareHandlerResponse 自动包装的格式
// GoFrame 自动包装为 { "code": 0, "message": "", "data": {...} }
// tRPC-Go 泛 HTTP 服务没有自动包装，需要手动定义
// -----------------------------------------------

// Response 统一 JSON 响应
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ─── 请求结构体 ─────────────────────────────────
// 对标 GoFrame: api/v1/todo.go 中的 Req/Res 结构体
// GoFrame 通过 g.Meta tag 自动绑定路由和校验
// tRPC-Go 版手动从 HTTP Request 解析，但保留相同的结构定义
// -----------------------------------------------

// TodoCreateReq 创建 Todo 请求
// 对标 GoFrame:
//
//	type TodoCreateReq struct {
//	    g.Meta    `path:"/todos" method:"post" tags:"Todo"`
//	    Title     string `v:"required|length:1,200" json:"title"`
//	    Completed bool   `json:"completed"`
//	}
type TodoCreateReq struct {
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

// TodoCreateRes 创建 Todo 响应
type TodoCreateRes struct {
	Id int `json:"id"`
}

// TodoListReq 列表查询请求
// 对标 GoFrame:
//
//	type TodoListReq struct {
//	    g.Meta `path:"/todos" method:"get" tags:"Todo"`
//	    Page   int `d:"1"  v:"min:1"   json:"page"`
//	    Size   int `d:"10" v:"max:100" json:"size"`
//	}
type TodoListReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

// TodoListRes 列表查询响应
type TodoListRes struct {
	List  []TodoItem `json:"list"`
	Total int        `json:"total"`
}

// TodoUpdateReq 更新 Todo 请求
// 对标 GoFrame:
//
//	type TodoUpdateReq struct {
//	    g.Meta    `path:"/todos/:id" method:"put" tags:"Todo"`
//	    Id        int    `v:"required" in:"path"`
//	    Title     string `v:"length:1,200" json:"title"`
//	    Completed *bool  `json:"completed"`
//	}
//
// 注意：Completed 使用 *bool 指针，区分"未传值"和"传了false"
type TodoUpdateReq struct {
	Id        int    `json:"-"` // 从 URL path 获取
	Title     string `json:"title"`
	Completed *bool  `json:"completed"`
}
