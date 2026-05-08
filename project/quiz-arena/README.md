# quiz-arena · Go 并发知识挑战平台

前后端分离架构的在线答题平台，专注于 Go 语言并发知识（Goroutine、Channel、WaitGroup、Mutex、Context 等）。

## 架构

```
quiz-arena/
├── backend/              # Go 后端（纯 API 服务 + CORS）
│   ├── main.go
│   └── go.mod
├── frontend/             # Vite + React + Tailwind CSS
│   ├── src/
│   │   ├── App.jsx
│   │   ├── api.js        # API 客户端封装
│   │   ├── main.jsx
│   │   ├── index.css     # Tailwind 入口
│   │   └── components/
│   │       ├── Header.jsx
│   │       ├── StartScreen.jsx
│   │       ├── QuizScreen.jsx
│   │       ├── Options.jsx
│   │       ├── Feedback.jsx
│   │       └── ResultScreen.jsx
│   ├── index.html
│   ├── vite.config.js    # Vite 配置 + API 代理
│   ├── tailwind.config.js
│   ├── postcss.config.js
│   └── package.json
├── run.sh                # 一键启动脚本
└── README.md
```

## API 接口

| 方法 | 路径 | 说明 | 请求体 | 响应 |
|------|------|------|--------|------|
| GET | `/api/health` | 健康检查 | - | `{status, total}` |
| GET | `/api/questions` | 获取题目列表（不含答案） | - | `Question[]` |
| POST | `/api/submit` | 提交单题答案 | `{id, choice}` | `{correct, answer, explain}` |

### Question 结构

```json
{
  "id": 1,
  "category": "Goroutine",
  "title": "下面代码的输出是？",
  "code": "func main() { ... }",
  "options": ["...", "..."],
  "total": 20,
  "index": 0
}
```

## 开发运行

### 1. 启动后端（端口 8090）

```bash
cd backend
go run main.go
# 或编译后运行
go build -o quiz-api . && ./quiz-api
```

### 2. 启动前端（端口 5173）

```bash
cd frontend
npm install       # 首次运行
npm run dev
```

Vite 开发服务器会自动将 `/api/*` 代理到后端 `http://localhost:8090`，前后端分离开发。

### 3. 一键启动

```bash
./run.sh
```

## 生产构建

```bash
# 前端构建（输出到 frontend/dist）
cd frontend && npm run build

# 后端编译
cd backend && go build -o quiz-api .

# 部署：前端静态产物可用任何 Web 服务器托管，
# 通过反向代理将 /api/* 转发到 Go 后端
```

## 技术栈

**后端**
- Go 1.21+ 标准库
- `net/http` 提供 REST API
- CORS 中间件支持跨域

**前端**
- Vite 5 构建工具
- React 18 组件化
- Tailwind CSS 3 原子化样式
- 原生 `fetch` API

## 知识点覆盖

| 分类 | 题目数 | 代表知识点 |
|------|-------|-----------|
| Goroutine | 3 | 启动、调度模型、闭包陷阱 |
| Channel | 4 | 缓冲、关闭、range、发送者责任 |
| Select | 2 | 随机选择、超时模式 |
| WaitGroup | 2 | 值传递陷阱、Add 时机 |
| Mutex | 2 | 不可重入、死锁避免 |
| Atomic | 1 | 无锁优势 |
| Context | 3 | cancel 释放、参数传递、级联取消 |
| 综合 | 3 | 核心哲学、sync.Once |
