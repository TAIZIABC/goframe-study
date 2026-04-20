# jsonfmt — JSON 格式化工具

从 stdin 或文件读取 JSON，输出格式化或压缩后的 JSON，支持自定义缩进，非法 JSON 精确报出错误位置。

## 使用方式

```bash
# 从 stdin 格式化
echo '{"a":1,"b":[2,3]}' | go run main.go

# 从文件读取
go run main.go data.json

# 指定 4 空格缩进
echo '{"a":1}' | go run main.go -indent 4

# 压缩模式
go run main.go -compact data.json

# 多文件批量处理
go run main.go a.json b.json c.json
```

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-indent` | `2` | 缩进空格数 |
| `-compact` | `false` | 压缩模式（去除所有空白） |

## 错误定位示例

```
✗ stdin: JSON 语法错误 (第 3 行, 第 9 列)
  错误: invalid character ',' looking for beginning of value
  位置: ..."a": 1,\n  "b": ,⚡\n  "c": 3...
```

## 功能特点

- 支持 **stdin** 和**文件**输入，可同时处理多个文件
- 自定义缩进空格数 / 压缩模式
- 非法 JSON 精确报出**行号、列号**
- `⚡` 标记出错位置，显示前后上下文
- 非零退出码（`exit 1`），方便脚本集成
