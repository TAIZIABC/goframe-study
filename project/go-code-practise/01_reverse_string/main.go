package main

import (
	"fmt"
	"slices"
	"strings"
)

// 题目：实现 ReverseString，反转字符串（支持中文等多字节字符）
// 示例：ReverseString("Hello, 世界") → "界世 ,olleH"

func ReverseString(s string) string {
	// 在此编写你的代码
	arr := strings.Split(s, "")
	slices.Reverse(arr)
	return strings.Join(arr, "")
}

func main() {
	fmt.Println(ReverseString("Hello, 世界")) // 界世 ,olleH
	fmt.Println(ReverseString("abcde"))     // edcba
	fmt.Println(ReverseString(""))          // ""
}
