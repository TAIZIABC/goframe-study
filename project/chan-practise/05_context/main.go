// Package main 演示 Context 的使用
//
// 知识点：
// 1. context.WithCancel：手动取消
// 2. context.WithTimeout：超时自动取消
// 3. context.WithDeadline：截止时间取消
// 4. context.WithValue：传递请求级别数据
// 5. Context 在 HTTP 服务和数据库查询中的最佳实践
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	fmt.Println("===== Context 示例 =====")

	cancelDemo()
	timeoutDemo()
	deadlineDemo()
	valueDemo()
	cascadeCancel()
	realWorldExample()
}

// cancelDemo 手动取消
func cancelDemo() {
	fmt.Println("\n--- 示例1: WithCancel 手动取消 ---")

	ctx, cancel := context.WithCancel(context.Background())

	// 模拟后台任务
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("  后台任务停止: %v\n", ctx.Err())
				return
			default:
				fmt.Println("  后台任务运行中...")
				time.Sleep(30 * time.Millisecond)
			}
		}
	}(ctx)

	time.Sleep(100 * time.Millisecond)
	cancel() // 发送取消信号
	time.Sleep(50 * time.Millisecond)
}

// timeoutDemo 超时自动取消
func timeoutDemo() {
	fmt.Println("\n--- 示例2: WithTimeout 超时控制 ---")

	// 设置 100ms 超时
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel() // ⚠️ 最佳实践：即使会自动取消，也要 defer cancel 释放资源

	// 模拟慢操作
	result := make(chan string, 1)
	go func() {
		time.Sleep(200 * time.Millisecond) // 超过超时时间
		result <- "操作完成"
	}()

	select {
	case res := <-result:
		fmt.Printf("  结果: %s\n", res)
	case <-ctx.Done():
		fmt.Printf("  ⏰ 超时: %v\n", ctx.Err())
	}
}

// deadlineDemo 截止时间
func deadlineDemo() {
	fmt.Println("\n--- 示例3: WithDeadline 截止时间 ---")

	deadline := time.Now().Add(80 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	// 查看截止时间
	if d, ok := ctx.Deadline(); ok {
		fmt.Printf("  截止时间: %v\n", d.Format("15:04:05.000"))
	}

	select {
	case <-time.After(200 * time.Millisecond):
		fmt.Println("  操作完成")
	case <-ctx.Done():
		fmt.Printf("  已超过截止时间: %v\n", ctx.Err())
	}
}

// valueDemo 传递请求级别数据
func valueDemo() {
	fmt.Println("\n--- 示例4: WithValue 传递数据 ---")

	type contextKey string
	const (
		userIDKey  contextKey = "userID"
		traceIDKey contextKey = "traceID"
	)

	// 逐层添加值
	ctx := context.Background()
	ctx = context.WithValue(ctx, userIDKey, "user-123")
	ctx = context.WithValue(ctx, traceIDKey, "trace-abc-456")

	// 在下游函数中获取值
	processRequest(ctx, userIDKey, traceIDKey)
}

func processRequest(ctx context.Context, userKey, traceKey fmt.Stringer) {
	userID := ctx.Value(userKey)
	traceID := ctx.Value(traceKey)
	fmt.Printf("  处理请求: userID=%v, traceID=%v\n", userID, traceID)
}

// cascadeCancel 级联取消：父 context 取消时，所有子 context 也被取消
func cascadeCancel() {
	fmt.Println("\n--- 示例5: 级联取消 ---")

	parent, parentCancel := context.WithCancel(context.Background())

	child1, _ := context.WithCancel(parent)
	child2, _ := context.WithCancel(parent)

	// 监听子 context
	go func() {
		<-child1.Done()
		fmt.Println("  子 Context 1 被取消")
	}()
	go func() {
		<-child2.Done()
		fmt.Println("  子 Context 2 被取消")
	}()

	fmt.Println("  取消父 Context...")
	parentCancel() // 取消父 context，子 context 也会被取消
	time.Sleep(50 * time.Millisecond)
}

// realWorldExample 模拟实际场景：并发请求多个API，任一超时则全部取消
func realWorldExample() {
	fmt.Println("\n--- 示例6: 实际场景 - 并发 API 调用 ---")

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	type apiResult struct {
		name string
		data string
		err  error
	}

	results := make(chan apiResult, 3)

	// 模拟并发调用多个 API
	callAPI := func(ctx context.Context, name string, delay time.Duration) {
		select {
		case <-time.After(delay):
			results <- apiResult{name: name, data: fmt.Sprintf("%s 返回数据", name)}
		case <-ctx.Done():
			results <- apiResult{name: name, err: ctx.Err()}
		}
	}

	go callAPI(ctx, "用户服务", 50*time.Millisecond)
	go callAPI(ctx, "订单服务", 100*time.Millisecond)
	go callAPI(ctx, "支付服务", 200*time.Millisecond) // 这个会超时

	for i := 0; i < 3; i++ {
		r := <-results
		if r.err != nil {
			fmt.Printf("  ❌ %s: %v\n", r.name, r.err)
		} else {
			fmt.Printf("  ✅ %s: %s\n", r.name, r.data)
		}
	}
}
