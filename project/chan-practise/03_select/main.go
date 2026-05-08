// Package main 演示 Select 语句的使用
//
// 知识点：
// 1. select 用于同时监听多个 channel 操作
// 2. 当多个 case 同时就绪时，随机选择一个执行
// 3. default 分支实现非阻塞操作
// 4. select + time.After 实现超时控制
// 5. 使用 select 实现优雅退出
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("===== Select 语句示例 =====")

	basicSelect()
	nonBlockingSelect()
	timeoutSelect()
	fanIn()
	gracefulShutdown()
}

// basicSelect 基础 select 用法
func basicSelect() {
	fmt.Println("\n--- 示例1: 基础 Select ---")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(10 * time.Millisecond)
		ch1 <- "来自 channel 1"
	}()

	go func() {
		time.Sleep(20 * time.Millisecond)
		ch2 <- "来自 channel 2"
	}()

	// 接收两次，每次选择先就绪的 channel
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Printf("  收到: %s\n", msg)
		case msg := <-ch2:
			fmt.Printf("  收到: %s\n", msg)
		}
	}
}

// nonBlockingSelect 非阻塞 select（使用 default）
func nonBlockingSelect() {
	fmt.Println("\n--- 示例2: 非阻塞 Select ---")

	ch := make(chan int, 1)

	// 非阻塞接收：channel 为空时走 default
	select {
	case val := <-ch:
		fmt.Printf("  收到: %d\n", val)
	default:
		fmt.Println("  channel 为空，没有数据可接收")
	}

	// 非阻塞发送：channel 满时走 default
	ch <- 1
	select {
	case ch <- 2:
		fmt.Println("  发送成功")
	default:
		fmt.Println("  channel 已满，发送被跳过")
	}
}

// timeoutSelect 使用 select + time.After 实现超时
func timeoutSelect() {
	fmt.Println("\n--- 示例3: 超时控制 ---")

	ch := make(chan string)

	go func() {
		time.Sleep(200 * time.Millisecond) // 模拟耗时操作
		ch <- "操作完成"
	}()

	select {
	case result := <-ch:
		fmt.Printf("  结果: %s\n", result)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("  ⏰ 操作超时！")
	}
}

// fanIn 扇入模式：将多个 channel 合并为一个
func fanIn() {
	fmt.Println("\n--- 示例4: 扇入模式 (Fan-In) ---")

	// 模拟多个数据源
	source := func(name string, delay time.Duration) <-chan string {
		ch := make(chan string)
		go func() {
			for i := 0; i < 3; i++ {
				time.Sleep(delay)
				ch <- fmt.Sprintf("%s: 数据 %d", name, i+1)
			}
			close(ch)
		}()
		return ch
	}

	ch1 := source("API-A", 10*time.Millisecond)
	ch2 := source("API-B", 15*time.Millisecond)

	// 合并两个 channel 的数据
	timeout := time.After(200 * time.Millisecond)
	received := 0
	for received < 6 {
		select {
		case msg, ok := <-ch1:
			if ok {
				fmt.Printf("  收到: %s\n", msg)
				received++
			}
		case msg, ok := <-ch2:
			if ok {
				fmt.Printf("  收到: %s\n", msg)
				received++
			}
		case <-timeout:
			fmt.Println("  总超时，退出")
			return
		}
	}
}

// gracefulShutdown 使用 select 实现优雅退出
func gracefulShutdown() {
	fmt.Println("\n--- 示例5: 优雅退出 ---")

	done := make(chan struct{}) // 退出信号
	work := make(chan int)

	// 工作 goroutine
	go func() {
		for {
			select {
			case <-done:
				fmt.Println("  Worker: 收到退出信号，清理资源...")
				return
			case val := <-work:
				fmt.Printf("  Worker: 处理任务 %d\n", val)
			}
		}
	}()

	// 发送几个任务
	for i := 1; i <= 3; i++ {
		work <- i
	}

	// 发送退出信号
	close(done)
	time.Sleep(50 * time.Millisecond)
	fmt.Println("  程序优雅退出")

	_ = rand.Int() // 消除 import 未使用警告
}
