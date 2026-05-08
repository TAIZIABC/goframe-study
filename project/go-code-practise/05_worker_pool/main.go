package main

import (
	"fmt"
	"sort"
	"sync"
)

// 题目：实现 Worker Pool
// WorkerPool(jobs []int, workers int) []int
// - jobs 中每个 int 代表需要计算平方的数字
// - 启动 workers 个 goroutine 并发处理
// - 返回所有结果（顺序不限）
// 示例：WorkerPool([]int{1,2,3,4,5}, 3) → [1,4,9,16,25]（顺序可不同）

func WorkerPool(jobs []int, workers int) []int {
	jobCn := make(chan int, len(jobs))
	res := make(chan int, len(jobs))
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobCn {
				res <- job * job
			}
		}()
	}
	for _, job := range jobs {
		jobCn <- job
	}
	close(jobCn)
	wg.Wait()
	close(res)
	result := make([]int, 0)
	for n := range res {
		result = append(result, n)
	}
	return result
}

func main() {
	result := WorkerPool([]int{1, 2, 3, 4, 5, 6, 7, 8}, 3)
	sort.Ints(result)
	fmt.Println(result) // [1 4 9 16 25 36 49 64]

	_ = sync.WaitGroup{}
}
