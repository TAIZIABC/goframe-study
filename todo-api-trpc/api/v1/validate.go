// api/v1/validate.go
// -----------------------------------------------
// 参数校验模块
// 对标 GoFrame: struct tag 自动校验 (v tag)
//
// GoFrame 通过声明式 tag 实现校验：
//   Title string `v:"required|length:1,200"`
//   Page  int    `d:"1" v:"min:1"`
//
// tRPC-Go 没有内置校验机制，需要手动实现
// 这里用独立的校验函数，保持与 GoFrame 版相同的校验规则
// -----------------------------------------------
package v1

import "fmt"

// ValidateCreateReq 校验创建请求
// 对标 GoFrame: v:"required|length:1,200"
func ValidateCreateReq(req *TodoCreateReq) error {
	if req.Title == "" {
		return fmt.Errorf("title 是必填项")
	}
	if len(req.Title) > 200 {
		return fmt.Errorf("title 长度需在 1-200 之间")
	}
	return nil
}

// ValidateListReq 校验并填充列表查询默认值
// 对标 GoFrame: d:"1" v:"min:1" 和 d:"10" v:"max:100"
func ValidateListReq(req *TodoListReq) error {
	if req.Page <= 0 {
		req.Page = 1 // 对标 GoFrame d:"1" 默认值
	}
	if req.Size <= 0 {
		req.Size = 10 // 对标 GoFrame d:"10" 默认值
	}
	if req.Size > 100 {
		return fmt.Errorf("size 最大值为 100")
	}
	return nil
}

// ValidateUpdateReq 校验更新请求
// 对标 GoFrame: v:"required"(id) + v:"length:1,200"(title)
func ValidateUpdateReq(req *TodoUpdateReq) error {
	if req.Id <= 0 {
		return fmt.Errorf("id 是必填项")
	}
	if req.Title != "" && len(req.Title) > 200 {
		return fmt.Errorf("title 长度需在 1-200 之间")
	}
	return nil
}
