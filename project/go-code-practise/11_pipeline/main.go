package main

import "fmt"

// 题目：实现 Pipeline 管道模式
//
// 1. Generate(nums ...int) <-chan int — 将数字发送到 channel 后关闭
// 2. Square(in <-chan int) <-chan int — 读入并平方后输出
// 3. Filter(in <-chan int, fn func(int) bool) <-chan int — 过滤
// 4. Collect(in <-chan int) []int — 收集所有结果
//
// 示例：
//   pipe := Square(Generate(1,2,3,4,5))
//   filtered := Filter(pipe, func(n int) bool { return n > 5 })
//   Collect(filtered) → [9, 16, 25]

func Generate(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func Square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func FilterChan(in <-chan int, fn func(int) bool) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if fn(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func Collect(in <-chan int) []int {
	var result []int
	for n := range in {
		result = append(result, n)
	}
	return result
}

func main() {
	pipe := Square(Generate(1, 2, 3, 4, 5))
	filtered := FilterChan(pipe, func(n int) bool { return n > 5 })
	result := Collect(filtered)
	fmt.Println(result) // [9 16 25]
}
