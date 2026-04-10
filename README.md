# 🚀 GoFrame 学习之路

> 专为前端开发者设计的 Go & GoFrame 交互式学习指南

[![GitHub](https://img.shields.io/badge/GitHub-TAIZIABC-blue?logo=github)](https://github.com/TAIZIABC/goframe-study)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

从 **JavaScript** 到 **Go**，以前端开发者熟悉的视角，系统学习 Go 语言与 GoFrame 框架。

![GoFrame 学习之路](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)
![GoFrame](https://img.shields.io/badge/GoFrame-企业级框架-00838F?style=for-the-badge)

---

## ✨ 特性

- 🎯 **JS ↔ Go 对比学习** — 用你最熟悉的 JavaScript 对比讲解，快速理解 Go 核心概念
- 🖥️ **交互式代码编辑器** — 内置 Go Playground，在线编写、运行 Go 代码，即时查看结果
- 📊 **学习进度追踪** — 自动记录已完成章节，可视化学习进度
- 🌙 **深色/浅色主题** — 支持主题切换，舒适阅读体验
- 📱 **响应式设计** — 完美适配桌面端和移动端
- 🗂️ **实战项目指南** — 5 个从简单到复杂的完整项目实现教程

## 📚 内容大纲

### 第一阶段：Go 语言基础（7 个核心章节）

| 章节 | 内容 | 说明 |
|------|------|------|
| 📦 变量与类型 | 静态类型、类型推断、零值机制 | 对比 JS/TS 类型系统 |
| ⚡ 函数与方法 | 多返回值、命名返回值、方法接收者 | 对比 JS 函数 |
| 🏗️ 结构体 | 结构体定义、组合、JSON 标签 | 对比 JS class |
| 🛡️ 错误处理 | error 接口、自定义错误、错误包装 | 对比 try/catch |
| 🔄 并发编程 | Goroutine、Channel、WaitGroup | 对比 async/await |
| 🔌 接口 | 隐式实现、空接口、类型断言 | 对比 TS interface |
| 🎯 指针 | 指针基础、引用传递、方法接收者选择 | JS 中无对应概念 |

### 第二阶段：GoFrame 框架（8 个核心章节）

| 章节 | 内容 |
|------|------|
| 🏗️ 项目结构 | GoFrame 工程化目录结构 vs Express 项目 |
| 🛣️ 路由注册 | 规范路由、分组路由、RESTful 设计 |
| 🔗 中间件 | 全局/路由中间件，对比 Express middleware |
| 📝 请求与响应 | 结构化参数绑定、数据校验、JSON 响应 |
| 💾 ORM 数据库 | Model 操作、链式查询、事务处理 |
| ✅ 数据校验 | 内置校验规则、自定义校验 |
| ⚙️ 配置管理 | YAML 配置、多环境管理 |
| 📋 日志管理 | 结构化日志、日志级别、文件输出 |

### 第三阶段：实战项目（5 个项目）

| 项目 | 难度 | 技术要点 |
|------|------|----------|
| ✅ Todo List API | 入门 | CRUD、路由注册、请求校验 |
| 🔐 用户认证系统 | 中级 | JWT、中间件、密码加密 |
| 📰 博客系统 | 中高级 | ORM 关联、分页、Redis 缓存、Swagger |
| 🛒 电商 API 平台 | 高级 | 库存扣减、订单状态机、支付、Docker 部署 |
| 💬 实时聊天应用 | 高级 | WebSocket、Hub 模式、Redis Pub/Sub、消息持久化 |

## 🚀 快速开始

### 在线访问

直接用浏览器打开 `index.html` 即可使用，无需安装任何依赖。

### 本地运行

```bash
# 克隆仓库
git clone https://github.com/TAIZIABC/goframe-study.git
cd goframe-study

# 方式一：直接打开
open index.html

# 方式二：使用本地服务器（推荐，支持代码在线运行）
python3 -m http.server 8080
# 然后访问 http://localhost:8080
```

## 🗂️ 项目结构

```
goframe-study/
├── index.html      # 页面结构
├── styles.css      # 样式（支持深色/浅色主题）
├── app.js          # 核心逻辑（课程数据、渲染、交互）
├── index.js        # 辅助脚本
├── todo-api/       # 实战项目：Todo List API（可独立运行）
│   ├── manifest/sql/init.sql  # 数据库初始化脚本
│   └── README.MD              # 项目说明与运行指南
├── .gitignore      # Git 忽略配置
└── README.md       # 项目说明
```

## 🎨 功能预览

### 交互式代码对比

左右对比展示 JavaScript 和 Go 的写法差异，直观理解语言特性。

### 在线代码编辑器

每个章节都内置 Go 代码编辑器，连接 Go Playground 实时编译运行。

### 实战项目指南

每个项目都包含完整的分步实现教程，包含数据库设计、业务逻辑、路由配置等。

## 🛠️ 技术栈

- **纯前端实现** — HTML + CSS + JavaScript，零依赖
- **Go Playground API** — 在线编译运行 Go 代码
- **LocalStorage** — 持久化学习进度和主题偏好
- **CSS 变量** — 实现深色/浅色主题切换
- **响应式布局** — Flexbox + Grid 适配多端

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m 'feat: 添加新特性'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 提交 Pull Request

## 📄 License

本项目采用 [MIT License](LICENSE) 开源协议。

---

<p align="center">
  <b>为前端开发者而建 ❤️</b><br>
  如果对你有帮助，请给个 ⭐ Star 支持！
</p>
