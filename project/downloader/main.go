package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Task struct {
	URL      string
	FileName string
	Index    int
}

type Result struct {
	Task    Task
	Success bool
	Size    int64
	Elapsed time.Duration
	Err     error
}

// progress 跟踪单个文件的下载进度
type progress struct {
	total   int64
	current int64
}

func main() {
	urlFile := flag.String("file", "", "URL 列表文件 (每行一个 URL)")
	outDir := flag.String("out", "./downloads", "下载输出目录")
	workers := flag.Int("workers", 3, "并发下载数")
	retries := flag.Int("retry", 2, "失败重试次数")
	timeout := flag.Duration("timeout", 60*time.Second, "单个文件下载超时")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: downloader [选项] [URL...]\n\n")
		fmt.Fprintf(os.Stderr, "批量文件下载工具，支持并发下载、自动重试、进度显示。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  downloader -file urls.txt -workers 5\n")
		fmt.Fprintf(os.Stderr, "  downloader https://example.com/a.zip https://example.com/b.zip\n")
		fmt.Fprintf(os.Stderr, "  downloader -file urls.txt -out ./data -retry 3 -timeout 120s\n")
	}
	flag.Parse()

	// 收集 URL
	var urls []string

	if *urlFile != "" {
		f, err := os.Open(*urlFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ 无法打开文件 %s: %v\n", *urlFile, err)
			os.Exit(1)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" && !strings.HasPrefix(line, "#") {
				urls = append(urls, line)
			}
		}
	}

	// 命令行直接传入的 URL
	urls = append(urls, flag.Args()...)

	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "✗ 没有要下载的 URL，使用 -file 指定文件或直接传入 URL")
		flag.Usage()
		os.Exit(1)
	}

	// 创建输出目录
	if err := os.MkdirAll(*outDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "✗ 创建目录失败: %v\n", err)
		os.Exit(1)
	}

	// 构建任务
	tasks := make([]Task, len(urls))
	for i, u := range urls {
		tasks[i] = Task{
			URL:      u,
			FileName: extractFileName(u, i),
			Index:    i + 1,
		}
	}

	fmt.Println()
	fmt.Printf("  下载任务: %d 个文件\n", len(tasks))
	fmt.Printf("  并发数:   %d\n", *workers)
	fmt.Printf("  重试次数: %d\n", *retries)
	fmt.Printf("  输出目录: %s\n", *outDir)
	fmt.Println("  " + strings.Repeat("─", 50))

	// 并发下载
	var completed int32
	totalTasks := int32(len(tasks))

	taskCh := make(chan Task, len(tasks))
	resultCh := make(chan Result, len(tasks))

	var wg sync.WaitGroup
	for i := 0; i < *workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			client := &http.Client{Timeout: *timeout}
			for task := range taskCh {
				result := downloadWithRetry(client, task, *outDir, *retries)
				done := atomic.AddInt32(&completed, 1)
				printResult(result, done, totalTasks)
				resultCh <- result
			}
		}()
	}

	startTime := time.Now()

	for _, t := range tasks {
		taskCh <- t
	}
	close(taskCh)

	wg.Wait()
	close(resultCh)

	// 汇总结果
	var successCount, failCount int
	var totalSize int64
	var failures []Result

	for r := range resultCh {
		if r.Success {
			successCount++
			totalSize += r.Size
		} else {
			failCount++
			failures = append(failures, r)
		}
	}

	elapsed := time.Since(startTime)
	fmt.Println("\n  " + strings.Repeat("─", 50))
	fmt.Printf("  下载完成  耗时 %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("  成功: \033[32m%d\033[0m  失败: \033[31m%d\033[0m  总大小: %s\n",
		successCount, failCount, humanSize(totalSize))

	if len(failures) > 0 {
		fmt.Println("\n  失败列表:")
		for _, f := range failures {
			fmt.Printf("    ✗ %s\n      %v\n", f.Task.URL, f.Err)
		}
	}
	fmt.Println()

	if failCount > 0 {
		os.Exit(1)
	}
}

func downloadWithRetry(client *http.Client, task Task, outDir string, maxRetries int) Result {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			fmt.Printf("  ↻ [%d] 重试 %d/%d: %s\n", task.Index, attempt, maxRetries, task.FileName)
			time.Sleep(time.Duration(attempt) * time.Second) // 递增等待
		}

		result := download(client, task, outDir)
		if result.Success {
			return result
		}
		lastErr = result.Err
	}

	return Result{Task: task, Success: false, Err: lastErr}
}

func download(client *http.Client, task Task, outDir string) Result {
	start := time.Now()

	resp, err := client.Get(task.URL)
	if err != nil {
		return Result{Task: task, Err: fmt.Errorf("请求失败: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{Task: task, Err: fmt.Errorf("HTTP %d", resp.StatusCode)}
	}

	filePath := filepath.Join(outDir, task.FileName)
	out, err := os.Create(filePath)
	if err != nil {
		return Result{Task: task, Err: fmt.Errorf("创建文件失败: %w", err)}
	}
	defer out.Close()

	// 带进度的写入
	total := resp.ContentLength
	pw := &progressWriter{
		index:    task.Index,
		fileName: task.FileName,
		total:    total,
		writer:   out,
	}

	written, err := io.Copy(pw, resp.Body)
	pw.finish()

	if err != nil {
		os.Remove(filePath)
		return Result{Task: task, Err: fmt.Errorf("下载中断: %w", err)}
	}

	return Result{
		Task:    task,
		Success: true,
		Size:    written,
		Elapsed: time.Since(start),
	}
}

// progressWriter 实时显示下载进度
type progressWriter struct {
	index    int
	fileName string
	total    int64
	current  int64
	writer   io.Writer
	lastPrint time.Time
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	pw.current += int64(n)

	if time.Since(pw.lastPrint) > 200*time.Millisecond {
		pw.lastPrint = time.Now()
		pw.printProgress()
	}

	return n, err
}

func (pw *progressWriter) printProgress() {
	if pw.total > 0 {
		pct := float64(pw.current) / float64(pw.total) * 100
		bar := progressBar(pct, 25)
		fmt.Printf("\r  ↓ [%d] %s %s %s/%s %.0f%%",
			pw.index, truncName(pw.fileName, 18), bar,
			humanSize(pw.current), humanSize(pw.total), pct)
	} else {
		fmt.Printf("\r  ↓ [%d] %s %s",
			pw.index, truncName(pw.fileName, 18), humanSize(pw.current))
	}
}

func (pw *progressWriter) finish() {
	fmt.Print("\r" + strings.Repeat(" ", 80) + "\r") // 清除进度行
}

func printResult(r Result, done, total int32) {
	if r.Success {
		fmt.Printf("  \033[32m✓\033[0m [%d/%d] %s  %s  %s\n",
			done, total, truncName(r.Task.FileName, 24),
			humanSize(r.Size), r.Elapsed.Round(time.Millisecond))
	} else {
		fmt.Printf("  \033[31m✗\033[0m [%d/%d] %s  %v\n",
			done, total, truncName(r.Task.FileName, 24), r.Err)
	}
}

func progressBar(pct float64, width int) string {
	filled := int(pct / 100 * float64(width))
	if filled > width {
		filled = width
	}
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", width-filled) + "]"
}

func extractFileName(rawURL string, index int) string {
	u, err := url.Parse(rawURL)
	if err == nil {
		name := filepath.Base(u.Path)
		if name != "" && name != "." && name != "/" {
			return name
		}
	}
	return fmt.Sprintf("download_%d", index+1)
}

func truncName(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s + strings.Repeat(" ", maxLen-len(s))
	}
	return s[:maxLen-3] + "..."
}

func humanSize(b int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case b >= GB:
		return fmt.Sprintf("%.2f GB", float64(b)/float64(GB))
	case b >= MB:
		return fmt.Sprintf("%.2f MB", float64(b)/float64(MB))
	case b >= KB:
		return fmt.Sprintf("%.2f KB", float64(b)/float64(KB))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
