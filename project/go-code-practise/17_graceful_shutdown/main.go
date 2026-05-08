package main

import (
	"context"
	"fmt"
	"time"
)

// 题目：实现优雅关闭的后台任务管理器
//
// TaskManager 可以启动后台 worker，收到 shutdown 信号后：
// 1. 停止接受新任务
// 2. 等待正在进行的任务完成（设置最大等待时间）
// 3. 超时则强制退出
//
// 接口：
//   NewTaskManager() *TaskManager
//   Submit(task func()) — 提交任务
//   Shutdown(timeout time.Duration) — 优雅关闭

type TaskManager struct {
	// 在此定义字段
}

func NewTaskManager(workerCount int) *TaskManager {
	// 在此实现：启动 workerCount 个 worker
	return nil
}

func (tm *TaskManager) Submit(task func()) {
	// 在此实现：提交任务到队列
}

func (tm *TaskManager) Shutdown(timeout time.Duration) {
	// 在此实现：优雅关闭
}

func main() {
	tm := NewTaskManager(3)

	// 提交 5 个任务
	for i := 0; i < 5; i++ {
		i := i
		tm.Submit(func() {
			time.Sleep(50 * time.Millisecond)
			fmt.Printf("task %d done\n", i)
		})
	}

	time.Sleep(10 * time.Millisecond)
	tm.Shutdown(2 * time.Second)
	fmt.Println("all done")

	_ = context.Background()
}
