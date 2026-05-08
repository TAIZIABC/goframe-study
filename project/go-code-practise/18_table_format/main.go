package main

import (
	"fmt"
	"strings"
)

// 题目：实现表格格式化输出
//
// FormatTable(headers []string, rows [][]string) string
// 输出对齐的 ASCII 表格：
//
// +--------+-----+-------+
// | Name   | Age | City  |
// +--------+-----+-------+
// | Alice  | 25  | BJ    |
// | Bob    | 30  | SH    |
// +--------+-----+-------+
//
// 要求：
// - 每列宽度取该列最宽的值（含 header）
// - 左对齐，右侧补空格

func FormatTable(headers []string, rows [][]string) string {
	// 计算每列最大宽度
	colWidths := make([]int, len(headers))
	for i, h := range headers {
		colWidths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// 构建分隔线 +--------+-----+
	var b strings.Builder
	buildSep := func() {
		b.WriteByte('+')
		for _, w := range colWidths {
			b.WriteString(strings.Repeat("-", w+2))
			b.WriteByte('+')
		}
		b.WriteByte('\n')
	}

	// 构建数据行 | Name   | Age |
	buildRow := func(cells []string) {
		b.WriteByte('|')
		for i, cell := range cells {
			b.WriteString(fmt.Sprintf(" %-*s |", colWidths[i], cell))
		}
		b.WriteByte('\n')
	}

	buildSep()
	buildRow(headers)
	buildSep()
	for _, row := range rows {
		buildRow(row)
	}
	buildSep()

	return b.String()
}

func main() {
	headers := []string{"Name", "Age", "City"}
	rows := [][]string{
		{"Alice", "25", "Beijing"},
		{"Bob", "30", "Shanghai"},
		{"Charlie", "28", "SZ"},
	}
	fmt.Println(FormatTable(headers, rows))
}
