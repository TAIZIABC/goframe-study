// Package main 演示 sync 包中的并发原语
//
// 知识点：
// 1. sync.WaitGroup：等待一组 goroutine 完成
// 2. sync.Mutex：互斥锁，保护共享资源
// 3. sync.RWMutex：读写锁，允许多个读操作并发
// 4. sync.Once：确保某个操作只执行一次（单例模式）
// 5. sync.Map：并发安全的 map
// 6. 竞态条件演示与修复
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	fmt.Println("===== Sync 并发原语示例 =====")

	waitGroupDemo()
	mutexDemo()
	rwMutexDemo()
	onceDemo()
	syncMapDemo()
	atomicDemo()
}

// waitGroupDemo 演示 WaitGroup 的使用
func waitGroupDemo() {
	fmt.Println("\n--- 示例1: WaitGroup ---")

	var wg sync.WaitGroup

	tasks := []string{"下载文件", "解析数据", "生成报告", "发送邮件", "清理缓存"}

	for _, task := range tasks {
		wg.Add(1)
		go func(t string) {
			defer wg.Done() // ✅ 使用 defer 确保 Done 一定会被调用
			fmt.Printf("  开始: %s\n", t)
			time.Sleep(20 * time.Millisecond)
			fmt.Printf("  完成: %s\n", t)
		}(task)
	}

	wg.Wait() // 阻塞，直到所有任务完成
	fmt.Println("  ✅ 所有任务完成")
}

// mutexDemo 演示互斥锁保护共享资源
func mutexDemo() {
	fmt.Println("\n--- 示例2: Mutex 互斥锁 ---")

	// ❌ 没有锁的情况：存在竞态条件
	var unsafeCount int
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			unsafeCount++ // 竞态条件！
		}()
	}
	wg.Wait()
	fmt.Printf("  无锁计数器（可能不正确）: %d\n", unsafeCount)

	// ✅ 使用 Mutex 保护
	var safeCount int
	var mu sync.Mutex
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			safeCount++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("  有锁计数器（一定正确）: %d\n", safeCount)
}

// rwMutexDemo 演示读写锁
// 适用场景：读多写少的情况，允许多个读操作并发执行
func rwMutexDemo() {
	fmt.Println("\n--- 示例3: RWMutex 读写锁 ---")

	var (
		data   = make(map[string]string)
		rwmu   sync.RWMutex
		wg     sync.WaitGroup
	)

	// 写操作（独占锁）
	write := func(key, value string) {
		defer wg.Done()
		rwmu.Lock()
		defer rwmu.Unlock()
		time.Sleep(5 * time.Millisecond) // 模拟写入耗时
		data[key] = value
		fmt.Printf("  写入: %s = %s\n", key, value)
	}

	// 读操作（共享锁，多个读可以并发）
	read := func(key string) {
		defer wg.Done()
		rwmu.RLock()
		defer rwmu.RUnlock()
		val := data[key]
		fmt.Printf("  读取: %s = %s\n", key, val)
	}

	// 先写入数据
	wg.Add(2)
	go write("name", "Go")
	go write("version", "1.22")
	wg.Wait()

	// 并发读取（不会互相阻塞）
	wg.Add(4)
	go read("name")
	go read("version")
	go read("name")
	go read("version")
	wg.Wait()
}

// onceDemo 演示 sync.Once（只执行一次）
func onceDemo() {
	fmt.Println("\n--- 示例4: sync.Once 单例 ---")

	var once sync.Once
	var wg sync.WaitGroup

	initConfig := func() {
		fmt.Println("  📋 配置初始化（只会执行一次）")
	}

	// 启动多个 goroutine 尝试初始化
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("  Goroutine %d 尝试初始化\n", id)
			once.Do(initConfig) // 只有第一个到达的会执行
		}(i)
	}
	wg.Wait()
}

// syncMapDemo 演示并发安全的 map
func syncMapDemo() {
	fmt.Println("\n--- 示例5: sync.Map ---")

	var m sync.Map
	var wg sync.WaitGroup

	// 并发写入
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", id)
			m.Store(key, id*10)
		}(i)
	}
	wg.Wait()

	// 遍历
	fmt.Println("  sync.Map 内容:")
	m.Range(func(key, value any) bool {
		fmt.Printf("    %s -> %v\n", key, value)
		return true // 返回 false 停止遍历
	})

	// LoadOrStore：不存在则存入，存在则返回已有值
	val, loaded := m.LoadOrStore("key_0", 999)
	fmt.Printf("  LoadOrStore key_0: val=%v, 已存在=%v\n", val, loaded)
}

// atomicDemo 演示原子操作（比 Mutex 更轻量）
func atomicDemo() {
	fmt.Println("\n--- 示例6: atomic 原子操作 ---")

	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1) // 原子自增
		}()
	}
	wg.Wait()
	fmt.Printf("  原子计数器: %d\n", atomic.LoadInt64(&counter))

	// CompareAndSwap：CAS 操作
	var flag int32
	swapped := atomic.CompareAndSwapInt32(&flag, 0, 1) // 如果是0则改为1
	fmt.Printf("  CAS 操作: swapped=%v, flag=%d\n", swapped, flag)
}
