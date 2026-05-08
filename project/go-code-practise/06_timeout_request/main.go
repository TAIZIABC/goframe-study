package main

import (
	"context"
	"fmt"
	"time"
)

// 题目：实现 FetchWithTimeout
// 模拟一个耗时操作（用 time.Sleep 模拟），如果超时则返回错误
// FetchWithTimeout(ctx context.Context, delay time.Duration) (string, error)
// - delay 模拟操作耗时
// - 如果 ctx 超时先到，返回 "", ctx.Err()
// - 否则返回 "data", nil

func FetchWithTimeout(ctx context.Context, delay time.Duration) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(delay):
		return "data", nil
	}
}

func main() {
	// 超时 100ms，操作耗时 50ms → 成功
	ctx1, cancel1 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel1()
	result, err := FetchWithTimeout(ctx1, 50*time.Millisecond)
	fmt.Println(result, err) // data <nil>

	// 超时 50ms，操作耗时 100ms → 超时
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	result, err = FetchWithTimeout(ctx2, 100*time.Millisecond)
	fmt.Println(result, err) // "" context deadline exceeded
}
