package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	nameFile := flag.String("file", "", "候选名单文件 (每行一个名字)")
	count := flag.Int("n", 1, "抽取人数")
	logFile := flag.String("log", "lottery.log", "日志文件路径")
	title := flag.String("title", "", "本次抽签标题 (如: 周会分享人)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: lottery -file <名单文件> -n <人数> [选项]\n\n")
		fmt.Fprintf(os.Stderr, "从候选名单中随机抽取指定数量的名字，结果去重并记录日志。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  lottery -file names.txt -n 3\n")
		fmt.Fprintf(os.Stderr, "  lottery -file names.txt -n 2 -title \"周会分享人\" -log record.log\n")
	}
	flag.Parse()

	if *nameFile == "" {
		fmt.Fprintln(os.Stderr, "✗ 请指定候选名单文件: -file <文件路径>")
		flag.Usage()
		os.Exit(1)
	}

	// 读取候选名单
	names, err := readNames(*nameFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "✗ 读取名单失败: %v\n", err)
		os.Exit(1)
	}

	if len(names) == 0 {
		fmt.Fprintln(os.Stderr, "✗ 名单为空")
		os.Exit(1)
	}

	if *count < 1 {
		fmt.Fprintln(os.Stderr, "✗ 抽取人数至少为 1")
		os.Exit(1)
	}

	if *count > len(names) {
		fmt.Fprintf(os.Stderr, "✗ 抽取人数 (%d) 超过候选人数 (%d)\n", *count, len(names))
		os.Exit(1)
	}

	// 随机抽取（Fisher-Yates 洗牌，天然去重）
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]string, len(names))
	copy(shuffled, names)
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	winners := shuffled[:*count]

	// 显示结果
	now := time.Now()
	titleStr := "抽签结果"
	if *title != "" {
		titleStr = *title
	}

	fmt.Println()
	fmt.Println("  ╔════════════════════════════════════════╗")
	fmt.Printf("  ║  🎲  %s\n", titleStr)
	fmt.Println("  ╠════════════════════════════════════════╣")
	fmt.Printf("  ║  候选人数: %d\n", len(names))
	fmt.Printf("  ║  抽取人数: %d\n", *count)
	fmt.Printf("  ║  时间: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Println("  ╠════════════════════════════════════════╣")
	for i, name := range winners {
		fmt.Printf("  ║  🎉 %d. %s\n", i+1, name)
	}
	fmt.Println("  ╚════════════════════════════════════════╝")
	fmt.Println()

	// 写入日志
	if err := appendLog(*logFile, now, titleStr, names, winners); err != nil {
		fmt.Fprintf(os.Stderr, "⚠ 日志写入失败: %v\n", err)
	} else {
		fmt.Printf("  📝 已记录到 %s\n\n", *logFile)
	}
}

// readNames 读取名单文件，自动去重和去空行
func readNames(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	seen := make(map[string]bool)
	var names []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name == "" || strings.HasPrefix(name, "#") {
			continue
		}
		if !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	return names, scanner.Err()
}

// appendLog 追加写入日志文件
func appendLog(path string, t time.Time, title string, candidates, winners []string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	var sb strings.Builder
	sb.WriteString("────────────────────────────────────────\n")
	sb.WriteString(fmt.Sprintf("标题: %s\n", title))
	sb.WriteString(fmt.Sprintf("时间: %s\n", t.Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("候选人数: %d\n", len(candidates)))
	sb.WriteString(fmt.Sprintf("抽取人数: %d\n", len(winners)))
	sb.WriteString(fmt.Sprintf("结果: %s\n", strings.Join(winners, ", ")))
	sb.WriteString("────────────────────────────────────────\n\n")

	_, err = f.WriteString(sb.String())
	return err
}
