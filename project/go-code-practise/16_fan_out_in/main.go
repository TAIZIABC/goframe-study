package main

import (
	"fmt"
	"sort"
	"sync"
)

// 题目：实现 Fan-Out/Fan-In 模式
//
// FanOut(input <-chan int, workers int) []<-chan int
// - 将 input 分发给 workers 个 channel
//
// FanIn(channels ...<-chan int) <-chan int
// - 将多个 channel 合并为一个
//
// 示例：输入 1~10，3 个 worker 各处理一部分，最终合并收集

func FanOut(input <-chan int, workers int) []<-chan int {
	channels := make([]<-chan int, workers)
	for i := 0; i < workers; i++ {
		ch := make(chan int)
		channels[i] = ch
		go func() {
			for v := range input {
				ch <- v
			}
			close(ch)
		}()
	}
	return channels
}

func FanIn(channels ...<-chan int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan int) {
			defer wg.Done()
			for v := range c {
				out <- v
			}
		}(ch)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	// 生成输入
	input := make(chan int, 10)
	for i := 1; i <= 10; i++ {
		input <- i
	}
	close(input)

	// Fan-Out 到 3 个 worker
	workers := FanOut(input, 3)

	// Fan-In 合并结果
	merged := FanIn(workers...)

	// 收集
	var results []int
	for v := range merged {
		results = append(results, v)
	}
	sort.Ints(results)
	fmt.Println(results) // [1 2 3 4 5 6 7 8 9 10]
}
