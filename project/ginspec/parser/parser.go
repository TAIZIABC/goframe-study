// parser/parser.go
// 解析 Go 源码，提取 Gin 路由注册和 Handler 函数注释
package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Route 表示一个解析出的路由
type Route struct {
	Method      string // GET, POST, PUT, DELETE...
	Path        string // /api/v1/users/:id
	HandlerName string // handler.CreateUser
	Summary     string // 从注释提取的摘要
	Description string // 从注释提取的描述
	Tags        []string
	SourceFile  string
	Line        int
}

// ParseProject 扫描项目目录，解析所有 .go 文件
func ParseProject(dir string) ([]Route, error) {
	var allRoutes []Route

	// 第一遍：收集所有函数的注释
	funcComments := make(map[string]commentInfo)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		comments := extractFuncComments(path)
		for k, v := range comments {
			funcComments[k] = v
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// 第二遍：扫描路由注册
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return err
		}
		routes := extractRoutes(path, funcComments)
		allRoutes = append(allRoutes, routes...)
		return nil
	})

	return allRoutes, err
}

type commentInfo struct {
	Summary     string
	Description string
	Tags        []string
}

// extractFuncComments 从文件中提取所有函数的注释
func extractFuncComments(filename string) map[string]commentInfo {
	result := make(map[string]commentInfo)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return result
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Doc == nil {
			continue
		}

		funcName := fn.Name.Name
		// 如果是方法 (有 receiver)，也记录
		if fn.Recv != nil && len(fn.Recv.List) > 0 {
			// 跳过 receiver 的方法名提取，直接用函数名
		}

		info := parseComment(fn.Doc.Text())
		result[funcName] = info
	}

	return result
}

// parseComment 解析注释文本，提取摘要、描述和标签
func parseComment(text string) commentInfo {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	info := commentInfo{}

	var descLines []string
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// 提取 @tag 标记
		if strings.HasPrefix(line, "@tag") || strings.HasPrefix(line, "@Tag") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				info.Tags = append(info.Tags, parts[1])
			}
			continue
		}

		// 第一行作为摘要
		if i == 0 && line != "" {
			info.Summary = line
			continue
		}

		if line != "" {
			descLines = append(descLines, line)
		}
	}

	info.Description = strings.Join(descLines, "\n")
	return info
}

// Gin HTTP 方法列表
var ginMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS", "Any"}

// 匹配路由注册: v1.GET("/users", handler.ListUsers) 或 r.POST("/users", CreateUser)
var routeRegex = regexp.MustCompile(
	`\.\s*(` + strings.Join(ginMethods, "|") + `)\s*\(\s*"([^"]+)"\s*,\s*([^,\)]+)`,
)

// 匹配路由组: r.Group("/api/v1")
var groupRegex = regexp.MustCompile(`\.Group\s*\(\s*"([^"]+)"`)

// extractRoutes 从文件中提取路由注册
func extractRoutes(filename string, funcComments map[string]commentInfo) []Route {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var routes []Route

	// 简单解析路由组前缀
	groupPrefixes := extractGroupPrefixes(content)

	for lineNum, line := range lines {
		matches := routeRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		method := strings.ToUpper(matches[1])
		if method == "ANY" {
			method = "*"
		}
		path := matches[2]
		handlerRef := strings.TrimSpace(matches[3])

		// 尝试匹配路由组前缀
		fullPath := resolveFullPath(line, path, groupPrefixes)

		// 提取 handler 函数名
		funcName := extractFuncName(handlerRef)

		// 查找注释
		info := funcComments[funcName]

		route := Route{
			Method:      method,
			Path:        fullPath,
			HandlerName: handlerRef,
			Summary:     info.Summary,
			Description: info.Description,
			Tags:        info.Tags,
			SourceFile:  filename,
			Line:        lineNum + 1,
		}

		// 如果没有 tag，从路径推断
		if len(route.Tags) == 0 {
			route.Tags = inferTags(fullPath)
		}

		routes = append(routes, route)
	}

	return routes
}

// extractGroupPrefixes 提取文件中所有 Group 定义及其变量名
func extractGroupPrefixes(content string) map[string]string {
	result := make(map[string]string)

	// 匹配: v1 := r.Group("/api/v1") 或 api := router.Group("/api")
	re := regexp.MustCompile(`(\w+)\s*[:=]+\s*\w+\.Group\s*\(\s*"([^"]+)"`)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, m := range matches {
		result[m[1]] = m[2]
	}

	// 匹配闭包形式: r.Group("/api/v1", func(...) { ... })
	// 这里简化处理，直接匹配所有 Group
	groupMatches := groupRegex.FindAllStringSubmatch(content, -1)
	if len(groupMatches) > 0 && len(result) == 0 {
		// 用 _default_ 作为 fallback key
		result["_default_"] = groupMatches[0][1]
	}

	return result
}

// resolveFullPath 根据调用者变量名解析完整路径
func resolveFullPath(line, path string, prefixes map[string]string) string {
	// 检查行中使用了哪个变量
	for varName, prefix := range prefixes {
		if varName == "_default_" {
			continue
		}
		if strings.Contains(line, varName+".") {
			return normalizePath(prefix + path)
		}
	}
	// fallback: 使用默认 group
	if prefix, ok := prefixes["_default_"]; ok {
		return normalizePath(prefix + path)
	}
	return normalizePath(path)
}

func normalizePath(p string) string {
	p = strings.ReplaceAll(p, "//", "/")
	// Gin 的 :id 参数转为 OpenAPI 的 {id}
	re := regexp.MustCompile(`:(\w+)`)
	p = re.ReplaceAllString(p, `{$1}`)
	return p
}

func extractFuncName(ref string) string {
	// "handler.CreateUser" → "CreateUser"
	parts := strings.Split(ref, ".")
	return parts[len(parts)-1]
}

func inferTags(path string) []string {
	// 从路径推断 tag: /api/v1/users/:id → Users
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i := len(parts) - 1; i >= 0; i-- {
		p := parts[i]
		if !strings.HasPrefix(p, "{") && p != "api" && p != "v1" && p != "v2" {
			return []string{strings.Title(p)}
		}
	}
	return []string{"Default"}
}
