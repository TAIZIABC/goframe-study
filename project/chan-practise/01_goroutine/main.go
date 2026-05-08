// Package main 演示 Goroutine 的创建与基本使用
//
// 知识点：
// 1. go 关键字启动 goroutine
// 2. goroutine 是轻量级线程，由 Go 运行时调度
// 3. 主 goroutine 退出后，其他 goroutine 也会被终止
// 4. 使用 sync.WaitGroup 等待所有 goroutine 完成
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	fmt.Println("===== Goroutine 基础示例 =====")

	// 示例1：简单启动 goroutine
	simpleGoroutine()

	// 示例2：启动多个 goroutine 并等待完成
	multipleGoroutines()

	// 示例3：查看 goroutine 数量
	goroutineCount()

	// 示例4：闭包陷阱
	closureTrap()
}

// simpleGoroutine 演示最简单的 goroutine 启动
func simpleGoroutine() {
	fmt.Println("\n--- 示例1: 简单 Goroutine ---")

	go func() {
		fmt.Println("  Hello from goroutine!")
	}()

	// 主 goroutine 需要等待，否则程序直接退出
	time.Sleep(100 * time.Millisecond)
	fmt.Println("  Hello from main!")
}

// multipleGoroutines 演示启动多个 goroutine 并用 WaitGroup 等待
func multipleGoroutines() {
	fmt.Println("\n--- 示例2: 多个 Goroutine ---")

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  Goroutine %d 正在执行\n", id)
			time.Sleep(time.Duration(id*10) * time.Millisecond)
			fmt.Printf("  Goroutine %d 执行完毕\n", id)
		}(i) // 注意：将 i 作为参数传入，避免闭包陷阱
	}

	wg.Wait()
	fmt.Println("  所有 goroutine 执行完毕")
}

// goroutineCount 演示如何查看当前 goroutine 数量
func goroutineCount() {
	fmt.Println("\n--- 示例3: Goroutine 数量 ---")
	fmt.Printf("  当前 goroutine 数量: %d\n", runtime.NumGoroutine())

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(50 * time.Millisecond)
		}()
	}

	fmt.Printf("  启动10个后 goroutine 数量: %d\n", runtime.NumGoroutine())
	wg.Wait()
	// 等待一小段时间让 goroutine 完全退出
	time.Sleep(10 * time.Millisecond)
	fmt.Printf("  等待完成后 goroutine 数量: %d\n", runtime.NumGoroutine())
}

// closureTrap 演示 goroutine 中常见的闭包陷阱
func closureTrap() {
	fmt.Println("\n--- 示例4: 闭包陷阱 ---")

	var wg sync.WaitGroup

	// ❌ 错误写法：所有 goroutine 可能打印相同的值
	fmt.Println("  错误写法（闭包捕获变量地址）:")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Printf("    i = %d\n", i) // i 可能都是 3
		}()
	}
	wg.Wait()

	// ✅ 正确写法1：通过参数传值
	fmt.Println("  正确写法1（参数传值）:")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			fmt.Printf("    i = %d\n", val)
		}(i)
	}
	wg.Wait()

	// ✅ 正确写法2：在循环内创建局部变量
	fmt.Println("  正确写法2（局部变量）:")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		i := i // 创建新的局部变量
		go func() {
			defer wg.Done()
			fmt.Printf("    i = %d\n", i)
		}()
	}
	wg.Wait()
}
