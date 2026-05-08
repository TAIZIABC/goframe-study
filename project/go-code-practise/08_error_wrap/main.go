package main

import (
	"errors"
	"fmt"
)

// 题目：实现自定义错误 + 错误链
//
// 1. 定义 AppError 结构体，包含 Code int 和 Message string
// 2. 实现 error 接口：Error() 返回 "[Code] Message"
// 3. 实现 Wrap(err error, code int, msg string) *AppError — 包装原始错误
// 4. 实现 Unwrap() error — 支持 errors.Is / errors.As
//
// 示例：
//   base := fmt.Errorf("connection refused")
//   e := Wrap(base, 500, "db error")
//   fmt.Println(e)                    // [500] db error
//   fmt.Println(errors.Is(e, base))   // true

type AppError struct {
	// 在此定义字段
}

func (e *AppError) Error() string {
	// 在此实现
	return ""
}

func (e *AppError) Unwrap() error {
	// 在此实现
	return nil
}

func Wrap(err error, code int, msg string) *AppError {
	// 在此实现
	return nil
}

func main() {
	base := fmt.Errorf("connection refused")
	e := Wrap(base, 500, "db error")
	fmt.Println(e)                  // [500] db error
	fmt.Println(errors.Is(e, base)) // true

	var appErr *AppError
	fmt.Println(errors.As(e, &appErr)) // true
}
