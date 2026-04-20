// generator/openapi.go
// 将解析出的路由转换为 OpenAPI 3.0 文档
package generator

import (
	"encoding/json"
	"sort"
	"strings"

	"ginspec/parser"

	"gopkg.in/yaml.v3"
)

// OpenAPI 3.0 结构体定义
type OpenAPI struct {
	OpenAPI string         `json:"openapi" yaml:"openapi"`
	Info    Info           `json:"info" yaml:"info"`
	Paths   map[string]any `json:"paths" yaml:"paths"`
	Tags    []Tag          `json:"tags,omitempty" yaml:"tags,omitempty"`
}

type Info struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
	Version     string `json:"version" yaml:"version"`
}

type Tag struct {
	Name string `json:"name" yaml:"name"`
}

type Operation struct {
	Summary     string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Tags        []string    `json:"tags,omitempty" yaml:"tags,omitempty"`
	OperationID string      `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	Responses   map[string]Response `json:"responses" yaml:"responses"`
}

type Parameter struct {
	Name     string `json:"name" yaml:"name"`
	In       string `json:"in" yaml:"in"`
	Required bool   `json:"required" yaml:"required"`
	Schema   Schema `json:"schema" yaml:"schema"`
}

type Schema struct {
	Type string `json:"type" yaml:"type"`
}

type Response struct {
	Description string `json:"description" yaml:"description"`
}

// Generate 生成 OpenAPI 文档
func Generate(routes []parser.Route, title, version string) *OpenAPI {
	doc := &OpenAPI{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:   title,
			Version: version,
		},
		Paths: make(map[string]any),
	}

	tagSet := make(map[string]bool)

	for _, route := range routes {
		path := route.Path

		// 获取或创建 path item
		pathItem, ok := doc.Paths[path]
		if !ok {
			pathItem = make(map[string]any)
			doc.Paths[path] = pathItem
		}
		pm := pathItem.(map[string]any)

		// 构建 operation
		op := Operation{
			Summary:     route.Summary,
			Description: route.Description,
			Tags:        route.Tags,
			OperationID: route.HandlerName,
			Parameters:  extractPathParams(path),
			Responses: map[string]Response{
				"200": {Description: "成功"},
			},
		}

		method := strings.ToLower(route.Method)
		if method == "*" {
			// Any → 注册所有方法
			for _, m := range []string{"get", "post", "put", "delete"} {
				pm[m] = op
			}
		} else {
			pm[method] = op
		}

		for _, t := range route.Tags {
			tagSet[t] = true
		}
	}

	// 排序 tags
	var tags []string
	for t := range tagSet {
		tags = append(tags, t)
	}
	sort.Strings(tags)
	for _, t := range tags {
		doc.Tags = append(doc.Tags, Tag{Name: t})
	}

	return doc
}

// extractPathParams 从路径中提取参数
func extractPathParams(path string) []Parameter {
	var params []Parameter
	parts := strings.Split(path, "/")
	for _, p := range parts {
		if strings.HasPrefix(p, "{") && strings.HasSuffix(p, "}") {
			name := strings.Trim(p, "{}")
			params = append(params, Parameter{
				Name:     name,
				In:       "path",
				Required: true,
				Schema:   Schema{Type: "string"},
			})
		}
	}
	return params
}

// ToJSON 输出 JSON 格式
func (doc *OpenAPI) ToJSON() ([]byte, error) {
	return json.MarshalIndent(doc, "", "  ")
}

// ToYAML 输出 YAML 格式
func (doc *OpenAPI) ToYAML() ([]byte, error) {
	return yaml.Marshal(doc)
}
