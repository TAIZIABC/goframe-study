# Go 编程练习题

20 道日常开发高频场景编程练习，每道题目是独立的 `main.go` 文件，包含题目描述、示例和代码框架。

## 题目列表

| # | 题目 | 难度 | 考点 |
|---|------|------|------|
| 01 | 反转字符串 | ⭐ | []rune 处理、Unicode |
| 02 | 两数之和 | ⭐ | map 哈希查找 |
| 03 | 切片去重 | ⭐ | map + 保序 |
| 04 | 并发安全计数器 | ⭐⭐ | sync.Mutex / atomic |
| 05 | Worker Pool | ⭐⭐ | goroutine 池、channel |
| 06 | 超时控制 | ⭐⭐ | context.WithTimeout、select |
| 07 | JSON 解析与构建 | ⭐ | encoding/json、struct tag |
| 08 | 自定义错误 + 错误链 | ⭐⭐ | error 接口、errors.Is/As、%w |
| 09 | 泛型 Filter/Map/Reduce | ⭐⭐ | Go 泛型 |
| 10 | LRU 缓存 | ⭐⭐⭐ | 数据结构设计 |
| 11 | Pipeline 管道模式 | ⭐⭐ | channel 串联、goroutine |
| 12 | 令牌桶限流器 | ⭐⭐ | time.Ticker、并发控制 |
| 13 | 指数退避重试 | ⭐⭐ | time.Sleep、错误处理 |
| 14 | 字符串处理工具 | ⭐⭐ | strings/unicode、命名转换 |
| 15 | 多条件结构体排序 | ⭐⭐ | sort.Slice、多级排序 |
| 16 | Fan-Out/Fan-In | ⭐⭐⭐ | 并发分发与合并 |
| 17 | 优雅关闭任务管理器 | ⭐⭐⭐ | context、sync.WaitGroup、信号 |
| 18 | 表格格式化输出 | ⭐⭐ | 字符串对齐、fmt.Sprintf |
| 19 | 事件发射器 | ⭐⭐ | 观察者模式、函数切片 |
| 20 | 中间件链 | ⭐⭐⭐ | 洋葱模型、闭包、高阶函数 |

## 使用方式

```bash
cd project/go-code-practise/01_reverse_string
# 在 main.go 中编写代码
go run main.go
```

每道题的 `main()` 函数中已包含测试用例，实现函数后运行即可验证。

## 建议练习顺序

1. 先做 ⭐ 题热身（01-03, 07）
2. 再做 ⭐⭐ 题巩固（04-06, 08-09, 11-15, 18-19）
3. 最后挑战 ⭐⭐⭐ 题（10, 16-17, 20）
