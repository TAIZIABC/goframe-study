package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	gmrender "github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"

	"md2wx/render"
)

func main() {
	input := flag.String("file", "", "Markdown 文件路径")
	output := flag.String("out", "", "输出 HTML 文件 (不指定则打印到终端)")
	preview := flag.Bool("preview", false, "生成完整 HTML 预览页面 (含 <html> 包裹)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: md2wx -file <markdown文件> [选项]\n\n")
		fmt.Fprintf(os.Stderr, "将 Markdown 转换为微信公众号编辑器兼容的 HTML（内联 CSS）。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  md2wx -file article.md -out article.html\n")
		fmt.Fprintf(os.Stderr, "  md2wx -file article.md -preview -out preview.html\n")
		fmt.Fprintf(os.Stderr, "  cat article.md | md2wx\n")
	}
	flag.Parse()

	// 读取输入
	var mdContent []byte
	var err error

	if *input != "" {
		mdContent, err = os.ReadFile(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ 读取文件失败: %v\n", err)
			os.Exit(1)
		}
	} else if flag.NArg() > 0 {
		mdContent, err = os.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ 读取文件失败: %v\n", err)
			os.Exit(1)
		}
	} else {
		// 从 stdin 读取
		mdContent, err = readStdin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "✗ 读取 stdin 失败: %v\n", err)
			os.Exit(1)
		}
	}

	if len(bytes.TrimSpace(mdContent)) == 0 {
		fmt.Fprintln(os.Stderr, "✗ 输入为空")
		os.Exit(1)
	}

	// 创建 goldmark，使用自定义渲染器
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithRenderer(
			gmrender.NewRenderer(
				gmrender.WithNodeRenderers(
					util.Prioritized(render.NewWxRenderer(), 100),
					util.Prioritized(render.NewWxTableRenderer(), 100),
				),
			),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(mdContent, &buf); err != nil {
		fmt.Fprintf(os.Stderr, "✗ 转换失败: %v\n", err)
		os.Exit(1)
	}

	result := buf.String()

	// 预览模式：包裹完整 HTML 页面
	if *preview {
		result = wrapPreviewHTML(result)
	}

	// 输出
	if *output != "" {
		if err := os.WriteFile(*output, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "✗ 写入失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "✓ 已输出到: %s\n", *output)
	} else {
		fmt.Print(result)
	}
}

func readStdin() ([]byte, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return nil, fmt.Errorf("没有输入数据，请用管道或指定 -file")
	}
	return os.ReadFile("/dev/stdin")
}

func wrapPreviewHTML(body string) string {
	return `<!DOCTYPE html>
<html><head><meta charset="UTF-8"><title>微信公众号预览</title>
<style>
body{max-width:680px;margin:40px auto;padding:0 20px;background:#f5f5f5}
.wx-content{background:#fff;padding:30px;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,.08)}
</style></head>
<body><div class="wx-content">` + body + `</div></body></html>`
}
