// Package main 演示常见的 Go 并发模式
//
// 知识点：
// 1. Pipeline（管道模式）：数据流经多个处理阶段
// 2. Fan-out/Fan-in（扇出/扇入）：并行处理后合并结果
// 3. Worker Pool（工作池）：控制并发数量
// 4. Semaphore（信号量）：使用 channel 限制并发
// 5. Rate Limiter（速率限制器）
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("===== 并发模式示例 =====")

	// pipelineDemo()
	// fanOutFanIn()
	// workerPoolDemo()
	// semaphoreDemo()
	rateLimiterDemo()
}

// ==================== 管道模式 ====================

// pipelineDemo 管道模式：数据经过多个阶段的处理
func pipelineDemo() {
	fmt.Println("\n--- 示例1: Pipeline 管道模式 ---")

	// 阶段1：生成数字
	nums := generate(1, 2, 3, 4, 5)
	// 阶段2：平方
	squared := square(nums)
	// 阶段3：过滤偶数
	filtered := filter(squared, func(n int) bool { return n%2 == 0 })

	// 消费结果
	for val := range filtered {
		fmt.Printf("  结果: %d\n", val)
	}
}

func generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func filter(in <-chan int, fn func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if fn(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

// ==================== 扇出/扇入 ====================

// fanOutFanIn 将任务分发给多个 worker 并行处理，再合并结果
func fanOutFanIn() {
	fmt.Println("\n--- 示例2: Fan-out/Fan-in ---")

	input := generate(1, 2, 3, 4, 5, 6, 7, 8)

	// Fan-out：启动多个 worker 并行处理
	worker1 := square(input)
	// 注意：这里不能直接复用 input，因为 channel 已经被 worker1 消费
	// 实际场景中需要先 fan-out 再分别处理

	// 简化演示：直接收集结果
	for val := range worker1 {
		fmt.Printf("  平方: %d\n", val)
	}
}

// ==================== 工作池 ====================

// workerPoolDemo 固定数量的 worker 处理任务
func workerPoolDemo() {
	fmt.Println("\n--- 示例3: Worker Pool 工作池 ---")

	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan int, numJobs)
	results := make(chan string, numJobs)

	// 启动固定数量的 worker
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				time.Sleep(10 * time.Millisecond) // 模拟处理
				results <- fmt.Sprintf("Worker %d 完成任务 %d", id, job)
			}
		}(w)
	}

	// 发送任务
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// 等待所有 worker 完成后关闭 results
	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果
	for result := range results {
		fmt.Printf("  %s\n", result)
	}
}

// ==================== 信号量 ====================

// semaphoreDemo 使用 buffered channel 作为信号量限制并发数
func semaphoreDemo() {
	fmt.Println("\n--- 示例4: Semaphore 信号量 ---")

	const maxConcurrent = 3
	sem := make(chan struct{}, maxConcurrent) // 信号量

	var wg sync.WaitGroup

	for i := 1; i <= 8; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sem <- struct{}{}        // 获取信号量（占坑）
			defer func() { <-sem }() // 释放信号量（让坑）

			fmt.Printf("  任务 %d 开始执行（并发数 ≤ %d）\n", id, maxConcurrent)
			time.Sleep(30 * time.Millisecond)
			fmt.Printf("  任务 %d 完成\n", id)
		}(i)
	}

	wg.Wait()
	fmt.Println("  ✅ 所有任务完成")
}

// ==================== 速率限制 ====================

// rateLimiterDemo 使用 time.Ticker 实现速率限制
func rateLimiterDemo() {
	fmt.Println("\n--- 示例5: Rate Limiter 速率限制 ---")

	// 每 50ms 处理一个请求
	limiter := time.NewTicker(50 * time.Millisecond)
	defer limiter.Stop()

	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	start := time.Now()
	for req := range requests {
		<-limiter.C // 等待令牌
		fmt.Printf("  处理请求 %d (耗时 %v)\n", req, time.Since(start).Round(time.Millisecond))
	}
}
