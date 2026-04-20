# wspush — WebSocket 消息推送服务

基于频道（channel）的 WebSocket 消息推送服务，客户端订阅频道，服务端通过 HTTP 接口推送消息，附带 Web 测试页面。

## 快速启动

```bash
go mod tidy
go run main.go -port 8090
```

浏览器打开 `http://localhost:8090` 进入测试页面。

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-port` | `8090` | 服务端口 |

## 接口说明

### WebSocket 连接

```
ws://localhost:8090/ws
```

**客户端发送命令：**

```json
// 订阅频道
{"action": "subscribe", "channel": "news"}

// 取消订阅
{"action": "unsubscribe", "channel": "news"}
```

**服务端推送消息格式：**

```json
{"channel": "news", "data": "Hello World!", "timestamp": "2026-04-17 10:45:00"}
```

### HTTP 推送接口

```bash
# 向频道推送消息
curl -X POST http://localhost:8090/api/publish \
  -H "Content-Type: application/json" \
  -d '{"channel": "news", "data": "Breaking News!"}'

# 查看统计信息
curl http://localhost:8090/api/stats
```

| 方法 | 路径 | 说明 |
|------|------|------|
| `WS` | `/ws` | WebSocket 连接 |
| `POST` | `/api/publish` | 推送消息到指定频道 |
| `GET` | `/api/stats` | 在线人数和频道统计 |
| `GET` | `/` | Web 测试页面 |

## 项目结构

```
wspush/
├── main.go              # 入口 + 内嵌 HTML 测试页面
├── hub/hub.go           # 频道管理中心（连接/订阅/广播）
├── handler/handler.go   # WebSocket handler + HTTP 接口
├── go.mod
└── README.md
```

## 测试页面功能

内嵌的 Web 测试页面提供：

- **连接/断开** WebSocket
- **订阅/取消** 任意频道
- **发送消息** 通过 HTTP API 推送
- **实时显示** 收到的频道消息
- 深色主题 UI

## 功能特点

- **频道模型**：客户端按需订阅，只收到已订阅频道的消息
- **并发安全**：Hub 单协程事件循环，避免锁竞争
- **缓冲区**：每个客户端 64 消息缓冲，防止慢客户端阻塞
- **自动清理**：客户端断开时自动从所有频道移除
- **终端日志**：实时显示连接/断开/订阅/推送事件
- **统计接口**：查看在线人数和各频道订阅数
