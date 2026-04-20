package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	indent := flag.Int("indent", 2, "缩进空格数")
	compact := flag.Bool("compact", false, "压缩模式（去除空白）")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: jsonfmt [选项] [文件...]\n\n")
		fmt.Fprintf(os.Stderr, "从文件或 stdin 读取 JSON 并格式化输出。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  echo '{\"a\":1}' | jsonfmt\n")
		fmt.Fprintf(os.Stderr, "  jsonfmt data.json\n")
		fmt.Fprintf(os.Stderr, "  jsonfmt -indent 4 data.json\n")
		fmt.Fprintf(os.Stderr, "  jsonfmt -compact data.json\n")
	}
	flag.Parse()

	indentStr := strings.Repeat(" ", *indent)

	// 确定输入源：文件参数 or stdin
	var inputs []namedReader
	if flag.NArg() > 0 {
		for _, path := range flag.Args() {
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "✗ 无法打开文件 %s: %v\n", path, err)
				os.Exit(1)
			}
			defer f.Close()
			inputs = append(inputs, namedReader{name: path, r: f})
		}
	} else {
		inputs = append(inputs, namedReader{name: "stdin", r: os.Stdin})
	}

	exitCode := 0
	for i, input := range inputs {
		if len(inputs) > 1 {
			if i > 0 {
				fmt.Println()
			}
			fmt.Fprintf(os.Stderr, "── %s ──\n", input.name)
		}

		data, err := io.ReadAll(input.r)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ 读取 %s 失败: %v\n", input.name, err)
			exitCode = 1
			continue
		}

		data = bytes.TrimSpace(data)
		if len(data) == 0 {
			fmt.Fprintf(os.Stderr, "✗ %s: 输入为空\n", input.name)
			exitCode = 1
			continue
		}

		// 先验证 JSON 合法性
		if err := validateJSON(data, input.name); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
			continue
		}

		// 格式化
		var buf bytes.Buffer
		if *compact {
			if err := json.Compact(&buf, data); err != nil {
				fmt.Fprintf(os.Stderr, "✗ %s: 压缩失败: %v\n", input.name, err)
				exitCode = 1
				continue
			}
		} else {
			if err := json.Indent(&buf, data, "", indentStr); err != nil {
				fmt.Fprintf(os.Stderr, "✗ %s: 格式化失败: %v\n", input.name, err)
				exitCode = 1
				continue
			}
		}

		buf.WriteByte('\n')
		os.Stdout.Write(buf.Bytes())
	}

	os.Exit(exitCode)
}

type namedReader struct {
	name string
	r    io.Reader
}

// validateJSON 校验 JSON 合法性，出错时定位行号和列号
func validateJSON(data []byte, name string) error {
	var v interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	if err := decoder.Decode(&v); err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			line, col := locateOffset(data, syntaxErr.Offset)
			context := extractContext(data, syntaxErr.Offset)
			return fmt.Errorf(
				"✗ %s: JSON 语法错误 (第 %d 行, 第 %d 列)\n"+
					"  错误: %s\n"+
					"  位置: ...%s...",
				name, line, col, syntaxErr.Error(), context,
			)
		}
		if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			line, col := locateOffset(data, unmarshalErr.Offset)
			return fmt.Errorf(
				"✗ %s: JSON 类型错误 (第 %d 行, 第 %d 列)\n  错误: %s",
				name, line, col, unmarshalErr.Error(),
			)
		}
		return fmt.Errorf("✗ %s: JSON 解析失败: %v", name, err)
	}
	return nil
}

// locateOffset 根据字节偏移量计算行号和列号
func locateOffset(data []byte, offset int64) (line, col int) {
	line = 1
	col = 1
	for i := int64(0); i < offset && i < int64(len(data)); i++ {
		if data[i] == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	return
}

// extractContext 提取出错位置前后的上下文片段
func extractContext(data []byte, offset int64) string {
	start := offset - 20
	if start < 0 {
		start = 0
	}
	end := offset + 20
	if end > int64(len(data)) {
		end = int64(len(data))
	}

	snippet := string(data[start:end])
	// 在出错位置插入标记
	markerPos := offset - start
	if markerPos >= 0 && markerPos <= int64(len(snippet)) {
		snippet = snippet[:markerPos] + "⚡" + snippet[markerPos:]
	}

	snippet = strings.ReplaceAll(snippet, "\n", "\\n")
	snippet = strings.ReplaceAll(snippet, "\t", "\\t")
	return snippet
}
