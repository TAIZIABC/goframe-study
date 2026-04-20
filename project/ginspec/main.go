package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"ginspec/generator"
	"ginspec/parser"
	"ginspec/swagger"
)

func main() {
	dir := flag.String("dir", ".", "Go 项目目录路径")
	output := flag.String("out", "", "输出文件路径 (如: api.json 或 api.yaml)")
	format := flag.String("format", "json", "输出格式: json 或 yaml")
	title := flag.String("title", "API Documentation", "API 文档标题")
	version := flag.String("version", "1.0.0", "API 版本号")
	serve := flag.Bool("serve", false, "启动 Swagger UI 服务")
	port := flag.Int("port", 8088, "Swagger UI 服务端口")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: ginspec [选项]\n\n")
		fmt.Fprintf(os.Stderr, "扫描 Gin 项目路由，自动生成 OpenAPI 3.0 文档。\n\n")
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  ginspec -dir ./my-api\n")
		fmt.Fprintf(os.Stderr, "  ginspec -dir ./my-api -out api.yaml -format yaml\n")
		fmt.Fprintf(os.Stderr, "  ginspec -dir ./my-api -serve -port 8088\n")
		fmt.Fprintf(os.Stderr, "  ginspec -dir ./my-api -title \"User API\" -version 2.0.0\n")
	}
	flag.Parse()

	// 也支持第一个参数直接传目录
	if flag.NArg() > 0 {
		*dir = flag.Arg(0)
	}

	fmt.Printf("\n  🔍 扫描项目: %s\n", *dir)

	// 解析路由
	routes, err := parser.ParseProject(*dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ 扫描失败: %v\n", err)
		os.Exit(1)
	}

	if len(routes) == 0 {
		fmt.Fprintln(os.Stderr, "  ⚠ 未发现 Gin 路由注册")
		os.Exit(0)
	}

	fmt.Printf("  ✓ 发现 %d 个路由\n\n", len(routes))

	// 打印路由列表
	fmt.Println("  " + strings.Repeat("─", 60))
	fmt.Printf("  %-8s %-28s %s\n", "方法", "路径", "Handler")
	fmt.Println("  " + strings.Repeat("─", 60))
	for _, r := range routes {
		summary := ""
		if r.Summary != "" {
			summary = "  // " + r.Summary
		}
		fmt.Printf("  \033[36m%-8s\033[0m %-28s %s%s\n", r.Method, r.Path, r.HandlerName, summary)
	}
	fmt.Println("  " + strings.Repeat("─", 60))
	fmt.Println()

	// 生成 OpenAPI 文档
	doc := generator.Generate(routes, *title, *version)

	// 输出到文件或 stdout
	var data []byte
	if strings.ToLower(*format) == "yaml" {
		data, err = doc.ToYAML()
	} else {
		data, err = doc.ToJSON()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ 生成文档失败: %v\n", err)
		os.Exit(1)
	}

	if *output != "" {
		if err := os.WriteFile(*output, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ 写入文件失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("  📄 文档已保存到: %s\n\n", *output)
	} else if !*serve {
		// 没有指定输出文件也没有 serve，打印到 stdout
		fmt.Println(string(data))
	}

	// 启动 Swagger UI
	if *serve {
		jsonData, _ := doc.ToJSON()
		fmt.Printf("  🌐 Swagger UI: http://localhost:%d\n", *port)
		fmt.Printf("  📋 OpenAPI Spec: http://localhost:%d/openapi.json\n", *port)
		fmt.Println("  按 Ctrl+C 停止\n")
		if err := swagger.Serve(*port, jsonData); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ 启动服务失败: %v\n", err)
			os.Exit(1)
		}
	}
}
