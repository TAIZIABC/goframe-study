package main

import "fmt"

// 题目：实现泛型函数 Filter、Map、Reduce
//
// 1. Filter[T any](s []T, fn func(T) bool) []T — 过滤
// 2. Map[T, U any](s []T, fn func(T) U) []U — 映射
// 3. Reduce[T, U any](s []T, init U, fn func(U, T) U) U — 归约
//
// 示例：
//   Filter([]int{1,2,3,4,5}, func(n int) bool { return n%2==0 }) → [2,4]
//   Map([]int{1,2,3}, func(n int) string { return fmt.Sprint(n*10) }) → ["10","20","30"]
//   Reduce([]int{1,2,3,4}, 0, func(acc, n int) int { return acc+n }) → 10

func Filter[T any](s []T, fn func(T) bool) []T {
	var result []T
	for _, v := range s {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func Map[T, U any](s []T, fn func(T) U) []U {
	result := make([]U, len(s))
	for i, v := range s {
		result[i] = fn(v)
	}
	return result
}

func Reduce[T, U any](s []T, init U, fn func(U, T) U) U {
	acc := init
	for _, v := range s {
		acc = fn(acc, v)
	}
	return acc
}

func main() {
	evens := Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 0 })
	fmt.Println(evens) // [2 4]

	strs := Map([]int{1, 2, 3}, func(n int) string { return fmt.Sprint(n * 10) })
	fmt.Println(strs) // [10 20 30]

	sum := Reduce([]int{1, 2, 3, 4}, 0, func(acc, n int) int { return acc + n })
	fmt.Println(sum) // 10
}
