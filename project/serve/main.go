package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	dir := flag.String("dir", ".", "要服务的目录路径")
	port := flag.Int("port", 8080, "监听端口号")
	host := flag.String("host", "0.0.0.0", "监听地址")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: serve [选项]\n\n")
		fmt.Fprintf(os.Stderr, "启动一个 HTTP 静态文件服务器，支持目录浏览。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  serve                          # 当前目录, 端口 8080\n")
		fmt.Fprintf(os.Stderr, "  serve -dir ./dist -port 3000   # 指定目录和端口\n")
	}
	flag.Parse()

	// 也支持第一个参数直接传目录
	if flag.NArg() > 0 {
		*dir = flag.Arg(0)
	}

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ 路径解析失败: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(absDir)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "✗ %s 不是一个有效的目录\n", absDir)
		os.Exit(1)
	}

	addr := fmt.Sprintf("%s:%d", *host, *port)

	// 带日志的文件服务
	handler := logMiddleware(http.FileServer(http.Dir(absDir)))

	fmt.Println()
	fmt.Println("  ┌──────────────────────────────────────────┐")
	fmt.Println("  │        静态文件服务器 (serve)             │")
	fmt.Println("  ├──────────────────────────────────────────┤")
	fmt.Printf("  │  📁 目录:  %s\n", absDir)
	fmt.Printf("  │  🔗 地址:  http://localhost:%d\n", *port)
	if *host == "0.0.0.0" {
		if ip := getLocalIP(); ip != "" {
			fmt.Printf("  │  🌐 局域网: http://%s:%d\n", ip, *port)
		}
	}
	fmt.Println("  │")
	fmt.Println("  │  按 Ctrl+C 停止服务")
	fmt.Println("  └──────────────────────────────────────────┘")
	fmt.Println()

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintf(os.Stderr, "✗ 服务启动失败: %v\n", err)
		os.Exit(1)
	}
}

// logMiddleware 请求日志中间件
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		statusColor := "\033[32m" // 绿色
		if rw.status >= 400 {
			statusColor = "\033[31m" // 红色
		} else if rw.status >= 300 {
			statusColor = "\033[33m" // 黄色
		}

		fmt.Printf("  %s %s%d\033[0m %s%-6s\033[0m %s  %s\n",
			time.Now().Format("15:04:05"),
			statusColor, rw.status,
			methodColor(r.Method), r.Method,
			r.URL.Path,
			duration.Round(time.Microsecond),
		)
	})
}

func methodColor(method string) string {
	switch method {
	case "GET":
		return "\033[36m" // 青色
	case "POST":
		return "\033[33m" // 黄色
	default:
		return "\033[0m"
	}
}

// responseWriter 包装 ResponseWriter 以捕获状态码
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// getLocalIP 获取本机局域网 IP
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip := ipNet.IP.String()
			if strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "172.") {
				return ip
			}
		}
	}
	return ""
}
