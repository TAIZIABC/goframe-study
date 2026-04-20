package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

type Result struct {
	Port     string
	Open     bool
	Duration time.Duration
	Err      error
}

func scanPort(host, port string, timeout time.Duration) Result {
	address := net.JoinHostPort(host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	duration := time.Since(start)

	if err != nil {
		return Result{Port: port, Open: false, Duration: duration, Err: err}
	}
	conn.Close()
	return Result{Port: port, Open: true, Duration: duration}
}

func main() {
	host := flag.String("host", "", "目标主机 (IP 或域名)")
	ports := flag.String("ports", "", "端口列表，逗号分隔 (如: 22,80,443,3306,8080)")
	timeout := flag.Duration("timeout", 3*time.Second, "单个端口连接超时时间")
	flag.Parse()

	if *host == "" || *ports == "" {
		fmt.Println("用法: portscan -host <主机> -ports <端口列表>")
		fmt.Println()
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("  portscan -host 127.0.0.1 -ports 22,80,443,3306,8080")
		fmt.Println("  portscan -host example.com -ports 80,443 -timeout 5s")
		os.Exit(1)
	}

	portList := strings.Split(*ports, ",")
	for i := range portList {
		portList[i] = strings.TrimSpace(portList[i])
	}

	fmt.Printf("扫描目标: %s\n", *host)
	fmt.Printf("端口列表: %s\n", strings.Join(portList, ", "))
	fmt.Printf("超时时间: %s\n", *timeout)
	fmt.Println(strings.Repeat("─", 50))

	var wg sync.WaitGroup
	results := make([]Result, len(portList))

	for i, port := range portList {
		wg.Add(1)
		go func(idx int, p string) {
			defer wg.Done()
			results[idx] = scanPort(*host, p, *timeout)
		}(i, port)
	}

	wg.Wait()

	// 按端口号排序输出
	sort.Slice(results, func(i, j int) bool {
		return results[i].Port < results[j].Port
	})

	openCount := 0
	for _, r := range results {
		status := "\033[31mclosed\033[0m"
		if r.Open {
			status = "\033[32mopen\033[0m"
			openCount++
		}
		fmt.Printf("  端口 %-6s  %-16s  耗时 %s\n", r.Port, status, r.Duration.Round(time.Microsecond))
	}

	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("扫描完成: %d 个端口, %d open, %d closed\n",
		len(results), openCount, len(results)-openCount)
}
