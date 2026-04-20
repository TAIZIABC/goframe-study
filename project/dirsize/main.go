package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Entry struct {
	Name  string
	Size  int64
	IsDir bool
}

// calcDirSize 递归计算目录总大小
func calcDirSize(path string) (int64, error) {
	var total int64
	err := filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			total += info.Size()
		}
		return nil
	})
	return total, err
}

// humanSize 将字节数转为人类可读格式
func humanSize(b int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)
	switch {
	case b >= TB:
		return fmt.Sprintf("%.2f TB", float64(b)/float64(TB))
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

func main() {
	dir := flag.String("dir", ".", "要扫描的目录路径")
	flag.Parse()

	// 也支持直接传参: dirsize /path/to/dir
	if flag.NArg() > 0 {
		*dir = flag.Arg(0)
	}

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}

	info, err := os.Stat(absDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "错误: %s 不是一个目录\n", absDir)
		os.Exit(1)
	}

	entries, err := os.ReadDir(absDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法读取目录: %v\n", err)
		os.Exit(1)
	}

	var items []Entry
	var totalSize int64

	for _, e := range entries {
		fullPath := filepath.Join(absDir, e.Name())
		var size int64

		if e.IsDir() {
			size, err = calcDirSize(fullPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "  警告: 计算 %s 大小失败: %v\n", e.Name(), err)
				continue
			}
		} else {
			fi, err := e.Info()
			if err != nil {
				fmt.Fprintf(os.Stderr, "  警告: 获取 %s 信息失败: %v\n", e.Name(), err)
				continue
			}
			size = fi.Size()
		}

		items = append(items, Entry{Name: e.Name(), Size: size, IsDir: e.IsDir()})
		totalSize += size
	}

	// 按大小降序排列
	sort.Slice(items, func(i, j int) bool {
		return items[i].Size > items[j].Size
	})

	// 输出
	fmt.Printf("目录: %s\n", absDir)
	fmt.Println(strings.Repeat("─", 56))
	fmt.Printf("  %-36s %10s  %s\n", "名称", "大小", "类型")
	fmt.Println(strings.Repeat("─", 56))

	for _, item := range items {
		typeMark := "📄"
		name := item.Name
		if item.IsDir {
			typeMark = "📁"
			name += "/"
		}
		fmt.Printf("  %-36s %10s  %s\n", name, humanSize(item.Size), typeMark)
	}

	fmt.Println(strings.Repeat("─", 56))
	fmt.Printf("  合计: %d 项, 总大小 %s\n", len(items), humanSize(totalSize))
}
