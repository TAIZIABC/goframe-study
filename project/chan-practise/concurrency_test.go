package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

// ==================== 竞态条件测试 ====================
// 运行命令：go test -race -v ./...
// -race 标志会检测竞态条件

// TestMutexCounter 测试互斥锁保护的计数器
func TestMutexCounter(t *testing.T) {
	var (
		counter int
		mu      sync.Mutex
		wg      sync.WaitGroup
	)

	const n = 10000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}
	wg.Wait()

	if counter != n {
		t.Errorf("期望 %d, 得到 %d", n, counter)
	}
}

// TestAtomicCounter 测试原子操作计数器
func TestAtomicCounter(t *testing.T) {
	var counter int64
	var wg sync.WaitGroup

	const n = 10000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}
	wg.Wait()

	if counter != int64(n) {
		t.Errorf("期望 %d, 得到 %d", n, counter)
	}
}

// TestChannelCommunication 测试 channel 通信正确性
func TestChannelCommunication(t *testing.T) {
	ch := make(chan int, 100)
	var wg sync.WaitGroup

	// 生产者
	const numItems = 100
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < numItems; i++ {
			ch <- i
		}
		close(ch)
	}()

	// 消费者
	var received int
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range ch {
			received++
		}
	}()

	wg.Wait()
	if received != numItems {
		t.Errorf("期望接收 %d 条消息, 实际 %d", numItems, received)
	}
}

// TestWorkerPool 测试工作池正确性
func TestWorkerPool(t *testing.T) {
	const (
		numWorkers = 5
		numJobs    = 50
	)

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// 启动 worker
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				results <- j * 2
			}
		}()
	}

	// 发送任务
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	// 验证结果数量
	count := 0
	for range results {
		count++
	}
	if count != numJobs {
		t.Errorf("期望 %d 个结果, 得到 %d", numJobs, count)
	}
}

// ==================== 性能基准测试 ====================

// BenchmarkMutexCounter Mutex 计数器性能
func BenchmarkMutexCounter(b *testing.B) {
	var counter int
	var mu sync.Mutex

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			mu.Lock()
			counter++
			mu.Unlock()
		}
	})
}

// BenchmarkAtomicCounter 原子操作计数器性能
func BenchmarkAtomicCounter(b *testing.B) {
	var counter int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			atomic.AddInt64(&counter, 1)
		}
	})
}

// BenchmarkRWMutexRead 读写锁读性能
func BenchmarkRWMutexRead(b *testing.B) {
	var data int
	var rwmu sync.RWMutex
	data = 42

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rwmu.RLock()
			_ = data
			rwmu.RUnlock()
		}
	})
}

// BenchmarkChannelPingPong Channel 通信性能
func BenchmarkChannelPingPong(b *testing.B) {
	ch := make(chan struct{})

	go func() {
		for {
			<-ch
			ch <- struct{}{}
		}
	}()

	for i := 0; i < b.N; i++ {
		ch <- struct{}{}
		<-ch
	}
}
