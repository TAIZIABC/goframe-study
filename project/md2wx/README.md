# md2wx — Markdown 转微信公众号 HTML

将 Markdown 文件转换为微信公众号编辑器兼容的 HTML 格式（所有 CSS 样式内联），输出可直接粘贴到公众号编辑器。

## 使用方式

```bash
# 转换文件
go run main.go -file article.md -out article.html

# 生成完整预览页面（含 HTML 包裹，可浏览器打开）
go run main.go -file article.md -preview -out preview.html

# 管道输入
cat article.md | go run main.go > output.html

# 直接传文件名
go run main.go article.md
```

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-file` | - | Markdown 文件路径 |
| `-out` | - | 输出 HTML 文件（不指定则输出到终端） |
| `-preview` | `false` | 生成完整 HTML 页面，可直接浏览器预览 |

## 支持的 Markdown 语法

| 语法 | 效果 |
|------|------|
| `# 标题` | H1-H6 标题（H1/H2 带底线） |
| `**粗体**` | 加粗 |
| `*斜体*` | 斜体 |
| `` `代码` `` | 行内代码（红色高亮） |
| ```` ```code``` ```` | 代码块（暗色背景） |
| `> 引用` | 引用块（灰色背景+左边线） |
| `- 列表` | 无序/有序列表 |
| `[链接](url)` | 链接（微信蓝色） |
| `![图片](url)` | 图片 |
| 表格 | 带边框表格 |
| `---` | 分隔线 |
| `~~删除线~~` | 删除线 |
| `- [x] 任务` | 任务列表 |

## 项目结构

```
md2wx/
├── main.go              # CLI 入口
├── render/
│   ├── render.go        # 微信 HTML 渲染器（内联 CSS）
│   └── table.go         # 表格渲染扩展
├── go.mod
└── README.md
```

## 为什么需要内联 CSS？

微信公众号编辑器会**过滤掉所有 `<style>` 标签和 `class` 属性**，只保留 `style` 内联样式。本工具的渲染器将所有样式直接写在 HTML 标签的 `style` 属性中，确保粘贴到公众号后样式不丢失。

## 使用流程

1. 用 Markdown 写好文章
2. 运行 `md2wx -file article.md -out article.html`
3. 用浏览器打开 `article.html`，全选复制
4. 粘贴到微信公众号编辑器，样式完美保留

## 功能特点

- **内联 CSS**：所有样式写在 style 属性，公众号兼容
- **代码高亮**：行内代码红色高亮，代码块暗色主题
- **表格支持**：带边框和表头背景色
- **引用块**：灰色背景 + 左侧蓝色边线
- **链接颜色**：微信蓝（#576b95）
- **预览模式**：`-preview` 生成完整页面，浏览器直接查看
- **管道支持**：支持 stdin 输入
