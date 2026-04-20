package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"gocron/api"
	"gocron/config"
	"gocron/scheduler"
)

func main() {
	cfgFile := flag.String("config", "tasks.yaml", "任务配置文件路径")
	httpPort := flag.Int("port", 8089, "HTTP 管理接口端口")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: gocron [选项]\n\n")
		fmt.Fprintf(os.Stderr, "基于 YAML 配置的定时任务调度器。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  gocron -config tasks.yaml -port 8089\n")
	}
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	// 创建调度器
	sched, err := scheduler.New(cfg.Tasks, func(format string, args ...interface{}) {
		fmt.Printf(format, args...)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("  ╔════════════════════════════════════════╗")
	fmt.Println("  ║        GoCron 定时任务调度器           ║")
	fmt.Println("  ╠════════════════════════════════════════╣")
	fmt.Printf("  ║  📋 配置: %s\n", *cfgFile)
	fmt.Printf("  ║  📊 任务数: %d\n", len(cfg.Tasks))
	fmt.Printf("  ║  🌐 管理面板: http://localhost:%d\n", *httpPort)
	fmt.Printf("  ║  📡 API: http://localhost:%d/api/tasks\n", *httpPort)
	fmt.Println("  ║  按 Ctrl+C 停止")
	fmt.Println("  ╠════════════════════════════════════════╣")

	tasks := sched.GetTasks()
	for _, t := range tasks {
		fmt.Printf("  ║  ⏱️  %-14s  %s\n", t.Name, t.Cron)
	}
	fmt.Println("  ╚════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("  " + strings.Repeat("─", 60))

	// 启动调度
	go sched.Start()

	// 启动 HTTP
	go func() {
		if err := api.StartAPI(*httpPort, sched); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ HTTP 启动失败: %v\n", err)
		}
	}()

	// 等待信号
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("\n  调度器正在停止...")
	sched.Stop()
	fmt.Println("  再见！")
}
