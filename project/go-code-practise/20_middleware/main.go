package main

import (
	"fmt"
	"time"
)

// 题目：实现中间件链（洋葱模型）
//
// 定义：
//   type Handler func(ctx *Context)
//   type Middleware func(next Handler) Handler
//   type Context struct { Data map[string]any }
//
// 实现：
//   Chain(middlewares ...Middleware) Middleware — 组合多个中间件为一个
//
// 示例：实现 Logger、Timer、Recovery 三个中间件
// Logger: 打印 "→ start" 和 "← end"
// Timer:  记录执行耗时
// Recovery: 捕获 panic

type Context struct {
	Data map[string]any
}

type Handler func(ctx *Context)
type Middleware func(next Handler) Handler

func Chain(middlewares ...Middleware) Middleware {
	return func(next Handler) Handler {
		// 从右往左包装：最左的中间件在最外层
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// Logger 中间件
func Logger() Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) {
			fmt.Println("→ start")
			next(ctx)
			fmt.Println("← end")
		}
	}
}

// Timer 中间件
func Timer() Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) {
			start := time.Now()
			next(ctx)
			fmt.Printf("  [timer] took %v\n", time.Since(start))
		}
	}
}

// Recovery 中间件 — 捕获 panic
func Recovery() Middleware {
	return func(next Handler) Handler {
		return func(ctx *Context) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("  [recovery] caught panic: %v\n", r)
				}
			}()
			next(ctx)
		}
	}
}

func main() {
	// 组合中间件
	combined := Chain(Logger(), Timer(), Recovery())

	// 最终业务 handler
	handler := combined(func(ctx *Context) {
		fmt.Println("  executing business logic")
		time.Sleep(10 * time.Millisecond)
		ctx.Data["result"] = "success"
	})

	ctx := &Context{Data: make(map[string]any)}
	handler(ctx)
	fmt.Println("result:", ctx.Data["result"])

	_ = time.Now()
}
