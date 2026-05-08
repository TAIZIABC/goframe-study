package main

import (
	"fmt"
	"strings"
	"unicode"
)

// 题目：实现一组字符串处理函数
//
// 1. CamelToSnake(s string) string — 驼峰转蛇形
//    "GoFrame" → "go_frame", "helloWorld" → "hello_world"
//
// 2. SnakeToCamel(s string) string — 蛇形转大驼峰
//    "go_frame" → "GoFrame", "hello_world" → "HelloWorld"
//
// 3. Truncate(s string, maxLen int) string — 截断字符串（按字符数）
//    超出 maxLen 时末尾加 "..."
//    Truncate("Hello, 世界", 7) → "Hello, ..."

func CamelToSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				b.WriteByte('_')
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func SnakeToCamel(s string) string {
	var b strings.Builder
	upper := true // 每个单词首字母大写
	for _, r := range s {
		if r == '_' {
			upper = true
			continue
		}
		if upper {
			b.WriteRune(unicode.ToUpper(r))
			upper = false
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func Truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

func main() {
	fmt.Println(CamelToSnake("GoFrame"))    // go_frame
	fmt.Println(CamelToSnake("helloWorld")) // hello_world

	fmt.Println(SnakeToCamel("go_frame"))    // GoFrame
	fmt.Println(SnakeToCamel("hello_world")) // HelloWorld

	fmt.Println(Truncate("Hello, 世界!", 7)) // Hello, ...
	fmt.Println(Truncate("Hi", 10))        // Hi
}
