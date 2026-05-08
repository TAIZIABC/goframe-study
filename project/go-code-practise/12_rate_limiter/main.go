package main

import (
	"fmt"
	"sync"
	"time"
)

// 题目：实现简易令牌桶限流器
//
// NewRateLimiter(rate int, interval time.Duration) *RateLimiter
// - rate: interval 时间内允许的最大请求数
// Allow() bool — 允许请求返回 true，否则返回 false
//
// 示例：
//   rl := NewRateLimiter(3, time.Second)  // 每秒最多3次
//   rl.Allow() // true
//   rl.Allow() // true
//   rl.Allow() // true
//   rl.Allow() // false（已用完）

type RateLimiter struct {
	mu       sync.Mutex
	rate     int           // 最大令牌数
	tokens   int           // 当前令牌数
	interval time.Duration // 补充间隔
	lastTime time.Time     // 上次补充时间
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		rate:     rate,
		tokens:   rate,
		interval: interval,
		lastTime: time.Now(),
	}
}

func (r *RateLimiter) Allow() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 计算经过的时间，补充令牌
	now := time.Now()
	elapsed := now.Sub(r.lastTime)
	if elapsed >= r.interval {
		// 经过了几个完整周期就补充几轮（最多补满）
		refill := int(elapsed / r.interval) * r.rate
		r.tokens += refill
		if r.tokens > r.rate {
			r.tokens = r.rate
		}
		r.lastTime = now
	}

	// 尝试消费一个令牌
	if r.tokens > 0 {
		r.tokens--
		return true
	}
	return false
}

func main() {
	rl := NewRateLimiter(3, time.Second)
	fmt.Println(rl.Allow()) // true
	fmt.Println(rl.Allow()) // true
	fmt.Println(rl.Allow()) // true
	fmt.Println(rl.Allow()) // false

	time.Sleep(time.Second)
	fmt.Println(rl.Allow()) // true（令牌已恢复）
}
