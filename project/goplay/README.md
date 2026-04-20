# goplay — 在线 Go 代码运行服务

接收 Go 代码片段，在隔离临时目录中编译运行，限制执行时间和内存，返回输出结果，附带 Web 代码编辑器。

## 快速启动

```bash
go run main.go -port 8092
```

浏览器打开 `http://localhost:8092` 进入代码编辑器。

## 参数

| 参数 | 默认值 | 说明 |
|------|--------|------|
| `-port` | `8092` | 服务端口 |

## API 接口

### POST /api/run

提交代码执行：

```bash
curl -X POST http://localhost:8092/api/run \
  -H "Content-Type: application/json" \
  -d '{"code":"package main\nimport \"fmt\"\nfunc main(){fmt.Println(\"Hello!\")}"}'
```

响应：

```json
{
  "stdout": "Hello!\n",
  "stderr": "",
  "exit_code": 0,
  "duration_ms": 2926,
  "success": true
}
```

## 项目结构

```
goplay/
├── main.go              # 入口 + HTTP 路由 + 内嵌编辑器页面
├── runner/
│   └── runner.go        # 代码沙箱（隔离目录/编译/运行/超时/安全检查）
├── go.mod
└── README.md
```

## Web 编辑器功能

内嵌暗色主题双栏编辑器：

- 左侧：代码编辑区（等宽字体、Tab 缩进、字符计数）
- 右侧：输出区（stdout 绿色、stderr 红色、错误橙色）
- **Ctrl+Enter** 快捷运行
- 运行时 Loading 动画

## 安全机制

| 措施 | 说明 |
|------|------|
| **临时目录隔离** | 每次运行创建独立目录，执行后自动清理 |
| **执行超时** | 10 秒限制，超时自动 kill |
| **内存限制** | `GOMEMLIMIT=128MiB` |
| **包黑名单** | 禁止 `os/exec`、`syscall`、`unsafe`、`net/http`、`plugin` |
| **代码大小** | 最大 100KB |
| **输出限制** | 最大 64KB，超出截断 |

## 功能特点

- 隔离临时目录，自动 `go mod init`，运行后 `os.RemoveAll` 清理
- `context.WithTimeout` 超时控制
- 编译错误精确报出行号（自动清理临时路径）
- 输出限流防止无限循环刷屏
- 内嵌 Web 编辑器，零依赖打开即用
