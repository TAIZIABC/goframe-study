package v1

import "github.com/gogf/gf/v2/frame/g"

// 创建 Todo
type TodoCreateReq struct {
    g.Meta    `path:"/todos" method:"post" tags:"Todo"`
    Title     string `v:"required|length:1,200" json:"title"`
    Completed bool   `json:"completed"`
}
type TodoCreateRes struct {
    Id int `json:"id"`
}

// 列表查询
type TodoListReq struct {
    g.Meta `path:"/todos" method:"get" tags:"Todo"`
    Page   int `d:"1"  v:"min:1"       json:"page"`
    Size   int `d:"10" v:"max:100"     json:"size"`
}
type TodoListRes struct {
    List  []TodoItem `json:"list"`
    Total int        `json:"total"`
}

// 更新 & 删除
type TodoUpdateReq struct {
    g.Meta    `path:"/todos/:id" method:"put" tags:"Todo"`
    Id        int    `v:"required" in:"path"`
    Title     string `v:"length:1,200" json:"title"`
    Completed *bool  `json:"completed"`
}
type TodoUpdateRes struct{}

type TodoDeleteReq struct {
    g.Meta `path:"/todos/:id" method:"delete" tags:"Todo"`
    Id     int `v:"required" in:"path"`
}
type TodoDeleteRes struct{}