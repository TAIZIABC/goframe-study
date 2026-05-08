# Go 并发练习项目

本项目包含 Go 语言核心并发知识点的完整示例，适合学习和练习。

## 项目结构

```
chan-practise/
├── 01_goroutine/     # Goroutine 创建与同步
│   └── main.go       # 基础用法、多goroutine、闭包陷阱
├── 02_channel/       # Channel 通信机制
│   └── main.go       # 无缓冲/有缓冲、单向channel、生产者消费者
├── 03_select/        # Select 多路复用
│   └── main.go       # 基础select、非阻塞、超时、扇入、优雅退出
├── 04_sync/          # 并发原语
│   └── main.go       # WaitGroup、Mutex、RWMutex、Once、sync.Map、atomic
├── 05_context/       # Context 上下文控制
│   └── main.go       # Cancel、Timeout、Deadline、Value、级联取消
├── 06_patterns/      # 常见并发模式
│   └── main.go       # Pipeline、Fan-out/Fan-in、Worker Pool、信号量、速率限制
├── concurrency_test.go  # 测试 & 基准测试
├── run.sh            # 运行脚本
├── go.mod
└── README.md
```

## 知识点覆盖

| 模块 | 知识点 |
|------|--------|
| 01_goroutine | go 关键字、WaitGroup 等待、闭包陷阱、goroutine 数量 |
| 02_channel | 无缓冲/有缓冲 channel、单向 channel、close/range、生产者消费者 |
| 03_select | 多路复用、非阻塞操作、超时控制、扇入模式、优雅退出 |
| 04_sync | Mutex、RWMutex、WaitGroup、sync.Once、sync.Map、atomic |
| 05_context | WithCancel、WithTimeout、WithDeadline、WithValue、级联取消 |
| 06_patterns | Pipeline、Fan-out/Fan-in、Worker Pool、Semaphore、Rate Limiter |

## 运行方式

```bash
# 运行单个示例
go run ./01_goroutine/
go run ./02_channel/
go run ./03_select/
go run ./04_sync/
go run ./05_context/
go run ./06_patterns/

# 使用运行脚本
chmod +x run.sh
./run.sh all       # 运行全部
./run.sh 1         # 运行 Goroutine 示例
./run.sh channel   # 运行 Channel 示例
./run.sh test      # 运行单元测试
./run.sh race      # 竞态条件检测
./run.sh bench     # 性能基准测试
```

## 测试

```bash
# 单元测试
go test -v .

# 竞态检测（重要！）
go test -race -v .

# 基准测试
go test -bench=. -benchmem .
```
