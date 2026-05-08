package main

import (
	"fmt"
	"sync"
)

// 题目：实现并发安全的计数器
// 要求：
// 1. Counter 结构体，包含 Inc()、Dec()、Value() int 方法
// 2. 多个 goroutine 并发调用时数据正确
// 测试：启动 1000 个 goroutine 并发 Inc()，最终 Value() 应为 1000

type Counter struct {
	// 在此定义字段
	mu    sync.Mutex
	value int
}

func (c *Counter) Inc() {
	// 在此实现
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

func (c *Counter) Dec() {
	// 在此实现
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value--
}

func (c *Counter) Value() int {
	// 在此实现
	return c.value
}

func main() {
	c := &Counter{}
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc()
		}()
	}
	wg.Wait()
	fmt.Println(c.Value()) // 1000
}
