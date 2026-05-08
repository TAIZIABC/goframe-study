package main

import "fmt"

// 题目：实现 Deduplicate，对整数切片去重并保持原始顺序
// 示例：Deduplicate([]int{3, 1, 2, 1, 3, 4, 2}) → [3, 1, 2, 4]

func Deduplicate(s []int) []int {
	// 在此编写你的代码
	if s == nil {
		return nil
	}
	m := make(map[int]bool)
	result := []int{}
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}

	return result
}

func main() {
	fmt.Println(Deduplicate([]int{3, 1, 2, 1, 3, 4, 2})) // [3 1 2 4]
	fmt.Println(Deduplicate([]int{1, 1, 1}))             // [1]
	fmt.Println(Deduplicate(nil))                        // []
}
