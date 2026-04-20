// scheduler/scheduler.go
// 核心调度器：并发执行、超时控制、执行日志
package scheduler

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"gocron/config"
	"gocron/cron"
)

// TaskStatus 任务状态
type TaskStatus string

const (
	StatusIdle    TaskStatus = "idle"
	StatusRunning TaskStatus = "running"
)

// ExecRecord 执行记录
type ExecRecord struct {
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
	Output    string        `json:"output"`
	Error     string        `json:"error,omitempty"`
	Success   bool          `json:"success"`
}

// Task 运行时任务
type Task struct {
	Config     config.TaskConfig `json:"config"`
	Schedule   *cron.Schedule    `json:"-"`
	Status     TaskStatus        `json:"status"`
	NextRun    time.Time         `json:"next_run"`
	LastRun    time.Time         `json:"last_run,omitempty"`
	RunCount   int               `json:"run_count"`
	FailCount  int               `json:"fail_count"`
	History    []ExecRecord      `json:"history"`
	mu         sync.RWMutex
}

// Scheduler 调度器
type Scheduler struct {
	tasks   []*Task
	stop    chan struct{}
	logFunc func(format string, args ...interface{})
	mu      sync.RWMutex
}

// New 创建调度器
func New(configs []config.TaskConfig, logFunc func(string, ...interface{})) (*Scheduler, error) {
	s := &Scheduler{
		stop:    make(chan struct{}),
		logFunc: logFunc,
	}

	for _, cfg := range configs {
		schedule, err := cron.Parse(cfg.Cron)
		if err != nil {
			return nil, fmt.Errorf("任务 [%s] cron 解析失败: %w", cfg.Name, err)
		}

		task := &Task{
			Config:   cfg,
			Schedule: schedule,
			Status:   StatusIdle,
			NextRun:  schedule.Next(time.Now()),
		}
		s.tasks = append(s.tasks, task)
	}

	return s, nil
}

// Start 启动调度
func (s *Scheduler) Start() {
	s.logFunc("  调度器启动，%d 个任务\n", len(s.tasks))

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			s.logFunc("  调度器已停止\n")
			return
		case now := <-ticker.C:
			now = now.Truncate(time.Second)
			for _, task := range s.tasks {
				task.mu.RLock()
				shouldRun := !task.NextRun.After(now) && task.Status == StatusIdle
				task.mu.RUnlock()

				if shouldRun {
					go s.executeTask(task)
				}
			}
		}
	}
}

// Stop 停止调度
func (s *Scheduler) Stop() {
	close(s.stop)
}

// GetTasks 获取所有任务状态（线程安全）
func (s *Scheduler) GetTasks() []TaskInfo {
	var infos []TaskInfo
	for _, t := range s.tasks {
		t.mu.RLock()
		info := TaskInfo{
			Name:      t.Config.Name,
			Cron:      t.Config.Cron,
			Command:   t.Config.Command,
			Timeout:   t.Config.Timeout.String(),
			Status:    string(t.Status),
			NextRun:   t.NextRun.Format("15:04:05"),
			RunCount:  t.RunCount,
			FailCount: t.FailCount,
		}
		if !t.LastRun.IsZero() {
			info.LastRun = t.LastRun.Format("15:04:05")
		}
		t.mu.RUnlock()
		infos = append(infos, info)
	}
	return infos
}

// GetHistory 获取任务执行历史
func (s *Scheduler) GetHistory(name string, limit int) []ExecRecord {
	for _, t := range s.tasks {
		if t.Config.Name == name {
			t.mu.RLock()
			defer t.mu.RUnlock()

			h := t.History
			if limit > 0 && limit < len(h) {
				h = h[len(h)-limit:]
			}

			// 倒序返回
			result := make([]ExecRecord, len(h))
			for i, r := range h {
				result[len(h)-1-i] = r
			}
			return result
		}
	}
	return nil
}

// TaskInfo HTTP 接口返回的任务信息
type TaskInfo struct {
	Name      string `json:"name"`
	Cron      string `json:"cron"`
	Command   string `json:"command"`
	Timeout   string `json:"timeout"`
	Status    string `json:"status"`
	NextRun   string `json:"next_run"`
	LastRun   string `json:"last_run,omitempty"`
	RunCount  int    `json:"run_count"`
	FailCount int    `json:"fail_count"`
}

func (s *Scheduler) executeTask(task *Task) {
	task.mu.Lock()
	task.Status = StatusRunning
	task.mu.Unlock()

	start := time.Now()
	record := ExecRecord{StartTime: start}

	s.logFunc("  \033[36m▶\033[0m [%s] %s 开始执行\n",
		start.Format("15:04:05"), task.Config.Name)

	// 带超时的命令执行
	ctx, cancel := context.WithTimeout(context.Background(), task.Config.Timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", task.Config.Command)
	output, err := cmd.CombinedOutput()

	record.EndTime = time.Now()
	record.Duration = record.EndTime.Sub(start)
	record.Output = strings.TrimSpace(string(output))

	if err != nil {
		record.Success = false
		if ctx.Err() == context.DeadlineExceeded {
			record.Error = fmt.Sprintf("超时 (%s)", task.Config.Timeout)
		} else {
			record.Error = err.Error()
		}
		s.logFunc("  \033[31m✗\033[0m [%s] %s 失败: %s (%s)\n",
			record.EndTime.Format("15:04:05"), task.Config.Name,
			record.Error, record.Duration.Round(time.Millisecond))
	} else {
		record.Success = true
		outputPreview := record.Output
		if len(outputPreview) > 60 {
			outputPreview = outputPreview[:60] + "..."
		}
		s.logFunc("  \033[32m✓\033[0m [%s] %s 完成 (%s) → %s\n",
			record.EndTime.Format("15:04:05"), task.Config.Name,
			record.Duration.Round(time.Millisecond), outputPreview)
	}

	task.mu.Lock()
	task.Status = StatusIdle
	task.LastRun = start
	task.RunCount++
	if !record.Success {
		task.FailCount++
	}
	task.NextRun = task.Schedule.Next(time.Now())

	// 保留最近 100 条历史
	task.History = append(task.History, record)
	if len(task.History) > 100 {
		task.History = task.History[len(task.History)-100:]
	}
	task.mu.Unlock()
}
