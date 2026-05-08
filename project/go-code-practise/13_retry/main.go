package main

import (
	"fmt"
	"time"
)

// 题目：实现带指数退避的重试函数
//
// Retry(fn func() error, maxRetries int, baseDelay time.Duration) error
// - 执行 fn，成功(nil)立即返回
// - 失败则等待 baseDelay * 2^i 后重试
// - 最多重试 maxRetries 次，全部失败返回最后一个错误
//
// 示例：前2次失败第3次成功，maxRetries=5 → 返回 nil

func Retry(fn func() error, maxRetries int, baseDelay time.Duration) error {
	var err error
	for i := 0; i <= maxRetries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		if i < maxRetries {
			time.Sleep(baseDelay * (1 << i)) // 指数退避：baseDelay * 2^i
		}
	}
	return err
}

func main() {
	attempt := 0
	err := Retry(func() error {
		attempt++
		if attempt < 3 {
			return fmt.Errorf("fail attempt %d", attempt)
		}
		return nil
	}, 5, 10*time.Millisecond)

	fmt.Println(err)     // <nil>
	fmt.Println(attempt) // 3
}
