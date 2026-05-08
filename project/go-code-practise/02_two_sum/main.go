package main

import "fmt"

// 题目：实现 TwoSum，找到切片中两个数之和等于 target 的下标
// 返回 [2]int{i, j}，保证 i < j。找不到返回 [2]int{-1, -1}
// 示例：TwoSum([]int{2, 7, 11, 15}, 9) → [0, 1]

func TwoSum(nums []int, target int) [2]int {
	// 在此编写你的代码
	m := make(map[int]int)
	for i, num := range nums {
		if _, ok := m[target-num]; ok {
			return [2]int{m[target-num], i}
		}
		m[num] = i
	}

	return [2]int{-1, -1}
}

func main() {
	fmt.Println(TwoSum([]int{2, 7, 11, 15}, 9)) // [0 1]
	fmt.Println(TwoSum([]int{3, 2, 4}, 6))      // [1 2]
	fmt.Println(TwoSum([]int{1, 2, 3}, 10))     // [-1 -1]
}
