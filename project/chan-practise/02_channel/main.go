// Package main 演示 Channel 的通信机制
//
// 知识点：
// 1. 无缓冲 channel：发送和接收必须同步（阻塞式）
// 2. 有缓冲 channel：缓冲区满时发送阻塞，空时接收阻塞
// 3. 单向 channel：限制 channel 只能发送或接收
// 4. channel 的关闭与 range 遍历
// 5. 使用 channel 实现生产者-消费者模式
package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("===== Channel 通信示例 =====")

	// unbufferedChannel()
	// bufferedChannel()
	// directionalChannel()
	// channelRange()
	producerConsumer()
}

// unbufferedChannel 无缓冲 channel 示例
// 无缓冲 channel 要求发送方和接收方同时就绪
func unbufferedChannel() {
	fmt.Println("\n--- 示例1: 无缓冲 Channel ---")

	ch := make(chan string) // 无缓冲 channel

	go func() {
		fmt.Println("  发送方: 准备发送数据...")
		ch <- "Hello, Channel!" // 阻塞，直到有接收方
		fmt.Println("  发送方: 数据已发送")
	}()

	time.Sleep(100 * time.Millisecond) // 模拟接收方延迟
	msg := <-ch                        // 接收数据
	fmt.Printf("  接收方: 收到 -> %s\n", msg)
}

// bufferedChannel 有缓冲 channel 示例
// 有缓冲 channel 在缓冲区未满时，发送不会阻塞
func bufferedChannel() {
	fmt.Println("\n--- 示例2: 有缓冲 Channel ---")

	ch := make(chan int, 3) // 缓冲大小为 3

	// 连续发送3个值不会阻塞
	ch <- 1
	ch <- 2
	ch <- 3
	fmt.Printf("  缓冲区: len=%d, cap=%d\n", len(ch), cap(ch))

	// 第4个会阻塞（除非有人接收）
	// ch <- 4 // 如果取消注释会导致死锁

	// 依次接收
	fmt.Printf("  接收: %d\n", <-ch)
	fmt.Printf("  接收: %d\n", <-ch)
	fmt.Printf("  接收: %d\n", <-ch)
}

// directionalChannel 单向 channel 示例
// 用于在函数签名中限制 channel 的使用方向，提高代码安全性
func directionalChannel() {
	fmt.Println("\n--- 示例3: 单向 Channel ---")

	ch := make(chan int)

	// send 只能往 channel 发送
	send := func(ch chan<- int, val int) {
		ch <- val
		fmt.Printf("  发送: %d\n", val)
	}

	// receive 只能从 channel 接收
	receive := func(ch <-chan int) int {
		val := <-ch
		fmt.Printf("  接收: %d\n", val)
		return val
	}

	go send(ch, 42)
	receive(ch)
}

// channelRange 使用 range 遍历 channel
// 关闭 channel 后，range 循环才会结束
func channelRange() {
	fmt.Println("\n--- 示例4: Channel 关闭与 Range ---")

	ch := make(chan int, 5)

	// 发送数据后关闭 channel
	go func() {
		for i := 1; i <= 5; i++ {
			ch <- i * i
		}
		close(ch) // ⚠️ 关闭 channel，通知接收方不再有数据
	}()

	// 使用 range 遍历，channel 关闭后自动退出
	for val := range ch {
		fmt.Printf("  收到平方数: %d\n", val)
	}

	// 从已关闭的 channel 接收会立即返回零值
	val, ok := <-ch
	fmt.Printf("  关闭后接收: val=%d, ok=%v\n", val, ok)
}

// producerConsumer 生产者-消费者模式
func producerConsumer() {
	fmt.Println("\n--- 示例5: 生产者-消费者模式 ---")

	jobs := make(chan int, 10)    // 任务队列
	results := make(chan int, 10) // 结果队列

	// 启动 3 个消费者（worker）
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// 生产者发送 9 个任务
	for j := 1; j <= 9; j++ {
		jobs <- j
	}
	close(jobs)

	// 收集所有结果
	for i := 0; i < 9; i++ {
		result := <-results
		fmt.Printf("  结果: %d\n", result)
	}
}

// worker 从 jobs 接收任务，将结果发送到 results
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("  Worker %d 处理任务 %d\n", id, j)
		time.Sleep(10 * time.Millisecond)
		results <- j * 2 // 模拟处理：结果为任务值的两倍
	}
}
